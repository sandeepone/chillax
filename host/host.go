package host

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/chillaxio/chillax/libmetrics"
	chillax_storage "github.com/chillaxio/chillax/storage"
)

func GetAllHosts(storage chillax_storage.Storer) ([]*ChillaxHost, error) {
	hosts, err := storage.List("/hosts")
	if err != nil {
		return nil, err
	}

	chosts := make([]*ChillaxHost, len(hosts))

	for i, host := range hosts {
		chosts[i] = NewChillaxHost(storage, host)
	}

	return chosts, nil
}

func GetCpuFromAllHosts(storage chillax_storage.Storer) (map[string]*libmetrics.CpuMetrics, error) {
	hosts, err := storage.List("/hosts")
	if err != nil {
		return nil, err
	}

	data := make(map[string]*libmetrics.CpuMetrics)

	for _, host := range hosts {
		chost := NewChillaxHost(storage, host)
		cpu, err := chost.LoadCpu()

		if err == nil {
			data[host] = cpu
		}
	}
	return data, err
}

func NewChillaxHost(storage chillax_storage.Storer, name string) *ChillaxHost {
	chost := &ChillaxHost{}
	chost.Name = name
	chost.storage = storage

	return chost
}

type ChillaxHost struct {
	Name       string
	CpuMetrics *libmetrics.CpuMetrics
	storage    chillax_storage.Storer
}

func (chost *ChillaxHost) SaveCpu() error {
	dataPath := fmt.Sprintf("/hosts/%v/metrics/cpu", chost.Name)
	cpu, err := libmetrics.NewCpuMetrics()
	if err != nil {
		return err
	}

	cpuToml, err := cpu.ToToml()

	if err != nil {
		return err
	}

	err = chost.storage.Update(dataPath, cpuToml)
	if err == nil {
		chost.CpuMetrics = cpu
	}

	return err
}

func (chost *ChillaxHost) LoadCpu() (*libmetrics.CpuMetrics, error) {
	cpuTomlBytes, err := chost.storage.Get(fmt.Sprintf("/hosts/%v/metrics/cpu", chost.Name))

	if err != nil {
		return nil, err
	}

	cpu, err := libmetrics.NewCpuMetrics()
	if err != nil {
		return nil, err
	}

	_, err = toml.Decode(string(cpuTomlBytes), cpu)
	if err != nil {
		return nil, err
	}

	chost.CpuMetrics = cpu

	return cpu, err
}
