// Package middlewares provides common middleware handlers.
package middlewares

import (
	"github.com/GeertJohan/go.rice"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"github.com/gorilla/context"
	"net/http"
)

func SetStorages(storages *chillax_storage.Storages) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			context.Set(req, "storages", storages)

			next.ServeHTTP(res, req)
		})
	}
}

func SetGoRice(gorice *rice.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			context.Set(req, "gorice", gorice)

			next.ServeHTTP(res, req)
		})
	}
}
