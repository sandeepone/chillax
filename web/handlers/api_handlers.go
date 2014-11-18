package handlers

import (
	"encoding/json"
	"fmt"
	chillax_metricskeeper "github.com/chillaxio/chillax/metricskeeper"
	chillax_proxy_backend "github.com/chillaxio/chillax/proxy/backend"
	chillax_storage "github.com/chillaxio/chillax/storage"
	chillax_web_pipelines "github.com/chillaxio/chillax/web/pipelines"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/peterbourgon/mergemap"
	"io/ioutil"
	"net/http"
	"strings"
)

func ApiStatsCpuJsonHandler(storage chillax_storage.Storer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			metrics, err := chillax_metricskeeper.LoadCpuFromAllHosts(storage)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			inJson, err := json.Marshal(metrics)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(inJson)
		}
	}
}

func ApiProxiesTomlHandler() func(http.ResponseWriter, *http.Request) {
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

func ApiProxiesRestartHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// tell chillax_web_handler to restart all endpoints and return error.
		}
	}
}

func ApiProxyTomlHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := gorilla_mux.Vars(r)
		proxyName := vars["name"]

		requestBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		proxyBackend, err := chillax_proxy_backend.LoadProxyBackendByName(proxyName)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if r.Method == "POST" || r.Method == "PUT" {
			proxyBackend, err = chillax_proxy_backend.UpdateProxyBackend(proxyBackend, requestBodyBytes)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			err = proxyBackend.Save()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

		} else if r.Method == "DELETE" {
			err = chillax_proxy_backend.DeleteProxyBackendByName(proxyName)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

		}
	}
}

func ApiProxyJsonHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := gorilla_mux.Vars(r)
		proxyName := vars["name"]

		proxyBackend, err := chillax_proxy_backend.LoadProxyBackendByName(proxyName)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if r.Method == "GET" {
			inJson, err := json.Marshal(proxyBackend)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(inJson)
		}
	}
}

func ApiProxyRestartHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// tell chillax_web_handler to restart 1 endpoint and return error.
		}
	}
}

func ApiPipelinesHandler() func(http.ResponseWriter, *http.Request) {
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

func ApiPipelinesRunHandler() func(http.ResponseWriter, *http.Request) {
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

func ApiPipelineRunHandler() func(http.ResponseWriter, *http.Request) {
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

			pipeline, err = mergePipelineBody(r, pipeline)
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

	pipeline, err := chillax_web_pipelines.NewPipelineGivenJsonBytes(requestBodyBytes)
	if err != nil {
		return nil, err
	}

	err = pipeline.Save()
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func mergePipelineBody(r *http.Request, pipeline *chillax_web_pipelines.Pipeline) (*chillax_web_pipelines.Pipeline, error) {
	requestBodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
	json.Unmarshal(requestBodyBytes, &body)

	pipeline.Body = mergemap.Merge(pipeline.Body, body)
	return pipeline, err
}
