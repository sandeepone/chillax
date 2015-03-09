package handlers

import (
	"encoding/json"
	chillax_dal "github.com/chillaxio/chillax/dal"
	"github.com/chillaxio/chillax/libhttp"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"github.com/gorilla/context"
	// "io/ioutil"
	"net/http"
)

func PostApiUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	storages := context.Get(r, "storages").(*chillax_storage.Storages)

	user, err := chillax_dal.NewUserGivenJson(storages, r.Body)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	err = user.Save()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(userJson)
}

func PostApiUsersLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	storages := context.Get(r, "storages").(*chillax_storage.Storages)

	user, err := chillax_dal.GetUserByEmailAndPasswordJson(storages, r.Body)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	session, _ := storages.Cookie.Get(r, "login-session")
	session.Values[user.ID] = user
	session.Save(r, w)

	userJson, err := json.Marshal(user)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(userJson)
}
