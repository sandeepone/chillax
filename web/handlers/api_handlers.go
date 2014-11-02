package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/chillaxio/chillax/libstring"
	"github.com/chillaxio/chillax/libtime"
	chillax_proxy_backend "github.com/chillaxio/chillax/proxy/backend"
	chillax_statskeeper "github.com/chillaxio/chillax/statskeeper"
	chillax_storage "github.com/chillaxio/chillax/storage"
	chillax_web_pipelines "github.com/chillaxio/chillax/web/pipelines"
	"github.com/peterbourgon/mergemap"
	"io/ioutil"
	"net/http"
	"strings"
)

func ApiStatsRequestsJsonHandler(storage chillax_storage.Storer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			params := r.URL.Query()
			endDateString := params["end"][0]
			durationString := params["duration"][0]

			endDate, err := libtime.ParseIsoString(endDateString)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			requestDataBytes, err := chillax_statskeeper.GetRequestDataDurationsAgo(endDate, durationString)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			requestDataBytesJsonBytes := libstring.SliceOfJsonBytesToJsonArrayBytes(requestDataBytes)

			w.Header().Set("Content-Type", "application/json")
			w.Write(requestDataBytesJsonBytes)
		}
	}
}

func ApiStatsRequestsCsvHandler(storage chillax_storage.Storer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			params := r.URL.Query()
			endDateString := params["end"][0]
			durationString := params["duration"][0]

			endDate, err := libtime.ParseIsoString(endDateString)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			requestData, err := chillax_statskeeper.GetRequestDataDurationsAgo(endDate, durationString)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Set("Content-Type", "text/csv")

			csvWriter := csv.NewWriter(w)

			// Header
			recordHeader := []string{
				"CurrentUnixNano",
				"Latency",
				"Method",
				"RemoteAddr",
				"URI",
				"UserAgent",
			}
			csvWriter.Write(recordHeader)

			for _, dataBytes := range requestData {
				data := make(map[string]interface{})
				json.Unmarshal(dataBytes, &data)

				var record []string
				record = append(record, fmt.Sprintf("%v", int64(data["CurrentUnixNano"].(float64))))
				record = append(record, fmt.Sprintf("%v", int64(data["Latency"].(float64))))
				record = append(record, data["Method"].(string))
				record = append(record, data["RemoteAddr"].(string))
				record = append(record, data["URI"].(string))
				record = append(record, data["UserAgent"].(string))
				csvWriter.Write(record)
			}
			csvWriter.Flush()
		}
	}
}

func ApiProxiesHandler() func(http.ResponseWriter, *http.Request) {
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
