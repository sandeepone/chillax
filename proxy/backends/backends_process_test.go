package backends

import (
    "os"
    "bufio"
    "testing"
    "io/ioutil"
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
        t.Errorf("backend.Command should exists. Backend: %v", backend)
    }
}

func TestCreateProcesses(t *testing.T) {
    backend := NewProcessProxyBackendForTest()

    err := backend.CreateProcesses()
    if err != nil {
        t.Errorf("Failed to create processes. Error: %v", err)
    }
}