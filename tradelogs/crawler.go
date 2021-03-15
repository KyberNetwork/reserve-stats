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
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/tokenrate"
)

const (

	// tradeExecute(address sender, address src, uint256 srcAmount, address destToken, uint256 destAmount, address destAddress)
	tradeExecuteEvent = "0x4ee2afc3e9f9e97f558641bdc31ff31e4f34a1aaa2390cffbd64ee9ac18dfbec"
)

// var defaultTimeout = 10 * time.Second

// var errUnknownLogTopic = errors.New("unknown log topic")

type tradeLogFetcher func(*big.Int, *big.Int, time.Duration) (*common.CrawlResult, error)

// NewCrawler create a new Crawler instance.
func NewCrawler(sugar *zap.SugaredLogger,
	client *ethclient.Client,
	broadcastClient broadcast.Interface,
	rateProvider tokenrate.ETHUSDRateProvider,
	addresses []ethereum.Address,
	sb deployment.VersionedStartingBlocks,
	etherscanClient *etherscan.Client) (*Crawler, error) {
	resolver, err := blockchain.NewBlockTimeResolver(sugar, client)
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

// func (crawler *Crawler) getTransactionReceipt(txHash ethereum.Hash, timeout time.Duration) (*types.Receipt, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()
// 	receipt, err := crawler.ethClient.TransactionReceipt(ctx, txHash)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return receipt, nil
// }

// func (crawler *Crawler) updateBasicInfo(log types.Log, tradeLog common.TradelogV4, timeout time.Duration) (common.TradelogV4, error) {
// 	var txSender ethereum.Address
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()
// 	tx, _, err := crawler.ethClient.TransactionByHash(ctx, log.TxHash)
// 	if err != nil {
// 		return tradeLog, err
// 	}
// 	txSender, err = crawler.ethClient.TransactionSender(ctx, tx, log.BlockHash, log.TxIndex)
// 	tradeLog.TxDetail.TxSender = txSender
// 	tradeLog.TxDetail.GasPrice = tx.GasPrice()

// 	return tradeLog, err
// }

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
	for index, tradeLog := range result.Trades {
		// TODO: in case we want to get this information later
		// var uid, ip, country string
		// uid, ip, country, err = crawler.broadcastClient.GetTxInfo(tradeLog.TransactionHash.Hex())
		// if err != nil {
		// 	return result, err
		// }
		// result.Trades[index].User.IP = ip
		// result.Trades[index].User.Country = country
		// result.Trades[index].User.UID = uid

		rate, err := crawler.rateProvider.USDRate(tradeLog.Timestamp)
		if err != nil {
			return nil, err
		}
		result.Trades[index].ETHUSDProvider = crawler.rateProvider.Name()
		result.Trades[index].ETHUSDRate = rate
	}
	return result, nil
}
