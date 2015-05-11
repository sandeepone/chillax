package dal

import (
	"fmt"
	"github.com/chillaxio/chillax/storage"
	"path"
	"strings"
	"testing"
)

func TestThingCreateGetDelete(t *testing.T) {
	storages, err := storage.NewStorages()
	if err != nil {
		t.Fatalf("Creating storages should not fail. Error: %v", err)
	}

	u, err := NewUser(storages, "didip", "password", "password")
	u.Save()

	thing, err := NewThing(storages, u.ID)
	if err != nil {
		t.Fatalf("Creating a thing should not fail. Error: %v", err)
	}

	year := thing.CreatedAt.Year()
	month := thing.CreatedAt.Month()
	day := thing.CreatedAt.Day()

	expectedPathPrefix := path.Join(
		fmt.Sprintf("%v", year),
		fmt.Sprintf("%02d", month),
		fmt.Sprintf("%02d", day),
		thing.ID,
	)

	if thing.Path == "" {
		t.Error("thing.Path should never be empty.")
	}
	if !strings.HasPrefix(thing.Path, expectedPathPrefix) {
		t.Error("thing.Path was set incorrectly.")
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

	err = storages.RemoveAll()
	if err != nil {
		t.Errorf("Remove all storage failed. Error: %v", err)
	}
}
