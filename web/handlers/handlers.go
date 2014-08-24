package handlers

import (
    "net/http"
    "github.com/GeertJohan/go.rice"
    "github.com/gorilla/mux"
)

func GorillaMuxRouteStaticDir(router *mux.Router, staticDirectory string) {
    box, err := rice.FindBox(staticDirectory)
    if err == nil {
        router.Handle(staticDirectory, http.FileServer(box.HTTPBox()))
    }
}
