package main

import (
    "net/http"
    chillax_web_handlers "github.com/didip/chillax/web/handlers"
    chillax_web_settings "github.com/didip/chillax/web/settings"
    chillax_proxy_muxproducer "github.com/didip/chillax/proxy/muxproducer"
)

func main() {
    settings, err := chillax_web_settings.NewServerSettings()

    if err != nil {
        panic(err)
    }

    muxProducer := chillax_proxy_muxproducer.NewMuxProducer(settings.ProxyHandlerTomls)

    muxProducer.CreateProxyBackends()
    muxProducer.StartProxyBackends()
    mux := muxProducer.GorillaMuxWithProxyBackends()

    mux.HandleFunc("/chillax/static", chillax_web_handlers.StaticDirHandler(settings.DefaultAssetsPath)).Methods("GET")

    mux.HandleFunc("/chillax/proxies", chillax_web_handlers.ProxiesHandler(muxProducer.ProxyHandlers)).Methods("GET")

    http.ListenAndServe(settings.HttpAddress(), mux)
}
