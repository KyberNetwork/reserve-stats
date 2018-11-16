package common

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

// ReserveRateEntry hold 4 float number represent necessary data for a rate entry
type ReserveRateEntry struct {
	BuyReserveRate  float64 `json:"buy_reserve_rate"`
	BuySanityRate   float64 `json:"buy_sanity_rate"`
	SellReserveRate float64 `json:"sell_reserve_rate"`
	SellSanityRate  float64 `json:"sell_sanity_rate"`
}

// NewReserveRateEntry returns new ReserveRateEntry from results of GetReserveRate method.
// The reserve rates are stored in following structure:
// - reserveRate: [sellReserveRate(index: 0)]-[buyReserveRate)(index: 0)]-[sellReserveRate(index: 1)]-[buyReserveRate)(index: 1)]...
// - sanityRate: [sellSanityRate(index: 0)]-[buySanityRate)(index: 0)]-[sellSanityRate(index: 1)]-[buySanityRate)(index: 1)]...
func NewReserveRateEntry(reserveRates, sanityRates []*big.Int, index int) ReserveRateEntry {
	return ReserveRateEntry{
		BuyReserveRate:  core.ETHToken.FromWei(reserveRates[index*2+1]),
		BuySanityRate:   core.ETHToken.FromWei(sanityRates[index*2+1]),
		SellReserveRate: core.ETHToken.FromWei(reserveRates[index*2]),
		SellSanityRate:  core.ETHToken.FromWei(sanityRates[index*2]),
	}
}

// ReserveRates hold all the pairs's rate for a particular reserve and metadata
type ReserveRates struct {
	Timestamp   time.Time                   `json:"timestamp"`
	BlockNumber uint64                      `json:"-"`
	Data        map[string]ReserveRateEntry `json:"data"`
	Reserve     string                      `json:"-"`
}

// MarshalJSON implements custom JSON marshaler for TradeLog to format timestamp in unix millis instead of RFC3339.
func (rr ReserveRates) MarshalJSON() ([]byte, error) {
	type AliasReserveRates ReserveRates
	return json.Marshal(struct {
		Timestamp uint64 `json:"timestamp"`
		AliasReserveRates
	}{
		AliasReserveRates: (AliasReserveRates)(rr),
		Timestamp:         timeutil.TimeToTimestampMs(rr.Timestamp),
	})
}
