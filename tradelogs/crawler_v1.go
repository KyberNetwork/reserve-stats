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

const (
	// use for crawler v4
	addReserveToStorageEvent = "0x50b2ce9e8f1a63ceaed262cc854dbf741b216e6429f7ba38403afbcdddc7f1ea"

	reserveRebateWalletSetEvent = "0x42cac9e63e37f62d5689493d04887a67fe3c68e1d3763c3f0890e1620a0465b3"

	feeDistributedEvent = "0x53e2e1b5ab64e0a76fcc6a932558eba265d4e58c512401a7d776ae0f8fc08994"

	kyberTradeEventV4 = "0x30bbea603a7b36858fe5e3ec6ba5ff59dde039d02120d758eacfaed01520577d"
)

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

func (crawler *Crawler) fillAddReserveToStorage(crResult *common.CrawlResult, log types.Log) error {
	reserve, err := crawler.kyberStorageContract.ParseAddReserveToStorage(log)
	if err != nil {
		return err
	}

	if len(crResult.Reserves) == 0 {
		crResult.Reserves = []common.Reserve{}
	}
	crResult.Reserves = append(crResult.Reserves, common.Reserve{
		Address:      reserve.Reserve,
		ReserveID:    reserve.ReserveId,
		ReserveType:  uint64(reserve.ReserveType),
		RebateWallet: reserve.RebateWallet,
		BlockNumber:  log.BlockNumber,
	})

	return nil
}

func (crawler *Crawler) fillRebateWalletSet(crResult *common.CrawlResult, log types.Log) error {
	reserve, err := crawler.kyberStorageContract.ParseReserveRebateWalletSet(log)
	if err != nil {
		return err
	}
	crResult.UpdateWallets = append(crResult.UpdateWallets, common.Reserve{
		ReserveID:    reserve.ReserveId,
		RebateWallet: reserve.RebateWallet,
		BlockNumber:  log.BlockNumber,
	})
	return nil
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
		default:
			return nil, errUnknownLogTopic
		}
	}

	return &result, err
}
