package handlers

import (
    "fmt"
    "net/http"
    "html/template"
    // "github.com/GeertJohan/go.rice"
    chillax_proxy_handler "github.com/didip/chillax/proxy/handler"
)

func StaticDirHandler(staticDirectory string) http.Handler {
    // box, _ := rice.FindBox(staticDirectory)
    // return http.FileServer(box.HTTPBox())

    fmt.Printf("staticDirectory: %v\n", staticDirectory)

    return http.FileServer(http.Dir(staticDirectory))
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