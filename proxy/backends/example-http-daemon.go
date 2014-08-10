package main

import (
    "fmt"
    "net/http"
)

//
// This is an HTTP daemon to test regular process of proxy/backends.
//
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, 世界")
    })
    http.ListenAndServe(":44444", nil) // TODO: Do i need to generate port dynamically here?
}