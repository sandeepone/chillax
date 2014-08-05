package dockerinventory

import (
    "fmt"
    "testing"
    chillax_storage "github.com/didip/chillax/storage"
)

func TestReserveLargestPortWhichIsTheDefault(t *testing.T) {
    dockerHost := "127.0.0.1:2375"
    storage    := chillax_storage.NewStorage()

    storage.Delete(fmt.Sprintf("/dockers/%v/used-ports/", dockerHost))

    port := ReservePort(dockerHost)

    if port != MAX_PORT {
        t.Errorf("port should equal to %v", MAX_PORT)
    }

    storage.Delete(fmt.Sprintf("/dockers/%v/used-ports/", dockerHost))

    usedPorts, _ := storage.List(fmt.Sprintf("/dockers/%v/used-ports", dockerHost))

    if len(usedPorts) != 0 {
        t.Errorf("Total used ports should be 0. Actual number: %v", len(usedPorts))
    }
}

func TestReserveGapPort(t *testing.T) {
    dockerHost := "127.0.0.1:2375"
    storage    := chillax_storage.NewStorage()

    storage.Delete(fmt.Sprintf("/dockers/%v/used-ports/", dockerHost))

    ReservePort(dockerHost)
    deleteme := ReservePort(dockerHost)
    ReservePort(dockerHost)

    usedPorts, _ := storage.List(fmt.Sprintf("/dockers/%v/used-ports", dockerHost))

    if len(usedPorts) != 3 {
        t.Errorf("Total used ports should be 3. Actual number: %v", len(usedPorts))
    }

    storage.Delete(fmt.Sprintf("/dockers/%v/used-ports/%v", dockerHost, deleteme))

    port := ReservePort(dockerHost)

    if port != deleteme {
        t.Errorf("ReservePort did not regenerate gap port. Port: %v", port)
    }
}