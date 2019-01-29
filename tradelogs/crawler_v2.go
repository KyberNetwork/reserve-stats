package tradelogs

import (
	"context"
	"errors"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/mathutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	//TradeExecute(address sender, address src, uint256 srcAmount, address destToken, uint256 destAmount, address destAddress)
	tradeExecuteEvent = "0xea9415385bae08fe9f6dc457b02577166790cde83bb18cc340aac6cb81b824de"
)

func (crawler *Crawler) fetchTradeLogV2(fromBlock, toBlock *big.Int, timeout time.Duration) ([]common.TradeLog, error) {
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
	for index := mathutil.MinUint64(uint64(len(receipt.Logs)-1), uint64(logIndex)); index >= 0; index-- {
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

func assembleTradeLogsReserveAddr(log common.TradeLog, sugar *zap.SugaredLogger) common.TradeLog {
	if blockchain.IsBurnable(log.SrcAddress) {
		if blockchain.IsBurnable(log.DestAddress) {
			if len(log.BurnFees) == 2 {
				log.SrcReserveAddress = log.BurnFees[0].ReserveAddress
				log.DstReserveAddress = log.BurnFees[1].ReserveAddress
			} else {
				sugar.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "2 burn fees (src-dst)")
			}
		} else {
			if len(log.BurnFees) == 1 {
				log.SrcReserveAddress = log.BurnFees[0].ReserveAddress
			} else {
				sugar.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "1 burn fees (src)")
			}
		}
	} else if blockchain.IsBurnable(log.DestAddress) {
		if len(log.BurnFees) == 1 {
			log.DstReserveAddress = log.BurnFees[0].ReserveAddress
		} else {
			sugar.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "1 burn fees (dst)")
		}
	} else if len(log.WalletFees) != 0 {
		if len(log.WalletFees) == 1 {
			log.SrcReserveAddress = log.WalletFees[0].ReserveAddress
		} else {
			log.SrcReserveAddress = log.WalletFees[0].ReserveAddress
			log.DstReserveAddress = log.WalletFees[1].ReserveAddress
		}
	}
	return log
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
				return nil, err
			}

			tradeLog = assembleTradeLogsReserveAddr(tradeLog, crawler.sugar)

			// when the tradelog does not contain burnfee and etherReceival event
			// get tx receipt to get reserve address
			if len(tradeLog.BurnFees) == 0 && blockchain.IsZeroAddress(tradeLog.SrcReserveAddress) {
				crawler.sugar.Debug("Tradelog have no burnfee and ethReceival event, fallback to get reserve address from tx receipt")
				tradeLog.SrcReserveAddress, err = crawler.getTransactionReceipt(tradeLog.TransactionHash, 10*time.Second, log.Index)
				if err != nil {
					return nil, err
				}
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
