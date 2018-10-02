package ethrate

// EthUSDRate contains method to get ETH/USD conversion rate.
type EthUSDRate interface {
	GetUSDRate(timepoint uint64) float64
}
