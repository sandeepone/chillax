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
	server.Paths["ApiProxies"] = "/chillax/api/proxies"
	server.Paths["ApiPipelines"] = "/chillax/api/pipelines"
	server.Paths["AdminProxies"] = "/chillax/admin/proxies"
	server.Paths["AdminStatic"] = "/chillax/admin/static/"

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
		s.Paths["ApiPipelines"]+"/{Id}/run",
		chillax_web_handlers.ApiPipelineRunHandler(s.Settings)).Methods("POST")

	// Admin Handlers
	mux.HandleFunc(
		s.Paths["AdminProxies"],
		chillax_web_handlers.AdminProxiesHandler(s.Settings, muxProducer.ProxyHandlers)).Methods("GET")

	staticHandler := http.StripPrefix(
		s.Paths["AdminStatic"],
		chillax_web_handlers.AdminStaticDirHandler(s.Settings.DefaultAssetsPath))

	mux.PathPrefix(s.Paths["AdminStatic"]).Handler(staticHandler)

	return mux
}
