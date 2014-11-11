package backend

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

func NewProcessProxyBackendForTest() *ProxyBackend {
	fileHandle, _ := os.Open("./example-process-backend.toml")
	bufReader := bufio.NewReader(fileHandle)
	definition, _ := ioutil.ReadAll(bufReader)
	backend, _ := NewProxyBackend(definition)
	return backend
}

func TestDeserializeProcessProxyBackendFromToml(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	backend := NewProcessProxyBackendForTest()

	if backend.Command == "" {
		t.Errorf("backend.Command should exists. Backend.Command: %v", backend.Command)
	}
	if len(backend.Process.Hosts) != 1 {
		t.Errorf("backend.Process.Hosts should contains 1 element. Backend.Process.Hosts: %v", backend.Process.Hosts)
	}
}

func TestStartRestartAndStopProcesses(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	backend := NewProcessProxyBackendForTest()

	errors := backend.CreateProcesses()
	if len(errors) > 0 {
		t.Fatalf("CreateProcesses should not return errors. Errors: %v", errors)
	}

	backendFromStorage, loadErr := LoadProxyBackendByName("test-process-backend")
	if loadErr != nil {
		t.Errorf("Unable to load proxy data from storage. loadErr: %v", loadErr)
	}

	if backend.Path != backendFromStorage.Path || backend.Command != backendFromStorage.Command {
		t.Errorf("Backend was improperly saved. backendFromStorage: %v", backendFromStorage)
	}

	errors = backend.StartProcesses()

	if len(errors) > 0 {
		t.Fatalf("StartProcesses returns errors. Errors: %v, backend.Process.Instances: %v", errors, backend.Process.Instances)
	}

	for _, err := range errors {
		if err != nil {
			t.Fatalf("Failed to start process. Error: %v", err)
		}
	}

	if len(backend.Process.Instances) != 2 {
		t.Errorf("Expected to start 2 processes, got: %v", len(backend.Process.Instances))
	}

	for _, instance := range backend.Process.Instances {
		if instance.ProcessWrapper == nil {
			t.Fatalf("Process was not started.")
		}
	}

	// errors = backend.RestartProcesses()

	// if len(errors) > 0 {
	// 	t.Fatalf("errors slice should be empty. Errors: %v", errors)
	// }

	// for _, err := range errors {
	// 	if err != nil {
	// 		t.Fatalf("Failed to restart process. Error: %v", err)
	// 	}
	// }

	errors = backend.StopProcesses()

	if len(errors) > 0 {
		t.Fatalf("errors slice should be empty. Errors: %v", errors)
	}

	for _, err := range errors {
		if err != nil {
			t.Fatalf("Failed to stop process. Error: %v", err)
		}
	}
}
