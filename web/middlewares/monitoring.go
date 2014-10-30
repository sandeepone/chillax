package middlewares

import (
	"fmt"
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
			gorilla_context.Delete(r, "BeginRequestTime")
			gorilla_context.Set(r, "RequestTime", time.Since(data.(time.Time)))

			fmt.Printf("%v %v took: %v\n", r.Method, r.RequestURI, time.Since(data.(time.Time)))
		}
	}
}
