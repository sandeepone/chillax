package backend

import (
    "os"
    "bufio"
    "testing"
    "io/ioutil"
    "github.com/didip/chillax/libtime"
)

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

    backend.StopAndRemoveDockerContainers()
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

func TestWatchDockerContainer(t *testing.T) {
    backend := NewDockerProxyBackendForTest()
    backend.Ping = "50ms"

    err := backend.CreateDockerContainers()

    if err != nil {
        t.Errorf("Failed to create Docker containers. Error: %v", err)
    }

    container1 := backend.Docker.Containers[0]
    containerJson, err := backend.InspectDockerContainer(container1)
    if containerJson.State.Running {
        t.Errorf("Container1 should not be running")
    }

    go backend.WatchDockerContainer(container1)

    libtime.SleepString("500ms")

    containerJson, err = backend.InspectDockerContainer(container1)
    if !containerJson.State.Running {
        t.Errorf("Container1 should be running")
    }

    err = backend.StopDockerContainer(container1)
    if err != nil {
        t.Errorf("Container1 should have been stopped. Error: %v", err)
    }

    libtime.SleepString("750ms")

    containerJson, err = backend.InspectDockerContainer(container1)
    if !containerJson.State.Running {
        t.Errorf("Container1 should still be running. Container state: %v", containerJson.State)
    }

    backend.StopAndRemoveDockerContainers()
}

