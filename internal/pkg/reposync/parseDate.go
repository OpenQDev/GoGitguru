package reposync

import "time"

func ParseDate(dateStr string) (time.Time, error) {
	const (
		DateFormat       = "2006-01-02"
		DefaultStartDate = "2020-01-01"
	)

	if dateStr == "" {
		dateStr = DefaultStartDate
	}
	return time.Parse(DateFormat, dateStr)
}
