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

	appname "github.com/KyberNetwork/reserve-stats/app-names"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/broadcast"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/tokenrate"
)

const (
	// executeTradeEvent is the topic of event
	// ExecuteTrade(address indexed sender, ERC20 src, ERC20 dest, uint actualSrcAmount, uint actualDestAmount).
	executeTradeEvent = "0x1849bd6a030a1bca28b83437fd3de96f3d27a5d172fa7e9c78e7b61468928a39"

	// tradeExecute(address sender, address src, uint256 srcAmount, address destToken, uint256 destAmount, address destAddress)
	// use for crawler v1 and v2
	tradeExecuteEvent = "0xea9415385bae08fe9f6dc457b02577166790cde83bb18cc340aac6cb81b824de"
)

var defaultTimeout = 10 * time.Second
var errUnknownLogTopic = errors.New("unknown log topic")

type tradeLogFetcher func(*big.Int, *big.Int, time.Duration) (*common.CrawlResult, error)

// NewCrawler create a new Crawler instance.
func NewCrawler(sugar *zap.SugaredLogger,
	client *ethclient.Client,
	broadcastClient broadcast.Interface,
	rateProvider tokenrate.ETHUSDRateProvider,
	addresses []ethereum.Address,
	sb deployment.VersionedStartingBlocks,
	etherscanClient *etherscan.Client,
	volumeExcludedReserves []ethereum.Address) (*Crawler, error) {
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
	sugar                 *zap.SugaredLogger
	ethClient             *ethclient.Client
	txTime                *blockchain.BlockTimeResolver
	broadcastClient       broadcast.Interface
	rateProvider          tokenrate.ETHUSDRateProvider
	addresses             []ethereum.Address
	startingBlocks        deployment.VersionedStartingBlocks
	volumeExludedReserves []ethereum.Address

	kyberStorageContract    *contracts.KyberStorage
	kyberFeeHandlerContract *contracts.KyberFeeHandler
	kyberNetworkContract    *contracts.KyberNetwork

	etherscanClient *etherscan.Client
	networkProxy    ethereum.Address
}

func logDataToExecuteTradeParams(data []byte) (ethereum.Address, ethereum.Address, ethereum.Hash, ethereum.Hash, error) {
	var srcAddr, desAddr ethereum.Address
	var srcAmount, desAmount ethereum.Hash

	if len(data) != 128 {
		err := errors.New("invalid trade data")
		return srcAddr, desAddr, srcAmount, desAmount, err
	}

	srcAddr = ethereum.BytesToAddress(data[0:32])
	desAddr = ethereum.BytesToAddress(data[32:64])
	srcAmount = ethereum.BytesToHash(data[64:96])
	desAmount = ethereum.BytesToHash(data[96:128])
	return srcAddr, desAddr, srcAmount, desAmount, nil
}

func fillExecuteTrade(tradeLog common.TradelogV4, logItem types.Log) (common.TradelogV4, error) {
	srcAddr, destAddr, srcAmount, destAmount, err := logDataToExecuteTradeParams(logItem.Data)
	if err != nil {
		return common.TradelogV4{}, err
	}
	tradeLog.TokenInfo.SrcAddress = srcAddr
	tradeLog.TokenInfo.DestAddress = destAddr
	tradeLog.SrcAmount = srcAmount.Big()
	tradeLog.DestAmount = destAmount.Big()

	tradeLog.TransactionHash = logItem.TxHash
	tradeLog.Index = logItem.Index
	tradeLog.User.UserAddress = ethereum.BytesToAddress(logItem.Topics[1].Bytes())
	tradeLog.BlockNumber = logItem.BlockNumber

	return tradeLog, nil
}

func isVolumeExcludedReserves(address ethereum.Address, volumeExcludedReserves []ethereum.Address) bool {
	for _, a := range volumeExcludedReserves {
		if address == a {
			return true
		}
	}
	return false
}

func (crawler *Crawler) calculateTradeAmount(t2ESrcAmounts, e2TSrcAmounts, t2ERates, e2TRates []*big.Int, srcToken, destToken ethereum.Address) (*big.Int, *big.Int) {
	srcAmount := big.NewInt(0)
	dstAmount := big.NewInt(0)
	if len(t2ESrcAmounts) != 0 {
		for _, amount := range t2ESrcAmounts {
			srcAmount = srcAmount.Add(srcAmount, amount)
		}
	} else {
		for _, amount := range e2TSrcAmounts {
			srcAmount = srcAmount.Add(srcAmount, amount)
		}
	}
	if len(e2TSrcAmounts) != 0 {
		for i, amount := range e2TSrcAmounts {
			dstAmount = dstAmount.Add(dstAmount, big.NewInt(0).Mul(amount, e2TRates[i]))
		}
	} else {
		for i, amount := range t2ESrcAmounts {
			dstAmount = dstAmount.Add(dstAmount, big.NewInt(0).Mul(amount, t2ERates[i]))
		}
	}
	return srcAmount, dstAmount
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

func (crawler *Crawler) getTransactionReceipt(txHash ethereum.Hash, timeout time.Duration) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	receipt, err := crawler.ethClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (crawler *Crawler) updateBasicInfo(log types.Log, tradeLog common.TradelogV4, timeout time.Duration) (common.TradelogV4, error) {
	var txSender ethereum.Address
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	tx, _, err := crawler.ethClient.TransactionByHash(ctx, log.TxHash)
	if err != nil {
		return tradeLog, err
	}
	txSender, err = crawler.ethClient.TransactionSender(ctx, tx, log.BlockHash, log.TxIndex)
	tradeLog.TxDetail.TxSender = txSender
	tradeLog.TxDetail.GasPrice = tx.GasPrice()

	return tradeLog, err
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
	for index, tradeLog := range result.Trades {
		var uid, ip, country string

		uid, ip, country, err = crawler.broadcastClient.GetTxInfo(tradeLog.TransactionHash.Hex())
		if err != nil {
			return result, err
		}
		result.Trades[index].User.IP = ip
		result.Trades[index].User.Country = country
		result.Trades[index].User.UID = uid

		if tradeLog.IsKyberSwap() {
			result.Trades[index].IntegrationApp = appname.KyberSwapAppName
		} else {
			result.Trades[index].IntegrationApp = appname.ThirdPartyAppName
		}

		rate, err := crawler.rateProvider.USDRate(tradeLog.Timestamp)
		if err != nil {
			return nil, err
		}
		result.Trades[index].ETHUSDProvider = crawler.rateProvider.Name()
		result.Trades[index].ETHUSDRate = rate
	}
	return result, nil
}
