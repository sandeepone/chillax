package host

import (
	"fmt"
	"github.com/chillaxio/chillax/libprocess"
	"github.com/chillaxio/chillax/libstring"
	"github.com/chillaxio/chillax/libtime"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func GetCpuTomlForTest(host string) ([]byte, error) {
	storage := chillax_storage.NewStorage()
	return storage.Get(fmt.Sprintf("/hosts/%v/metrics/cpu", host))
}

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

func TestGetAllHosts(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	host := "127.0.0.1"
	storage := chillax_storage.NewStorage()

	storage.Delete("/hosts")

	chost := NewChillaxHost(storage, host)

	err := chost.SaveCpu()
	if err != nil {
		t.Errorf("Saving CPU data should work. err: %v", err)
	}

	chillaxHosts, err := GetAllHosts(storage)
	if err != nil {
		t.Errorf("Unable to get chillax hosts. Error: %v", err)
	}

	if chillaxHosts[0].Name != host {
		t.Errorf("chillax host was not retrieved correctly.")
	}

	// Cleanup
	storage.Delete("/hosts")
}

func TestSaveAndLoadCpu(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	host := "127.0.0.1"
	storage := chillax_storage.NewStorage()

	storage.Delete(fmt.Sprintf("/hosts/%v/metrics/cpu", host))

	_, err := GetCpuTomlForTest(host)
	if err == nil {
		t.Errorf("Cpu data should be nil.")
	}

	chost := NewChillaxHost(storage, host)

	err = chost.SaveCpu()
	if err != nil {
		t.Errorf("Saving CPU data should work. err: %v", err)
	}

	_, err = GetCpuTomlForTest(host)
	if err != nil {
		t.Errorf("Cpu data should not be nil.")
	}

	cpuFromStorage, err := chost.LoadCpu()

	if (chost.CpuMetrics.LoadAverages[0] != cpuFromStorage.LoadAverages[0]) || (chost.CpuMetrics.NumCpu != cpuFromStorage.NumCpu) || (chost.CpuMetrics.LoadAveragesPerCpu[0] != cpuFromStorage.LoadAveragesPerCpu[0]) {
		t.Errorf("Cpu data was not saved properly. chost.CpuMetrics: %v, cpuFromStorage: %v", chost.CpuMetrics, cpuFromStorage)
	}

	// Cleanup
	storage.Delete(fmt.Sprintf("/hosts/%v/metrics/cpu", host))
}

func TestGetCpuFromAllHosts(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	host := "127.0.0.1"
	storage := chillax_storage.NewStorage()

	storage.Delete(fmt.Sprintf("/hosts/%v/metrics/cpu", host))

	chost := NewChillaxHost(storage, host)
	err := chost.SaveCpu()
	if err != nil {
		t.Errorf("Saving CPU data should work. err: %v", err)
	}

	cpusFromStorage, err := GetCpuFromAllHosts(storage)

	for h, cpuData := range cpusFromStorage {
		if h == host {
			if (chost.CpuMetrics.LoadAverages[0] != cpuData.LoadAverages[0]) || (chost.CpuMetrics.NumCpu != cpuData.NumCpu) || (chost.CpuMetrics.LoadAveragesPerCpu[0] != cpuData.LoadAveragesPerCpu[0]) {
				t.Errorf("Cpu data was not saved properly. chost.CpuMetrics: %v, cpusFromStorage: %v", chost.CpuMetrics, cpuData)
			}
		}
	}

	// Cleanup
	storage.Delete(fmt.Sprintf("/hosts/%v/metrics/cpu", host))
}

func TestLsofPort(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	cmd := exec.Command("python", "-m", "SimpleHTTPServer", "33456")
	cmd.Start()

	libtime.SleepString("1200ms")

	output, _ := libprocess.LsofPort(33456)
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
	chost := NewChillaxHost(storage, dockerHost)

	storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/", dockerHost))

	chost.ReservePort()
	CheckLengthOfUsedPortsForTest(t, dockerHost, 1)

	chost.ReservePort()
	CheckLengthOfUsedPortsForTest(t, dockerHost, 2)
}

func TestCleanReservedPortsShouldNotCareOfHostProtocol(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	dockerUri := "tcp://127.0.0.1:2375"
	dockerHost := libstring.HostWithoutPort(dockerUri)
	storage := chillax_storage.NewStorage()
	chost := NewChillaxHost(storage, dockerHost)

	storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/", dockerHost))

	chost.ReservePort()
	err := chost.CleanReservedPorts()
	if err != nil {
		t.Errorf("Unable to clean reserved port. Error: %v", err)
	}
	CheckLengthOfUsedPortsForTest(t, dockerHost, 0)

	chost.ReservePort()
	err = chost.CleanReservedPorts()
	if err != nil {
		t.Errorf("Unable to clean reserved port. Error: %v", err)
	}
	CheckLengthOfUsedPortsForTest(t, dockerHost, 0)
}

func TestReserveLargestPortWhichIsTheDefault(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	dockerHost := "127.0.0.1"
	storage := chillax_storage.NewStorage()
	chost := NewChillaxHost(storage, dockerHost)

	storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/", dockerHost))

	CheckLengthOfUsedPortsForTest(t, dockerHost, 0)

	port := chost.ReservePort()

	if port != MAX_PORT {
		t.Errorf("port should equal to %v", MAX_PORT)
	}

	CheckLengthOfUsedPortsForTest(t, dockerHost, 1)

	// Since we are not actually using the port, CleanReservedPorts will delete port correctly.
	err := chost.CleanReservedPorts()
	if err != nil {
		t.Errorf("CleanReservedPorts should work. Error: %v", err)
	}

	CheckLengthOfUsedPortsForTest(t, dockerHost, 0)
}

func TestReserveGapPort(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	dockerHost := "127.0.0.1"
	storage := chillax_storage.NewStorage()
	chost := NewChillaxHost(storage, dockerHost)

	storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/", dockerHost))

	CheckLengthOfUsedPortsForTest(t, dockerHost, 0)

	chost.ReservePort()
	deleteme := chost.ReservePort()
	chost.ReservePort()

	CheckLengthOfUsedPortsForTest(t, dockerHost, 3)

	storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/%v", dockerHost, deleteme))

	CheckLengthOfUsedPortsForTest(t, dockerHost, 2)

	port := chost.ReservePort()

	if port != deleteme {
		t.Errorf("ReservePort did not regenerate gap port. Port: %v", port)
	}

	// Since we are not actually using the port, CleanReservedPorts will delete port correctly.
	chost.CleanReservedPorts()

	CheckLengthOfUsedPortsForTest(t, dockerHost, 0)
}
