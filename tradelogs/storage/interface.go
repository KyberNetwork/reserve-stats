package storage

import (
	"fmt"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres"
)

const (
	// DbEngineFlag flag option
	DBEngineFlag    = "db-engine"
	defaultDBEngine = "influx"
	// InfluxDbEngine influxdb
	InfluxDBEngine = "influx"
	// PostgresDbEngine postgres db
	PostgresDBEngine = "postgres"

	// PostgresDefaultDb default db name when choosing Postgres
	PostgresDefaultDB = "reserve_stats"
)

// Interface represent a storage for TradeLogs data
type Interface interface {
	LoadTradeLogsByTxHash(tx ethereum.Hash) ([]common.TradeLog, error)
	LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error)
	GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddrs []ethereum.Address) (map[ethereum.Address]map[string]float64, error)
	GetAssetVolume(token ethereum.Address, fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error)
	GetReserveVolume(rsvAddr ethereum.Address, token ethereum.Address, fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error)
	GetMonthlyVolume(rsvAddr ethereum.Address, from, to time.Time) (map[uint64]*common.VolumeStats, error)
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
	GetTokenSymbol(address string) (string, error)
	UpdateTokens(tokenAddresses, symbols []string) error
	GetStats(from, to time.Time) (common.StatsResponse, error)
	GetTopTokens(from, to time.Time, limit uint64) (common.TopTokens, error)
	GetTopIntegrations(from, to time.Time, limit uint64) (common.TopIntegrations, error)
	GetTopReserves(from, to time.Time, limit uint64) (common.TopReserves, error)
	GetNotTwittedTrades(from, to time.Time) ([]common.BigTradeLog, error)
	SaveBigTrades(bigVolume float32, fromBlock uint64) error
	UpdateBigTradesTwitted(trades []uint64) error
}

// NewCliFlags return dbEngine flag option
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   DBEngineFlag,
			Usage:  "db engine to write trade logs, pls select influx or postgres",
			EnvVar: "DB_ENGINE",
			Value:  defaultDBEngine,
		},
	}
}

// KNCAddressFromContext return knc address by deployment mode
func KNCAddressFromContext(c *cli.Context) ethereum.Address {
	deploymentMode := deployment.MustGetDeploymentFromContext(c)
	switch deploymentMode {
	case deployment.Ropsten:
		return ethereum.HexToAddress("0x4e470dc7321e84ca96fcaedd0c8abcebbaeb68c6")
	default:
		return ethereum.HexToAddress("0xdd974d5c2e2928dea5f71b9825b8b646686bd200")
	}
}

// NewStorageInterfaceFromContext return new storage interface
func NewStorageInterfaceFromContext(sugar *zap.SugaredLogger, c *cli.Context, tokenAmountFormatter blockchain.TokenAmountFormatterInterface) (Interface, error) {
	dbEngine := c.String(DBEngineFlag)
	kncAddr := KNCAddressFromContext(c)
	switch dbEngine {
	case InfluxDBEngine:
		influxClient, err := influxdb.NewClientFromContext(c)
		if err != nil {
			return nil, err
		}

		influxStorage, err := influx.NewInfluxStorage(
			sugar,
			common.DatabaseName,
			influxClient,
			tokenAmountFormatter,
			kncAddr,
		)
		if err != nil {
			return nil, err
		}
		return influxStorage, nil
	case PostgresDBEngine:
		db, err := libapp.NewDBFromContext(c)
		if err != nil {
			return nil, err
		}
		postgresStorage, err := postgres.NewTradeLogDB(sugar, db, tokenAmountFormatter, kncAddr)
		if err != nil {
			sugar.Errorw("failed to initiate postgres storage", "error", err)
			return nil, err
		}
		return postgresStorage, nil
	default:
		return nil, fmt.Errorf("invalid db engine: %q", dbEngine)
	}
}
