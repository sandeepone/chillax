package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	chillax_storage "github.com/chillaxio/chillax/storage"
	gorilla_context "github.com/gorilla/context"
	"net/http"
	"time"
)

func BeginRequestTimerMiddleware() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gorilla_context.Set(r, "BeginRequestTime", time.Now())
	}
}

func RecordRequestTimerMiddleware(storage chillax_storage.Storer) func(http.ResponseWriter, *http.Request) {
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

					fields := logrus.Fields{
						"CurrentTime":     fmt.Sprintf(`"%v"`, currentTime.String()),
						"CurrentUnixNano": currentTime.UnixNano(),
						"Method":          r.Method,
						"URI":             r.RequestURI,
						"RemoteAddr":      r.RemoteAddr,
						"UserAgent":       fmt.Sprintf(`"%v"`, r.UserAgent()),
						"Latency":         latency,
					}
					logger.(*logrus.Logger).WithFields(fields).Info(fmt.Sprintf(`"%v %v took %v"`, r.Method, r.RequestURI, latency))

					datetime := time.Unix(0, currentTime.UnixNano())
					dataPath := fmt.Sprintf(
						"/logs/requests/%v/%d/%v/%v/%v/%v-%v",
						datetime.Year(), datetime.Month(), datetime.Day(), datetime.Hour(), datetime.Minute(),
						currentTime.UnixNano(), r.Method)

					fields["CurrentTime"] = currentTime.String()
					fields["UserAgent"] = r.UserAgent()

					inBytes, _ := json.Marshal(fields)
					storage.Create(dataPath, inBytes)
				}
			}()
		}
	}
}
