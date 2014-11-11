package portkeeper

import (
	"fmt"
	"github.com/chillaxio/chillax/libstring"
	"github.com/chillaxio/chillax/libtime"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func GetUsedPortsForTest(dockerHost string) []string {
	storage := chillax_storage.NewStorage()
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
	os.Setenv("CHILLAX_ENV", "test")

	cmd := exec.Command("python", "-m", "SimpleHTTPServer", "33456")
	cmd.Start()

	libtime.SleepString("1200ms")

	output, _ := LsofPort(33456)
	if !strings.Contains(string(output), ":33456") {
		t.Errorf("lsof should found that port 33456 is taken. Output: %v", string(output))
	}

	cmd.Process.Kill()
}

func TestReservePortShouldNotCareOfHostProtocol(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	dockerUri := "tcp://127.0.0.1:2375"
	dockerHost := libstring.HostWithoutPort(dockerUri)
	storage := chillax_storage.NewStorage()

	storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/", dockerHost))

	ReservePort(dockerUri)
	CheckLengthOfUsedPortsForTest(t, dockerHost, 1)

	ReservePort(dockerHost)
	CheckLengthOfUsedPortsForTest(t, dockerHost, 2)
}

func TestCleanReservedPortsShouldNotCareOfHostProtocol(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	dockerUri := "tcp://127.0.0.1:2375"
	dockerHost := libstring.HostWithoutPort(dockerUri)
	storage := chillax_storage.NewStorage()

	storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/", dockerHost))

	ReservePort(dockerUri)
	CleanReservedPorts(dockerUri)
	CheckLengthOfUsedPortsForTest(t, dockerHost, 0)

	ReservePort(dockerHost)
	CleanReservedPorts(dockerHost)
	CheckLengthOfUsedPortsForTest(t, dockerHost, 0)
}

func TestReserveLargestPortWhichIsTheDefault(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	dockerHost := "127.0.0.1"
	storage := chillax_storage.NewStorage()

	storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/", dockerHost))

	CheckLengthOfUsedPortsForTest(t, dockerHost, 0)

	port := ReservePort(dockerHost)

	if port != MAX_PORT {
		t.Errorf("port should equal to %v", MAX_PORT)
	}

	CheckLengthOfUsedPortsForTest(t, dockerHost, 1)

	// Since we are not actually using the port, CleanReservedPorts will delete port correctly.
	err := CleanReservedPorts(dockerHost)
	if err != nil {
		t.Errorf("CleanReservedPorts should work. Error: %v", err)
	}

	CheckLengthOfUsedPortsForTest(t, dockerHost, 0)
}

func TestReserveGapPort(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	dockerHost := "127.0.0.1"
	storage := chillax_storage.NewStorage()

	storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/", dockerHost))

	CheckLengthOfUsedPortsForTest(t, dockerHost, 0)

	ReservePort(dockerHost)
	deleteme := ReservePort(dockerHost)
	ReservePort(dockerHost)

	CheckLengthOfUsedPortsForTest(t, dockerHost, 3)

	storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/%v", dockerHost, deleteme))

	CheckLengthOfUsedPortsForTest(t, dockerHost, 2)

	port := ReservePort(dockerHost)

	if port != deleteme {
		t.Errorf("ReservePort did not regenerate gap port. Port: %v", port)
	}

	// Since we are not actually using the port, CleanReservedPorts will delete port correctly.
	CleanReservedPorts(dockerHost)

	CheckLengthOfUsedPortsForTest(t, dockerHost, 0)
}
