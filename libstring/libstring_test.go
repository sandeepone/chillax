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
    if StripProtocol("tcp://localhost:2375") != "localhost:2375" {
        t.Errorf("Failed to strip protocol: %v", StripProtocol("tcp://localhost:2375"))
    }
}

func TestSplitDockerPorts(t *testing.T) {
    ports := "127.0.0.1:99999:80/tcp"
    hostIp, hostPort, containerPort := SplitDockerPorts(ports)

    if hostIp != "127.0.0.1" {
        t.Errorf("Failed to split docker ports. hostIp: %v, hostPort: %v, containerPort: %v", hostIp, hostPort, containerPort)
    }
    if hostPort != "99999" {
        t.Errorf("Failed to split docker ports. hostIp: %v, hostPort: %v, containerPort: %v", hostIp, hostPort, containerPort)
    }
    if containerPort != "80/tcp" {
        t.Errorf("Failed to split docker ports. hostIp: %v, hostPort: %v, containerPort: %v", hostIp, hostPort, containerPort)
    }

    ports = "99999:80/tcp"
    hostIp, hostPort, containerPort = SplitDockerPorts(ports)

    if hostIp != "" {
        t.Errorf("Failed to split docker ports. hostIp: %v, hostPort: %v, containerPort: %v", hostIp, hostPort, containerPort)
    }
    if hostPort != "99999" {
        t.Errorf("Failed to split docker ports. hostIp: %v, hostPort: %v, containerPort: %v", hostIp, hostPort, containerPort)
    }
    if containerPort != "80/tcp" {
        t.Errorf("Failed to split docker ports. hostIp: %v, hostPort: %v, containerPort: %v", hostIp, hostPort, containerPort)
    }
}

