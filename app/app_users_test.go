package app

import (
	"bytes"
	chillax_dal "github.com/chillaxio/chillax/dal"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestCreateUserFromApi(t *testing.T) {
	// ---- Setup ----
	chillax, err := NewChillax()
	if err != nil {
		t.Fatalf("Failed to create chillax app. Error: %v", err)
	}

	middle, err := chillax.Middlewares()
	if err != nil {
		t.Fatalf("Failed to create middlewares. Error: %v", err)
	}

	go http.ListenAndServe(":18000", middle)

	// ---- Setup ----

	data := `{"Email": "didip@example.com", "Password": "abc123"}`

	resp, err := http.Post("http://localhost:18000/api/users", "application/json", bytes.NewBuffer([]byte(data)))
	if err != nil {
		t.Errorf("Failed to create user via API. Error: %v", err)
	}
	defer resp.Body.Close()

	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to response body via API. Error: %v", err)
	}

	if !strings.Contains(string(bodyContent), "ID") {
		t.Errorf("Failed to response body via API. Body: %v", string(bodyContent))
	}

	chillax.storages.RemoveAll()
}

func TestLoginFromApi(t *testing.T) {
	// ---- Setup ----
	chillax, err := NewChillax()
	if err != nil {
		t.Fatalf("Failed to create chillax app. Error: %v", err)
	}

	middle, err := chillax.Middlewares()
	if err != nil {
		t.Fatalf("Failed to create middlewares. Error: %v", err)
	}

	go http.ListenAndServe(":18000", middle)

	_, err = chillax_dal.NewUser(chillax.storages, "didip@example.com", "abc123")
	if err != nil {
		t.Errorf("Failed to create user. Error: %v", err)
	}

	// ---- Setup ----

	data := `{"Email": "didip@example.com", "Password": "abc123"}`

	resp, err := http.Post("http://localhost:18000/api/users/login", "application/json", bytes.NewBuffer([]byte(data)))
	if err != nil {
		t.Errorf("Failed to login via API. Error: %v", err)
	}
	defer resp.Body.Close()

	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to response body via API. Error: %v", err)
	}

	if !strings.Contains(string(bodyContent), "ID") {
		t.Errorf("Failed to response body via API. Body: %v", string(bodyContent))
	}

	chillax.storages.RemoveAll()
}
