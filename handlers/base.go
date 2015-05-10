// Package handlers provides request handlers.
package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func getIdFromPath(w http.ResponseWriter, r *http.Request) (string, error) {
	userId := mux.Vars(r)["id"]
	if userId == "" {
		return "", errors.New("user id cannot be empty.")
	}

	return userId, nil
}
