package handlers

import (
	"encoding/json"
	"github.com/GeertJohan/go.rice"
	chillax_dal "github.com/chillaxio/chillax/dal"
	"github.com/chillaxio/chillax/libhttp"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"github.com/gorilla/context"
	"html/template"
	// "io/ioutil"
	"errors"
	"net/http"
)

func GetSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	gorice := context.Get(r, "gorice").(*rice.Config)

	box, err := gorice.FindBox("users-templates")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	templateString, err := box.String("signup.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl, err := template.New("signup").Parse(templateString)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, nil)
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	gorice := context.Get(r, "gorice").(*rice.Config)

	box, err := gorice.FindBox("users-templates")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	templateString, err := box.String("login.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl, err := template.New("login").Parse(templateString)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, nil)
}

func PostApiUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	storages := context.Get(r, "storages").(*chillax_storage.Storages)

	existingUser, err := chillax_dal.GetUserByEmailAndPasswordJson(storages, r.Body)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	if existingUser != nil {
		err = errors.New("User already exists.")
		libhttp.HandleErrorJson(w, err)
		return
	}

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
	session.Values["user"] = user
	session.Save(r, w)

	userJson, err := json.Marshal(user)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(userJson)
}

func GetApiUsersLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	storages := context.Get(r, "storages").(*chillax_storage.Storages)

	session, _ := storages.Cookie.Get(r, "login-session")

	user := session.Values["user"]
	delete(session.Values, "user")
	session.Save(r, w)

	userJson, err := json.Marshal(user)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(userJson)
}
