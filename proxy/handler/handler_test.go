package handler

import (
    "os"
    "fmt"
    "bufio"
    "testing"
    "io/ioutil"
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
