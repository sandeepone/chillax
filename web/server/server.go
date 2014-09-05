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
		AdminProxiesPath: "/chillax/proxies",
		AdminStaticPath:  "/chillax/static/",
		Settings:         settings,
		Server: &graceful.Server{
			Timeout: requestTimeoutOnRestart,
			Server:  &http.Server{Addr: settings.HttpAddress()},
		},
	}

	return server, err
}

type Server struct {
	*graceful.Server

	Settings         *chillax_web_settings.ServerSettings
	AdminProxiesPath string
	AdminStaticPath  string
}

func (s *Server) NewGorillaMux() *gorilla_mux.Router {
	muxProducer := chillax_proxy_muxproducer.NewMuxProducer(s.Settings.ProxyHandlerTomls)

	muxProducer.CreateProxyBackends()
	muxProducer.StartProxyBackends()
	mux := muxProducer.GorillaMuxWithProxyBackends()

	mux.HandleFunc(
		s.AdminProxiesPath,
		chillax_web_handlers.ProxiesHandler(s.Settings, muxProducer.ProxyHandlers)).Methods("GET")

	staticHandler := http.StripPrefix(
		s.AdminStaticPath,
		chillax_web_handlers.StaticDirHandler(s.Settings.DefaultAssetsPath))

	mux.PathPrefix(s.AdminStaticPath).Handler(staticHandler)

	return mux
}
