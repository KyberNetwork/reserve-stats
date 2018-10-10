package common

import (
	"math/big"
	"time"
)

// BigToFloat converts a big int to float according to its number of decimal digits
// Example:
// - BigToFloat(1100, 3) = 1.1
// - BigToFloat(1100, 2) = 11
// - BigToFloat(1100, 5) = 0.11
func BigToFloat(b *big.Int, decimal int64) float64 {
	f := new(big.Float).SetInt(b)
	power := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10), big.NewInt(decimal), nil,
	))
	res := new(big.Float).Quo(f, power)
	result, _ := res.Float64()
	return result
}

// TimestampMsToTime turn a uint64 timestamp in millisecond to a golang time object
func TimestampMsToTime(ms uint64) time.Time {
	return time.Unix(0, int64(ms)*int64(time.Millisecond))
}

func TimeToTimestampMs(t time.Time) uint64 {
	return uint64(t.UnixNano() / int64(time.Millisecond))
}
