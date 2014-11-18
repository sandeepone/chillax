package host

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/chillaxio/chillax/libmetrics"
	"github.com/chillaxio/chillax/libnumber"
	"github.com/chillaxio/chillax/libprocess"
	"github.com/chillaxio/chillax/libstring"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"strconv"
	"strings"
	"sync"
)

const MAX_PORT = 65535

var ModuleMutex = &sync.Mutex{}

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

func (chost *ChillaxHost) ReservePort() int {
	host := libstring.HostWithoutPort(chost.Name)

	ModuleMutex.Lock()

	usedPorts, _ := chost.storage.List(fmt.Sprintf("/hosts/%v/used-ports", host))

	var reservedPort int

	if len(usedPorts) == 0 {
		reservedPort = MAX_PORT
		chost.storage.Create(fmt.Sprintf("/hosts/%v/used-ports/%v", host, reservedPort), make([]byte, 0))

		ModuleMutex.Unlock()
		return reservedPort
	}

	usedIntPorts := make([]int, len(usedPorts))
	for index, port := range usedPorts {
		usedIntPorts[index], _ = strconv.Atoi(port)
	}

	newSmallestPort := usedIntPorts[0] - 1
	firstGapPort := libnumber.FirstGapIntSlice(usedIntPorts)
	reservedPort = libnumber.LargestInt([]int{firstGapPort, newSmallestPort})

	chost.storage.Create(fmt.Sprintf("/hosts/%v/used-ports/%v", host, reservedPort), make([]byte, 0))

	ModuleMutex.Unlock()

	return reservedPort
}

func (chost *ChillaxHost) CleanReservedPorts() error {
	var err error

	host := libstring.HostWithoutPort(chost.Name)

	usedPorts, err := chost.storage.List(fmt.Sprintf("/hosts/%v/used-ports", host))
	if err != nil || (len(usedPorts) == 0) {
		return err
	}

	lsofOutput, _ := libprocess.LsofAll()
	lsofOutputString := string(lsofOutput)

	if lsofOutputString != "" {
		for _, port := range usedPorts {
			chunk := fmt.Sprintf(":%v ", port)

			if !strings.Contains(lsofOutputString, chunk) {
				err = chost.storage.Delete(fmt.Sprintf("/hosts/%v/used-ports/%v", host, port))
			}
		}
	}
	return err
}
