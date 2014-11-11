package metricskeeper

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

func TestSaveAndLoadCpu(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	host := "127.0.0.1"
	storage := chillax_storage.NewStorage()

	storage.Delete(fmt.Sprintf("/hosts/%v/metrics/cpu", host))

	_, err := GetCpuTomlForTest(host)
	if err == nil {
		t.Errorf("Cpu data should be nil.")
	}

	cpu, err := SaveCpu(storage, host)
	if err != nil {
		t.Errorf("Saving CPU data should work. err: %v", err)
	}

	_, err = GetCpuTomlForTest(host)
	if err != nil {
		t.Errorf("Cpu data should not be nil.")
	}

	cpuFromStorage, err := LoadCpu(storage, host)

	if (cpu.LoadAverages[0] != cpuFromStorage.LoadAverages[0]) || (cpu.NumCpu != cpuFromStorage.NumCpu) || (cpu.LoadAveragesPerCpu[0] != cpuFromStorage.LoadAveragesPerCpu[0]) {
		t.Errorf("Cpu data was not saved properly. cpu: %v, cpuFromStorage: %v", cpu, cpuFromStorage)
	}

	// Cleanup
	storage.Delete(fmt.Sprintf("/hosts/%v/metrics/cpu", host))
}

func TestLoadCpuFromAllHosts(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	host := "127.0.0.1"
	storage := chillax_storage.NewStorage()

	storage.Delete(fmt.Sprintf("/hosts/%v/metrics/cpu", host))

	cpu, err := SaveCpu(storage, host)
	if err != nil {
		t.Errorf("Saving CPU data should work. err: %v", err)
	}

	cpusFromStorage, err := LoadCpuFromAllHosts(storage)

	for h, cpuData := range cpusFromStorage {
		if h == host {
			if (cpu.LoadAverages[0] != cpuData.LoadAverages[0]) || (cpu.NumCpu != cpuData.NumCpu) || (cpu.LoadAveragesPerCpu[0] != cpuData.LoadAveragesPerCpu[0]) {
				t.Errorf("Cpu data was not saved properly. cpu: %v, cpusFromStorage: %v", cpu, cpuData)
			}
		}
	}

	// Cleanup
	storage.Delete(fmt.Sprintf("/hosts/%v/metrics/cpu", host))
}
