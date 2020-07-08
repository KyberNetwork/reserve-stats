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

func (crawler *Crawler) fetchTradeLogV4(fromBlock, toBlock *big.Int, timeout time.Duration) (*common.CrawlResult, error) {
	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(feeDistributedEvent),
			ethereum.HexToHash(kyberTradeEventV4),
			ethereum.HexToHash(addReserveToStorageEvent),
			ethereum.HexToHash(reserveRebateWalletSetEvent),
		},
	}

	typeLogs, err := crawler.fetchLogsWithTopics(fromBlock, toBlock, timeout, topics)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch log by topic")
	}

	return crawler.assembleTradeLogsV4(typeLogs)
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

func (crawler *Crawler) fillFeeDistributed(tradelog common.TradelogV4, log types.Log) (common.TradelogV4, error) {
	var (
		logger = crawler.sugar.With("func", caller.GetCurrentFunctionName())
	)
	fee, err := crawler.kyberFeeHandlerContract.ParseFeeDistributed(log)
	if err != nil {
		logger.Errorw("failed to parse fee distributed event", "error", err)
		return tradelog, err
	}

	tradelog.Fees = append(tradelog.Fees,
		common.TradelogFee{
			PlatformFee:               fee.PlatformFeeWei,
			PlatformWallet:            fee.PlatformWallet,
			Burn:                      fee.BurnAmtWei,
			Rebate:                    fee.RebateWei,
			Reward:                    fee.RewardWei,
			RebateWallets:             fee.RebateWallets,
			RebatePercentBpsPerWallet: fee.RebatePercentBpsPerWallet,
			Index:                     log.Index,
		})
	tradelog.WalletAddress = fee.PlatformWallet
	tradelog.WalletName = WalletAddrToName(fee.PlatformWallet)
	return tradelog, nil
}

func (crawler *Crawler) assembleTradeLogsV4(eventLogs []types.Log) (*common.CrawlResult, error) {
	var (
		logger   = crawler.sugar.With("func", caller.GetCurrentFunctionName())
		result   common.CrawlResult
		tradeLog common.TradelogV4
		err      error
	)

	for _, log := range eventLogs {
		if log.Removed {
			continue // Removed due to chain reorg
		}

		if len(log.Topics) == 0 {
			return &result, errors.New("log item has no topic")
		}

		topic := log.Topics[0]
		switch topic.Hex() {
		case addReserveToStorageEvent:
			if err = crawler.fillAddReserveToStorage(&result, log); err != nil {
				return &result, err
			}
		case reserveRebateWalletSetEvent:
			if err = crawler.fillRebateWalletSet(&result, log); err != nil {
				return &result, err
			}
		case feeDistributedEvent:
			if tradeLog, err = crawler.fillFeeDistributed(tradeLog, log); err != nil {
				return nil, err
			}
		case kyberTradeEventV4:
			tradeLog.Version = 4 // tradelog version 4
			if tradeLog, err = crawler.fillKyberTradeV4(tradeLog, log, crawler.volumeExludedReserves); err != nil {
				return nil, err
			}
			receipt, err := crawler.getTransactionReceipt(tradeLog.TransactionHash, defaultTimeout)
			if err != nil {
				logger.Errorw("failed to get transaction receipt", "error", err)
				return nil, errors.Wrapf(err, "failed to get transaction receipt tx: %v", tradeLog.TransactionHash)
			}
			tradeLog.TxDetail.GasUsed = receipt.GasUsed
			if tradeLog.Timestamp, err = crawler.txTime.Resolve(log.BlockNumber); err != nil {
				return nil, errors.Wrapf(err, "failed to resolve timestamp by block_number %v", log.BlockNumber)
			}
			tradeLog, err = crawler.updateBasicInfoV4(log, tradeLog, defaultTimeout)
			if err != nil {
				return &result, errors.Wrap(err, "could not update trade log basic info")
			}
			tradeLog.TxDetail.TransactionFee = big.NewInt(0).Mul(tradeLog.TxDetail.GasPrice, big.NewInt(int64(tradeLog.TxDetail.GasUsed)))
			crawler.sugar.Infow("gathered new trade log", "trade_log", tradeLog)
			// one trade only has one and only ExecuteTrade event
			result.Trades = append(result.Trades, tradeLog)
			tradeLog = common.TradelogV4{}
		default:
			return nil, errUnknownLogTopic
		}
	}

	return &result, nil
}
