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

	u, err := NewUser(storages, "didip@example.com", "password", "password")
	if err != nil {
		t.Errorf("Creating user should not fail. Error: %v", err)
	}
	if u == nil {
		t.Error("Creating user should not fail.")
	}
	if u.storages == nil {
		t.Error("storages should not be empty.")
	}
	if u.bucketName == "" {
		t.Error("bucketName should not be empty.")
	}
	if u.ID == "" {
		t.Error("User ID should not be empty.")
	}

	err = storages.RemoveAll()
	if err != nil {
		t.Fatalf("Wiping storage should work. Error: %v", err)
	}
}

func TestHashedPassword(t *testing.T) {
	storages, err := chillax_storage.NewStorages()
	if err != nil {
		t.Fatalf("Creating storages should not fail. Error: %v", err)
	}

	u, err := NewUser(storages, "didip@example.com", "password", "password")
	if err != nil {
		t.Errorf("Creating user should not fail. Error: %v", err)
	}

	if u.Password == "password" {
		t.Fatal("Hashing password should work.")
	}

	storages.RemoveAll()
}

func TestValidateBeforeSave(t *testing.T) {
	storages, err := chillax_storage.NewStorages()
	if err != nil {
		t.Fatalf("Creating storages should not fail. Error: %v", err)
	}

	u, err := NewUser(storages, "didip@example.com", "password", "password")
	if err != nil {
		t.Errorf("Creating user should not fail. Error: %v", err)
	}

	err = u.ValidateBeforeSave()
	if err != nil {
		t.Fatalf("Validation should pass because Email or Password is not empty. Error: %v", err)
	}

	storages.RemoveAll()
}

func TestUserSave(t *testing.T) {
	storages, err := chillax_storage.NewStorages()
	if err != nil {
		t.Fatalf("Creating storages should not fail. Error: %v", err)
	}

	u, err := NewUser(storages, "didip@example.com", "password", "password")
	if err != nil {
		t.Errorf("Creating user should not fail. Error: %v", err)
	}

	err = u.Save()
	if err != nil {
		t.Fatalf("Saving user should work because Email and Password is not empty.")
	}

	storages.RemoveAll()
}

func TestGetUserById(t *testing.T) {
	storages, err := chillax_storage.NewStorages()
	if err != nil {
		t.Fatalf("Creating storages should not fail. Error: %v", err)
	}

	u, err := NewUser(storages, "didip@example.com", "password", "password")
	if err != nil {
		t.Errorf("Creating user should not fail. Error: %v", err)
	}

	err = u.Save()
	if err != nil {
		t.Fatalf("Saving user should work because Name and Password is not empty.")
	}

	userFromStorage, err := GetUserById(storages, u.ID)
	if err != nil {
		t.Fatalf("Getting user should work. Error: %v", err)
	}

	if u.ID != userFromStorage.ID {
		t.Error("Got the wrong user.")
	}
	if u.Email != userFromStorage.Email {
		t.Error("Got the wrong user.")
	}
	if u.Password != userFromStorage.Password {
		t.Error("Got the wrong user.")
	}

	storages.RemoveAll()
}

func TestGetUserByEmailAndPassword(t *testing.T) {
	storages, err := chillax_storage.NewStorages()
	if err != nil {
		t.Fatalf("Creating storages should not fail. Error: %v", err)
	}

	u, err := NewUser(storages, "didip@example.com", "password", "password")
	if err != nil {
		t.Errorf("Creating user should not fail. Error: %v", err)
	}

	err = u.Save()
	if err != nil {
		t.Fatalf("Saving user should work because Name and Password is not empty.")
	}

	userFromStorage, err := GetUserByEmailAndPassword(storages, "didip@example.com", "password")
	if err != nil {
		t.Fatalf("Getting user should work. Error: %v", err)
	}

	if u.ID != userFromStorage.ID {
		t.Errorf("Got the wrong user. userFromStorage.ID: %v, userFromStorage.Email: %v, userFromStorage.Password: %v", userFromStorage.ID, userFromStorage.Email, userFromStorage.Password)
	}
	if u.Email != userFromStorage.Email {
		t.Errorf("Got the wrong user. userFromStorage.ID: %v, userFromStorage.Email: %v, userFromStorage.Password: %v", userFromStorage.ID, userFromStorage.Email, userFromStorage.Password)
	}
	if u.Password != userFromStorage.Password {
		t.Errorf("Got the wrong user. userFromStorage.ID: %v, userFromStorage.Email: %v, userFromStorage.Password: %v", userFromStorage.ID, userFromStorage.Email, userFromStorage.Password)
	}

	storages.RemoveAll()
}
