package server

import (
	"os"
	"testing"
)

func TestNewServerWithoutConfigPath(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	server, err := NewServer()

	if err != nil {
		t.Errorf("Should be able to create server without CONFIG_PATH. Error: %v", err)
	}
	if server.Paths["AdminProxies"] == "" {
		t.Errorf("server.AdminProxiesPath should not be empty")
	}
	if server.Paths["AdminStatic"] == "" {
		t.Errorf("server.AdminStaticPath should not be empty")
	}
}
