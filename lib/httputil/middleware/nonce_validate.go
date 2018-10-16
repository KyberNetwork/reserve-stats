package middleware

import "time"

const maxTimeGapMillis = 30000 // 30 secs

// UnixMillis return current Unix timestamp in milliseconds
func UnixMillis() uint64 {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return uint64(timestamp)
}

// NonceValidator interface for checking if a nonce is valid or not
type NonceValidator interface {
	IsValid(int64) bool
}

// nonceValidatorByTime checking validate by time range
type nonceValidatorByTime struct {
	// timeGap is max time different between client submit timestamp
	// and server time that considered valid. The time precision is millisecond.
	timeGap uint64
}

// newValidateNonceByTime return nonceValidatorByTime with default value
func newValidateNonceByTime() nonceValidatorByTime {
	return nonceValidatorByTime{
		timeGap: maxTimeGapMillis,
	}
}

// IsValid return nonce is valid or not by time range
func (v nonceValidatorByTime) IsValid(nonce int64) bool {
	serverTime := UnixMillis()
	difference := nonce - int64(serverTime)
	if difference < -maxTimeGapMillis || difference > maxTimeGapMillis {
		return false
	}
	return true
}
