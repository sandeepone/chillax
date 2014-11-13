package server

import (
	"os"
	"testing"
)

func TestNewServerWithoutConfigPath(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	_, err := NewServer()

	if err != nil {
		t.Errorf("Should be able to create server without CONFIG_PATH. Error: %v", err)
	}
}
