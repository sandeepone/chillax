package muxproducer

import (
    "os"
    "testing"
    "github.com/didip/chillax/libtime"
)


func NewMuxProducerForTest(t *testing.T) *MuxProducer {
    os.Setenv("DEFAULT_PROXY_BACKENDS_DIR", "./example-default-backend-dir")
    muxObj, err := NewMuxProducer()

    if err != nil {
        t.Errorf("Failed to create HTTP muxObj cleanly. Error: %v", err)
    }

    return muxObj
}

func TestMuxProducerStartStopBackends(t *testing.T) {
    muxObj := NewMuxProducerForTest(t)

    errors := muxObj.CreateProxyBackends()
    for _, err := range errors {
        if err != nil {
            t.Errorf("Failed to create backends. Error: %v", err)
        }
    }

    errors = muxObj.StartProxyBackends()
    for _, err := range errors {
        if err != nil {
            t.Errorf("Failed to start backends. Error: %v", err)
        }
    }

    libtime.SleepString("250ms")

    errors = muxObj.StopProxyBackends()
    for _, err := range errors {
        if err != nil {
            t.Errorf("Failed to stop backends. Error: %v", err)
        }
    }
}