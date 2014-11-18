package backend

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"strings"
)

const DOCKER_TIMEOUT = uint(5)

func LoadProxyBackendByName(proxyName string) (*ProxyBackend, error) {
	storage := chillax_storage.NewStorage()

	definition, err := storage.Get(fmt.Sprintf("/proxies/%v", proxyName))

	if err != nil {
		return nil, err
	}

	return NewProxyBackend(definition)
}

func DeleteProxyBackendByName(proxyName string) error {
	storage := chillax_storage.NewStorage()

	err := storage.Delete(fmt.Sprintf("/proxies/%v", proxyName))

	return err
}

func NewProxyBackend(tomlBytes []byte) (*ProxyBackend, error) {
	backend := &ProxyBackend{}
	backend.Numprocs = 1

	storage := chillax_storage.NewStorage()

	_, err := toml.Decode(string(tomlBytes), backend)
	if err != nil {
		return nil, err
	}

	backend.storage = storage

	return backend, err
}

func UpdateProxyBackend(backend *ProxyBackend, tomlBytes []byte) (*ProxyBackend, error) {
	_, err := toml.Decode(string(tomlBytes), backend)
	if err != nil {
		return nil, err
	}

	err = backend.Save()
	if err != nil {
		return nil, err
	}

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
	Process  *ProxyBackendProcessConfig `json:",omitempty"`
	Docker   *ProxyBackendDockerConfig  `json:",omitempty"`
	storage  chillax_storage.Storer     `json:"-"`
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
		err = pb.storage.Create(fmt.Sprintf("/proxies/%v", pb.ProxyName()), inBytes)
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
