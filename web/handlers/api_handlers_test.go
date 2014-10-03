package handlers

import (
	"bufio"
	"bytes"
	chillax_storage "github.com/didip/chillax/storage"
	chillax_web_settings "github.com/didip/chillax/web/settings"
	gorilla_mux "github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func NewProxyBackendDefinitionForTest() []byte {
	fileHandle, _ := os.Open("../../proxy/backend/example-process-backend.toml")
	bufReader := bufio.NewReader(fileHandle)
	definition, _ := ioutil.ReadAll(bufReader)
	return definition
}

func NewServerSettingsForTest() *chillax_web_settings.ServerSettings {
	settings, _ := chillax_web_settings.NewServerSettings()
	return settings
}

func TestApiProxies(t *testing.T) {
	// ---- Setup ----
	settings := NewServerSettingsForTest()

	mux := gorilla_mux.NewRouter()
	mux.HandleFunc(
		"/chillax/api/proxies",
		ApiProxiesHandler(settings)).Methods("POST")

	go http.ListenAndServe(":18000", mux)

	definition := NewProxyBackendDefinitionForTest()

	storage := chillax_storage.NewStorage()

	storage.Delete("/proxies/")

	// ---- Setup ----

	// Get existing number of proxies
	proxies, err := storage.List("/proxies")
	prevProxiesLength := len(proxies)

	req, err := http.NewRequest("POST", "http://localhost:18000/chillax/api/proxies", bytes.NewBuffer(definition))
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
