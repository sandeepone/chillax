package multiplexer

import (
	"github.com/chillaxio/chillax/libtime"
	chillax_storage "github.com/chillaxio/chillax/storage"
	chillax_web_settings "github.com/chillaxio/chillax/web/settings"
	"os"
	"path/filepath"
	"testing"
)

func NewMuxFactoryForTest(t *testing.T) *MuxFactory {
	fullpath, _ := filepath.Abs("../../examples/configs/proxy-handlers")
	os.Setenv("PROXY_HANDLERS_PATH", fullpath)

	settings, _ := chillax_web_settings.NewServerSettings()
	mp := NewMuxFactory(chillax_storage.NewStorage(), settings.ProxyHandlerTomls)

	return mp
}

func TestMuxFactoryStartStopBackends(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	mp := NewMuxFactoryForTest(t)

	errors := mp.CreateProxyBackends()
	for _, err := range errors {
		if err != nil {
			t.Errorf("Failed to create backends. Error: %v", err)
		}
	}

	errors = mp.StartProxyBackends()
	for _, err := range errors {
		if err != nil {
			t.Errorf("Failed to start backends. Error: %v", err)
		}
	}

	libtime.SleepString("250ms")

	errors = mp.StopProxyBackends()
	for _, err := range errors {
		if err != nil {
			t.Errorf("Failed to stop backends. Error: %v", err)
		}
	}
}
