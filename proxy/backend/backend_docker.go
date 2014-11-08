package backend

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/chillaxio/chillax/libstring"
	chillax_portkeeper "github.com/chillaxio/chillax/portkeeper"
	dockerclient "github.com/fsouza/go-dockerclient"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type ProxyBackendDockerConfig struct {
	Tag         string
	Hosts       []string
	Ports       []string
	HttpPortEnv string `toml:httpportenv`
	Containers  []ProxyBackendDockerContainerConfig
}

type ProxyBackendDockerContainerConfig struct {
	Id       string
	Tag      string
	Env      []string
	Host     string
	Ports    []string
	MapPorts map[string]int
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
		client, err := pb.NewDockerClient(dockerUri)
		if err == nil {
			dockers[dockerUri] = client
		}
	}

	return dockers
}

func (pb *ProxyBackend) NewDockerClient(dockerUri string) (client *dockerclient.Client, err error) {
	if os.Getenv("DOCKER_CERT_PATH") != "" {
		cert := path.Join(os.Getenv("DOCKER_CERT_PATH"), "cert.pem")
		key := path.Join(os.Getenv("DOCKER_CERT_PATH"), "key.pem")
		ca := path.Join(os.Getenv("DOCKER_CERT_PATH"), "ca.pem")

		client, err = dockerclient.NewTLSClient(dockerUri, cert, key, ca)
	} else {
		client, err = dockerclient.NewClient(dockerUri)
	}

	return client, err
}

func (pb *ProxyBackend) CreateDockerContainerOptions(publiclyAvailablePort int) *dockerclient.CreateContainerOptions {
	containerOpts := &dockerclient.CreateContainerOptions{}
	containerOpts.Name = fmt.Sprintf("%v-%v", pb.ProxyName(), publiclyAvailablePort)

	containerOpts.Config = &dockerclient.Config{}
	containerOpts.Config.Image = pb.Docker.Tag
	containerOpts.Config.Env = pb.Env
	containerOpts.Config.AttachStdout = true
	containerOpts.Config.AttachStderr = true
	containerOpts.Config.ExposedPorts = make(map[dockerclient.Port]struct{})

	for _, port := range pb.Docker.Ports {
		containerOpts.Config.ExposedPorts[dockerclient.Port(port)] = struct{}{}
	}

	return containerOpts
}

func (pb *ProxyBackend) StartDockerContainerOptions(containerPorts []string) *dockerclient.HostConfig {
	config := &dockerclient.HostConfig{}
	// config.ContainerIDFile = "/etc/cidfile"
	config.PortBindings = make(map[dockerclient.Port][]dockerclient.PortBinding)

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
	client, err := pb.NewDockerClient(dockerHost)
	if err != nil {
		return err
	}

	apiImages, err := client.ListImages(false)

	for _, apiImage := range apiImages {
		for _, remoteRepository := range apiImage.RepoTags {
			if remoteRepository == pb.Docker.Tag {
				return nil
			}
		}
	}

	tagParts := strings.Split(pb.Docker.Tag, ":")
	repository := tagParts[0]
	tag := "latest"

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
	if numDockers < 1 {
		return nil
	}

	pb.Docker.Containers = make([]ProxyBackendDockerContainerConfig, pb.Numprocs)

	for i := 0; i < pb.Numprocs; i++ {
		dockerHostsIndex := i

		if i >= numDockers {
			dockerHostsIndex = i % numDockers
		}

		dockerHost := pb.Docker.Hosts[dockerHostsIndex]

		// Pull docker image first.
		err = pb.PullDockerImage(dockerHost)
		if err != nil {
			return err
		}

		containerConfig, err := pb.CreateDockerContainer(dockerHost)
		if err != nil {
			return err
		}

		pb.Docker.Containers[i] = containerConfig

		err = pb.Save()
		if err != nil {
			return err
		}
	}
	return err
}

func (pb *ProxyBackend) NewProxyBackendDockerContainerConfig(dockerHost string) ProxyBackendDockerContainerConfig {
	containerConfig := ProxyBackendDockerContainerConfig{}

	containerConfig.Tag = pb.Docker.Tag
	containerConfig.Env = pb.Env
	containerConfig.Host = dockerHost
	containerConfig.Ports = make([]string, len(pb.Docker.Ports))
	containerConfig.MapPorts = make(map[string]int)

	containerConfig.MapPorts[pb.Docker.HttpPortEnv] = chillax_portkeeper.ReservePort(containerConfig.Host)

	return containerConfig
}

func (pb *ProxyBackend) CreateDockerContainer(dockerHost string) (ProxyBackendDockerContainerConfig, error) {
	containerConfig := pb.NewProxyBackendDockerContainerConfig(dockerHost)

	publiclyAvailablePort := containerConfig.MapPorts[pb.Docker.HttpPortEnv]

	for index, backendPort := range pb.Docker.Ports {
		containerConfig.Ports[index] = fmt.Sprintf("%v:%v", publiclyAvailablePort, backendPort)
	}

	client, err := pb.NewDockerClient(containerConfig.Host)
	if err != nil {
		return containerConfig, err
	}

	containerOpts := pb.CreateDockerContainerOptions(publiclyAvailablePort)
	container, err := client.CreateContainer(*containerOpts)

	if err != nil {
		return containerConfig, err
	}

	containerConfig.Id = container.ID

	err = pb.Save()

	return containerConfig, err
}

func (pb *ProxyBackend) StartDockerContainers() []error {
	errChan := make(chan error, len(pb.Docker.Containers))
	errors := make([]error, len(pb.Docker.Containers))

	for _, containerConfig := range pb.Docker.Containers {
		go func(containerConfig ProxyBackendDockerContainerConfig) {
			errChan <- pb.StartDockerContainer(containerConfig)
		}(containerConfig)
	}

	for i := 0; i < len(pb.Docker.Containers); i++ {
		err := <-errChan
		errors[i] = err
	}
	return errors
}

func (pb *ProxyBackend) InspectDockerContainer(containerConfig ProxyBackendDockerContainerConfig) (*dockerclient.Container, error) {
	client, err := pb.NewDockerClient(containerConfig.Host)
	if err != nil {
		return nil, err
	}

	return client.InspectContainer(containerConfig.Id)
}

func (pb *ProxyBackend) StartDockerContainer(containerConfig ProxyBackendDockerContainerConfig) error {
	client, err := pb.NewDockerClient(containerConfig.Host)
	if err != nil {
		return err
	}

	err = client.StartContainer(containerConfig.Id, pb.StartDockerContainerOptions(containerConfig.Ports))
	return err
}

func (pb *ProxyBackend) StopAndRemoveDockerContainers() error {
	for _, containerConfig := range pb.Docker.Containers {
		client, err := pb.NewDockerClient(containerConfig.Host)
		if err != nil {
			return err
		}

		client.StopContainer(containerConfig.Id, DOCKER_TIMEOUT)

		client.RemoveContainer(dockerclient.RemoveContainerOptions{ID: containerConfig.Id})
	}
	return nil
}

func (pb *ProxyBackend) StopAndRemoveDockerContainer(containerConfig ProxyBackendDockerContainerConfig) error {
	client, err := pb.NewDockerClient(containerConfig.Host)
	if err != nil {
		return err
	}

	err = client.StopContainer(containerConfig.Id, DOCKER_TIMEOUT)
	if err != nil {
		return err
	}

	err = client.RemoveContainer(dockerclient.RemoveContainerOptions{ID: containerConfig.Id})
	return err
}

func (pb *ProxyBackend) StopDockerContainers() []error {
	var errors []error

	for i, containerConfig := range pb.Docker.Containers {
		client, err := pb.NewDockerClient(containerConfig.Host)
		if err != nil {
			errors[i] = err
		} else {
			errors[i] = client.StopContainer(containerConfig.Id, DOCKER_TIMEOUT)
		}
	}
	return errors
}

func (pb *ProxyBackend) StopDockerContainer(containerConfig ProxyBackendDockerContainerConfig) error {
	client, err := pb.NewDockerClient(containerConfig.Host)
	if err != nil {
		return err
	}

	return client.StopContainer(containerConfig.Id, DOCKER_TIMEOUT)
}

func (pb *ProxyBackend) RestartDockerContainers() error {
	for _, containerConfig := range pb.Docker.Containers {
		client, err := pb.NewDockerClient(containerConfig.Host)
		if err != nil {
			return err
		}

		return client.RestartContainer(containerConfig.Id, DOCKER_TIMEOUT)
	}
	return nil
}

func (pb *ProxyBackend) RestartDockerContainer(containerConfig ProxyBackendDockerContainerConfig) error {
	client, err := pb.NewDockerClient(containerConfig.Host)
	if err != nil {
		return err
	}

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
	if err != nil {
		return err
	}

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
