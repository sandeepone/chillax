package server

import (
    "os"
    "path"
    "path/filepath"
    "bufio"
    "io/ioutil"

    "github.com/zenazn/goji"
    "github.com/didip/chillax/libenv"
    chillax_proxy_handler "github.com/didip/chillax/proxy/handler"
)

func NewHttpServer() (*HttpServer, error) {
    server := &HttpServer{}

    err := server.LoadProxyHandlersFromConfig()
    if err != nil { return server, err }

    return server, err
}

type HttpServer struct {
    ProxyHandlers []*chillax_proxy_handler.ProxyHandler
}

func (h *HttpServer) LoadProxyHandlersFromConfig() error {
    defaultProxyBackendsDir := libenv.EnvWithDefault("DEFAULT_PROXY_BACKENDS_DIR", "")

    if defaultProxyBackendsDir != "" {
        files, err := filepath.Glob(path.Join(defaultProxyBackendsDir, "*.toml"))
        if err != nil { return err }

        h.ProxyHandlers = make([]*chillax_proxy_handler.ProxyHandler, len(files))

        for i, fullFilename := range files {
            fileHandle, err := os.Open(fullFilename)

            if err != nil { return err }

            bufReader       := bufio.NewReader(fileHandle)
            definition, err := ioutil.ReadAll(bufReader)

            if err != nil { return err }

            h.ProxyHandlers[i] = chillax_proxy_handler.NewProxyHandler(definition)
        }
    }
    return nil
}

func (h *HttpServer) ReloadProxyHandlers() {
    h.ProxyHandlers = chillax_proxy_handler.NewProxyHandlers()
}

func (h *HttpServer) CreateProxyBackends() []error {
    errors := make([]error, 0)

    for _, handler := range h.ProxyHandlers {
        err := handler.CreateBackends()
        if err != nil { errors = append(errors, err) }
    }

    return errors
}

func (h *HttpServer) StartProxyBackends() []error {
    errors := make([]error, 0)

    for _, handler := range h.ProxyHandlers {
        errs := handler.StartBackends()
        if errs != nil { errors = append(errors, errs...) }
    }

    return errors
}

func (h *HttpServer) StopProxyBackends() []error {
    errors := make([]error, 0)

    for _, handler := range h.ProxyHandlers {
        errs := handler.StopBackends()
        if errs != nil { errors = append(errors, errs...) }
    }

    return errors
}

func (h *HttpServer) Serve() {
    goji.Serve()
}

