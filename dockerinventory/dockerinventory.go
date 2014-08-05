package dockerinventory

import (
    "fmt"
    "strconv"
    "github.com/didip/chillax/libstring"
    "github.com/didip/chillax/libnumber"
    chillax_storage "github.com/didip/chillax/storage"
)

const MAX_PORT = 65536

func ReservePort(dockerHost string) int {
    dockerHost    = libstring.StripProtocol(dockerHost)
    store        := chillax_storage.NewStorage()
    usedPorts, _ := store.List(fmt.Sprintf("/dockers/%v/used-ports", dockerHost))

    var reservedPort int

    if len(usedPorts) == 0 {
        reservedPort = MAX_PORT
        store.Create(fmt.Sprintf("/dockers/%v/used-ports/%v", dockerHost, reservedPort), make([]byte, 0))
        return reservedPort
    }

    usedIntPorts := make([]int, len(usedPorts))
    for index, port := range usedPorts {
        usedIntPorts[index], _ = strconv.Atoi(port)
    }

    newSmallestPort := usedIntPorts[0] - 1
    firstGapPort    := libnumber.FirstGapIntSlice(usedIntPorts)
    reservedPort     = libnumber.LargestInt([]int{firstGapPort, newSmallestPort})

    store.Create(fmt.Sprintf("/dockers/%v/used-ports/%v", dockerHost, reservedPort), make([]byte, 0))
    return reservedPort
}
