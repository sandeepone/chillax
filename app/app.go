package app

import (
	"github.com/carbocation/interpose"
	"github.com/chillaxio/chillax/handlers"
	"github.com/chillaxio/chillax/middlewares"
	"github.com/chillaxio/chillax/storage"
	gorilla_mux "github.com/gorilla/mux"
	"net/http"
)

// NewChillax is the constructor for Chillax struct.
func NewChillax() (*Chillax, error) {
	storages, err := storage.NewStorages()
	if err != nil {
		return nil, err
	}

	chillax := &Chillax{storages}
	return chillax, nil
}

type Chillax struct {
	storages *storage.Storages
}

func (chillax *Chillax) Middlewares() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetStorages(chillax.storages))

	middle.UseHandler(chillax.Mux())

	return middle, nil
}

func (chillax *Chillax) Mux() *gorilla_mux.Router {
	MustLogin := middlewares.MustLogin

	router := gorilla_mux.NewRouter()

	router.Handle("/", MustLogin(http.HandlerFunc(handlers.GetHome))).Methods("GET")

	router.HandleFunc("/signup", handlers.GetSignup).Methods("GET")
	router.HandleFunc("/signup", handlers.PostSignup).Methods("POST")
	router.HandleFunc("/login", handlers.GetLogin).Methods("GET")
	router.HandleFunc("/login", handlers.PostLogin).Methods("POST")
	router.HandleFunc("/logout", handlers.GetLogout).Methods("GET")

	router.Handle("/users/{id:[0-9]+}", MustLogin(http.HandlerFunc(handlers.PostPutDeleteUsersID))).Methods("POST", "PUT", "DELETE")

	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return router
}
