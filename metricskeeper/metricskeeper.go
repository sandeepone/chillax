package metricskeeper

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/chillaxio/chillax/libmetrics"
	"github.com/chillaxio/chillax/libtime"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"os"
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

func SaveCpuAsync(storage chillax_storage.Storer, intervalString string) {
	host, _ := os.Hostname()

	go func() {
		for {
			SaveCpu(storage, host)
			libtime.SleepString(intervalString)
		}
	}()
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

func LoadCpuFromAllHosts(storage chillax_storage.Storer) ([]*libmetrics.CpuMetrics, error) {
	hosts, err := storage.List("/hosts")
	if err != nil {
		return nil, err
	}

	data := make([]*libmetrics.CpuMetrics, 0)
	for _, host := range hosts {
		cpu, err := LoadCpu(storage, host)
		if err == nil {
			data = append(data, cpu)
		}
	}
	return data, err
}
