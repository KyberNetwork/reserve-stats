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
