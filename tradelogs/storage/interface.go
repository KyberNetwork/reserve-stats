package storage

import (
	"fmt"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres"
)

const (
	DbEngineFlag     = "db-engine"
	defaultDbEngine  = "influx"
	InfluxDbEngine   = "influx"
	PostgresDbEngine = "postgres"

	PostgresDefaultDb = "reserve_stats"
)

// Interface represent a storage for TradeLogs data
type Interface interface {
	LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error)
	GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddrs []ethereum.Address) (map[ethereum.Address]map[string]float64, error)
	GetAssetVolume(token ethereum.Address, fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error)
	GetReserveVolume(rsvAddr ethereum.Address, token ethereum.Address, fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error)
	GetAggregatedWalletFee(reserveAddr, walletAddr, freq string,
		fromTime, toTime time.Time, timezone int8) (map[uint64]float64, error)
	GetTradeSummary(from, to time.Time, timezone int8) (map[uint64]*common.TradeSummary, error)
	GetUserVolume(userAddr ethereum.Address, from, to time.Time, freq string) (map[uint64]common.UserVolume, error)
	GetUserList(from, to time.Time) ([]common.UserInfo, error)
	GetWalletStats(fromTime, toTime time.Time, walletAddr string, timezone int8) (map[uint64]common.WalletStats, error)
	GetCountryStats(countryCode string, from, to time.Time, timezone int8) (map[uint64]*common.CountryStats, error)
	GetTokenHeatmap(asset ethereum.Address, from, to time.Time, timezone int8) (map[string]common.Heatmap, error)
	GetIntegrationVolume(fromTime, toTime time.Time) (map[uint64]*common.IntegrationVolume, error)
	LastBlock() (int64, error)
	SaveTradeLogs(logs []common.TradeLog) error
}

func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   DbEngineFlag,
			Usage:  "db engine to write trade logs, pls select influx or postgres",
			EnvVar: "DB_ENGINE",
			Value:  defaultDbEngine,
		},
	}
}

func NewStorageInterfaceFromContext(sugar *zap.SugaredLogger, c *cli.Context, tokenAmountFormatter blockchain.TokenAmountFormatterInterface) (Interface, error) {
	dbEngine := c.String(DbEngineFlag)
	switch dbEngine {
	case InfluxDbEngine:
		influxClient, err := influxdb.NewClientFromContext(c)
		if err != nil {
			return nil, err
		}

		influxStorage, err := influx.NewInfluxStorage(
			sugar,
			common.DatabaseName,
			influxClient,
			tokenAmountFormatter,
		)
		if err != nil {
			return nil, err
		}
		return influxStorage, nil
	case PostgresDbEngine:
		db, err := libapp.NewDBFromContext(c)
		if err != nil {
			return nil, err
		}
		postgresStorage, err := postgres.NewTradeLogDB(sugar, db, tokenAmountFormatter)
		if err != nil {
			sugar.Infow("error", err)
			return nil, err
		}
		return postgresStorage, nil
	default:
		return nil, fmt.Errorf("invalid db engine: %q", dbEngine)
	}
}
