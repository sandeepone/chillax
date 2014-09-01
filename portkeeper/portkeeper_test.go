package portkeeper

import (
    "fmt"
    "strings"
    "os/exec"
    "testing"
    "github.com/didip/chillax/libtime"
    "github.com/didip/chillax/libstring"
    chillax_storage "github.com/didip/chillax/storage"
)

func GetUsedPortsForTest(dockerHost string) []string {
    storage      := chillax_storage.NewStorage()
    usedPorts, _ := storage.List(fmt.Sprintf("/hosts/%v/used-ports", dockerHost))
    return usedPorts
}

func CheckLengthOfUsedPortsForTest(t *testing.T, dockerHost string, expectation int) {
    usedPorts := GetUsedPortsForTest(dockerHost)

    if len(usedPorts) != expectation {
        t.Errorf("Total used ports should be %v. Used ports: %v", expectation, usedPorts)
    }
}

func TestLsofPort(t *testing.T) {
    cmd := exec.Command("python", "-m", "SimpleHTTPServer", "33456")
    cmd.Start()

    libtime.SleepString("1s")

    output, _ := LsofPort(33456)
    if !strings.Contains(string(output), ":33456") {
        t.Errorf("lsof should found that port 33456 is taken. Output: %v", string(output))
    }

    cmd.Process.Kill()
}

func TestReserveLargestPortWhichIsTheDefault(t *testing.T) {
    dockerUri  := "tcp://127.0.0.1:2375"
    dockerHost := libstring.HostWithoutPort(dockerUri)
    storage    := chillax_storage.NewStorage()

    storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/", dockerHost))

    CheckLengthOfUsedPortsForTest(t, dockerHost, 0)

    port := ReservePort(dockerUri)

    if port != MAX_PORT {
        t.Errorf("port should equal to %v", MAX_PORT)
    }

    CheckLengthOfUsedPortsForTest(t, dockerHost, 1)

    storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/", dockerHost))

    CheckLengthOfUsedPortsForTest(t, dockerHost, 0)
}

func TestReserveGapPort(t *testing.T) {
    dockerUri  := "tcp://127.0.0.1:2375"
    dockerHost := libstring.HostWithoutPort(dockerUri)
    storage    := chillax_storage.NewStorage()

    storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/", dockerHost))

    CheckLengthOfUsedPortsForTest(t, dockerHost, 0)

    ReservePort(dockerUri)
    deleteme := ReservePort(dockerHost)
    ReservePort(dockerUri)

    CheckLengthOfUsedPortsForTest(t, dockerHost, 3)

    storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/%v", dockerHost, deleteme))

    CheckLengthOfUsedPortsForTest(t, dockerHost, 2)

    port := ReservePort(dockerUri)

    if port != deleteme {
        t.Errorf("ReservePort did not regenerate gap port. Port: %v", port)
    }

    storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/", dockerHost))

    CheckLengthOfUsedPortsForTest(t, dockerHost, 0)
}