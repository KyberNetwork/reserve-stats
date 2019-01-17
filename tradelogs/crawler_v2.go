package tradelogs

import (
	"errors"
	"math/big"
	"time"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
