package tradelogs

import (
	"context"
	"fmt"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	// tradeExecute(address sender, address src, uint256 srcAmount, address destToken, uint256 destAmount, address destAddress)
	tradeExecuteEvent = "0x4ee2afc3e9f9e97f558641bdc31ff31e4f34a1aaa2390cffbd64ee9ac18dfbec"
)

func (crawler *Crawler) fetchTradeLog(fromBlock, toBlock *big.Int, timeout time.Duration) (*common.CrawlResult, error) {
	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(tradeExecuteEvent), // should be TradeExecuteEvent from reserve contract
		},
	}

	typeLogs, err := crawler.fetchLogsWithTopics(fromBlock, toBlock, timeout, topics)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch log by topic")
	}

	return crawler.assembleTradeLogs(typeLogs)
}

func (crawler *Crawler) assembleTradeLogs(eventLogs []types.Log) (*common.CrawlResult, error) {
	var (
		logger   = crawler.sugar.With("func", caller.GetCurrentFunctionName())
		result   common.CrawlResult
		tradeLog common.Tradelog
		err      error
	)
	logger.Info("assemble tradelogs")
	for _, log := range eventLogs {
		if log.Removed {
			continue // Removed due to chain reorg
		}

		if len(log.Topics) == 0 {
			return &result, errors.New("log item has no topic")
		}

		topic := log.Topics[0]
		switch topic.Hex() {
		case tradeExecuteEvent:
			tradeLog.Version = 1

			if tradeLog.Timestamp, err = crawler.txTime.Resolve(log.BlockNumber); err != nil {
				return nil, errors.Wrapf(err, "failed to resolve timestamp by block_number %v", log.BlockNumber)
			}

			if tradeLog, err = crawler.fillExecuteTrade(tradeLog, log); err != nil {
				return nil, errors.Wrapf(err, "failed to fill execute trade tx: %v", tradeLog.TransactionHash)
			}

			receipt, err := crawler.getTransactionReceipt(tradeLog.TransactionHash, defaultTimeout)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get transaction receipt tx: %v", tradeLog.TransactionHash)
			}
			tradeLog.TxDetail.GasUsed = receipt.GasUsed

			// set tradeLog.EthAmount
			if tradeLog.TokenInfo.SrcAddress == blockchain.USDTAddr {
				tradeLog.EthAmount = tradeLog.SrcAmount
			} else if tradeLog.TokenInfo.DestAddress == blockchain.USDTAddr {
				tradeLog.EthAmount = tradeLog.DestAmount
			}
			tradeLog.OriginalEthAmount = tradeLog.EthAmount // some case EthAmount

			tradeLog, err = crawler.updateBasicInfo(log, tradeLog, defaultTimeout)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to update basic info tx: %v", tradeLog.TransactionHash)
			}
			tradeLog.TxDetail.TransactionFee = big.NewInt(0).Mul(tradeLog.TxDetail.GasPrice, big.NewInt(int64(tradeLog.TxDetail.GasUsed)))

			crawler.sugar.Infow("gathered new trade log", "trade_log", tradeLog)

			result.Trades = append(result.Trades, tradeLog)
			tradeLog = common.Tradelog{}

		default:
			return nil, fmt.Errorf("unknown topic")
		}
	}

	return &result, err
}

func (crawler *Crawler) updateBasicInfo(log types.Log, tradeLog common.Tradelog, timeout time.Duration) (common.Tradelog, error) {
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
	tradeLog.User.UserAddress = txSender

	return tradeLog, err
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

func (crawler *Crawler) fillExecuteTrade(tradeLog common.Tradelog, logItem types.Log) (common.Tradelog, error) {

	executeTrade, err := crawler.reserveContract.ParseTradeExecute(logItem)
	if err != nil {
		return tradeLog, err
	}

	tradeLog.TokenInfo.SrcAddress = executeTrade.Src
	tradeLog.TokenInfo.DestAddress = executeTrade.DestToken
	tradeLog.SrcAmount = executeTrade.SrcAmount
	tradeLog.DestAmount = executeTrade.DestAmount

	tradeLog.TransactionHash = logItem.TxHash
	tradeLog.Index = logItem.Index
	tradeLog.BlockNumber = logItem.BlockNumber
	return tradeLog, nil
}
