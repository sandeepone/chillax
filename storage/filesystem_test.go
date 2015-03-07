package storage

import (
	"testing"
)

func TestFileSystemCreateGetDelete(t *testing.T) {
	storage := NewFileSystem("111")

	err := storage.Create("/dostuff", []byte("dostuff"))
	if err != nil {
		t.Errorf("Create should not fail. Error: %v", err)
	}

	data, err := storage.Get("/dostuff")
	if err != nil {
		t.Errorf("Get should not fail. Error: %v", err)
	}
	if string(data) != "dostuff" {
		t.Errorf("Get should not fail.")
	}

	err = storage.Delete("/dostuff")
	if err != nil {
		t.Errorf("Delete should not fail. Error: %v", err)
	}
}
