package common

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

// ReserveRateEntry hold 4 float number represent necessary data for a rate entry
type ReserveRateEntry struct {
	BuyReserveRate  float64 `json:"buy_reserve_rate"`
	BuySanityRate   float64 `json:"buy_sanity_rate"`
	SellReserveRate float64 `json:"sell_reserve_rate"`
	SellSanityRate  float64 `json:"sell_sanity_rate"`
}

// LastRate store last rate
type LastRate struct {
	FromBlock uint64 `json:"from_block"`
	ToBlock   uint64 `json:"to_block"`
	Rate      *ReserveRateEntry
}

// NewReserveRateEntry returns new ReserveRateEntry from results of GetReserveRate method.
// The reserve rates are stored in following structure:
// - reserveRate: [sellReserveRate(index: 0)]-[buyReserveRate)(index: 0)]-[sellReserveRate(index: 1)]-[buyReserveRate)(index: 1)]...
// - sanityRate: [sellSanityRate(index: 0)]-[buySanityRate)(index: 0)]-[sellSanityRate(index: 1)]-[buySanityRate)(index: 1)]...
func NewReserveRateEntry(reserveRates, sanityRates []*big.Int, index int) ReserveRateEntry {
	return ReserveRateEntry{
		BuyReserveRate:  fromWeiETH(reserveRates[index*2+1]),
		BuySanityRate:   fromWeiETH(sanityRates[index*2+1]),
		SellReserveRate: fromWeiETH(reserveRates[index*2]),
		SellSanityRate:  fromWeiETH(sanityRates[index*2]),
	}
}

func fromWeiETH(amount *big.Int) float64 {
	const ethDecimal = 18
	if amount == nil {
		return 0
	}
	f := new(big.Float).SetInt(amount)
	power := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10), big.NewInt(ethDecimal), nil,
	))
	res := new(big.Float).Quo(f, power)
	result, _ := res.Float64()
	return result
}

// ReserveRates hold all the pairs's rate for a particular reserve and metadata
type ReserveRates struct {
	Timestamp time.Time        `json:"timestamp"`
	FromBlock uint64           `json:"from_block"`
	ToBlock   uint64           `json:"to_block"`
	Rates     ReserveRateEntry `json:"rates"`
}

// MarshalJSON implements custom JSON marshaler for ReserveRates to format timestamp in unix millis instead of RFC3339.
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

// UnmarshalJSON implements custom JSON unmarshaler for ReserveRates to format timestamp in unix millis instead of RFC3339.
func (rr *ReserveRates) UnmarshalJSON(data []byte) error {
	type AliasReserveRates ReserveRates
	decoded := new(struct {
		Timestamp uint64 `json:"timestamp"`
		AliasReserveRates
	})

	if err := json.Unmarshal(data, decoded); err != nil {
		return err
	}
	rr.Timestamp = timeutil.TimestampMsToTime(decoded.Timestamp)
	rr.FromBlock = decoded.FromBlock
	rr.ToBlock = decoded.ToBlock
	rr.Rates = decoded.Rates
	return nil
}
