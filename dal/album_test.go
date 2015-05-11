package dal

import (
	"github.com/chillaxio/chillax/storage"
	"testing"
)

func TestCreateAndDeleteUserAlbums(t *testing.T) {
	storages, err := storage.NewStorages()
	if err != nil {
		t.Fatalf("Creating storages should not fail. Error: %v", err)
	}

	u, err := NewUser(storages, "didip@example.com", "password", "password")
	if err != nil {
		t.Errorf("Creating user should not fail. Error: %v", err)
	}

	err = u.CreateAlbum("programming")
	if err != nil {
		t.Errorf("Creating a album should not fail. Error: %v", err)
	}

	album := u.GetAlbumByName("programming")
	if album == nil {
		t.Error("Album 'programming' should exist.")
	}

	err = u.DeleteAlbum("programming")
	if err != nil {
		t.Errorf("Deleting a album should not fail. Error: %v", err)
	}

	err = u.DeleteAlbum("aaa")
	if err != nil {
		t.Errorf("Deleting non-existing album should not fail. Error: %v", err)
	}

	err = storages.RemoveAll()
	if err != nil {
		t.Errorf("Remove all storage failed. Error: %v", err)
	}
}
