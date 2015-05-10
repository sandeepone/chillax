// Package middlewares provides common middleware handlers.
package middlewares

import (
	"github.com/chillaxio/chillax/storage"
	"github.com/gorilla/context"
	"net/http"
)

func SetStorages(storages *storage.Storages) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			context.Set(req, "storages", storages)

			next.ServeHTTP(res, req)
		})
	}
}

// MustLogin is a middleware that checks existence of current user.
func MustLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		storages := context.Get(req, "storages").(*storage.Storages)
		session, _ := storages.Cookie.Get(req, "chillax-session")
		userRowInterface := session.Values["user"]

		if userRowInterface == nil {
			http.Redirect(res, req, "/login", 301)
			return
		}

		next.ServeHTTP(res, req)
	})
}
