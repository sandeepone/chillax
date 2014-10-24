package handlers

import (
	"fmt"
	chillax_proxy_handler "github.com/chillaxio/chillax/proxy/handler"
	chillax_web_pipelines "github.com/chillaxio/chillax/web/pipelines"
	chillax_web_settings "github.com/chillaxio/chillax/web/settings"
	chillax_web_templates_admin "github.com/chillaxio/chillax/web/templates/admin"
	"net/http"
)

func AdminBaseHandler(settings *chillax_web_settings.ServerSettings) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, chillax_web_templates_admin.NewAdminBase().String())
	}
}

func AdminProxiesHandler(settings *chillax_web_settings.ServerSettings, proxyHandlers []*chillax_proxy_handler.ProxyHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			ProxyHandlers []*chillax_proxy_handler.ProxyHandler
		}{
			proxyHandlers,
		}
		t, err := chillax_web_templates_admin.NewAdminProxies().Parse()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		t.Execute(w, data)
	}
}

func AdminPipelinesHandler(settings *chillax_web_settings.ServerSettings) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pipelines, err := chillax_web_pipelines.AllPipelines()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := struct {
			Pipelines []*chillax_web_pipelines.Pipeline
		}{
			pipelines,
		}
		t, err := chillax_web_templates_admin.NewAdminPipelines().Parse()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		t.Execute(w, data)
	}
}
