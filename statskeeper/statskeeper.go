package statskeeper

import (
	"encoding/json"

	"fmt"
	"github.com/Sirupsen/logrus"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"time"
)

func SaveRequest(currentTime time.Time, fields logrus.Fields) {
	datetime := time.Unix(0, currentTime.UnixNano())
	dataPath := fmt.Sprintf(
		"/logs/requests/%v/%d/%v/%v/%v/%v-%v",
		datetime.Year(), datetime.Month(), datetime.Day(), datetime.Hour(), datetime.Minute(),
		currentTime.UnixNano(), fields["Method"])

	inBytes, err := json.Marshal(fields)
	if err == nil {
		storage := chillax_storage.NewStorage()
		storage.Create(dataPath, inBytes)
	}
}

// func GetRequestJsonBetween(begin, end string) []string {}
