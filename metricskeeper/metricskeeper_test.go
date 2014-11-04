package metricskeeper

import (
	"fmt"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"testing"
)

func GetCpuTomlForTest(host string) ([]byte, error) {
	storage := chillax_storage.NewStorage()
	return storage.Get(fmt.Sprintf("/hosts/%v/metrics/cpu", host))
}

func TestSaveAndLoadCpu(t *testing.T) {
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
	host := "127.0.0.1"
	storage := chillax_storage.NewStorage()

	storage.Delete(fmt.Sprintf("/hosts/%v/metrics/cpu", host))

	cpu, err := SaveCpu(storage, host)
	if err != nil {
		t.Errorf("Saving CPU data should work. err: %v", err)
	}

	cpusFromStorage, err := LoadCpuFromAllHosts(storage)

	if (cpu.LoadAverages[0] != cpusFromStorage[0].LoadAverages[0]) || (cpu.NumCpu != cpusFromStorage[0].NumCpu) || (cpu.LoadAveragesPerCpu[0] != cpusFromStorage[0].LoadAveragesPerCpu[0]) {
		t.Errorf("Cpu data was not saved properly. cpu: %v, cpusFromStorage: %v", cpu, cpusFromStorage)
	}

	if len(cpusFromStorage) != 1 {
		t.Errorf("Should get only 1 cpu data. cpusFromStorage: %v", cpusFromStorage)
	}

	// Cleanup
	storage.Delete(fmt.Sprintf("/hosts/%v/metrics/cpu", host))
}
