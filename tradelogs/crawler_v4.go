package tradelogs

import (
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	// use for crawler v4
	addReserveToStorageEvent = "0x4649526e2876a69a4439244e5d8a32a6940a44a92b5390fdde1c22a26cc54004"

	reserveRebateWalletSetEvent = "0x42cac9e63e37f62d5689493d04887a67fe3c68e1d3763c3f0890e1620a0465b3"

	feeDistributedEvent = "0x53e2e1b5ab64e0a76fcc6a932558eba265d4e58c512401a7d776ae0f8fc08994"

	kyberTradeEventV4 = "0x30bbea603a7b36858fe5e3ec6ba5ff59dde039d02120d758eacfaed01520577d"
)

func init() {
	var err error
	networkABI, err = abi.JSON(strings.NewReader(contracts.NetworkProxyABI))
	if err != nil {
		panic(err)
	}
}
func (crawler *Crawler) fetchTradeLogV4(fromBlock, toBlock *big.Int, timeout time.Duration) ([]common.TradelogV4, error) {
	var result []common.TradelogV4

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

	result, err = crawler.assembleTradeLogsV4(typeLogs)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// AddReserveToStorage object
type AddReserveToStorage struct {
	Reserve      ethereum.Address `json:"reserve"`
	ReserveID    [32]byte         `json:"reserve_id"`
	RebateWallet ethereum.Address `json:"rebate_wallet"`
	BlockNumber  uint64           `json:"block_number"`
}

func (crawler *Crawler) fillAddReserveToStorage(tradelog common.TradelogV4, log types.Log) (common.TradelogV4, error) {
	reserve, err := crawler.kyberStorageContract.ParseAddReserveToStorage(log)
	if err != nil {
		return tradelog, err
	}

	if len(tradelog.Reserves) == 0 {
		tradelog.Reserves = []common.Reserve{}
	}
	tradelog.Reserves = append(tradelog.Reserves, common.Reserve{
		Address:      reserve.Reserve,
		ReserveID:    reserve.ReserveId,
		RebateWallet: reserve.RebateWallet,
		BlockNumber:  log.BlockNumber,
	})

	return tradelog, nil
}

// UpdateRebateWallet object
type UpdateRebateWallet struct {
	RebateWallet ethereum.Address `json:"rebate_wallet"`
	ReserveID    [32]byte         `json:"reserve_id"`
	BlockNumber  uint64           `json:"block_number"`
}

func (crawler *Crawler) fillRebateWalletSet(tradelog common.TradelogV4, log types.Log) (common.TradelogV4, error) {
	reserve, err := crawler.kyberStorageContract.ParseReserveRebateWalletSet(log)
	if err != nil {
		return tradelog, err
	}
	tradelog.Reserves = append(tradelog.Reserves, common.Reserve{
		ReserveID:    reserve.ReserveId,
		RebateWallet: reserve.RebateWallet,
		BlockNumber:  log.BlockNumber,
	})
	return tradelog, nil
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
			PlatformFee:    fee.PlatformFeeWei,
			PlatformWallet: fee.PlatformWallet,
			Burn:           fee.BurnAmtWei,
			Rebate:         fee.RebateWei,
			Reward:         fee.RewardWei,
		})
	return tradelog, nil
}

func (crawler *Crawler) assembleTradeLogsV4(eventLogs []types.Log) ([]common.TradelogV4, error) {
	var (
		logger   = crawler.sugar.With("func", caller.GetCurrentFunctionName())
		result   []common.TradelogV4
		tradeLog common.TradelogV4
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
		case feeDistributedEvent:
			if tradeLog, err = crawler.fillFeeDistributed(tradeLog, log); err != nil {
				return nil, err
			}
		case addReserveToStorageEvent:
			if tradeLog, err = crawler.fillAddReserveToStorage(tradeLog, log); err != nil {
				return result, err
			}
		case reserveRebateWalletSetEvent:
			if tradeLog, err = crawler.fillRebateWalletSet(tradeLog, log); err != nil {
				return result, err
			}
		case kyberTradeEventV4:
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
				return result, errors.Wrap(err, "could not update trade log basic info")
			}
			tradeLog.TxDetail.TransactionFee = big.NewInt(0).Mul(tradeLog.TxDetail.GasPrice, big.NewInt(int64(tradeLog.TxDetail.GasUsed)))
			crawler.sugar.Infow("gathered new trade log", "trade_log", tradeLog)
			// one trade only has one and only ExecuteTrade event
			result = append(result, tradeLog)
			tradeLog = common.TradelogV4{}
		default:
			return nil, errUnknownLogTopic
		}
	}

	return result, nil
}
