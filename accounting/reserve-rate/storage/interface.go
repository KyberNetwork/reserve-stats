package storage

import (
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
)

//Interface defines a set of interface for reserve rate storage, which can be implemented by any DB
type Interface interface {
	UpdateRatesRecords(lastblockdaily.BlockInfo, map[string]map[string]common.ReserveRateEntry) error
	UpdateETHUSDPrice(blockInfo lastblockdaily.BlockInfo, ethusdRate float64) error
	GetLastResolvedBlockInfo() (lastblockdaily.BlockInfo, error)
}
