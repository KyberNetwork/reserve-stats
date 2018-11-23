package storage

import (
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// Interface represent a storage for TradeLogs data
type Interface interface {
	LastBlock() (int64, error)
	SaveTradeLogs(logs []common.TradeLog) error
	LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error)
	GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddrs []ethereum.Address) (map[ethereum.Address]map[string]float64, error)
	GetAssetVolume(token core.Token, fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error)
	GetReserveVolume(rsvAddr ethereum.Address, token core.Token, fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error)
	GetAggregatedWalletFee(reserveAddr, walletAddr, freq string,
		fromTime, toTime time.Time, timezone int8) (map[uint64]float64, error)
	GetTradeSummary(from, to time.Time, timezone int8) (map[uint64]*common.TradeSummary, error)
	GetUserVolume(userAddr ethereum.Address, from, to time.Time, freq string) (map[uint64]common.UserVolume, error)
	GetUserList(from, to time.Time, timezone int8) ([]common.UserInfo, error)
	GetWalletStats(fromTime, toTime time.Time, walletAddr string, timezone int8) (map[uint64]common.WalletStats, error)
	GetCountryStats(countryCode string, from, to time.Time, timezone int8) (map[uint64]*common.CountryStats, error)
	GetTokenHeatmap(asset core.Token, from, to time.Time, timezone int8) (map[string]common.Heatmap, error)
}
