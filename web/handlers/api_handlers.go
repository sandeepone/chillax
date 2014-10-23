package handlers

import (
	"fmt"
	chillax_proxy_backend "github.com/didip/chillax/proxy/backend"
	chillax_web_pipelines "github.com/didip/chillax/web/pipelines"
	chillax_web_settings "github.com/didip/chillax/web/settings"
	"io/ioutil"
	"net/http"
	"strings"
)

func ApiProxiesHandler(settings *chillax_web_settings.ServerSettings) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {

		} else if r.Method == "POST" {
			requestBodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			proxyBackend, err := chillax_proxy_backend.NewProxyBackend(requestBodyBytes)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			err = proxyBackend.Save()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
	}
}

func ApiPipelinesHandler(settings *chillax_web_settings.ServerSettings) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {

		} else if r.Method == "POST" {
			pipeline, err := savePipeline(w, r)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf(`{"Id": "%v"}`, pipeline.Id)))
		}
	}
}

func ApiPipelinesRunHandler(settings *chillax_web_settings.ServerSettings) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {

		} else if r.Method == "POST" {
			pipeline, err := savePipeline(w, r)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			runInstance := pipeline.RunWithCrashRequeue()
			if runInstance.ErrorMessage != "" {
				http.Error(w, runInstance.ErrorMessage, 500)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf(`{"Id": "%v"}`, pipeline.Id)))
		}
	}
}

func ApiPipelineRunHandler(settings *chillax_web_settings.ServerSettings) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {

		} else if r.Method == "POST" {
			pathChunk := strings.Split(r.URL.Path, "/")

			pipelineId := pathChunk[len(pathChunk)-2]

			pipeline, err := chillax_web_pipelines.PipelineById(pipelineId)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			runInstance := pipeline.RunWithCrashRequeue()
			if runInstance.ErrorMessage != "" {
				http.Error(w, runInstance.ErrorMessage, 500)
				return
			}
		}
	}
}

//
// Private functions
//
func savePipeline(w http.ResponseWriter, r *http.Request) (*chillax_web_pipelines.Pipeline, error) {
	requestBodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	pipeline, err := chillax_web_pipelines.NewPipeline(string(requestBodyBytes))
	if err != nil {
		return nil, err
	}

	err = pipeline.Save()
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}
