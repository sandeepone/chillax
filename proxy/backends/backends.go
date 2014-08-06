package backends

import (
    "fmt"
    "bytes"
    "strings"
    "github.com/BurntSushi/toml"
    chillax_storage "github.com/didip/chillax/storage"
)

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
    Containers []*ProxyBackendDockerContainerConfig
}

type ProxyBackendDockerContainerConfig struct {
    Id    string
    Tag   string
    Env   []string
    Host  string
    Ports []string
}

func (pb *ProxyBackend) ProxyName() (string) {
    name := strings.Replace(pb.Path, "/", "", 1)
    return strings.Replace(name, "/", "-", -1)
}

func (pb *ProxyBackend) Serialize() ([]byte, error) {
    var buffer bytes.Buffer
    err := toml.NewEncoder(&buffer).Encode(pb)

    return buffer.Bytes(), err
}

func (pb *ProxyBackend) Save() (error) {
    inBytes, err := pb.Serialize()

    if err == nil {
        err = pb.Storage.Create(fmt.Sprintf("/proxies/%v", pb.ProxyName()), inBytes)
    }
    return err
}

// func (pb *ProxyBackend) IsDocker() (bool) {}

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

// //
// // Docker container
// //
// func (pb *ProxyBackend) StartDockerContainer() (error) {}

// func (pb *ProxyBackend) StopDockerContainer() (error) {}

// func (pb *ProxyBackend) RestartDockerContainer() (error) {}

// func (pb *ProxyBackend) WatchDockerContainer() (error) {}
