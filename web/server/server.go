package server

import (
    "net/http"
    gorilla_mux "github.com/gorilla/mux"
    chillax_web_handlers "github.com/didip/chillax/web/handlers"
    chillax_web_settings "github.com/didip/chillax/web/settings"
    chillax_proxy_muxproducer "github.com/didip/chillax/proxy/muxproducer"
)

func NewServer() (*Server, error) {
    var err error

    server := &Server{}
    server.AdminProxiesPath = "/chillax/proxies"
    server.AdminStaticPath  = "/chillax/static/"

    server.Settings, err = chillax_web_settings.NewServerSettings()
    return server, err
}

type Server struct {
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

func (s *Server) Serve(handler http.Handler) error {
    return http.ListenAndServe(s.Settings.HttpAddress(), handler)
}
