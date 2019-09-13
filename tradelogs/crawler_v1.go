package tradelogs

import (
	"context"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/mathutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	// Transfer (index_topic_1 address from, index_topic_2 address to, uint256 value)
	transferEvent = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

	internalNetworkAddrV1 = "0x964F35fAe36d75B1e72770e244F6595B68508CF5"
)

func (crawler *Crawler) fetchTradeLogV1(fromBlock, toBlock *big.Int, timeout time.Duration) ([]common.TradeLog, error) {
	var result []common.TradeLog
	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(executeTradeEvent),
			ethereum.HexToHash(burnFeeEvent),
			ethereum.HexToHash(feeToWalletEvent),
			ethereum.HexToHash(etherReceivalEvent),
		},
	}

	typeLogs, err := crawler.fetchLogsWithTopics(fromBlock, toBlock, timeout, topics)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch log by topic")
	}

	result, err = crawler.assembleTradeLogsV1(typeLogs)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (crawler *Crawler) getTransactionReceiptV1(tradeLog common.TradeLog, timeout time.Duration, logIndex uint,
	shouldGetReserveAddr bool) (common.TradeLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	receipt, err := crawler.ethClient.TransactionReceipt(ctx, tradeLog.TransactionHash)
	if err != nil {
		return tradeLog, err
	}

	for index := mathutil.MinInt64(int64(len(receipt.Logs)-1), int64(logIndex)); index >= 0; index-- {
		log := receipt.Logs[index]
		if len(log.Topics) == 0 {
			continue
		}
		topic := log.Topics[0].Hex()
		switch {
		case topic == transferEvent && blockchain.IsZeroAddress(tradeLog.ReceiverAddress):
			if len(log.Topics) != 3 {
				return tradeLog, errors.New("invalid transfer event topics data")
			}
			tradeLog.ReceiverAddress = ethereum.BytesToAddress(log.Topics[2].Bytes())
			if !shouldGetReserveAddr {
				break
			}
		case shouldGetReserveAddr && topic == tradeExecuteEvent:
			tradeLog.SrcReserveAddress = log.Address
			if !blockchain.IsZeroAddress(tradeLog.ReceiverAddress) {
				break
			}
		}
	}
	return tradeLog, nil
}

func (crawler *Crawler) getInternalTransaction(tradeLog common.TradeLog) (common.TradeLog, error) {
	blockInt := int(tradeLog.BlockNumber)
	internalTxs, err := crawler.etherscanClient.InternalTxByAddress(internalNetworkAddrV1, &blockInt, &blockInt, 0, 0, false)
	if err != nil {
		switch {
		case blockchain.IsEtherscanNotransactionFound(err):
			crawler.sugar.Warnw("internal transaction not found on etherscan", "err", err,
				"tx_hash", tradeLog.TransactionHash, "block_number", tradeLog.BlockNumber)
			return tradeLog, nil
		case blockchain.IsEtherscanRateLimit(err):
			crawler.sugar.Warnw("failed to get internal transaction", "err", err,
				"tx_hash", tradeLog.TransactionHash, "block_number", tradeLog.BlockNumber)
			return tradeLog, err
		default:
			return tradeLog, err
		}
	}
	for _, tx := range internalTxs {
		txHash := ethereum.HexToHash(tx.Hash)
		fromAddr := ethereum.HexToAddress(tx.From)
		if txHash == tradeLog.TransactionHash && fromAddr == ethereum.HexToAddress(internalNetworkAddrV1) {
			tradeLog.ReceiverAddress = ethereum.HexToAddress(tx.To)
			return tradeLog, nil
		}
	}
	return tradeLog, nil
}

func (crawler *Crawler) assembleTradeLogsV1(eventLogs []types.Log) ([]common.TradeLog, error) {
	var (
		result   []common.TradeLog
		tradeLog common.TradeLog
		err      error
	)

	for _, log := range eventLogs {
		if log.Removed {
			continue // Removed due to chain reorg
		}

		if len(log.Topics) == 0 {
			return result, errors.New("log item has no topic")
		}

		tradeLog.TxSender, err = crawler.getTxSender(log, defaultTimeout)
		if err != nil {
			return result, errors.New("could not get trade log sender")
		}

		topic := log.Topics[0]
		switch topic.Hex() {
		case feeToWalletEvent:
			if tradeLog, err = fillWalletFees(tradeLog, log); err != nil {
				return nil, err
			}
		case burnFeeEvent:
			if tradeLog, err = fillBurnFees(tradeLog, log); err != nil {
				return nil, err
			}
		case etherReceivalEvent:
			if tradeLog, err = fillEtherReceival(tradeLog, log); err != nil {
				return nil, err
			}
		case executeTradeEvent:
			if tradeLog, err = fillExecuteTrade(tradeLog, log); err != nil {
				return nil, err
			}
			if tradeLog.Timestamp, err = crawler.txTime.Resolve(log.BlockNumber); err != nil {
				return nil, errors.Wrapf(err, "failed to resolve timestamp by block_number %v", log.BlockNumber)
			}

			tradeLog = assembleTradeLogsReserveAddr(tradeLog, crawler.sugar)

			// when the tradelog does not contain burnfee and etherReceival event
			// or trasaction is not from token to ether
			// get tx receipt to get reserve address, receiver address
			shouldGetReserveAddr := len(tradeLog.BurnFees) == 0 && blockchain.IsZeroAddress(tradeLog.SrcReserveAddress)
			if shouldGetReserveAddr || tradeLog.DestAddress != blockchain.ETHAddr {
				crawler.sugar.Debug("trade logs  has no burn fee, no ethReceival event, no wallet fee, getting reserve address from tx receipt")
				tradeLog, err = crawler.getTransactionReceiptV1(tradeLog, 10*time.Second, log.Index, shouldGetReserveAddr)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get info from transaction receipt")
				}
			}

			if tradeLog.DestAddress == blockchain.ETHAddr {
				// if transaction is from token to eth
				// get internal transaction to detect receiver address
				crawler.sugar.Debugw("get internal transaction", "tx_hash", tradeLog.TransactionHash.Hex(), "block_number", tradeLog.BlockNumber)
				tradeLog, err = crawler.getInternalTransaction(tradeLog)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get internal transaction")
				}
			}

			// set tradeLog.EthAmount
			if tradeLog.SrcAddress == blockchain.ETHAddr {
				tradeLog.EthAmount = tradeLog.SrcAmount
				tradeLog.TradeVolume = tradeLog.SrcAmount
			} else if tradeLog.DestAddress == blockchain.ETHAddr {
				tradeLog.EthAmount = tradeLog.DestAmount
				tradeLog.TradeVolume = tradeLog.DestAmount
			}
			crawler.sugar.Infow("gathered new trade log", "trade_log", tradeLog)
			// one trade only has one and only ExecuteTrade event
			result = append(result, tradeLog)
			tradeLog = common.TradeLog{}
		default:
			return nil, errUnknownLogTopic
		}
	}

	return result, nil
}
