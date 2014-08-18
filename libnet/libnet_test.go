package libnet

import (
    "os"
    "net"
    "testing"
)

func TestRemoteToLocalHostEquality(t *testing.T) {
    hostname, _ := os.Hostname()
    addrs, _    := net.InterfaceAddrs()

    if !RemoteToLocalHostEquality(hostname) {
        t.Errorf("Local hostname should == to local hostname. Interface addresses: %v, Hostname: %v", addrs, hostname)
    }
    if !RemoteToLocalHostEquality("127.0.0.1") {
        t.Errorf("127.0.0.1 should == to local hostname. Interface addresses: %v, Hostname: %v", addrs, hostname)
    }
    if RemoteToLocalHostEquality("blah-i-do-not-exist.local") {
        t.Errorf("Local hostname should != to local hostname. Interface addresses: %v, Hostname: %v", addrs, hostname)
    }
}
