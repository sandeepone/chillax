package handler

import(
    "os"
    "net"
    "net/url"
    "net/http"
    "net/http/httputil"
    "strconv"
    "github.com/didip/chillax/libstring"
    chillax_storage "github.com/didip/chillax/storage"
    chillax_proxy_backend "github.com/didip/chillax/proxy/backend"
)

func NewProxyHandlers() []*ProxyHandler {
    storage       := chillax_storage.NewStorage()
    handlers      := make([]*ProxyHandler, 0)
    allProxies, _ := storage.List("/proxies")

    for index, proxyName := range allProxies {
        definition, _  := storage.Get("/proxies/" + proxyName)
        handlers[index] = NewProxyHandler(definition)
    }
    return handlers
}

func NewProxyHandler(tomlBytes []byte) *ProxyHandler {
    handler := &ProxyHandler{}
    handler.Backend = chillax_proxy_backend.NewProxyBackend(tomlBytes)

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

func (ph *ProxyHandler) RealIP(r *http.Request) string {
    host, _, _ := net.SplitHostPort(r.RemoteAddr)
    return host
}

func (ph *ProxyHandler) BackendHosts() []string {
    var backendHosts []string

    if ph.Backend.IsDocker() {
        backendHosts = make([]string, len(ph.Backend.Docker.Containers))

        for index, container := range ph.Backend.Docker.Containers {
            backendHosts[index] = libstring.HostWithoutPort(container.Host) + ":" + strconv.Itoa(container.MapPorts[ph.Backend.Process.HttpPortEnv])
        }
    } else {
        backendHosts = make([]string, len(ph.Backend.Process.Instances))
        hostname, _ := os.Hostname()

        for index, instance := range ph.Backend.Process.Instances {
            backendHosts[index] = hostname + ":" + strconv.Itoa(instance.MapPorts[ph.Backend.Process.HttpPortEnv])
        }
    }
    return backendHosts
}

func (ph *ProxyHandler) Function() func(http.ResponseWriter, *http.Request) {
    url       := &url.URL{}
    url.Scheme = "http"
    url.Path   = ph.Backend.Path

    if ph.Backend.IsDocker() {
    } else {
    }

    return func(w http.ResponseWriter, r *http.Request) {
        r.URL.Path = "/"
        r.Header.Add("X-Real-IP", ph.RealIP(r))

        proxy := httputil.NewSingleHostReverseProxy(url)
        proxy.ServeHTTP(w, r)
    }
}
