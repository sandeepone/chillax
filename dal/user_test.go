package dal

import (
	chillax_storage "github.com/chillaxio/chillax/storage"
	"testing"
)

func TestNewUser(t *testing.T) {
	storages, err := chillax_storage.NewStorages()
	if err != nil {
		t.Fatalf("Creating storages should not fail. Error: %v", err)
	}
	u := NewUser(storages)
	if u == nil {
		t.Error("Creating user should not fail.")
	}
	if u.storages == nil {
		t.Error("storages should not be empty.")
	}
	if u.bucketName == "" {
		t.Error("bucketName should not be empty.")
	}
	storages.KeyValue.Close()
}

func TestSaveAndGetUser(t *testing.T) {
	storages, err := chillax_storage.NewStorages()
	if err != nil {
		t.Fatalf("Creating storages should not fail. Error: %v", err)
	}

	u := NewUser(storages)

	err = u.Save()
	if err != nil {
		t.Errorf("Saving user should work. Error: %v", err)
	}

	if u.ID == "" {
		t.Error("User ID should not be empty.")
	}

	userFromStorage, err := GetUserById(storages, u.ID)
	if err != nil {
		t.Fatalf("Getting user should work. Error: %v", err)
	}
	if u.ID != userFromStorage.ID {
		t.Error("Got the wrong user.")
	}
	storages.KeyValue.Close()
}
