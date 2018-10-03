package utils

import (
	"math"
	"math/big"
	"time"
)

// FloatToBigInt converts a float to a big int with specific decimal
// Example:
// - FloatToBigInt(1, 4) = 10000
// - FloatToBigInt(1.234, 4) = 12340
func FloatToBigInt(amount float64, decimal int64) *big.Int {
	// 6 is our smallest precision
	if decimal < 6 {
		return big.NewInt(int64(amount * math.Pow10(int(decimal))))
	}
	result := big.NewInt(int64(amount * math.Pow10(6)))
	return result.Mul(result, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimal-6), nil))
}

// EthToWei converts Gwei as a float to Wei as a big int
func EthToWei(n float64) *big.Int {
	return FloatToBigInt(n, 18)
}

//GetTimepoint return current timepoint by millisecond
func GetTimepoint() uint64 {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return uint64(timestamp)
}
