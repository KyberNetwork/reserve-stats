package common

import "time"

// ReserveTokenRateEntry is a  map[ETH-tokenID]ratesEntry
type ReserveTokenRateEntry map[string]ReserveRateEntry

// ReserveRateEntry hold 4 float number represent necessary data for a rate entry
type ReserveRateEntry struct {
	BuyReserveRate  float64
	BuySanityRate   float64
	SellReserveRate float64
	SellSanityRate  float64
}

// ReserveRates hold all the pairs's rate for a particular reserve and metadata
type ReserveRates struct {
	Timestamp   time.Time
	ReturnTime  time.Time
	BlockNumber int64
	Data        ReserveTokenRateEntry
}
