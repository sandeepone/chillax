package main

import (
    "net/http"
    "github.com/didip/chillax/proxy/muxproducer"
)

func main() {
    mp, _ := muxproducer.NewMuxProducer()

    mp.CreateProxyBackends()
    mp.StartProxyBackends()
    mux := mp.GorillaMuxWithProxyBackends()

    http.ListenAndServe(":8080", mux)
}
