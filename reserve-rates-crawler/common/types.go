package common

import (
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"math/big"
	"time"
)

// ReserveRateEntry hold 4 float number represent necessary data for a rate entry
type ReserveRateEntry struct {
	BuyReserveRate  float64
	BuySanityRate   float64
	SellReserveRate float64
	SellSanityRate  float64
}

// NewReserveRateEntry returns new ReserveRateEntry from results of GetReserveRate method.
// The reserve rates are stored in following structure:
// - reserveRate: [sellReserveRate(index: 0)]-[buyReserveRate)(index: 0)]-[sellReserveRate(index: 1)]-[buyReserveRate)(index: 1)]...
// - sanityRate: [sellSanityRate(index: 0)]-[buySanityRate)(index: 0)]-[sellSanityRate(index: 1)]-[buySanityRate)(index: 1)]...
func NewReserveRateEntry(reserveRates, sanityRates []*big.Int, index int) ReserveRateEntry {
	return ReserveRateEntry{
		BuyReserveRate:  core.ETHToken.FormatAmount(reserveRates[index*2+1]),
		BuySanityRate:   core.ETHToken.FormatAmount(sanityRates[index*2+1]),
		SellReserveRate: core.ETHToken.FormatAmount(reserveRates[index*2]),
		SellSanityRate:  core.ETHToken.FormatAmount(sanityRates[index*2]),
	}
}

// ReserveRates hold all the pairs's rate for a particular reserve and metadata
type ReserveRates struct {
	Timestamp   time.Time
	BlockNumber uint64
	Data        map[string]ReserveRateEntry
}
