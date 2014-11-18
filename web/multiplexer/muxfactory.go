package multiplexer

import (
	chillax_proxy_handler "github.com/chillaxio/chillax/proxy/handler"
	chillax_storage "github.com/chillaxio/chillax/storage"
	gorilla_mux "github.com/gorilla/mux"
)

// Constructor for MuxFactory.
func NewMuxFactory(storage chillax_storage.Storer, proxyHandlerTomlsFromConfig [][]byte) *MuxFactory {
	mf := &MuxFactory{}

	mf.LoadProxyHandlersFromStorage(storage)
	mf.LoadProxyHandlersFromConfig(proxyHandlerTomlsFromConfig)

	return mf
}

// MuxFactory is responsible for creating new mux and starting proxy backends.
type MuxFactory struct {
	ProxyHandlersFromConfig []*chillax_proxy_handler.ProxyHandler
	ProxyHandlers           []*chillax_proxy_handler.ProxyHandler
}

// LoadProxyHandlersFromConfig loads proxies data from config.
func (mf *MuxFactory) LoadProxyHandlersFromConfig(proxyHandlerTomls [][]byte) {
	mf.ProxyHandlersFromConfig = chillax_proxy_handler.NewProxyHandlers(proxyHandlerTomls)
}

// LoadProxyHandlersFromStorage loads proxies data from storage.
func (mf *MuxFactory) LoadProxyHandlersFromStorage(storage chillax_storage.Storer) {
	proxyHandlers := chillax_proxy_handler.NewProxyHandlersFromStorage(storage)
	mf.ProxyHandlers = proxyHandlers
}

// CreateAndStartBackends create and start new backends as needed per numprocs.
func (mf *MuxFactory) CreateAndStartBackends() []error {
	if len(mf.ProxyHandlers) == 0 {
		mf.ProxyHandlers = mf.ProxyHandlersFromConfig
	}

	errors := make([]error, 0)

	for _, handler := range mf.ProxyHandlers {
		errs := handler.CreateBackends()
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}

		if len(errors) == 0 {
			errs = handler.StartBackends()
			if errs != nil {
				errors = append(errors, errs...)
			}
		}
	}

	return errors
}

func (mf *MuxFactory) StopProxyBackends() []error {
	errors := make([]error, 0)

	for _, handler := range mf.ProxyHandlers {
		errs := handler.StopBackends()
		if errs != nil {
			errors = append(errors, errs...)
		}
	}

	return errors
}

func (mf *MuxFactory) GorillaMuxWithProxyBackends() *gorilla_mux.Router {
	mux := gorilla_mux.NewRouter()

	for _, handler := range mf.ProxyHandlers {
		if handler.Backend.Domain != "" {
			mux.Host(handler.Backend.Domain).Subrouter().HandleFunc(handler.Backend.Path, handler.Function())
		} else {
			mux.HandleFunc(handler.Backend.Path, handler.Function())
		}
	}
	return mux
}
