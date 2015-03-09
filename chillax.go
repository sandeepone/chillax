package main

import (
	"github.com/carbocation/interpose"
	chillax_handlers "github.com/chillaxio/chillax/handlers"
	chillax_middlewares "github.com/chillaxio/chillax/middlewares"
	chillax_storage "github.com/chillaxio/chillax/storage"
	gorilla_mux "github.com/gorilla/mux"
	"net/http"
	"os"
)

// NewChillax is the constructor for Chillax struct.
func NewChillax() (*Chillax, error) {
	storages, err := chillax_storage.NewStorages()
	if err != nil {
		return nil, err
	}

	chillax := &Chillax{storages}
	return chillax, nil
}

type Chillax struct {
	storages *chillax_storage.Storages
}

func (chillax *Chillax) middlewares(storages *chillax_storage.Storages) (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(chillax_middlewares.SetStorages(storages))

	middle.UseHandler(chillax.mux())

	return middle, nil
}

func (chillax *Chillax) mux() *gorilla_mux.Router {
	router := gorilla_mux.NewRouter()

	router.HandleFunc("/api/users", chillax_handlers.PostApiUsers).Methods("POST")
	router.HandleFunc("/api/users/login", chillax_handlers.PostApiUsersLogin).Methods("POST")

	return router
}

func main() {
	chillax, err := NewChillax()
	if err != nil {
		println(err)
		os.Exit(1)
	}

	middle, err := chillax.middlewares(chillax.storages)
	if err != nil {
		println(err)
		os.Exit(1)
	}

	http.ListenAndServe(":3333", middle)
}
