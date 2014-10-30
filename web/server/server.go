package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/carbocation/interpose"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/stretchr/graceful"

	"github.com/chillaxio/chillax/portkeeper"
	chillax_web_handlers "github.com/chillaxio/chillax/web/handlers"
	chillax_web_middlewares "github.com/chillaxio/chillax/web/middlewares"
	chillax_web_multiplexer "github.com/chillaxio/chillax/web/multiplexer"
	chillax_web_pingers "github.com/chillaxio/chillax/web/pingers"
	chillax_web_pipelines "github.com/chillaxio/chillax/web/pipelines"
	chillax_web_settings "github.com/chillaxio/chillax/web/settings"
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
	server.Paths["AdminProxy"] = server.Paths["AdminProxies"] + "/{Name}"
	server.Paths["AdminPipelines"] = server.Paths["AdminPrefix"] + "/pipelines"
	server.Paths["AdminPipeline"] = server.Paths["AdminPipelines"] + "/{Id}"

	server.Handler = server.NewGorillaMux()
	server.Middleware = server.NewInterposeMiddleware()

	server.SetDefaultMiddlewaresAfterInitialize()

	return server, err
}

// Server struct
type Server struct {
	*graceful.Server

	Settings   *chillax_web_settings.ServerSettings
	Middleware *interpose.Middleware
	Paths      map[string]string
}

func (s *Server) NewInterposeMiddleware() *interpose.Middleware {
	return interpose.New()
}

func (s *Server) SetDefaultMiddlewaresAfterInitialize() {
	s.Middleware.UseHandler(http.HandlerFunc(chillax_web_middlewares.ServerNameMiddleware(s.Settings)))
	s.Middleware.UseHandler(http.HandlerFunc(chillax_web_middlewares.BeginRequestTimerMiddleware(s.Settings)))
}

func (s *Server) SetDefaultMiddlewaresBeforeHttpServe() {
	s.Middleware.UseHandler(http.HandlerFunc(chillax_web_middlewares.RecordRequestTimerMiddleware(s.Settings)))
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
		s.Paths["AdminProxy"],
		chillax_web_handlers.AdminProxyHandler(s.Settings)).Methods("GET")

	mux.HandleFunc(
		s.Paths["AdminPipelines"],
		chillax_web_handlers.AdminPipelinesHandler(s.Settings)).Methods("GET")

	mux.HandleFunc(
		s.Paths["AdminPipeline"],
		chillax_web_handlers.AdminPipelineHandler(s.Settings)).Methods("GET")

	return mux
}

// BeforeListenAndServeGeneric runs background responsibilities.
func (s *Server) BeforeListenAndServeGeneric() {
	s.RunAllInProgressPipelinesAsync()
	s.CheckProxiesAsync()

	// Clean reserved ports every 5 minutes.
	// This value is hard-coded for now.
	s.CleanReservedPortsAsync("5m")

	// Wrap mux inside middleware before launching server.
	s.SetDefaultMiddlewaresBeforeHttpServe()
	s.Middleware.UseHandler(s.Handler)
	s.Handler = s.Middleware
}

// ListenAndServeGeneric runs the server.
func (s *Server) ListenAndServeGeneric() {
	s.BeforeListenAndServeGeneric()

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

func (s *Server) CleanReservedPortsAsync(sleepString string) {
	portkeeper.CleanReservedPortsAsync(sleepString)
}
