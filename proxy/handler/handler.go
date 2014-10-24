package handler

import (
	"github.com/chillaxio/chillax/libstring"
	chillax_proxy_backend "github.com/chillaxio/chillax/proxy/backend"
	chillax_proxy_selectors "github.com/chillaxio/chillax/proxy/selectors"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func NewProxyHandlers() []*ProxyHandler {
	storage := chillax_storage.NewStorage()
	handlers := make([]*ProxyHandler, 0)
	allProxies, _ := storage.List("/proxies")

	for index, proxyName := range allProxies {
		definition, _ := storage.Get("/proxies/" + proxyName)
		handlers[index] = NewProxyHandler(definition)
	}
	return handlers
}

func NewProxyHandler(tomlBytes []byte) *ProxyHandler {
	backend, _ := chillax_proxy_backend.NewProxyBackend(tomlBytes)

	handler := &ProxyHandler{}
	handler.Backend = backend

	return handler
}

type ProxyHandler struct {
	Backend *chillax_proxy_backend.ProxyBackend
}

func (ph *ProxyHandler) CreateBackends() error {
	var err error

	if ph.Backend.IsDocker() {
		err = ph.Backend.CreateDockerContainers()
	} else {
		err = ph.Backend.CreateProcesses()
	}
	return err
}

func (ph *ProxyHandler) StartBackends() []error {
	var errors []error

	if ph.Backend.IsDocker() {
		errors = ph.Backend.StartDockerContainers()
	} else {
		errors = ph.Backend.StartProcesses()
	}
	return errors
}

func (ph *ProxyHandler) StopBackends() []error {
	var errors []error

	if ph.Backend.IsDocker() {
		errors = ph.Backend.StopDockerContainers()
	} else {
		errors = ph.Backend.StopProcesses()
	}
	return errors
}

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
