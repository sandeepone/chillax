package handler

import (
    "os"
    "fmt"
    "bufio"
    "testing"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "github.com/didip/chillax/libtime"
)


func NewProxyHandlerForTest() *ProxyHandler {
    fileHandle, _ := os.Open("../backend/example-process-backend.toml")
    bufReader     := bufio.NewReader(fileHandle)
    definition, _ := ioutil.ReadAll(bufReader)
    handler       := NewProxyHandler(definition)
    return handler
}

func TestBackendHosts(t *testing.T) {
    handler := NewProxyHandlerForTest()
    handler.CreateBackends()

    if handler.Backend == nil {
        t.Errorf("handler.Backend should exists. handler.Backend: %v", handler.Backend)
    }

    if handler.Backend.Process == nil {
        t.Errorf("handler.Backend.Process should exists. handler.Backend.Process: %v", handler.Backend.Process)
    }

    if len(handler.BackendHosts()) != 1 {
        t.Errorf("handler.BackendHosts should exists. handler.BackendHosts: %v", handler.BackendHosts)
    }

    instance1 := handler.Backend.Process.Instances[0]

    if handler.BackendHosts()[0] != fmt.Sprintf("127.0.0.1:%v", instance1.MapPorts[handler.Backend.Process.HttpPortEnv]) {
        t.Errorf("handler.BackendHosts()[0] is incorrect. handler.BackendHosts()[0]: %v", handler.BackendHosts()[0])
    }
}

func TestChooseBackendHost(t *testing.T) {
    handler := NewProxyHandlerForTest()
    handler.CreateBackends()

    host := handler.ChooseBackendHost()

    if host != handler.BackendHosts()[0] {
        t.Errorf("handler.BackendHosts()[0] should always be chosen. host: %v", host)
    }
}

func TestProxyHandlerFunction(t *testing.T) {
    handler := NewProxyHandlerForTest()
    handler.CreateBackends()

    errors := handler.StartBackends()
    for _, err := range errors {
        if err != nil {
            t.Errorf("Failed to start process. Error: %v", err)
        }
    }

    libtime.SleepString("500ms")

    server := httptest.NewServer(http.HandlerFunc(handler.Function()))
    defer server.Close()

    response, err := http.Get(server.URL)
    if err != nil {
        t.Fatalf("Unable to hit server.URL endpoint. Error: %v. URL: %v", err, server.URL)
    }
    defer response.Body.Close()

    content, err := ioutil.ReadAll(response.Body)
    if err != nil || len(content) == 0 {
        t.Errorf("Unable to read content of the endpoint. Error: %v, Content: %v", err, string(content))
    }

    if response.StatusCode != 200 {
        t.Errorf("response.StatusCode should == 200. Response: %v", response.StatusCode)
    }

    directResponse, err := http.Get("http://" + handler.BackendHosts()[0])
    if err != nil {
        t.Fatalf("Unable to get direct response to backend. Error: %v", err)
    }
    directResponseContent, _ := ioutil.ReadAll(directResponse.Body)
    directResponse.Body.Close()

    if string(content) != string(directResponseContent) {
        t.Errorf("Content is not what we expect. Content: %v, Expected Content: %v", string(content), string(directResponseContent))
    }

    errors = handler.StopBackends()
    for _, err := range errors {
        if err != nil {
            t.Errorf("Failed to stop backend process. Error: %v", err)
        }
    }
}