package app

import (
	"github.com/carbocation/interpose"
	chillax_handlers "github.com/chillaxio/chillax/handlers"
	chillax_middlewares "github.com/chillaxio/chillax/middlewares"
	chillax_storage "github.com/chillaxio/chillax/storage"
	gorilla_mux "github.com/gorilla/mux"
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

func (chillax *Chillax) Middlewares() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(chillax_middlewares.SetStorages(chillax.storages))

	middle.UseHandler(chillax.Mux())

	return middle, nil
}

func (chillax *Chillax) Mux() *gorilla_mux.Router {
	router := gorilla_mux.NewRouter()

	router.HandleFunc("/api/users", chillax_handlers.PostApiUsers).Methods("POST")
	router.HandleFunc("/api/users/login", chillax_handlers.PostApiUsersLogin).Methods("POST")

	return router
}
