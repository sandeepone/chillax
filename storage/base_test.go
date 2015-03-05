package storage

import (
	"github.com/chillaxio/chillax/libstring"
	"os"
	"testing"
)

func TestNewStoragesWithDefault(t *testing.T) {
	_, err := NewStorages()
	if err != nil {
		t.Fatalf("Creating storages should not fail. Error: %v", err)
	}
	if _, err := os.Stat(libstring.ExpandTildeAndEnv("~/chillax/kv-db")); os.IsNotExist(err) {
		t.Fatal("Default key-value db file should exist.")
	}
}
