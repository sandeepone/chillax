package statskeeper

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"strconv"
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

// Assumptions:
// - Never get data older than this year.
func GetRequestDataPathsDurationsAgo(endDatetime time.Time, durationString string) ([]string, error) {
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return nil, err
	}

	startDatetime := endDatetime.Add(duration)

	storage := chillax_storage.NewStorage()

	dataGlobPaths := make([]string, 0)

	// Grab list of days within duration.
	days, err := storage.List(fmt.Sprintf("/logs/requests/%v/%d", startDatetime.Year(), startDatetime.Month()))
	if err != nil {
		return nil, err
	}

	sameMonthDataGlobPathsFetch := func(startDatetime time.Time, endDatetime time.Time, days []string) {
		if startDatetime.Month() == endDatetime.Month() {
			for _, day := range days {
				dayInt, err := strconv.Atoi(day)
				if err == nil && dayInt >= startDatetime.Day() && dayInt <= endDatetime.Day() {
					dataGlobPaths = append(dataGlobPaths, fmt.Sprintf("/logs/requests/%v/%d/%v/**/**/*", startDatetime.Year(), startDatetime.Month(), dayInt))
				}
			}
		}
	}

	if startDatetime.Month() != endDatetime.Month() {
		// Grab list of months within duration.
		months, err := storage.List(fmt.Sprintf("/logs/requests/%v", startDatetime.Year()))
		if err != nil {
			return nil, err
		}

		for _, month := range months {
			monthInt, err := strconv.Atoi(month)
			if err == nil && monthInt >= int(startDatetime.Month()) && monthInt < int(endDatetime.Month()) {
				dataGlobPaths = append(dataGlobPaths, fmt.Sprintf("/logs/requests/%v/%d/**/**/*", startDatetime.Year(), monthInt))
			}
			sameMonthDataGlobPathsFetch(startDatetime, endDatetime, days)
		}
	} else {
		sameMonthDataGlobPathsFetch(startDatetime, endDatetime, days)
	}

	dataPaths := make([]string, 0)
	for _, glob := range dataGlobPaths {
		paths, err := storage.ListRecursive(glob)
		if err != nil {
			return nil, err
		}
		if err == nil {
			dataPaths = append(dataPaths, paths...)
		}
	}

	return dataPaths, nil
}

func GetRequestDataDurationsAgo(endDatetime time.Time, durationString string) ([][]byte, error) {
	dataPaths, err := GetRequestDataPathsDurationsAgo(endDatetime, durationString)
	if err != nil {
		return nil, err
	}

	dataSlice := make([][]byte, len(dataPaths))
	storage := chillax_storage.NewStorage()

	for i, dataPath := range dataPaths {
		data, err := storage.Get(dataPath)
		if err != nil {
			return dataSlice, err
		}
		dataSlice[i] = data
	}
	return dataSlice, nil
}
