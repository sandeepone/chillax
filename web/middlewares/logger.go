package middlewares

import (
	"github.com/Sirupsen/logrus"
	gorilla_context "github.com/gorilla/context"
	"net/http"
)

func SetLoggerMiddleware(logger *logrus.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gorilla_context.Set(r, "Logger", logger)
	}
}
