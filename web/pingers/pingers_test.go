package pingers

import (
	"strings"
	"testing"
	"time"

	chillax_storage "github.com/chillaxio/chillax/storage"
)

func TestNewPinger(t *testing.T) {
	pinger := NewPinger("http://localhost:8080/chillax/admin")

	if pinger.Method != "GET" {
		t.Error("Default method should be POST.")
	}
	if pinger.TimeoutString != "1s" {
		t.Error("Default TimeoutString should be 1 second.")
	}
	if pinger.FailCount != 0 {
		t.Error("Default FailCount should be 0.")
	}
	if pinger.FailMax != 10 {
		t.Error("Default FailMax should be 10.")
	}
}

func TestPingerIsUp(t *testing.T) {
	pinger := NewPinger("http://localhost:8080/chillax/admin")

	isUp, err := pinger.IsUp()

	if err != nil && !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("Pinger should indicates connection refuse. Error: %v", err)
	}

	if isUp {
		t.Error("Pinger should indicates that endpoint is down.")
	}
}

func TestPingerFailCount(t *testing.T) {
	chillax_storage.NewStorage().Delete("/pingers")

	pinger := NewPinger("http://localhost:8080/chillax/admin")

	pinger.IsUp()

	if pinger.FailCount != 1 {
		t.Errorf("Pinger.FailCount should increase by 1. pinger.FailCount: %v", pinger.FailCount)
	}
}

func TestNewPingerGroup(t *testing.T) {
	pg := NewPingerGroup([]string{"http://localhost:8080/chillax/admin"})

	if pg.SleepTime != 1*time.Minute {
		t.Error("Default SleepTime should be 1 minute.")
	}
	if len(pg.Pingers) != 1 {
		t.Error("There should be 1 URL in list of Pingers.")
	}
}

func TestPingerGroupSave(t *testing.T) {
	pg := NewPingerGroup([]string{"http://localhost:8080/chillax/admin"})

	for uri, pinger := range pg.Pingers {
		pg.IsOnePingerUp(uri, pinger)
	}

	err := pg.Save()
	if err != nil {
		t.Errorf("Unable to serialize and save PingerGroup. Error: %v", err)
	}
}
