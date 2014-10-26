package handler

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/chillaxio/chillax/libstring"
	chillax_proxy_backend "github.com/chillaxio/chillax/proxy/backend"
	chillax_proxy_selectors "github.com/chillaxio/chillax/proxy/selectors"
	chillax_storage "github.com/chillaxio/chillax/storage"
	chillax_web_pingers "github.com/chillaxio/chillax/web/pingers"
)

// NewProxyHandlers is constructor for all ProxyHandler.
func NewProxyHandlers() []*ProxyHandler {
	var handlers []*ProxyHandler

	storage := chillax_storage.NewStorage()
	allProxies, _ := storage.List("/proxies")

	for index, proxyName := range allProxies {
		definition, _ := storage.Get("/proxies/" + proxyName)
		handlers[index] = NewProxyHandler(definition)
	}
	return handlers
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

// PingData returns ping data per host.
func (ph *ProxyHandler) PingData() map[string]bool {
	storage := chillax_storage.NewStorage()
	pingData := make(map[string]bool)

	hosts, err := storage.List("/pingers")
	if err != nil {
		return pingData
	}

	for _, host := range hosts {
		// Set defaults
		pingData[host] = false

		pg, err := chillax_web_pingers.LoadPingerGroupFromStorage(host)
		if err == nil {
			for uri := range pg.Pingers {
				if strings.Contains(uri, ph.Backend.Path) {
					pingData[host] = pg.PingersCheck[uri]
					continue
				}
			}
		}
	}
	return pingData
}

// CreateBackends instantiate all backends.
func (ph *ProxyHandler) CreateBackends() error {
	var err error

	if ph.Backend.IsDocker() {
		err = ph.Backend.CreateDockerContainers()
	} else {
		err = ph.Backend.CreateProcesses()
	}
	return err
}

// StartBackends start all backends.
func (ph *ProxyHandler) StartBackends() []error {
	var errors []error

	if ph.Backend.IsDocker() {
		errors = ph.Backend.StartDockerContainers()
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
