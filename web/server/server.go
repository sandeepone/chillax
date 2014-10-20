package server

import (
	chillax_proxy_muxproducer "github.com/didip/chillax/proxy/muxproducer"
	chillax_web_handlers "github.com/didip/chillax/web/handlers"
	chillax_web_settings "github.com/didip/chillax/web/settings"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/stretchr/graceful"
	"net/http"
	"time"
)

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
	server.Paths["ApiPipelineRun"] = server.Paths["ApiPrefix"] + "/pipelines/{Id}/run"

	server.Paths["AdminPrefix"] = "/chillax/admin"
	server.Paths["AdminProxies"] = server.Paths["AdminPrefix"] + "/proxies"
	server.Paths["AdminPipelines"] = server.Paths["AdminPrefix"] + "/pipelines"

	server.Handler = server.NewGorillaMux()

	return server, err
}

type Server struct {
	*graceful.Server

	Settings *chillax_web_settings.ServerSettings
	Paths    map[string]string
}

func (s *Server) NewGorillaMux() *gorilla_mux.Router {
	muxProducer := chillax_proxy_muxproducer.NewMuxProducer(s.Settings.ProxyHandlerTomls)

	muxProducer.CreateProxyBackends()
	muxProducer.StartProxyBackends()
	mux := muxProducer.GorillaMuxWithProxyBackends()

	// API Handlers
	mux.HandleFunc(
		s.Paths["ApiProxies"],
		chillax_web_handlers.ApiProxiesHandler(s.Settings)).Methods("POST")

	mux.HandleFunc(
		s.Paths["ApiPipelines"],
		chillax_web_handlers.ApiPipelinesHandler(s.Settings)).Methods("POST")

	mux.HandleFunc(
		s.Paths["ApiPipelineRun"],
		chillax_web_handlers.ApiPipelineRunHandler(s.Settings)).Methods("POST")

	// Admin Handlers
	mux.HandleFunc(
		s.Paths["AdminPrefix"],
		chillax_web_handlers.AdminBaseHandler(s.Settings)).Methods("GET")

	mux.HandleFunc(
		s.Paths["AdminProxies"],
		chillax_web_handlers.AdminProxiesHandler(s.Settings, muxProducer.ProxyHandlers)).Methods("GET")

	mux.HandleFunc(
		s.Paths["AdminPipelines"],
		chillax_web_handlers.AdminPipelinesHandler(s.Settings)).Methods("GET")

	return mux
}

func (s *Server) ListenAndServeGeneric() {
	if s.Settings.CertFile != "" && s.Settings.KeyFile != "" {
		s.ListenAndServeTLS(s.Settings.CertFile, s.Settings.KeyFile)
	} else {
		s.ListenAndServe()
	}
}
