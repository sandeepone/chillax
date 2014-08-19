package server

import (
    "os"
    "testing"
    "github.com/didip/chillax/libtime"
)


func NewServerForTest(t *testing.T) *HttpServer {
    os.Setenv("DEFAULT_PROXY_BACKENDS_DIR", "./example-default-backend-dir")
    server, err := NewHttpServer()

    if err != nil {
        t.Errorf("Failed to create HTTP server cleanly. Error: %v", err)
    }

    return server
}

func TestServerStartStopBackends(t *testing.T) {
    server := NewServerForTest(t)

    errors := server.CreateProxyBackends()
    for _, err := range errors {
        if err != nil {
            t.Errorf("Failed to create backends. Error: %v", err)
        }
    }

    errors = server.StartProxyBackends()
    for _, err := range errors {
        if err != nil {
            t.Errorf("Failed to start backends. Error: %v", err)
        }
    }

    libtime.SleepString("250ms")

    errors = server.StopProxyBackends()
    for _, err := range errors {
        if err != nil {
            t.Errorf("Failed to stop backends. Error: %v", err)
        }
    }
}