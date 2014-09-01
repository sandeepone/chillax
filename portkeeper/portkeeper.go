package portkeeper

import (
    "fmt"
    "os/exec"
    "strconv"
    "sync"
    "github.com/didip/chillax/libstring"
    "github.com/didip/chillax/libnumber"
    chillax_storage "github.com/didip/chillax/storage"
)

const MAX_PORT = 65535

func LsofPort(port int) ([]byte, error) {
    return exec.Command("lsof", "-i", fmt.Sprintf(":%v", port)).Output()
}

func ReservePort(host string) int {
    host   = libstring.HostWithoutPort(host)
    store := chillax_storage.NewStorage()
    mutex := &sync.Mutex{}

    mutex.Lock()

    usedPorts, _ := store.List(fmt.Sprintf("/hosts/%v/used-ports", host))

    var reservedPort int

    if len(usedPorts) == 0 {
        reservedPort = MAX_PORT
        store.Create(fmt.Sprintf("/hosts/%v/used-ports/%v", host, reservedPort), make([]byte, 0))
        return reservedPort
    }

    usedIntPorts := make([]int, len(usedPorts))
    for index, port := range usedPorts {
        usedIntPorts[index], _ = strconv.Atoi(port)
    }

    newSmallestPort := usedIntPorts[0] - 1
    firstGapPort    := libnumber.FirstGapIntSlice(usedIntPorts)
    reservedPort     = libnumber.LargestInt([]int{firstGapPort, newSmallestPort})

    store.Create(fmt.Sprintf("/hosts/%v/used-ports/%v", host, reservedPort), make([]byte, 0))

    mutex.Unlock()

    return reservedPort
}
