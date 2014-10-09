package muxproducer

import (
	chillax_proxy_handler "github.com/didip/chillax/proxy/handler"
	gorilla_mux "github.com/gorilla/mux"
)

func NewMuxProducer(proxyHandlerTomls [][]byte) *MuxProducer {
	mp := &MuxProducer{}

	mp.LoadProxyHandlersFromConfig(proxyHandlerTomls)

	return mp
}

type MuxProducer struct {
	ProxyHandlers []*chillax_proxy_handler.ProxyHandler
}

func (mp *MuxProducer) LoadProxyHandlersFromConfig(proxyHandlerTomls [][]byte) {
	mp.ProxyHandlers = make([]*chillax_proxy_handler.ProxyHandler, len(proxyHandlerTomls))

	for i, definition := range proxyHandlerTomls {
		mp.ProxyHandlers[i] = chillax_proxy_handler.NewProxyHandler(definition)
	}
}

func (mp *MuxProducer) ReloadProxyHandlers() {
	mp.ProxyHandlers = chillax_proxy_handler.NewProxyHandlers()
}

func (mp *MuxProducer) CreateProxyBackends() []error {
	errors := make([]error, 0)

	for _, handler := range mp.ProxyHandlers {
		err := handler.CreateBackends()
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (mp *MuxProducer) StartProxyBackends() []error {
	errors := make([]error, 0)

	for _, handler := range mp.ProxyHandlers {
		errs := handler.StartBackends()
		if errs != nil {
			errors = append(errors, errs...)
		}
	}

	return errors
}

func (mp *MuxProducer) StopProxyBackends() []error {
	errors := make([]error, 0)

	for _, handler := range mp.ProxyHandlers {
		errs := handler.StopBackends()
		if errs != nil {
			errors = append(errors, errs...)
		}
	}

	return errors
}

func (mp *MuxProducer) GorillaMuxWithProxyBackends() *gorilla_mux.Router {
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
