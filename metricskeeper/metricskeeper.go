package metricskeeper

import (
	"fmt"
	"github.com/chillaxio/chillax/libmetrics"
	chillax_storage "github.com/chillaxio/chillax/storage"
)

func SaveCpu(storage chillax_storage.Storer, host string) error {
	dataPath := fmt.Sprintf("/hosts/%v/metrics/cpu", host)
	cpu := libmetrics.NewCpuMetrics()
	cpuToml, err := cpu.ToToml()

	if err != nil {
		return err
	}

	return storage.Create(dataPath, cpuToml)
}
