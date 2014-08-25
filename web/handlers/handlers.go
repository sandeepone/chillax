package handlers

import (
    "net/http"
    "html/template"
    "github.com/GeertJohan/go.rice"
    gorilla_mux "github.com/gorilla/mux"
    chillax_proxy_handler "github.com/didip/chillax/proxy/handler"
)

func StaticDirHandler(staticDirectory string) func(http.ResponseWriter, *http.Request) {
    box, _ := rice.FindBox(staticDirectory)
    return func(w http.ResponseWriter, r *http.Request) {
        http.FileServer(box.HTTPBox())
    }
}

func ProxiesHandler(proxyHandlers []*chillax_proxy_handler.ProxyHandler) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        data := struct {
            ProxyHandlers []*chillax_proxy_handler.ProxyHandler
        }{
            proxyHandlers,
        }
        t, _ := template.ParseFiles("/Users/didip/projects/go/src/github.com/didip/chillax/web/default-assets/server-templates/proxies/list.html")
        t.Execute(w, data)
    }
}