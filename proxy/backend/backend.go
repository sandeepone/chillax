package backend

import (
    "os"
    "fmt"
    "time"
    "bytes"
    "strings"
    "strconv"
    "errors"
    "github.com/BurntSushi/toml"
    "github.com/didip/chillax/libstring"
    "github.com/didip/chillax/libprocess"
    dockerclient "github.com/fsouza/go-dockerclient"
    chillax_storage "github.com/didip/chillax/storage"
    chillax_portkeeper "github.com/didip/chillax/portkeeper"
)

const DOCKER_TIMEOUT = uint(5)

func NewProxyBackend(tomlBytes []byte) *ProxyBackend {
    backend := &ProxyBackend{}
    backend.Numprocs = 1

    storage := chillax_storage.NewStorage()

    toml.Decode(string(tomlBytes), backend)

    backend.Storage = storage
    return backend
}

type ProxyBackend struct {
    Path            string
    Command         string
    Numprocs        int
    Delay           string
    Ping            string
    Env             []string
    Storage         chillax_storage.Storer
    Process         *ProxyBackendProcessConfig
    Docker          *ProxyBackendDockerConfig
}

// Process data structure
type ProxyBackendProcessConfig struct {
    HttpPortEnv string `toml:httpportenv`
    Instances   []ProxyBackendProcessInstanceConfig
}

type ProxyBackendProcessInstanceConfig struct {
    Command        string
    Delay          string
    Ping           string
    Env            []string
    ProcessWrapper *libprocess.ProcessWrapper
}


// Docker data structure
type ProxyBackendDockerConfig struct {
    Tag        string
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


//
// Regular processes
//
func (pb *ProxyBackend) NewProcessWrapper(command string) *libprocess.ProcessWrapper {
    pw := &libprocess.ProcessWrapper{
        Name:       pb.ProxyName(),
        Command:    command,
        StopDelay:  pb.Delay,
        StartDelay: pb.Delay,
        Ping:       pb.Ping,
    }
    pw.SetDefaults()
    return pw
}

func (pb *ProxyBackend) NewProxyBackendProcessInstanceConfig(httpPort int) ProxyBackendProcessInstanceConfig {
    pbpi := ProxyBackendProcessInstanceConfig{}

    pbpi.Command = pb.Command
    pbpi.Command = strings.Replace(pbpi.Command, fmt.Sprintf("$%v", pb.Process.HttpPortEnv), strconv.Itoa(httpPort), -1)
    pbpi.Delay   = pb.Delay
    pbpi.Ping    = pb.Ping
    pbpi.Env     = pb.Env
    pbpi.Env     = append(pbpi.Env, fmt.Sprintf("%v=%v", pb.Process.HttpPortEnv, httpPort))

    pbpi.ProcessWrapper = pb.NewProcessWrapper(pbpi.Command)

    return pbpi
}

func (pb *ProxyBackend) CreateProcesses() error {
    if pb.Process == nil { return errors.New("[process] section is missing.") }

    for i := 0; i < pb.Numprocs; i++ {
        hostname, _ := os.Hostname()
        newPort     := chillax_portkeeper.ReservePort(hostname)

        pb.Process.Instances[i] = pb.NewProxyBackendProcessInstanceConfig(newPort)

        err := pb.Save()
        if err != nil { return err }
    }
    return nil
}

func (pb *ProxyBackend) StartProcesses() []error {
    errChan    := make(chan error, pb.Numprocs)
    errorSlice := make([]error, pb.Numprocs)

    if pb.Process == nil { return errorSlice }

    for _, env := range pb.Env {
        envParts := strings.Split(env, "=")
        os.Setenv(envParts[0], envParts[1])
    }

    for i := 0; i < pb.Numprocs; i++ {
        go func(i int) {
            err := pb.Process.Instances[i].ProcessWrapper.StartAndWatch()
            if err == nil {
                err = pb.Save()
            }
            errChan <- err
        }(i)
    }

    for i := 0; i < pb.Numprocs; i++ {
        err := <- errChan
        errorSlice[i] = err
    }
    close(errChan)

    return errorSlice
}

func (pb *ProxyBackend) StopProcesses() []error {
    errorSlice := make([]error, pb.Numprocs)

    if pb.Process == nil { return errorSlice }

    errChan := make(chan error, pb.Numprocs)

    for i, instance := range pb.Process.Instances {
        go func(i int) {
            if instance.ProcessWrapper == nil {
                errChan <- errors.New("Process has not been started.")
            } else {
                err := instance.ProcessWrapper.Stop()

                if err == nil {
                    err = pb.Save()
                }
                errChan <- err
            }
        }(i)
    }

    for i := 0; i < pb.Numprocs; i++ {
        err := <- errChan
        errorSlice[i] = err
    }
    close(errChan)

    return errorSlice
}

func (pb *ProxyBackend) RestartProcesses() []error {
    errorSlice := make([]error, pb.Numprocs)

    if pb.Process == nil { return errorSlice }

    errChan := make(chan error, len(pb.Process.Instances))

    for i, instance := range pb.Process.Instances {
        go func(i int) {
            if instance.ProcessWrapper == nil {
                errChan <- errors.New("Process has not been started.")
            } else {
                err := instance.ProcessWrapper.RestartAndWatch()

                if err == nil {
                    err = pb.Save()
                }
                errChan <- err
            }
        }(i)
    }

    for i := 0; i < pb.Numprocs; i++ {
        err := <- errChan
        errorSlice[i] = err
    }
    close(errChan)

    return errorSlice
}


//
// Docker containers
//
func (pb *ProxyBackend) CreateDockerContainerOptions(publiclyAvailablePort int) *dockerclient.CreateContainerOptions {
    containerOpts := &dockerclient.CreateContainerOptions{}
    containerOpts.Name = fmt.Sprintf("%v-%v", pb.ProxyName(), publiclyAvailablePort)

    containerOpts.Config              = &dockerclient.Config{}
    containerOpts.Config.Image        = pb.Docker.Tag
    containerOpts.Config.Env          = pb.Env
    containerOpts.Config.AttachStdout = true
    containerOpts.Config.AttachStderr = true
    containerOpts.Config.ExposedPorts = make(map[dockerclient.Port]struct{})

    for _, port := range pb.Docker.Ports {
        containerOpts.Config.ExposedPorts[dockerclient.Port(port)] = struct {}{}
    }

    return containerOpts
}

func (pb *ProxyBackend) StartDockerContainerOptions(containerPorts []string) *dockerclient.HostConfig {
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

func (pb *ProxyBackend) PullDockerImage(dockerHost string) error {
    client, err := dockerclient.NewClient(dockerHost)
    if err != nil { return err }

    apiImages, err := client.ListImages(false)

    for _, apiImage := range apiImages {
        for _, remoteRepository := range apiImage.RepoTags {
            if remoteRepository == pb.Docker.Tag {
                return nil
            }
        }
    }

    tagParts   := strings.Split(pb.Docker.Tag, ":")
    repository := tagParts[0]
    tag        := "latest"

    if len(tagParts) >= 2 {
        tag = tagParts[1]
    }

    return client.PullImage(
        dockerclient.PullImageOptions{Repository: repository, Tag: tag},
        dockerclient.AuthConfiguration{},
    )
}

func (pb *ProxyBackend) CreateDockerContainers() error {
    var err error

    numDockers := len(pb.Docker.Hosts)
    if numDockers < 1 { return nil }

    pb.Docker.Containers = make([]ProxyBackendDockerContainerConfig, pb.Numprocs)

    for i := 0; i < pb.Numprocs; i++ {
        dockerHostsIndex := i

        if i >= numDockers {
            dockerHostsIndex = i % numDockers
        }

        dockerHost := pb.Docker.Hosts[dockerHostsIndex]

        // Pull docker image first.
        err = pb.PullDockerImage(dockerHost)
        if err != nil { return err }

        containerConfig, err := pb.CreateDockerContainer(dockerHost)
        if err != nil { return err }

        pb.Docker.Containers[i] = containerConfig

        err = pb.Save()
        if err != nil { return err }
    }
    return err
}

func (pb *ProxyBackend) CreateDockerContainer(dockerHost string) (ProxyBackendDockerContainerConfig, error) {
    containerConfig := ProxyBackendDockerContainerConfig{}

    containerConfig.Tag   = pb.Docker.Tag
    containerConfig.Env   = pb.Env
    containerConfig.Host  = dockerHost
    containerConfig.Ports = make([]string, len(pb.Docker.Ports))

    publiclyAvailablePort := chillax_portkeeper.ReservePort(containerConfig.Host)

    for index, backendPort := range pb.Docker.Ports {
        containerConfig.Ports[index] = fmt.Sprintf("%v:%v", publiclyAvailablePort, backendPort)
    }

    client, err := dockerclient.NewClient(containerConfig.Host)
    if err != nil { return containerConfig, err }

    containerOpts  := pb.CreateDockerContainerOptions(publiclyAvailablePort)
    container, err := client.CreateContainer(*containerOpts)

    if err != nil { return containerConfig, err }

    containerConfig.Id = container.ID

    return containerConfig, err
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

    err = client.StartContainer(containerConfig.Id, pb.StartDockerContainerOptions(containerConfig.Ports))
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

func (pb *ProxyBackend) WatchDockerContainer(containerConfig ProxyBackendDockerContainerConfig) error {
    delayTime, err := time.ParseDuration(pb.Ping)
    if err != nil { return err }

    inspectErrCounter := 0

    for {
        err = pb.InspectAndStartDockerContainer(containerConfig)

        if err != nil {
            if inspectErrCounter > 10 {
                inspectErrCounter = 0
                return err
            }
            inspectErrCounter++
        }

        time.Sleep(delayTime)
    }
    return nil
}
