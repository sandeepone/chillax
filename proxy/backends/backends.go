package backends

import (
    "fmt"
    "bytes"
    "strings"
    "github.com/BurntSushi/toml"
    dockerclient "github.com/fsouza/go-dockerclient"
    "github.com/didip/chillax/libstring"
    chillax_storage "github.com/didip/chillax/storage"
    chillax_dockerinventory "github.com/didip/chillax/dockerinventory"
)

const DOCKER_TIMEOUT = uint(5)

func NewProxyBackend(tomlBytes []byte) *ProxyBackend {
    backend := &ProxyBackend{}
    storage := chillax_storage.NewStorage()

    toml.Decode(string(tomlBytes), backend)

    backend.Storage = storage
    return backend
}

type ProxyBackend struct {
    Storage chillax_storage.Storer
    Path    string
    Command string
    Delay   string
    Ping    string
    Docker  *ProxyBackendDockerConfig
}

type ProxyBackendDockerConfig struct {
    Tag        string
    Env        []string
    Numprocs   int
    Hosts      []string
    Ports      []string
    Containers []ProxyBackendDockerContainerConfig
}

type ProxyBackendDockerContainerConfig struct {
    Id    string
    Tag   string
    Env   []string
    Host  string
    Ports []string
}

func (pb *ProxyBackend) ProxyName() string {
    name := strings.Replace(pb.Path, "/", "", 1)
    return strings.Replace(name, "/", "-", -1)
}

func (pb *ProxyBackend) Serialize() ([]byte, error) {
    var buffer bytes.Buffer
    err := toml.NewEncoder(&buffer).Encode(pb)

    return buffer.Bytes(), err
}

func (pb *ProxyBackend) Save() error {
    inBytes, err := pb.Serialize()

    if err == nil {
        err = pb.Storage.Create(fmt.Sprintf("/proxies/%v", pb.ProxyName()), inBytes)
    }
    return err
}

func (pb *ProxyBackend) IsDocker() bool {
    if pb.Docker != nil && pb.Docker.Tag != "" && len(pb.Docker.Hosts) > 0 {
        return true
    }
    return false
}

func (pb *ProxyBackend) NewDockerClients() map[string]*dockerclient.Client {
    dockers := make(map[string]*dockerclient.Client)

    for _, dockerUri := range pb.Docker.Hosts {
        client, err := dockerclient.NewClient(dockerUri)

        if err == nil {
            dockers[dockerUri] = client
        }
    }

    return dockers
}


// // protocol can be: TCP or HTTP
// func (pb *ProxyBackend) Ping(protocol string) (bool) {}

// func (pb *ProxyBackend) PingDocker() (bool) {}


// func (pb *ProxyBackend) Start() (error) {}

// func (pb *ProxyBackend) Stop() (error) {}

// func (pb *ProxyBackend) Restart() (error) {}

// func (pb *ProxyBackend) Watch() (error) {}

// //
// // Regular process
// //
// func (pb *ProxyBackend) StartProcess() (error) {}

// func (pb *ProxyBackend) StopProcess() (error) {}

// func (pb *ProxyBackend) RestartProcess() (error) {}

// func (pb *ProxyBackend) WatchProcess() (error) {}

//
// Docker container
//
func (pb *ProxyBackend) ContainerIds() []string {
    ids := make([]string, len(pb.Docker.Containers))

    for index, containerConfig := range pb.Docker.Containers {
        ids[index] = containerConfig.Id
    }
    return ids
}

func (pb *ProxyBackend) CreateDockerContainerOptions(publiclyAvailablePort int) *dockerclient.CreateContainerOptions {
    containerOpts := &dockerclient.CreateContainerOptions{}
    containerOpts.Name = fmt.Sprintf("%v-%v", pb.ProxyName(), publiclyAvailablePort)

    containerOpts.Config              = &dockerclient.Config{}
    containerOpts.Config.Image        = pb.Docker.Tag
    containerOpts.Config.Env          = pb.Docker.Env
    containerOpts.Config.AttachStdout = true
    containerOpts.Config.AttachStderr = true
    containerOpts.Config.ExposedPorts = make(map[dockerclient.Port]struct{})

    for _, port := range pb.Docker.Ports {
        containerOpts.Config.ExposedPorts[dockerclient.Port(port)] = struct {}{}
    }

    return containerOpts
}

func (pb *ProxyBackend) StartContainerOptions(containerPorts []string) *dockerclient.HostConfig {
    config := &dockerclient.HostConfig{}
    config.ContainerIDFile = "/etc/cidfile"
    config.PortBindings    = make(map[dockerclient.Port][]dockerclient.PortBinding)

    for _, ports := range containerPorts {
        hostIp, hostPort, containerPort := libstring.SplitDockerPorts(ports)

        hostPortBinding := &dockerclient.PortBinding{}
        hostPortBinding.HostPort = hostPort

        if hostIp != "" {
            hostPortBinding.HostIp = hostIp
        }

        config.PortBindings[dockerclient.Port(containerPort)] = append(config.PortBindings[dockerclient.Port(containerPort)], *hostPortBinding)
    }

    return config
}

func (pb *ProxyBackend) CreateDockerContainers() error {
    var err error

    numDockers := len(pb.Docker.Hosts)
    if numDockers < 1 { return nil }

    pb.Docker.Containers = make([]ProxyBackendDockerContainerConfig, pb.Docker.Numprocs)

    dockerClients := pb.NewDockerClients()

    for i := 0; i < pb.Docker.Numprocs; i++ {
        dockerHostsIndex := i

        if i >= numDockers {
            dockerHostsIndex = i % numDockers
        }

        pb.Docker.Containers[i] = ProxyBackendDockerContainerConfig{}

        pb.Docker.Containers[i].Tag   = pb.Docker.Tag
        pb.Docker.Containers[i].Env   = pb.Docker.Env
        pb.Docker.Containers[i].Host  = pb.Docker.Hosts[dockerHostsIndex]
        pb.Docker.Containers[i].Ports = make([]string, len(pb.Docker.Ports))

        publiclyAvailablePort := chillax_dockerinventory.ReservePort(pb.Docker.Containers[i].Host)

        for index, backendPort := range pb.Docker.Ports {
            pb.Docker.Containers[i].Ports[index] = fmt.Sprintf("%v:%v", publiclyAvailablePort, backendPort)
        }

        containerOpts  := pb.CreateDockerContainerOptions(publiclyAvailablePort)
        container, err := dockerClients[pb.Docker.Containers[i].Host].CreateContainer(*containerOpts)
        if err != nil { return err }

        pb.Docker.Containers[i].Id = container.ID

        err = pb.Save()
        if err != nil { return err }
    }
    return err
}

func (pb *ProxyBackend) StartDockerContainers() []error {
    errChan := make(chan error, len(pb.Docker.Containers))
    errors  := make([]error, len(pb.Docker.Containers))

    for _, containerConfig := range pb.Docker.Containers {
        go func(containerConfig ProxyBackendDockerContainerConfig) {
            errChan <- pb.StartDockerContainer(containerConfig)
        } (containerConfig);
    }

    for i := 0; i < len(pb.Docker.Containers); i++ {
        err := <- errChan
        errors[i] = err
    }
    return errors
}

func (pb *ProxyBackend) InspectDockerContainer(containerConfig ProxyBackendDockerContainerConfig) (*dockerclient.Container, error) {
    client, err := dockerclient.NewClient(containerConfig.Host)
    if err != nil { return nil, err }

    return client.InspectContainer(containerConfig.Id)
}

func (pb *ProxyBackend) StartDockerContainer(containerConfig ProxyBackendDockerContainerConfig) error {
    client, err := dockerclient.NewClient(containerConfig.Host)
    if err != nil { return err }

    err = client.StartContainer(containerConfig.Id, pb.StartContainerOptions(containerConfig.Ports))
    return err
}

func (pb *ProxyBackend) StopAndRemoveDockerContainers() error {
    for _, containerConfig := range pb.Docker.Containers {
        client, err := dockerclient.NewClient(containerConfig.Host)
        if err != nil { return err }

        client.StopContainer(containerConfig.Id, DOCKER_TIMEOUT)

        client.RemoveContainer(dockerclient.RemoveContainerOptions{ID: containerConfig.Id})
    }
    return nil
}

func (pb *ProxyBackend) StopAndRemoveDockerContainer(containerConfig ProxyBackendDockerContainerConfig) error {
    client, err := dockerclient.NewClient(containerConfig.Host)
    if err != nil { return err }

    err = client.StopContainer(containerConfig.Id, DOCKER_TIMEOUT)
    if err != nil { return err }

    err = client.RemoveContainer(dockerclient.RemoveContainerOptions{ID: containerConfig.Id})
    return err
}

func (pb *ProxyBackend) StopDockerContainers() error {
    for _, containerConfig := range pb.Docker.Containers {
        client, err := dockerclient.NewClient(containerConfig.Host)
        if err != nil { return err }

        return client.StopContainer(containerConfig.Id, DOCKER_TIMEOUT)
    }
    return nil
}

func (pb *ProxyBackend) StopDockerContainer(containerConfig ProxyBackendDockerContainerConfig) error {
    client, err := dockerclient.NewClient(containerConfig.Host)
    if err != nil { return err }

    return client.StopContainer(containerConfig.Id, DOCKER_TIMEOUT)
}

func (pb *ProxyBackend) RestartDockerContainers() error {
    for _, containerConfig := range pb.Docker.Containers {
        client, err := dockerclient.NewClient(containerConfig.Host)
        if err != nil { return err }

        return client.RestartContainer(containerConfig.Id, DOCKER_TIMEOUT)
    }
    return nil
}

func (pb *ProxyBackend) RestartDockerContainer(containerConfig ProxyBackendDockerContainerConfig) error {
    client, err := dockerclient.NewClient(containerConfig.Host)
    if err != nil { return err }

    return client.RestartContainer(containerConfig.Id, DOCKER_TIMEOUT)
}

func (pb *ProxyBackend) InspectAndStartDockerContainer(containerConfig ProxyBackendDockerContainerConfig) error {
    jsonData, err := pb.InspectDockerContainer(containerConfig)
    if err == nil && !jsonData.State.Running {
        err = pb.StartDockerContainer(containerConfig)
    }
    return err
}

// func (pb *ProxyBackend) WatchDockerContainer() error {}
