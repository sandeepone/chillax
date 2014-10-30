package portkeeper

import (
	"fmt"
	"github.com/chillaxio/chillax/libnumber"
	"github.com/chillaxio/chillax/libstring"
	"github.com/chillaxio/chillax/libtime"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const MAX_PORT = 65535

var ModuleMutex = &sync.Mutex{}

func LsofAll() ([]byte, error) {
	return exec.Command("lsof", "-i").Output()
}

func LsofPort(port int) ([]byte, error) {
	return exec.Command("lsof", "-i", fmt.Sprintf(":%v", port)).Output()
}

func ReservePort(host string) int {
	host = libstring.HostWithoutPort(host)
	store := chillax_storage.NewStorage()

	ModuleMutex.Lock()

	usedPorts, _ := store.List(fmt.Sprintf("/hosts/%v/used-ports", host))

	var reservedPort int

	if len(usedPorts) == 0 {
		reservedPort = MAX_PORT
		store.Create(fmt.Sprintf("/hosts/%v/used-ports/%v", host, reservedPort), make([]byte, 0))

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

	store.Create(fmt.Sprintf("/hosts/%v/used-ports/%v", host, reservedPort), make([]byte, 0))

	ModuleMutex.Unlock()

	return reservedPort
}

func CleanReservedPorts(host string) error {
	var err error

	host = libstring.HostWithoutPort(host)
	store := chillax_storage.NewStorage()

	usedPorts, err := store.List(fmt.Sprintf("/hosts/%v/used-ports", host))
	if err != nil || (len(usedPorts) == 0) {
		return err
	}

	lsofOutput, _ := LsofAll()
	lsofOutputString := string(lsofOutput)

	if lsofOutputString != "" {
		for _, port := range usedPorts {
			chunk := fmt.Sprintf(":%v ", port)

			if !strings.Contains(lsofOutputString, chunk) {
				err = store.Delete(fmt.Sprintf("/hosts/%v/used-ports/%v", host, port))
			}
		}
	}
	return err
}

func CleanReservedPortsAsync(sleepString string) {
	hostname, _ := os.Hostname()

	go func(hostname string) {
		CleanReservedPorts(hostname)
		libtime.SleepString(sleepString)
	}(hostname)
}
