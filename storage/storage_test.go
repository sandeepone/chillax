package storage

import (
	"github.com/didip/chillax/libenv"
	"os/user"
	"testing"
)

func TestDefaultStorageType(t *testing.T) {
	storageType := libenv.EnvWithDefault("STORAGE_TYPE", "FileSystem")
	if storageType != "FileSystem" {
		t.Error("Default storageType should equal to FileSystem")
	}
}

func TestRootFileSystem(t *testing.T) {
	currentUser, _ := user.Current()

	storage := NewStorage()

	if storage.GetRoot() != currentUser.HomeDir+"/chillax" {
		t.Error("Root of FileSystem storage should be located at $HOME/chillax")
	}
}
