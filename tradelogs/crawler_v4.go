package tradelogs

import (
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	// use for crawler v4
	addReserveToStorageEvent = "0x4649526e2876a69a4439244e5d8a32a6940a44a92b5390fdde1c22a26cc54004"

	// use for crawler v4
	reserveRebateWalletSetEvent = "0x42cac9e63e37f62d5689493d04887a67fe3c68e1d3763c3f0890e1620a0465b3"

	//
	feeDistributedEvent = ""
)

func init() {
	var err error
	networkABI, err = abi.JSON(strings.NewReader(contracts.NetworkProxyABI))
	if err != nil {
		panic(err)
	}
}
func (crawler *Crawler) fetchTradeLogV4(fromBlock, toBlock *big.Int, timeout time.Duration) ([]common.TradeLog, error) {
	var result []common.TradeLog

	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(feeDistributedEvent),
			ethereum.HexToHash(kyberTradeEvent),
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

// type tradeWithHintV4Param struct {
// 	Src               ethereum.Address
// 	SrcAmount         *big.Int
// 	Dest              ethereum.Address
// 	DestAddress       ethereum.Address
// 	MaxDestAmount     *big.Int
// 	MinConversionRate *big.Int
// 	WalletID          ethereum.Address `abi:"walletId"`
// 	Hint              []byte
// }

func fillAddReserveToStorage(log types.Log) error {
	// TODO: get
	return nil
	//
}

func (crawler *Crawler) assembleTradeLogsV4(eventLogs []types.Log) ([]common.TradeLog, error) {
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
		case kyberTradeEvent:
			if tradeLog, err = fillKyberTradeV3(tradeLog, log, crawler.volumeExludedReserves); err != nil {
				return nil, err
			}
			receipt, err := crawler.getTransactionReceipt(tradeLog.TransactionHash, defaultTimeout)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get transaction receipt tx: %v", tradeLog.TransactionHash)
			}
			tradeLog.GasUsed = receipt.GasUsed
			if tradeLog.Timestamp, err = crawler.txTime.Resolve(log.BlockNumber); err != nil {
				return nil, errors.Wrapf(err, "failed to resolve timestamp by block_number %v", log.BlockNumber)
			}
			tradeLog, err = crawler.updateBasicInfo(log, tradeLog, defaultTimeout)
			if err != nil {
				return result, errors.Wrap(err, "could not update trade log basic info")
			}
			tradeLog.TransactionFee = big.NewInt(0).Mul(tradeLog.GasPrice, big.NewInt(int64(tradeLog.GasUsed)))
			crawler.sugar.Infow("gathered new trade log", "trade_log", tradeLog)
			// one trade only has one and only ExecuteTrade event
			result = append(result, tradeLog)
			tradeLog = common.TradeLog{}
		case addReserveToStorageEvent:
			if err := fillAddReserveToStorage(log); err != nil {
				return result, err
			}
		case reserveRebateWalletSetEvent:
		default:
			return nil, errUnknownLogTopic
		}
	}

	return result, nil
}
