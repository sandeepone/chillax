package backends

import (
    "os"
    "bufio"
    "testing"
    "io/ioutil"
    dockerclient "github.com/fsouza/go-dockerclient"
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
    if backend.Docker.Env[0] != "HTTP_PORT=8080" {
        t.Errorf("backend.Docker.Env[0] should == HTTP_PORT=8080. Backend.Docker.Env: %s", backend.Docker.Env)
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

func TestIsDocker(t *testing.T) {
    backend := NewDockerProxyBackendForTest()

    if !backend.IsDocker() {
        t.Errorf("Backend1 should be docker")
    }

    fileHandle2, _ := os.Open("./example-process-backend.toml")
    bufReader2     := bufio.NewReader(fileHandle2)
    definition2, _ := ioutil.ReadAll(bufReader2)
    backend2       := NewProxyBackend(definition2)

    if backend2.IsDocker() {
        t.Errorf("Backend2 must not be docker")
    }
}

func TestContainerIds(t *testing.T) {
    backend := NewDockerProxyBackendForTest()

    if backend.ContainerIds()[0] != "abc123" {
        t.Errorf("backend.ContainerIds[0] should == abc123")
    }
    if backend.ContainerIds()[1] != "abc12456" {
        t.Errorf("backend.ContainerIds[1] should == abc12456")
    }
}

func TestCreateDockerContainerOptions(t *testing.T) {
    backend := NewDockerProxyBackendForTest()

    publiclyAvailablePort  := 65536
    createContainerOptions := backend.CreateDockerContainerOptions(publiclyAvailablePort)

    if createContainerOptions.Name == "" {
        t.Errorf("createContainerOptions.Name should not be empty string. Actually: %v", createContainerOptions.Name)
    }
    if len(createContainerOptions.Config.ExposedPorts) != len(backend.Docker.Ports) {
        t.Errorf("Number of ExposedPorts per container should == backend.Docker.Ports. Actually: %v", createContainerOptions.Config.ExposedPorts)
    }
}

func TestCreateDockerContainers(t *testing.T) {
    backend := NewDockerProxyBackendForTest()
    err     := backend.CreateDockerContainers()

    if err != nil {
        t.Errorf("Failed to create Docker containers. Error: %v", err)
    }

    dockerHost := backend.Docker.Hosts[0]

    for _, containerId := range backend.ContainerIds() {
        backend.NewDockerClients()[dockerHost].RemoveContainer(dockerclient.RemoveContainerOptions{ID: containerId})
    }
}

func TestStartStopRestartAndRemoveOneDockerContainer(t *testing.T) {
    backend := NewDockerProxyBackendForTest()

    err := backend.CreateDockerContainers()

    if err != nil {
        t.Errorf("Failed to create Docker containers. Error: %v", err)
    }

    container1 := backend.Docker.Containers[0]
    container2 := backend.Docker.Containers[1]

    err = backend.StartDockerContainer(container1)

    if err != nil {
        t.Errorf("Failed to start Docker container. Error: %v", err)
    }

    err = backend.StopDockerContainer(container1)

    if err != nil {
        t.Errorf("Failed to stop Docker container. Error: %v", err)
    }

    err = backend.RestartDockerContainer(container1)

    if err != nil {
        t.Errorf("Failed to restart Docker container. Error: %v", err)
    }

    backend.StopAndRemoveDockerContainer(container1)
    backend.StopAndRemoveDockerContainer(container2)
}

func TestStartMultipleDockerContainers(t *testing.T) {
    backend := NewDockerProxyBackendForTest()

    backend.CreateDockerContainers()

    errs := backend.StartDockerContainers()

    if errs[0] != nil || errs[1] != nil {
        t.Errorf("Failed to start Docker containers. Errors: %v", errs)
    }

    container1 := backend.Docker.Containers[0]
    container2 := backend.Docker.Containers[1]

    containerJson1, err := backend.InspectDockerContainer(container1)
    containerJson2, err := backend.InspectDockerContainer(container2)

    if err != nil {
        t.Errorf("Failed to inspect Docker container. JSON: %v, Error: %v", containerJson1, err)
    }

    if (containerJson1.ID != container1.Id) || (containerJson2.ID != container2.Id) {
        t.Errorf("ID must match between container JSON and containerConfig")
    }

    if (!containerJson1.State.Running) || (!containerJson2.State.Running) {
        t.Errorf("Container JSON must indicate state: Running")
    }

    backend.StopAndRemoveDockerContainers()
}

func TestInspectAndRestartDockerContainer(t *testing.T) {
    backend := NewDockerProxyBackendForTest()

    err := backend.CreateDockerContainers()

    if err != nil {
        t.Errorf("Failed to create Docker containers. Error: %v", err)
    }

    container1 := backend.Docker.Containers[0]

    containerJson, err := backend.InspectDockerContainer(container1)
    if containerJson.State.Running {
        t.Errorf("Container1 should not be running")
    }

    err = backend.InspectAndStartDockerContainer(container1)
    if err != nil {
        t.Errorf("Failed to inspect and restart Docker container. Error: %v", err)
    }

    containerJson, err = backend.InspectDockerContainer(container1)
    if !containerJson.State.Running {
        t.Errorf("Container1 should be running")
    }

    backend.StopAndRemoveDockerContainers()
}


