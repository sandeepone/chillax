package handler

import (
	"github.com/chillaxio/chillax/libstring"
	chillax_proxy_backend "github.com/chillaxio/chillax/proxy/backend"
	chillax_proxy_selectors "github.com/chillaxio/chillax/proxy/selectors"
	chillax_storage "github.com/chillaxio/chillax/storage"
	chillax_web_pingers "github.com/chillaxio/chillax/web/pingers"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

// LoadProxyHandlerByName creates ProxyHandler given proxy backend name.
func LoadProxyHandlerByName(proxyName string) (*ProxyHandler, error) {
	backend, err := chillax_proxy_backend.LoadProxyBackendByName(proxyName)
	if err != nil {
		return nil, err
	}

	handler := &ProxyHandler{}
	handler.Backend = backend

	return handler, nil
}

// NewProxyHandlers is constructor for all ProxyHandlers given as TOML.
func NewProxyHandlers(proxyHandlerTomls [][]byte) []*ProxyHandler {
	proxyHandlers := make([]*ProxyHandler, len(proxyHandlerTomls))

	for i, definition := range proxyHandlerTomls {
		proxyHandlers[i] = NewProxyHandler(definition)
	}
	return proxyHandlers
}

// NewProxyHandlersFromStorage loads proxies data from storage.
func NewProxyHandlersFromStorage(storage chillax_storage.Storer) []*ProxyHandler {
	proxyNames, err := storage.List("/proxies")
	proxyHandlers := make([]*ProxyHandler, 0)

	if err == nil {
		for _, proxyName := range proxyNames {
			proxyHandlerToml, err := storage.Get("/proxies/" + proxyName)
			if err == nil {
				handler := NewProxyHandler(proxyHandlerToml)
				proxyHandlers = append(proxyHandlers, handler)
			}
		}
	}
	return proxyHandlers
}

// NewProxyHandler is constructor for one ProxyHandler.
func NewProxyHandler(tomlBytes []byte) *ProxyHandler {
	backend, _ := chillax_proxy_backend.NewProxyBackend(tomlBytes)

	handler := &ProxyHandler{}
	handler.Backend = backend

	return handler
}

// ProxyHandler is a struct that represents 1 proxy endpoint.
type ProxyHandler struct {
	Backend *chillax_proxy_backend.ProxyBackend
}

// PingBool returns ping data per proxy.BackendPath for all hosts.
func (ph *ProxyHandler) PingBool() map[string]bool {
	return chillax_web_pingers.PingBoolGivenProxyPath(ph.Backend.Path)
}

// PingLastCheck returns ping last check Unix Nano per host.
func (ph *ProxyHandler) PingLastCheck(host string) int64 {
	return chillax_web_pingers.PingLastCheckGivenProxyPathAndHost(ph.Backend.Path, host)
}

// CreateBackends instantiate all backends.
func (ph *ProxyHandler) CreateBackends() []error {
	var errors []error

	if ph.Backend.IsDocker() {
		errors = ph.Backend.CreateDockerContainers()
	} else {
		errors = ph.Backend.CreateProcesses()
	}
	return errors
}

// StartBackends start all backends.
func (ph *ProxyHandler) StartBackends() []error {
	var errors []error

	if ph.Backend.IsDocker() {
		errors = ph.Backend.InspectAndStartDockerContainers()
		if len(errors) == 0 {
			ph.Backend.WatchDockerContainers()
		}
	} else {
		errors = ph.Backend.StartProcesses()
	}
	return errors
}

// StopBackends start all backends.
func (ph *ProxyHandler) StopBackends() []error {
	var errors []error

	if ph.Backend.IsDocker() {
		errors = ph.Backend.StopDockerContainers()
	} else {
		errors = ph.Backend.StopProcesses()
	}
	return errors
}

// RestartBackends rolling restarts all backends.
func (ph *ProxyHandler) RestartBackends() []error {
	var errors []error

	if ph.Backend.IsDocker() {
		errors = ph.Backend.CreateStartAndStopRemoveOldDockerContainers()
	} else {
		errors = ph.Backend.StopProcesses()
		if len(errors) == 0 {
			errors = ph.Backend.StartProcesses()
		}
	}
	return errors
}

// RealIP extracts actual client IP address.
func (ph *ProxyHandler) RealIP(r *http.Request) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func (ph *ProxyHandler) BackendHosts() []string {
	var backendHosts []string

	if ph.Backend.IsDocker() {
		backendHosts = make([]string, len(ph.Backend.Docker.Containers))

		for index, container := range ph.Backend.Docker.Containers {
			backendHosts[index] = libstring.HostWithoutPort(container.Host) + ":" + strconv.Itoa(container.MapPorts[ph.Backend.Docker.HttpPortEnv])
		}
	} else {
		backendHosts = make([]string, len(ph.Backend.Process.Instances))

		for index, instance := range ph.Backend.Process.Instances {
			backendHosts[index] = libstring.HostWithoutPort(instance.Host) + ":" + strconv.Itoa(instance.MapPorts[ph.Backend.Process.HttpPortEnv])
		}
	}
	return backendHosts
}

func (ph *ProxyHandler) ChooseBackendHost() string {
	backendHosts := ph.BackendHosts()
	selector := chillax_proxy_selectors.NewRoundRobin(backendHosts)
	return selector.Choose()
}

func (ph *ProxyHandler) BackendUrl() *url.URL {
	url := &url.URL{}
	url.Scheme = "http"
	url.Path = "/"
	url.Host = ph.ChooseBackendHost()
	return url
}

func (ph *ProxyHandler) Function() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("X-Real-IP", ph.RealIP(r))

		proxy := httputil.NewSingleHostReverseProxy(ph.BackendUrl())
		proxy.ServeHTTP(w, r)
	}
}
