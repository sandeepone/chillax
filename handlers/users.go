package handlers

import (
	"encoding/json"
	"errors"
	"github.com/GeertJohan/go.rice"
	chillax_dal "github.com/chillaxio/chillax/dal"
	"github.com/chillaxio/chillax/libhttp"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"github.com/gorilla/context"
	"html/template"
	"io/ioutil"
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

	tmpl := template.New("signup")

	for _, filename := range []string{"login-signup-parent.html.tmpl", "signup.html.tmpl"} {

		templateString, err := box.String(filename)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}

		tmpl, err = tmpl.Parse(templateString)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}
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

	tmpl := template.New("signup")

	for _, filename := range []string{"login-signup-parent.html.tmpl", "login.html.tmpl"} {

		templateString, err := box.String(filename)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}

		tmpl, err = tmpl.Parse(templateString)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}
	}

	tmpl.Execute(w, nil)
}

func PostApiUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	storages := context.Get(r, "storages").(*chillax_storage.Storages)

	jsonPayload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	existingUser, err := chillax_dal.GetUserByEmailAndPasswordJson(storages, jsonPayload)
	if err != nil && err.Error() != "Failed to get user." {
		libhttp.HandleErrorJson(w, err)
		return
	}

	if existingUser != nil {
		err = errors.New("User already exists.")
		libhttp.HandleErrorJson(w, err)
		return
	}

	user, err := chillax_dal.NewUserGivenJson(storages, jsonPayload)
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

	jsonPayload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	user, err := chillax_dal.GetUserByEmailAndPasswordJson(storages, jsonPayload)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	session, _ := storages.Cookie.Get(r, "chillax-session")
	session.Values["user"] = user

	err = session.Save(r, w)
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

func GetApiUsersLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	storages := context.Get(r, "storages").(*chillax_storage.Storages)

	session, _ := storages.Cookie.Get(r, "chillax-session")

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
