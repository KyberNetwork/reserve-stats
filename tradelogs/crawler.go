package tradelogs

import (
	"context"
	"math/big"
	"time"

	ether "github.com/ethereum/go-ethereum"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nanmu42/etherscan-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/broadcast"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/tokenrate"
)

var defaultTimeout = 10 * time.Second

type tradeLogFetcher func(*big.Int, *big.Int, time.Duration) (*common.CrawlResult, error)

// NewCrawler create a new Crawler instance.
func NewCrawler(sugar *zap.SugaredLogger,
	client *ethclient.Client,
	broadcastClient broadcast.Interface,
	rateProvider tokenrate.ETHUSDRateProvider,
	addresses []ethereum.Address,
	sb deployment.VersionedStartingBlocks,
	reserveAddress ethereum.Address,
	etherscanClient *etherscan.Client) (*Crawler, error) {
	resolver, err := blockchain.NewBlockTimeResolver(sugar, client)
	if err != nil {
		return nil, err
	}

	reserveContract, err := contracts.NewReserve(reserveAddress, client)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		sugar:           sugar,
		ethClient:       client,
		txTime:          resolver,
		broadcastClient: broadcastClient,
		rateProvider:    rateProvider,
		addresses:       addresses,
		startingBlocks:  sb,
		reserveContract: reserveContract,
		etherscanClient: etherscanClient,
	}, nil
}

// Crawler gets trade logs on KyberNetwork on blockchain, adding the
// information about USD equivalent on each trade.
type Crawler struct {
	sugar           *zap.SugaredLogger
	ethClient       *ethclient.Client
	txTime          *blockchain.BlockTimeResolver
	broadcastClient broadcast.Interface
	rateProvider    tokenrate.ETHUSDRateProvider
	addresses       []ethereum.Address
	startingBlocks  deployment.VersionedStartingBlocks

	reserveContract *contracts.Reserve

	etherscanClient *etherscan.Client
}

func (crawler *Crawler) fetchLogsWithTopics(fromBlock, toBlock *big.Int, timeout time.Duration, topics [][]ethereum.Hash) ([]types.Log, error) {
	query := ether.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: crawler.addresses,
		Topics:    topics,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return crawler.ethClient.FilterLogs(ctx, query)

}

// GetTradeLogs returns trade logs from KyberNetwork.
func (crawler *Crawler) GetTradeLogs(fromBlock, toBlock *big.Int, timeout time.Duration) (*common.CrawlResult, error) {
	var (
		result  *common.CrawlResult
		fetchFn tradeLogFetcher
	)

	fetchFn = crawler.fetchTradeLog

	result, err := fetchFn(fromBlock, toBlock, timeout)
	if err != nil {
		return result, errors.Wrapf(err, "failed to fetch trade logs fromBlock: %v toBlock:%v", fromBlock, toBlock)
	}
	if result == nil {
		return result, nil
	}
	return result, nil
}
