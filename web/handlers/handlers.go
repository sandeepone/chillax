package handlers

import (
    "net/http"
    "html/template"
    "github.com/GeertJohan/go.rice"
    gorilla_mux "github.com/gorilla/mux"
    chillax_proxy_handler "github.com/didip/chillax/proxy/handler"
    chillax_proxy_muxproducer "github.com/didip/chillax/proxy/muxproducer"
)

func GorillaMuxRouteStaticDir(router *gorilla_mux.Router, staticDirectory string) {
    box, err := rice.FindBox(staticDirectory)
    if err == nil {
        router.Handle(staticDirectory, http.FileServer(box.HTTPBox()))
    }
}

func ProxiesHandler(mp *chillax_proxy_muxproducer.MuxProducer) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        proxyHandlers := make([]*chillax_proxy_handler.ProxyHandler, len(mp.ProxyHandlers))
        copy(proxyHandlers, mp.ProxyHandlers)

        data := struct {
            ProxyHandlers []*chillax_proxy_handler.ProxyHandler
        }{
            proxyHandlers,
        }
        t, _ := template.ParseFiles("/Users/didip/projects/go/src/github.com/didip/chillax/web/default-assets/server-templates/proxies/list.html")
        t.Execute(w, data)
    }
}