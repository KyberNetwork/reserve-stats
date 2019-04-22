package timeutil

import "time"

// TimestampMsToTime turn a uint64 timestamp in millisecond to a golang time object
func TimestampMsToTime(ms uint64) time.Time {
	return time.Unix(0, int64(ms)*int64(time.Millisecond))
}

// TimeToTimestampMs turn a golang time object into uint64 timestamp in millisecond
func TimeToTimestampMs(t time.Time) uint64 {
	return uint64(t.UnixNano() / int64(time.Millisecond))
}

// UnixMilliSecond return current timestamp in millisecond
func UnixMilliSecond() uint64 {
	return TimeToTimestampMs(time.Now())
}

// Midnight returns midnight of given timestamp.
// Truncate(24*time.Hour) only works for UTC based.
// https://github.com/golang/go/commit/e1ced3219506938daf404bb2373333cd3352f350
func Midnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// TimestampSecondToNs convert time stamp in second to nanosecond
func TimestampSecondToNs(ts uint64) uint64 {
	const secondToNs = 1000000
	return ts * secondToNs
}
