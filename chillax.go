package main

import (
	"encoding/gob"
	chillax_app "github.com/chillaxio/chillax/app"
	chillax_dal "github.com/chillaxio/chillax/dal"
	"net/http"
	"os"
)

func main() {
	gob.Register(&chillax_dal.User{})

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
