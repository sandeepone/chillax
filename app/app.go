package app

import (
	"github.com/GeertJohan/go.rice"
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

	gorice := &rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateFS},
	}

	chillax := &Chillax{storages, gorice}
	return chillax, nil
}

type Chillax struct {
	storages *chillax_storage.Storages
	gorice   *rice.Config
}

func (chillax *Chillax) Middlewares() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(chillax_middlewares.SetStorages(chillax.storages))
	middle.Use(chillax_middlewares.SetGoRice(chillax.gorice))

	middle.UseHandler(chillax.Mux())

	return middle, nil
}

func (chillax *Chillax) Mux() *gorilla_mux.Router {
	router := gorilla_mux.NewRouter()

	router.HandleFunc("/", chillax_handlers.GetDashboard).Methods("GET")

	router.HandleFunc("/signup", chillax_handlers.GetSignup).Methods("GET")
	router.HandleFunc("/login", chillax_handlers.GetLogin).Methods("GET")

	router.HandleFunc("/api/users", chillax_handlers.PostApiUsers).Methods("POST")
	router.HandleFunc("/api/users/login", chillax_handlers.PostApiUsersLogin).Methods("POST")

	return router
}
