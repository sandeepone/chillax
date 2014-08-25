package muxproducer

import (
    "os"
    "testing"
    "path/filepath"
    "github.com/didip/chillax/libtime"
    chillax_web_settings "github.com/didip/chillax/web/settings"
)


func NewMuxProducerForTest(t *testing.T) *MuxProducer {
    fullpath, _ := filepath.Abs("../../examples/configs/proxy-handlers")
    os.Setenv("PROXY_HANDLERS_PATH", fullpath)

    settings := chillax_web_settings.NewServerSettings()
    mp       := NewMuxProducer(settings.ProxyHandlerTomls)

    return mp
}

func TestMuxProducerStartStopBackends(t *testing.T) {
    mp := NewMuxProducerForTest(t)

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