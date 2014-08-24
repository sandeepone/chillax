package muxproducer

import (
    "os"
    "path"
    "path/filepath"
    "bufio"
    "io/ioutil"
    "github.com/didip/chillax/libenv"
    gorilla_mux "github.com/gorilla/mux"
    chillax_proxy_handler "github.com/didip/chillax/proxy/handler"
)

func NewMuxProducer() (*MuxProducer, error) {
    mp := &MuxProducer{}

    err := mp.LoadProxyHandlersFromConfig()
    if err != nil { return mp, err }

    return mp, err
}

type MuxProducer struct {
    ProxyHandlers []*chillax_proxy_handler.ProxyHandler
}

func (mp *MuxProducer) LoadProxyHandlersFromConfig() error {
    defaultProxyBackendsDir := libenv.EnvWithDefault("DEFAULT_PROXY_BACKENDS_DIR", "")

    if defaultProxyBackendsDir != "" {
        files, err := filepath.Glob(path.Join(defaultProxyBackendsDir, "*.toml"))
        if err != nil { return err }

        mp.ProxyHandlers = make([]*chillax_proxy_handler.ProxyHandler, len(files))

        for i, fullFilename := range files {
            fileHandle, err := os.Open(fullFilename)

            if err != nil { return err }

            bufReader       := bufio.NewReader(fileHandle)
            definition, err := ioutil.ReadAll(bufReader)

            if err != nil { return err }

            mp.ProxyHandlers[i] = chillax_proxy_handler.NewProxyHandler(definition)
        }
    }
    return nil
}

func (mp *MuxProducer) ReloadProxyHandlers() {
    mp.ProxyHandlers = chillax_proxy_handler.NewProxyHandlers()
}

func (mp *MuxProducer) CreateProxyBackends() []error {
    errors := make([]error, 0)

    for _, handler := range mp.ProxyHandlers {
        err := handler.CreateBackends()
        if err != nil { errors = append(errors, err) }
    }

    return errors
}

func (mp *MuxProducer) StartProxyBackends() []error {
    errors := make([]error, 0)

    for _, handler := range mp.ProxyHandlers {
        errs := handler.StartBackends()
        if errs != nil { errors = append(errors, errs...) }
    }

    return errors
}

func (mp *MuxProducer) StopProxyBackends() []error {
    errors := make([]error, 0)

    for _, handler := range mp.ProxyHandlers {
        errs := handler.StopBackends()
        if errs != nil { errors = append(errors, errs...) }
    }

    return errors
}

func (mp *MuxProducer) GorillaMuxWithProxyBackends() *gorilla_mux.Router {
    mux := gorilla_mux.NewRouter()

    for _, handler := range mp.ProxyHandlers {
        mux.HandleFunc(handler.Backend.Path, handler.Function())
    }
    return mux
}
