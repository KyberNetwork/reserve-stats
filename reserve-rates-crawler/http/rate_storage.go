package http

import (
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type rateStorage interface {
	GetRatesByTimePoint(rsvAddr ethereum.Address, fromTime, toTime uint64) (map[uint64]common.ReserveRates, error)
}
