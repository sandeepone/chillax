package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	chillax_proxy_backend "github.com/chillaxio/chillax/proxy/backend"
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
		ApiProxyJsonHandler()).Methods("GET", "DELETE")

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

func TestCreateGetUpdateAndDeleteProxyFromApi(t *testing.T) {
	// ---- Setup ----
	os.Setenv("CHILLAX_ENV", "test")

	mux := NewGorillaMuxForTest()

	go http.ListenAndServe(":18000", mux)

	definition := NewProxyBackendDefinitionForTest()

	// ---- Setup ----

	// Create a proxy
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

	// Get a proxy
	resp, err = http.Get("http://localhost:18000/chillax/api/proxies/test-chillax-api-proxies-post.json")
	if err != nil {
		t.Errorf("Fail to perform GET request. Error: %v", err)
	}
	defer resp.Body.Close()

	// Update a proxy with new definition via PUT
	backend, _ := chillax_proxy_backend.NewProxyBackend(definition)
	backend.Numprocs = 5
	definition, _ = backend.Serialize()

	req, err = http.NewRequest("PUT", "http://localhost:18000/chillax/api/proxies/test-chillax-api-proxies-post.json", bytes.NewBuffer(definition))
	if err != nil {
		t.Errorf("Fail to create PUT request. Error: %v", err)
	}

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Fail to perform PUT request. Error: %v", err)
	}
	defer resp.Body.Close()

	// DELETE a proxy
	req, err = http.NewRequest("DELETE", "http://localhost:18000/chillax/api/proxies/test-chillax-api-proxies-post.json", nil)
	if err != nil {
		t.Errorf("Fail to create DELETE request. Error: %v", err)
	}

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Fail to perform DELETE request. Error: %v", err)
	}
	defer resp.Body.Close()

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
