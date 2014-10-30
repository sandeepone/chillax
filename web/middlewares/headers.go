package middlewares

import (
	"net/http"
)

func SetServerNameMiddleware() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Server-Name", "Chillax HTTP Server")
	}
}
