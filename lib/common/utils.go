package common

import "time"

// TimestampMsToTime turn a uint64 timestamp in millisecond to a golang time object
func TimestampMsToTime(ms uint64) time.Time {
	return time.Unix(0, int64(ms)*int64(time.Millisecond))
}
