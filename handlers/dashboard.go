package handlers

import (
	"github.com/GeertJohan/go.rice"
	chillax_dal "github.com/chillaxio/chillax/dal"
	"github.com/chillaxio/chillax/libhttp"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"github.com/gorilla/context"
	"html/template"
	"net/http"
)

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	gorice := context.Get(r, "gorice").(*rice.Config)

	storages := context.Get(r, "storages").(*chillax_storage.Storages)

	session, _ := storages.Cookie.Get(r, "chillax-session")

	currentUserInterface := session.Values["user"]
	if currentUserInterface == nil {
		http.Redirect(w, r, "/login", 301)
		return
	}

	currentUser := currentUserInterface.(*chillax_dal.User)

	box, err := gorice.FindBox("dashboard-templates")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl := template.New("dashboard")

	for _, filename := range []string{"dashboard-base.html.tmpl"} {

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

	tmpl.Execute(w, currentUser)
}
