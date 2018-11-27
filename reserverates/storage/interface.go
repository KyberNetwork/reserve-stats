package storage

import (
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

//ReserveRatesStorage defines a set of interface for reserve rate storage, which can be implemented by any DB
type ReserveRatesStorage interface {
	UpdateRatesRecords(uint64, map[string]map[string]common.ReserveRateEntry) error
	GetRatesByTimePoint(addrs []ethereum.Address, fromTime, toTime uint64) (map[string]map[string][]common.ReserveRates, error)
	LastBlock() (int64, error)
}
