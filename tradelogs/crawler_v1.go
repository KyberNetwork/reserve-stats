package tradelogs

import (
	"fmt"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

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
		logger = crawler.sugar.With("func", caller.GetCurrentFunctionName())
		result common.CrawlResult
		err    error
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
			// TODO
			logger.Info("handle trade execute event")
		default:
			return nil, fmt.Errorf("unknown topic")
		}
	}

	return &result, err
}
