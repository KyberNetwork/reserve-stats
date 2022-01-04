package storage

import (
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres"
)

const (
	// PostgresDBEngine postgres db
	PostgresDBEngine = "postgres"

	// PostgresDefaultDB default db name when choosing Postgres
	PostgresDefaultDB = "reserve_stats"
)

// Interface represent a storage for TradeLogs data
type Interface interface {
	LoadTradeLogsByTxHash(tx ethereum.Hash) ([]common.TradelogV4, error)
	LoadTradeLogs(from, to time.Time) ([]common.TradelogV4, error)
	LastBlock() (int64, error)
	SaveTradeLogs(log *common.CrawlResult) error
	GetTokenSymbol(address string) (string, error)
	UpdateTokens(tokenAddresses, symbols []string) error
	GetStats(from, to time.Time) (common.StatsResponse, error)
	GetTopTokens(from, to time.Time, limit uint64) (common.TopTokens, error)
	GetTopIntegrations(from, to time.Time, limit uint64) (common.TopIntegrations, error)
	GetTopReserves(from, to time.Time, limit uint64) (common.TopReserves, error)
	GetNotTwittedTrades(from, to time.Time) ([]common.BigTradeLog, error)
	SaveBigTrades(bigVolume float32, fromBlock uint64) error
	UpdateBigTradesTwitted(trades []uint64) error
	GetTokenInfo() ([]common.TokenInfo, error)
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
	kncAddr := KNCAddressFromContext(c)
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
}
