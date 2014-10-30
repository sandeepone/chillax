package middlewares

import (
	"net/http"

	chillax_web_settings "github.com/chillaxio/chillax/web/settings"
)

func ServerNameMiddleware(settings *chillax_web_settings.ServerSettings) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Server-Name", "Chillax HTTP Server")
	}
}
