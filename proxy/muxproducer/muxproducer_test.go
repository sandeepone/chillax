package muxproducer

import (
    "os"
    "testing"
    "github.com/didip/chillax/libtime"
)


func NewMuxProducerForTest(t *testing.T) *MuxProducer {
    os.Setenv("DEFAULT_PROXY_BACKENDS_DIR", "./example-default-backend-dir")
    mp, err := NewMuxProducer()

    if err != nil {
        t.Errorf("Failed to create HTTP mp cleanly. Error: %v", err)
    }

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