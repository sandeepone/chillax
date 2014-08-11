package libstring

import (
    "os"
    "testing"
)

func TestGuid(t *testing.T) {
    if Guid() == "" {
        t.Errorf("Failed to generate GUID")
    }
}

func TestNormalizeLocalhost(t *testing.T) {
    hostname, _ := os.Hostname()

    if NormalizeLocalhost("tcp://localhost:2375") != ("tcp://" + hostname + ":2375") {
        t.Errorf("Failed to normalize localhost: %v", NormalizeLocalhost("tcp://localhost:2375"))
    }
}

func TestStripProtocol(t *testing.T) {
    hostname, _ := os.Hostname()

    if StripProtocol("tcp://localhost:2375") != (hostname + ":2375") {
        t.Errorf("Failed to strip protocol: %v", StripProtocol("tcp://localhost:2375"))
    }
}


