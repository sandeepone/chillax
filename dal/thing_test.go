package dal

import (
	chillax_storage "github.com/chillaxio/chillax/storage"
	"testing"
)

func TestThingCreateGetDelete(t *testing.T) {
	storages, err := chillax_storage.NewStorages()
	if err != nil {
		t.Fatalf("Creating storages should not fail. Error: %v", err)
	}

	u, err := NewUser(storages, "didip", "password")
	u.Save()

	thing, err := NewThing(storages, u.ID)
	if err != nil {
		t.Fatalf("Creating a thing should not fail. Error: %v", err)
	}

	if thing.Path == "" {
		t.Error("thing.Path should never be empty.")
	}

	err = thing.Save([]byte("dostuff"))
	if err != nil {
		t.Errorf("Saving should not fail. Error: %v", err)
	}

	data, err := thing.GetContent()
	if err != nil {
		t.Errorf("Get should not fail. Error: %v", err)
	}
	if string(data) != "dostuff" {
		t.Errorf("Get should not fail.")
	}

	err = thing.Delete()
	if err != nil {
		t.Errorf("Delete should not fail. Error: %v", err)
	}

	storages.RemoveAll()
}
