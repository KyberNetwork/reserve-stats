package storage

import (
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	lastblockdaily "github.com/KyberNetwork/reserve-stats/lib/lastblockdaily/common"
)

//Interface defines a set of interface for reserve rate storage, which can be implemented by any DB
type Interface interface {
	UpdateRatesRecords(lastblockdaily.BlockInfo, map[string]map[string]float64, float64) error
	GetLastResolvedBlockInfo(ethereum.Address) (lastblockdaily.BlockInfo, error)
	GetETHUSDRates(from, to time.Time) (AccountingReserveRates, error)
	GetRates(from, to time.Time) (map[string]AccountingReserveRates, error)
}
