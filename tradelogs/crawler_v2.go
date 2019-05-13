package tradelogs

import (
	"context"
	"errors"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/mathutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	// KyberTrade (address srcAddress, address srcToken, uint256 srcAmount, address destAddress, address destToken, uint256 destAmount)
	kyberTradeEventV2 = "0x1c8399ecc5c956b9cb18c820248b10b634cca4af308755e07cd467655e8ec3c7"
)

func (crawler *Crawler) fetchTradeLogV2(fromBlock, toBlock *big.Int, timeout time.Duration) ([]common.TradeLog, error) {
	var result []common.TradeLog
	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(kyberTradeEventV2),
			ethereum.HexToHash(burnFeeEvent),
			ethereum.HexToHash(feeToWalletEvent),
			ethereum.HexToHash(etherReceivalEvent),
		},
	}

	typeLogs, err := crawler.fetchLogsWithTopics(fromBlock, toBlock, timeout, topics)
	if err != nil {
		return nil, err
	}

	result, err = crawler.assembleTradeLogsV2(typeLogs)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (crawler *Crawler) getTransactionReceipt(txHash ethereum.Hash, timeout time.Duration, logIndex uint) (ethereum.Address, error) {
	var (
		reserveAddr ethereum.Address
	)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	receipt, err := crawler.ethClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		return reserveAddr, err
	}
	for index := mathutil.MinUint64(uint64(len(receipt.Logs)-1), uint64(logIndex)); ; index-- {
		log := receipt.Logs[index]
		for _, topic := range log.Topics {
			if topic == ethereum.HexToHash(tradeExecuteEvent) {
				reserveAddr = log.Address
				break
			}
		}
		if !blockchain.IsZeroAddress(reserveAddr) {
			break
		}
	}
	return reserveAddr, nil
}

func (crawler *Crawler) assembleTradeLogsV2(eventLogs []types.Log) ([]common.TradeLog, error) {
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
		case kyberTradeEventV2:
			if tradeLog, err = fillKyberTradeV2(tradeLog, log); err != nil {
				return nil, err
			}
			if tradeLog.Timestamp, err = crawler.txTime.Resolve(log.BlockNumber); err != nil {
				return nil, err
			}

			tradeLog = assembleTradeLogsReserveAddr(tradeLog, crawler.sugar)

			// when the tradelog does not contain burnfee and etherReceival event
			// get tx receipt to get reserve address
			if len(tradeLog.BurnFees) == 0 && blockchain.IsZeroAddress(tradeLog.SrcReserveAddress) {
				crawler.sugar.Debug("trade logs has no burn fee, no ethReceival event, no wallet fee, getting reserve address from tx receipt")
				tradeLog.SrcReserveAddress, err = crawler.getTransactionReceipt(tradeLog.TransactionHash, defaultTimeout, log.Index)
				if err != nil {
					return nil, err
				}
			}
			// set tradeLog.EthAmount
			if tradeLog.SrcAddress == blockchain.ETHAddr {
				tradeLog.EthAmount = tradeLog.SrcAmount
			} else if tradeLog.DestAddress == blockchain.ETHAddr {
				tradeLog.EthAmount = tradeLog.DestAmount
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
