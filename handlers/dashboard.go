package handlers

import (
	chillax_dal "github.com/chillaxio/chillax/dal"
	"github.com/chillaxio/chillax/libhttp"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"github.com/gorilla/context"
	"html/template"
	"net/http"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	storages := context.Get(r, "storages").(*chillax_storage.Storages)

	session, _ := storages.Cookie.Get(r, "chillax-session")

	currentUserInterface := session.Values["user"]
	if currentUserInterface == nil {
		http.Redirect(w, r, "/login", 301)
		return
	}

	currentUser := currentUserInterface.(*chillax_dal.User)

	data := struct {
		CurrentUser *chillax_dal.User
	}{
		currentUser,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}
