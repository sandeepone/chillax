package main

import (
    "net/http"
    chillax_web_handlers "github.com/didip/chillax/web/handlers"
    chillax_proxy_muxproducer "github.com/didip/chillax/proxy/muxproducer"
)

func main() {
    mp, _ := chillax_proxy_muxproducer.NewMuxProducer()

    mp.CreateProxyBackends()
    mp.StartProxyBackends()
    mux := mp.GorillaMuxWithProxyBackends()

    staticDir := "./web/static"
    chillax_web_handlers.GorillaMuxRouteStaticDir(mux, staticDir)

    http.ListenAndServe(":8080", mux)
}
