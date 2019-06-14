package postgrestorage

import (
	ethereum "github.com/ethereum/go-ethereum/common"
	"time"
)

//TODO implement this
func (tldb *TradeLogDB) GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddrs []ethereum.Address) (map[ethereum.Address]map[string]float64, error) {
	return nil, nil
}
