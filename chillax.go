package main

import (
    "net/http"
    chillax_web_handlers "github.com/didip/chillax/web/handlers"
    chillax_web_settings "github.com/didip/chillax/web/settings"
    chillax_proxy_muxproducer "github.com/didip/chillax/proxy/muxproducer"
)

func main() {
    settings    := chillax_web_settings.NewServerSettings()
    muxProducer := chillax_proxy_muxproducer.NewMuxProducer(settings.ProxyHandlerTomls)

    muxProducer.CreateProxyBackends()
    muxProducer.StartProxyBackends()
    mux := muxProducer.GorillaMuxWithProxyBackends()

    chillax_web_handlers.GorillaMuxRouteStaticDir(mux, settings.DefaultAssetsPath)

    mux.HandleFunc("/proxies", chillax_web_handlers.ProxiesHandler(muxProducer)).Methods("GET")

    http.ListenAndServe(settings.HttpAddress(), mux)
}
