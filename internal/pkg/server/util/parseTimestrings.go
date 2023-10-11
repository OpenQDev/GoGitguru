package util

import (
	"time"
)

// ISO 8601 format
const TIME_LAYOUT = "2006-01-02T15:04:05Z"

func ParseTimestrings(timeStrings ...string) ([]time.Time, error) {
	parsedTimes := make([]time.Time, len(timeStrings))

	for i, timeString := range timeStrings {
		parsedTime, err := time.Parse(TIME_LAYOUT, timeString)
		if err != nil {
			return nil, err
		}
		parsedTimes[i] = parsedTime
	}

	return parsedTimes, nil
}
