package multiplexer

import (
	chillax_proxy_handler "github.com/chillaxio/chillax/proxy/handler"
	gorilla_mux "github.com/gorilla/mux"
)

// Constructor for MuxFactory.
func NewMuxFactory(proxyHandlerTomls [][]byte) *MuxFactory {
	mp := &MuxFactory{}

	mp.LoadProxyHandlersFromConfig(proxyHandlerTomls)

	return mp
}

// MuxFactory is responsible for creating new mux.
type MuxFactory struct {
	// Complete list of all proxy paths.
	ProxyHandlers []*chillax_proxy_handler.ProxyHandler
}

func (mp *MuxFactory) LoadProxyHandlersFromConfig(proxyHandlerTomls [][]byte) {
	mp.ProxyHandlers = make([]*chillax_proxy_handler.ProxyHandler, len(proxyHandlerTomls))

	for i, definition := range proxyHandlerTomls {
		mp.ProxyHandlers[i] = chillax_proxy_handler.NewProxyHandler(definition)
	}
}

func (mp *MuxFactory) ReloadProxyHandlers() {
	mp.ProxyHandlers = chillax_proxy_handler.NewProxyHandlers()
}

func (mp *MuxFactory) CreateProxyBackends() []error {
	errors := make([]error, 0)

	for _, handler := range mp.ProxyHandlers {
		err := handler.CreateBackends()
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (mp *MuxFactory) StartProxyBackends() []error {
	errors := make([]error, 0)

	for _, handler := range mp.ProxyHandlers {
		errs := handler.StartBackends()
		if errs != nil {
			errors = append(errors, errs...)
		}
	}

	return errors
}

func (mp *MuxFactory) StopProxyBackends() []error {
	errors := make([]error, 0)

	for _, handler := range mp.ProxyHandlers {
		errs := handler.StopBackends()
		if errs != nil {
			errors = append(errors, errs...)
		}
	}

	return errors
}

func (mp *MuxFactory) GorillaMuxWithProxyBackends() *gorilla_mux.Router {
	mux := gorilla_mux.NewRouter()

	for _, handler := range mp.ProxyHandlers {
		if handler.Backend.Domain != "" {
			mux.Host(handler.Backend.Domain).Subrouter().HandleFunc(handler.Backend.Path, handler.Function())
		} else {
			mux.HandleFunc(handler.Backend.Path, handler.Function())
		}
	}
	return mux
}
