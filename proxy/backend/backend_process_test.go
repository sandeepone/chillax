package backend

import (
    "os"
    "bufio"
    "testing"
    "io/ioutil"
    "github.com/didip/chillax/libtime"
)


func NewProcessProxyBackendForTest() *ProxyBackend {
    fileHandle, _ := os.Open("./example-process-backend.toml")
    bufReader     := bufio.NewReader(fileHandle)
    definition, _ := ioutil.ReadAll(bufReader)
    backend       := NewProxyBackend(definition)
    return backend
}

func TestDeserializeProcessProxyBackendFromToml(t *testing.T) {
    backend := NewProcessProxyBackendForTest()

    if backend.Command == "" {
        t.Errorf("backend.Command should exists. Backend.Command: %v", backend.Command)
    }
}

func TestCreateProcesses(t *testing.T) {
    backend := NewProcessProxyBackendForTest()

    err := backend.CreateProcesses()
    if err != nil {
        t.Errorf("Failed to create processes. Error: %v", err)
    }
}

func TestStartStopProcesses(t *testing.T) {
    backend := NewProcessProxyBackendForTest()
    backend.CreateProcesses()

    go func() {
        errors := backend.StartProcesses()
        for _, err := range errors {
            if err != nil {
                t.Errorf("Failed to start process. Error: %v", err)
            }
        }
    }()

    libtime.SleepString("5s")

    go func() {
        errors := backend.StopProcesses()
        t.Errorf("Failed to stop process. Errors: %v", errors)
        for _, err := range errors {
            if err != nil {
                t.Errorf("Failed to stop process. Error: %v", err)
            }
        }
    }()
}