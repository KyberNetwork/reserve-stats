package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
)

//Interface defines a set of interface for reserve rate storage, which can be implemented by any DB
type Interface interface {
	// TODO: should we merge UpdateETHUSDPrice and UpdateETHUSDPrice?
	UpdateRatesRecords(lastblockdaily.BlockInfo, map[string]map[string]float64) error
	UpdateETHUSDPrice(blockInfo lastblockdaily.BlockInfo, ethUSDRate float64) error
	GetLastResolvedBlockInfo() (lastblockdaily.BlockInfo, error)
	// TODO: should we merge GetETHUSDRates and GetRates?
	GetETHUSDRates(from, to time.Time) (AccountingReserveRates, error)
	GetRates(from, to time.Time) (map[string]AccountingReserveRates, error)
}
