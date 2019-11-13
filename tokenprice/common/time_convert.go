package common

import "time"

const (
	timeLayout = "2006-01-02"
)

// DateStringToTime convert date string to time with format timeLayout
func DateStringToTime(date string) (time.Time, error) {
	return time.Parse(timeLayout, date)
}

// TimeToDateString convert time to date string with format timeLayout
func TimeToDateString(t time.Time) string {
	return t.Format(timeLayout)
}
