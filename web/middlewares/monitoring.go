package middlewares

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	gorilla_context "github.com/gorilla/context"
	"net/http"
	"time"
)

func BeginRequestTimerMiddleware() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gorilla_context.Set(r, "BeginRequestTime", time.Now())
	}
}

func RecordRequestTimerMiddleware() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, ok := gorilla_context.GetOk(r, "BeginRequestTime")
		if ok {
			latency := time.Since(data.(time.Time))

			gorilla_context.Delete(r, "BeginRequestTime")
			gorilla_context.Set(r, "RequestTime", latency)

			go func() {
				logger, ok := gorilla_context.GetOk(r, "Logger")
				if ok {
					currentTime := time.Now()

					logger.(*logrus.Logger).WithFields(logrus.Fields{
						"CurrentTime":     fmt.Sprintf(`"%v"`, currentTime.String()),
						"CurrentUnixNano": currentTime.UnixNano(),
						"Method":          r.Method,
						"URI":             r.RequestURI,
						"RemoteAddr":      r.RemoteAddr,
						"UserAgent":       fmt.Sprintf(`"%v"`, r.UserAgent()),
						"Latency":         latency,
					}).Info(fmt.Sprintf(`"%v %v took %v"`, r.Method, r.RequestURI, latency))
				}
			}()
		}
	}
}
