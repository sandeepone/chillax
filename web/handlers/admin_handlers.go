package handlers

import (
	"github.com/GeertJohan/go.rice"
	chillax_proxy_handler "github.com/didip/chillax/proxy/handler"
	chillax_web_settings "github.com/didip/chillax/web/settings"
	"html/template"
	"net/http"
	"path/filepath"
)

func AdminStaticDirHandler(staticDirectory string) http.Handler {
	box, _ := rice.FindBox(staticDirectory)
	return http.FileServer(box.HTTPBox())
}

func AdminProxiesHandler(settings *chillax_web_settings.ServerSettings, proxyHandlers []*chillax_proxy_handler.ProxyHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			ProxyHandlers []*chillax_proxy_handler.ProxyHandler
		}{
			proxyHandlers,
		}
		templatePath := filepath.Join(settings.DefaultAssetsPath, "server-templates/proxies/list.html")
		t, _ := template.ParseFiles(templatePath)
		t.Execute(w, data)
	}
}
