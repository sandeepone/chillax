package handlers

import (
	"github.com/GeertJohan/go.rice"
	"github.com/chillaxio/chillax/libhttp"
	"github.com/gorilla/context"
	"html/template"
	"net/http"
)

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	gorice := context.Get(r, "gorice").(*rice.Config)

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

	tmpl.Execute(w, nil)
}
