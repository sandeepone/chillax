package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	chillax_storage "github.com/chillaxio/chillax/storage"
	gorilla_mux "github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func NewProxyBackendDefinitionForTest() []byte {
	fileHandle, _ := os.Open("./tests-data/serialized-process-backend.toml")
	bufReader := bufio.NewReader(fileHandle)
	definition, _ := ioutil.ReadAll(bufReader)
	return definition
}

func NewPipelineJsonDefinitionForTest() []byte {
	fileHandle, _ := os.Open("./tests-data/serialized-pipeline.json")
	bufReader := bufio.NewReader(fileHandle)
	definition, _ := ioutil.ReadAll(bufReader)
	return definition
}

func NewGorillaMuxForTest() *gorilla_mux.Router {
	mux := gorilla_mux.NewRouter()

	mux.HandleFunc(
		"/chillax/api/proxies.toml",
		ApiProxiesTomlHandler()).Methods("POST")

	mux.HandleFunc(
		"/chillax/api/proxies/restart",
		ApiProxiesRestartHandler()).Methods("POST")

	mux.HandleFunc(
		"/proxies/{name}.toml",
		ApiProxyTomlHandler()).Methods("PUT", "POST", "DELETE")

	mux.HandleFunc(
		"/proxies/{name}.json",
		ApiProxyJsonHandler()).Methods("GET")

	mux.HandleFunc(
		"/chillax/api/proxy/{name}/restart",
		ApiProxyRestartHandler()).Methods("POST")

	mux.HandleFunc(
		"/chillax/api/pipelines",
		ApiPipelinesHandler()).Methods("POST")

	mux.HandleFunc(
		"/chillax/api/pipelines/run",
		ApiPipelinesRunHandler()).Methods("POST")

	mux.HandleFunc(
		"/chillax/api/pipelines/{Id}/run",
		ApiPipelineRunHandler()).Methods("POST")

	return mux
}

func TestApiProxies(t *testing.T) {
	// ---- Setup ----
	os.Setenv("CHILLAX_ENV", "test")

	mux := NewGorillaMuxForTest()

	go http.ListenAndServe(":18000", mux)

	definition := NewProxyBackendDefinitionForTest()

	storage := chillax_storage.NewStorage()

	storage.Delete("/proxies/")

	// ---- Setup ----

	// Get existing number of proxies
	proxies, err := storage.List("/proxies")
	prevProxiesLength := len(proxies)

	req, err := http.NewRequest("POST", "http://localhost:18000/chillax/api/proxies.toml", bytes.NewBuffer(definition))
	if err != nil {
		t.Errorf("Fail to create POST request. Error: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Fail to perform POST request. Error: %v", err)
	}
	defer resp.Body.Close()

	// Get new number of proxies
	proxies, err = storage.List("/proxies")
	currentProxiesLength := len(proxies)

	if currentProxiesLength <= prevProxiesLength {
		t.Errorf("proxy definition was not saved correctly. prevProxiesLength: %v, currentProxiesLength: %v", prevProxiesLength, currentProxiesLength)
	}
}

func TestApiPipelinesAndRun(t *testing.T) {
	// ---- Setup ----
	os.Setenv("CHILLAX_ENV", "test")

	mux := NewGorillaMuxForTest()

	go http.ListenAndServe(":18001", mux)

	definition := NewPipelineJsonDefinitionForTest()

	storage := chillax_storage.NewStorage()

	storage.Delete("/pipelines/")

	// ---- Setup ----

	// Get existing number of pipelines
	pipelines, err := storage.List("/pipelines")
	prevPipelinesLength := len(pipelines)

	// Create a pipeline
	req, err := http.NewRequest("POST", "http://localhost:18001/chillax/api/pipelines", bytes.NewBuffer(definition))
	if err != nil {
		t.Errorf("Fail to create POST request. Error: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Fail to perform POST request. Error: %v", err)
	}
	defer resp.Body.Close()

	// Get new number of pipelines
	pipelines, err = storage.List("/pipelines")
	currentPipelinesLength := len(pipelines)

	if currentPipelinesLength <= prevPipelinesLength {
		t.Errorf("pipeline definition was not saved correctly. prevPipelinesLength: %v, currentPipelinesLength: %v", prevPipelinesLength, currentPipelinesLength)
	}

	//
	// Read JSON payload after saving a pipeline and check its Id.
	//
	jsonBytes, _ := ioutil.ReadAll(resp.Body)

	data := make(map[string]string)
	json.Unmarshal(jsonBytes, &data)

	if data["Id"] == "" {
		t.Errorf("POST /pipelines did not return Id. jsonBytes: %v", string(jsonBytes))
	}

	//
	// Run pipeline
	//
	req, err = http.NewRequest(
		"POST",
		"http://localhost:18001/chillax/api/pipelines/"+data["Id"]+"/run",
		bytes.NewBuffer([]byte(`{"Brotato": true}`)))
	if err != nil {
		t.Errorf("Fail to create POST request. Error: %v", err)
	}

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Fail to perform POST request. Error: %v", err)
	}
	defer resp.Body.Close()

}

func TestApiPipelinesRun(t *testing.T) {
	// ---- Setup ----
	os.Setenv("CHILLAX_ENV", "test")

	mux := NewGorillaMuxForTest()

	go http.ListenAndServe(":18002", mux)

	definition := NewPipelineJsonDefinitionForTest()

	storage := chillax_storage.NewStorage()

	storage.Delete("/pipelines/")

	// ---- Setup ----

	// Get existing number of pipelines
	pipelines, err := storage.List("/pipelines")
	prevPipelinesLength := len(pipelines)

	req, err := http.NewRequest("POST", "http://localhost:18002/chillax/api/pipelines/run", bytes.NewBuffer(definition))
	if err != nil {
		t.Errorf("Fail to create POST request. Error: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Fail to perform POST request. Error: %v", err)
	}
	defer resp.Body.Close()

	// Get new number of pipelines
	pipelines, err = storage.List("/pipelines")
	currentPipelinesLength := len(pipelines)

	if currentPipelinesLength <= prevPipelinesLength {
		t.Errorf("pipeline definition was not saved correctly. prevPipelinesLength: %v, currentPipelinesLength: %v", prevPipelinesLength, currentPipelinesLength)
	}

	//
	// Read JSON payload after saving a pipeline and check its Id.
	//
	jsonBytes, _ := ioutil.ReadAll(resp.Body)

	data := make(map[string]string)
	json.Unmarshal(jsonBytes, &data)

	if data["Id"] == "" {
		t.Errorf("POST /pipelines/run did not return Id. jsonBytes: %v", string(jsonBytes))
	}
}
