package common

import "time"

const (
	timeLayout = "2006-01-02"
)

// YYYYMMDDToTime convert date string in format YYYY-MM-DD to time
func YYYYMMDDToTime(date string) (time.Time, error) {
	return time.Parse(timeLayout, date)
}

// TimeToYYYYMMDD convert time to date string with format YYYY-MM-DD
func TimeToYYYYMMDD(t time.Time) string {
	return t.Format(timeLayout)
}
