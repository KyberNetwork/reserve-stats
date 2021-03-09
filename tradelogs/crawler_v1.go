package tradelogs

import (
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// const (
// use for crawler v4
// kyberTradeEventV4 = "0x30bbea603a7b36858fe5e3ec6ba5ff59dde039d02120d758eacfaed01520577d"
// )

func (crawler *Crawler) fetchTradeLog(fromBlock, toBlock *big.Int, timeout time.Duration) (*common.CrawlResult, error) {
	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(""), // should be TradeExecuteEvent from reserve contract
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

		// topic := log.Topics[0]
		// switch topic.Hex() {
		// default:
		// 	return nil, errUnknownLogTopic
		// }
	}

	return &result, err
}
