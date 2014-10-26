package server

import (
	"fmt"
	"net/http"
	"time"

	chillax_web_handlers "github.com/chillaxio/chillax/web/handlers"
	chillax_web_multiplexer "github.com/chillaxio/chillax/web/multiplexer"
	chillax_web_pingers "github.com/chillaxio/chillax/web/pingers"
	chillax_web_pipelines "github.com/chillaxio/chillax/web/pipelines"
	chillax_web_settings "github.com/chillaxio/chillax/web/settings"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/stretchr/graceful"
)

// NewServer is the constructor for creating Chillax server.
func NewServer() (*Server, error) {
	settings, err := chillax_web_settings.NewServerSettings()
	if err != nil {
		return nil, err
	}

	requestTimeoutOnRestart, err := time.ParseDuration(settings.RequestTimeoutOnRestart)
	if err != nil {
		settings.RequestTimeoutOnRestart = "3s"
		requestTimeoutOnRestart, _ = time.ParseDuration(settings.RequestTimeoutOnRestart)
	}

	server := &Server{
		Settings: settings,
		Server: &graceful.Server{
			Timeout: requestTimeoutOnRestart,
			Server:  &http.Server{Addr: settings.HttpAddress()},
		},
	}

	server.Paths = make(map[string]string)

	server.Paths["ApiPrefix"] = "/chillax/api"
	server.Paths["ApiProxies"] = server.Paths["ApiPrefix"] + "/proxies"
	server.Paths["ApiPipelines"] = server.Paths["ApiPrefix"] + "/pipelines"
	server.Paths["ApiPipelinesRun"] = server.Paths["ApiPrefix"] + "/pipelines/run"
	server.Paths["ApiPipelineRun"] = server.Paths["ApiPrefix"] + "/pipelines/{Id}/run"

	server.Paths["AdminPrefix"] = "/chillax/admin"
	server.Paths["AdminProxies"] = server.Paths["AdminPrefix"] + "/proxies"
	server.Paths["AdminPipelines"] = server.Paths["AdminPrefix"] + "/pipelines"

	server.Handler = server.NewGorillaMux()

	return server, err
}

// Server struct
type Server struct {
	*graceful.Server

	Settings *chillax_web_settings.ServerSettings
	Paths    map[string]string
}

// NewGorillaMux creates a multiplexer will all the correct endpoints as well as admin pages.
func (s *Server) NewGorillaMux() *gorilla_mux.Router {
	muxFactory := chillax_web_multiplexer.NewMuxFactory(s.Settings.ProxyHandlerTomls)

	muxFactory.CreateProxyBackends()
	muxFactory.StartProxyBackends()
	mux := muxFactory.GorillaMuxWithProxyBackends()

	// API Handlers
	mux.HandleFunc(
		s.Paths["ApiProxies"],
		chillax_web_handlers.ApiProxiesHandler(s.Settings)).Methods("POST")

	mux.HandleFunc(
		s.Paths["ApiPipelines"],
		chillax_web_handlers.ApiPipelinesHandler(s.Settings)).Methods("POST")

	mux.HandleFunc(
		s.Paths["ApiPipelinesRun"],
		chillax_web_handlers.ApiPipelinesRunHandler(s.Settings)).Methods("POST")

	mux.HandleFunc(
		s.Paths["ApiPipelineRun"],
		chillax_web_handlers.ApiPipelineRunHandler(s.Settings)).Methods("POST")

	// Admin Handlers
	mux.HandleFunc(
		s.Paths["AdminPrefix"],
		chillax_web_handlers.AdminBaseHandler(s.Settings)).Methods("GET")

	mux.HandleFunc(
		s.Paths["AdminProxies"],
		chillax_web_handlers.AdminProxiesHandler(s.Settings, muxFactory.ProxyHandlers)).Methods("GET")

	mux.HandleFunc(
		s.Paths["AdminPipelines"],
		chillax_web_handlers.AdminPipelinesHandler(s.Settings)).Methods("GET")

	return mux
}

// ListenAndServeGeneric runs the server.
func (s *Server) ListenAndServeGeneric() {
	if s.Settings.CertFile != "" && s.Settings.KeyFile != "" {
		s.ListenAndServeTLS(s.Settings.CertFile, s.Settings.KeyFile)
	} else {
		s.ListenAndServe()
	}
}

// RunAllInProgressPipelinesAsync loads all in-progress pipelines.
// This method is used in the event of server crash.
func (s *Server) RunAllInProgressPipelinesAsync() {
	numGoroutinesForCrashedInProgressPipelines := 50 // Hard-coded for now
	chillax_web_pipelines.RunAllInProgressPipelinesAsync(numGoroutinesForCrashedInProgressPipelines)
}

// CheckProxiesAsync hits every proxy endpoints.
func (s *Server) CheckProxiesAsync() {
	muxFactory := chillax_web_multiplexer.NewMuxFactory(s.Settings.ProxyHandlerTomls)

	proxyUris := make([]string, len(muxFactory.ProxyHandlers))
	for i, proxyHandler := range muxFactory.ProxyHandlers {
		if proxyHandler.Backend.Domain != "" {
			proxyUris[i] = fmt.Sprintf("http://%v%v", proxyHandler.Backend.Domain, proxyHandler.Backend.Path)
		} else {
			proxyUris[i] = fmt.Sprintf("http://localhost:%v%v", s.Settings.HttpPort, proxyHandler.Backend.Path)
		}
	}

	pingerGroup := chillax_web_pingers.NewPingerGroup(proxyUris)
	pingerGroup.IsUpAsync()
}
