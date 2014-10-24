package backend

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"strings"
)

const DOCKER_TIMEOUT = uint(5)

func NewProxyBackend(tomlBytes []byte) (*ProxyBackend, error) {
	backend := &ProxyBackend{}
	backend.Numprocs = 1

	storage := chillax_storage.NewStorage()

	_, err := toml.Decode(string(tomlBytes), backend)
	if err != nil {
		return nil, err
	}

	backend.Storage = storage

	return backend, err
}

type ProxyBackend struct {
	Domain   string
	Path     string
	Command  string
	Numprocs int
	Delay    string
	Ping     string
	Env      []string
	Storage  chillax_storage.Storer
	Process  *ProxyBackendProcessConfig
	Docker   *ProxyBackendDockerConfig
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

func (pb *ProxyBackend) UpNumprocs() int {
	if pb.IsDocker() {
		return len(pb.Docker.Containers)
	} else {
		return len(pb.Process.Instances)
	}
	return 0
}
