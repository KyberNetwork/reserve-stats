package ethrate

import "time"

// EthUSDRate contains method to get ETH/USD conversion rate.
type EthUSDRate interface {
	GetUSDRate(timepoint time.Time) float64
}
