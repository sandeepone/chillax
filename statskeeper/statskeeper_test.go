package statskeeper

import (
	"strings"
	"testing"
	"time"
)

func TestGetRequestDataPathsDurationsAgo(t *testing.T) {
	dataPaths, err := GetRequestDataPathsDurationsAgo(time.Now(), "-5h")
	if err != nil {
		t.Errorf("Unable to get dataPaths. err: %v", err)
	}
	if len(dataPaths) <= 0 {
		t.Errorf("Unable to get dataPaths. dataPaths: %v", dataPaths)
	}
}

func TestGetRequestDataDurationsAgo(t *testing.T) {
	data, err := GetRequestDataDurationsAgo(time.Now(), "-5h")
	if err != nil {
		t.Errorf("Unable to get data. err: %v", err)
	}
	if len(data) <= 0 {
		t.Error("Unable to get data.")
	}

	firstData := data[0]
	firstDataString := string(firstData)

	for _, key := range []string{"CurrentUnix", "CurrentUnixNano", "Latency"} {
		if !strings.Contains(firstDataString, key) {
			t.Errorf("Bad key on data. key: %v, firstData: %v", key, firstDataString)
		}
	}
}
