package storage

import (
	"fmt"
	"github.com/chillaxio/chillax/libenv"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestDefaultStorageType(t *testing.T) {
	storageType := libenv.EnvWithDefault("STORAGE_TYPE", "FileSystem")
	if storageType != "FileSystem" {
		t.Error("Default storageType should equal to FileSystem")
	}
}

func TestRootFileSystemWithDefaultEnvironment(t *testing.T) {
	currentUser, _ := user.Current()
	chillaxEnv := "development"

	storage := NewStorage()

	if storage.GetRoot() != filepath.Join(currentUser.HomeDir, fmt.Sprintf("chillax-%v", chillaxEnv)) {
		t.Errorf("Root of FileSystem storage should be located at $HOME/chillax-%v", chillaxEnv)
	}
}

func TestRootFileSystemWithTestEnvironment(t *testing.T) {
	currentUser, _ := user.Current()
	chillaxEnv := "test"

	os.Setenv("CHILLAX_ENV", chillaxEnv)

	storage := NewStorage()

	if storage.GetRoot() != filepath.Join(currentUser.HomeDir, fmt.Sprintf("chillax-%v", chillaxEnv)) {
		t.Errorf("Root of FileSystem storage should be located at $HOME/chillax-%v", chillaxEnv)
	}
}
