package storage

import (
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type ReserveRatesStorage interface {
	UpdateRatesRecords(rateRecords map[string]common.ReserveRates) error
	GetRatesByTimePoint(rsvAddr ethereum.Address, fromTime, toTime uint64) ([]common.ReserveRates, error)
}
