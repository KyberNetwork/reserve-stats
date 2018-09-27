package common

import "time"

// GetTimepoint return current timestamp in millisecond.
func GetTimepoint() uint64 {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return uint64(timestamp)
}
