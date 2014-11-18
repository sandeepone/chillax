package host

import (
	"fmt"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"os"
	"testing"
)

func GetCpuTomlForTest(host string) ([]byte, error) {
	storage := chillax_storage.NewStorage()
	return storage.Get(fmt.Sprintf("/hosts/%v/metrics/cpu", host))
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
