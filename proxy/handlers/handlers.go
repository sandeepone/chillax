package chillax_proxy_handlers

import (
    "github.com/BurntSushi/toml"
)

func NewProxyHandler(tomlBytes []byte) *ProxyHandler {
    handler := &ProxyHandler{}
    storage := NewStorage()

    toml.Decode(string(tomlBytes), handler)

    handler.storage = storage
    return handler
}

type ProxyHandler struct {}

func (ph *ProxyHandler) Serialize() ([]byte, error) {}

func (ph *ProxyHandler) IsDocker() (bool) {}

// protocol can be: TCP or HTTP
func (ph *ProxyHandler) Ping(protocol string) (bool) {}


func (ph *ProxyHandler) Start() (error) {}

func (ph *ProxyHandler) Stop() (error) {}

func (ph *ProxyHandler) Restart() (error) {}

func (ph *ProxyHandler) Watch() (error) {}

//
// Regular process
//
func (ph *ProxyHandler) StartProcess() (error) {}

func (ph *ProxyHandler) StopProcess() (error) {}

func (ph *ProxyHandler) RestartProcess() (error) {}

func (ph *ProxyHandler) WatchProcess() (error) {}

//
// Docker container
//
func (ph *ProxyHandler) StartDockerContainer() (error) {}

func (ph *ProxyHandler) StopDockerContainer() (error) {}

func (ph *ProxyHandler) RestartDockerContainer() (error) {}

func (ph *ProxyHandler) WatchDockerContainer() (error) {}
