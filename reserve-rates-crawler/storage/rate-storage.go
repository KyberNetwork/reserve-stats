package storage

import (
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/crawler"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type ReserveRatesStorage interface {
	UpdateRatesRecords(rateRecords map[string]crawler.ReserveRates) error
	GetRatesByTimePoint(rsvAddr ethereum.Address, fromTime, toTime int64) ([]crawler.ReserveRates, error)
}
