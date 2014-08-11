package backends

import (
    "os"
    "bufio"
    "testing"
    "io/ioutil"
)

func NewDockerProxyBackendForTest() *ProxyBackend {
    fileHandle, _ := os.Open("./example-docker-backend.toml")
    bufReader     := bufio.NewReader(fileHandle)
    definition, _ := ioutil.ReadAll(bufReader)
    backend       := NewProxyBackend(definition)
    return backend
}

func TestDeserializeFromToml(t *testing.T) {
    backend := NewDockerProxyBackendForTest()

    if backend.Path != "/path/to/scraper" {
        t.Errorf("backend.Path should be /path/to/scraper. Backend: %s", backend)
    }
    if backend.Command != "" {
        t.Errorf("backend.Command should not exists. Backend: %s", backend)
    }
    if backend.Numprocs != 2 {
        t.Errorf("backend.Numprocs should == 2. Backend: %s", backend)
    }
    if backend.Delay != "1m" {
        t.Errorf("backend.Command should == 1m. Backend: %s", backend)
    }
    if backend.Ping != "30s" {
        t.Errorf("backend.Command should == 30s. Backend: %s", backend)
    }

    if backend.Docker.Tag != "didip/go-urldownloader:latest" {
        t.Errorf("backend.Docker.Tag should == didip/go-urldownloader:latest. Backend.Docker: %s", backend.Docker)
    }
    if backend.Env[0] != "HTTP_PORT=8080" {
        t.Errorf("backend.Docker.Env[0] should == HTTP_PORT=8080. Backend.Env: %s", backend.Env)
    }
    if backend.Docker.Hosts[0] != "tcp://127.0.0.1:2375" {
        t.Errorf("backend.Docker.Hosts[0] should == tcp://127.0.0.1:2375. Backend.Docker.Hosts: %s", backend.Docker.Hosts)
    }
    if backend.Docker.Ports[0] != "8080/tcp" {
        t.Errorf("backend.Docker.Ports[0] should == 8080/tcp. Backend.Docker.Ports: %s", backend.Docker.Ports)
    }

    if backend.Docker.Containers[0].Id != "abc123" {
        t.Errorf("backend.Docker.Containers[0].Id should == abc123. Backend.Docker.Containers[0]: %s", backend.Docker.Containers[0])
    }
    if backend.Docker.Containers[0].Tag != "didip/go-urldownloader:latest" {
        t.Errorf("backend.Docker.Containers[0].Tag should == didip/go-urldownloader:latest. Backend.Docker.Containers[0]: %s", backend.Docker.Containers[0])
    }
    if len(backend.Docker.Containers[0].Env) != 4 {
        t.Errorf("backend.Docker.Containers[0].Env should contains 4 items. Backend.Docker.Containers[0]: %s", backend.Docker.Containers[0])
    }
    if backend.Docker.Containers[0].Host != "tcp://127.0.0.1:2375" {
        t.Errorf("backend.Docker.Containers[0].Host should == tcp://127.0.0.1:2375. Backend.Docker.Containers[0]: %s", backend.Docker.Containers[0])
    }
    if backend.Docker.Containers[0].Ports[0] != "65000:8080/tcp" {
        t.Errorf("backend.Docker.Containers[0].Ports[0] should == 65000:8080/tcp. Backend.Docker.Containers[0]: %s", backend.Docker.Containers[0])
    }

    if backend.Docker.Containers[1].Id != "abc12456" {
        t.Errorf("backend.Docker.Containers[1].Id should == abc12456. Backend.Docker.Containers[1]: %s", backend.Docker.Containers[1])
    }

    if backend.Docker.Containers[0].Tag != backend.Docker.Containers[1].Tag {
        t.Errorf("backend.Docker.Containers[0].Tag should == backend.Docker.Containers[1].Tag")
    }
    if backend.Docker.Containers[0].Env[0] != backend.Docker.Containers[1].Env[0] {
        t.Errorf("backend.Docker.Containers[0].Env[0] should == backend.Docker.Containers[1].Env[0]")
    }
}

func TestSerializeFromToml(t *testing.T) {
    backend := NewDockerProxyBackendForTest()

    _, err := backend.Serialize()

    if err != nil {
        t.Errorf("Failed to serialize backend")
    }
}

func TestSave(t *testing.T) {
    backend := NewDockerProxyBackendForTest()

    err := backend.Save()

    if err != nil {
        t.Errorf("Failed to save backend")
    }
}

