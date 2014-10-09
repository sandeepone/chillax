package main

import (
	chillax_web_server "github.com/didip/chillax/web/server"
)

func main() {
	server, err := chillax_web_server.NewServer()
	if err != nil {
		panic(err)
	}

	server.ListenAndServeGeneric()
}
