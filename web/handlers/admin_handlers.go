package handlers

import (
	"fmt"
	"net/http"
	"strings"

	chillax_proxy_backend "github.com/chillaxio/chillax/proxy/backend"
	chillax_proxy_handler "github.com/chillaxio/chillax/proxy/handler"
	chillax_web_pipelines "github.com/chillaxio/chillax/web/pipelines"
	chillax_web_templates_admin "github.com/chillaxio/chillax/web/templates/admin"
)

// AdminBaseHandler renders HTML for /admin/base
func AdminBaseHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, chillax_web_templates_admin.NewAdminBase().String())
	}
}

// AdminStatsHandler renders HTML for /admin/stats
func AdminStatsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			DurationString string
		}{
			"-336h",
		}
		t, err := chillax_web_templates_admin.NewAdminStats().Parse()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		t.Execute(w, data)
	}
}

// AdminProxiesHandler renders HTML for /admin/proxies
func AdminProxiesHandler(proxyHandlers []*chillax_proxy_handler.ProxyHandler) func(http.ResponseWriter, *http.Request) {
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

// AdminProxyHandler renders HTML for /admin/proxies/{Name}
func AdminProxyHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pathChunk := strings.Split(r.URL.Path, "/")
		proxyName := pathChunk[len(pathChunk)-1]

		proxyBackend, err := chillax_proxy_backend.LoadProxyBackendByName(proxyName)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := struct {
			ProxyBackend *chillax_proxy_backend.ProxyBackend
		}{
			proxyBackend,
		}
		t, err := chillax_web_templates_admin.NewAdminProxy().Parse()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		t.Execute(w, data)
	}
}

// AdminPipelinesHandler renders HTML for /admin/pipelines
func AdminPipelinesHandler() func(http.ResponseWriter, *http.Request) {
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

// AdminPipelineHandler renders HTML for /admin/pipeline/{Id}
func AdminPipelineHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pathChunk := strings.Split(r.URL.Path, "/")
		pipelineId := pathChunk[len(pathChunk)-1]

		pipeline, err := chillax_web_pipelines.PipelineById(pipelineId)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := struct {
			Pipeline *chillax_web_pipelines.Pipeline
		}{
			pipeline,
		}
		t, err := chillax_web_templates_admin.NewAdminPipeline().Parse()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		t.Execute(w, data)
	}
}
