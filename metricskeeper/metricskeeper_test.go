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

func TestReserveLargestPortWhichIsTheDefault(t *testing.T) {
	host := "127.0.0.1"
	storage := chillax_storage.NewStorage()

	storage.Delete(fmt.Sprintf("/hosts/%v/metrics/cpu", host))

	_, err := GetCpuTomlForTest(host)
	if err == nil {
		t.Errorf("Cpu data should be nil.")
	}

	err = SaveCpu(storage, host)
	if err != nil {
		t.Errorf("Saving CPU data should work. err: %v", err)
	}

	_, err = GetCpuTomlForTest(host)
	if err != nil {
		t.Errorf("Cpu data should not be nil.")
	}

	// Cleanup
	storage.Delete(fmt.Sprintf("/hosts/%v/metrics/cpu", host))
}
