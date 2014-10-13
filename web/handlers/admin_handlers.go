package handlers

import (
	chillax_proxy_handler "github.com/didip/chillax/proxy/handler"
	chillax_web_settings "github.com/didip/chillax/web/settings"
	chillax_web_templates_admin "github.com/didip/chillax/web/templates/admin"
	"net/http"
)

func AdminProxiesHandler(settings *chillax_web_settings.ServerSettings, proxyHandlers []*chillax_proxy_handler.ProxyHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			ProxyHandlers []*chillax_proxy_handler.ProxyHandler
		}{
			proxyHandlers,
		}
		t, err := chillax_web_templates_admin.NewProxies().Parse()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		t.Execute(w, data)
	}
}
