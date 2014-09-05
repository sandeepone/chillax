package libtime

import (
	"time"
)

func SleepString(definition string) error {
	delayTime, err := time.ParseDuration(definition)
	if err != nil {
		return err
	}

	time.Sleep(delayTime)
	return nil
}
