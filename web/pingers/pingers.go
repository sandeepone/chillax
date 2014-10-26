package pingers

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	chillax_storage "github.com/chillaxio/chillax/storage"

	"github.com/franela/goreq"
)

// NewPinger is the constructor for Pinger
func NewPinger(uri string) *Pinger {
	p := &Pinger{}
	p.Uri = uri
	p.Method = "GET"
	p.TimeoutString = "1s"
	p.FailMax = 10
	return p
}

// Pinger checks endpoint using GET request.
type Pinger struct {
	goreq.Request

	// Default is "1s"
	TimeoutString string
	FailCount     int

	// Default is 10
	FailMax int
}

// IsUp checks if endpoint's status code == 200.
func (p *Pinger) IsUp() (bool, error) {
	timeoutTime, err := time.ParseDuration(p.TimeoutString)
	if err != nil {
		return false, err
	}

	p.Timeout = timeoutTime

	res, err := p.Do()
	if err != nil {
		p.FailCount++
		return false, err
	}

	if res.StatusCode != 200 {
		p.FailCount++
		return false, nil
	}

	return true, nil
}

// BelowsFailMax checks if FailCount is still below FailMax
func (p *Pinger) BelowsFailMax() bool {
	return p.FailCount < p.FailMax
}

// NewPingerGroup is the constructor for PingerGroup
func NewPingerGroup(uris []string) *PingerGroup {
	pg := &PingerGroup{}

	pg.SleepTime = 1 * time.Minute

	pg.Pingers = make(map[string]*Pinger)
	for _, uri := range uris {
		pg.Pingers[uri] = NewPinger(uri)
	}

	pg.PingersCheck = make(map[string]bool)

	return pg
}

// PingerGroup is a collection of Pingers.
type PingerGroup struct {
	SleepTime    time.Duration
	Pingers      map[string]*Pinger
	PingersCheck map[string]bool
}

// IsUpAsync checks all endpoints in their own goroutines.
func (pg *PingerGroup) IsUpAsync() {
	for uri, pinger := range pg.Pingers {
		go func(uri string, pinger *Pinger) {
			sleepTime := pg.SleepTime

			for {
				time.Sleep(sleepTime)

				isUp, _ := pinger.IsUp()

				pg.PingersCheck[uri] = isUp

				// Sleeps longer if pinger exceeds FailMax
				if !pinger.BelowsFailMax() {
					sleepTime = sleepTime * 2
				}

				// Return sleepTime back to original if endpoint is finally up
				// or sleepTime is greater than 2 hours
				if isUp || sleepTime > (2*time.Hour) {
					sleepTime = pg.SleepTime
				}

				// Serialize and save checks data to storage.
				go func(pg *PingerGroup) {
					pg.Save()
				}(pg)
			}
		}(uri, pinger)
	}
}

// Serialize current checks data.
func (pg *PingerGroup) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(pg.PingersCheck)

	return buffer.Bytes(), err
}

// Save current checks data to storage.
func (pg *PingerGroup) Save() error {
	hostname, _ := os.Hostname()

	inBytes, err := pg.Serialize()
	if err != nil {
		return err
	}

	return chillax_storage.NewStorage().Update(fmt.Sprintf("/pingers/%v", hostname), inBytes)
}
