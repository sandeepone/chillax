package main

import (
	chillax_app "github.com/chillaxio/chillax/app"
	"net/http"
	"os"
)

func main() {
	chillax, err := chillax_app.NewChillax()
	if err != nil {
		println(err)
		os.Exit(1)
	}

	middle, err := chillax.Middlewares()
	if err != nil {
		println(err)
		os.Exit(1)
	}

	http.ListenAndServe(":3333", middle)
}
