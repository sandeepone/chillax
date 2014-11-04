package metricskeeper

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/chillaxio/chillax/libmetrics"
	chillax_storage "github.com/chillaxio/chillax/storage"
)

func SaveCpu(storage chillax_storage.Storer, host string) (*libmetrics.CpuMetrics, error) {
	dataPath := fmt.Sprintf("/hosts/%v/metrics/cpu", host)
	cpu := libmetrics.NewCpuMetrics()
	cpuToml, err := cpu.ToToml()

	if err != nil {
		return nil, err
	}

	err = storage.Create(dataPath, cpuToml)

	return cpu, err
}

func LoadCpu(storage chillax_storage.Storer, host string) (*libmetrics.CpuMetrics, error) {
	cpuTomlBytes, err := storage.Get(fmt.Sprintf("/hosts/%v/metrics/cpu", host))

	if err != nil {
		return nil, err
	}

	cpu := libmetrics.NewCpuMetrics()

	_, err = toml.Decode(string(cpuTomlBytes), cpu)
	if err != nil {
		return nil, err
	}

	return cpu, err
}
