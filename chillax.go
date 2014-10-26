package main

import (
	chillax_web_server "github.com/chillaxio/chillax/web/server"
)

func main() {
	server, err := chillax_web_server.NewServer()
	if err != nil {
		panic(err)
	}

	server.RunAllInProgressPipelinesAsync()
	server.CheckProxiesAsync()
	server.ListenAndServeGeneric()
}
