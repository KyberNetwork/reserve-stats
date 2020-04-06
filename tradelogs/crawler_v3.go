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

var (
	networkABI abi.ABI
)

func init() {
	var err error
	networkABI, err = abi.JSON(strings.NewReader(contracts.NetworkProxyABI))
	if err != nil {
		panic(err)
	}
}
func (crawler *Crawler) fetchTradeLogV3(fromBlock, toBlock *big.Int, timeout time.Duration) ([]common.TradeLog, error) {
	var result []common.TradeLog

	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(burnFeeEvent),
			ethereum.HexToHash(feeToWalletEvent),
			ethereum.HexToHash(kyberTradeEvent),
		},
	}

	typeLogs, err := crawler.fetchLogsWithTopics(fromBlock, toBlock, timeout, topics)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch log by topic")
	}

	result, err = crawler.assembleTradeLogsV3(typeLogs)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type tradeWithHintParam struct {
	Src               ethereum.Address
	SrcAmount         *big.Int
	Dest              ethereum.Address
	DestAddress       ethereum.Address
	MaxDestAmount     *big.Int
	MinConversionRate *big.Int
	WalletID          ethereum.Address `abi:"walletId"`
	Hint              []byte
}

func decodeTradeWithHintParam(data []byte) (tradeWithHintParam, error) { // decode txInput method signature
	if len(data) < 4 {
		return tradeWithHintParam{}, errors.New("input data not valid")
	}
	// recover Method from signature and ABI
	method, err := networkABI.MethodById(data[0:4])
	if err != nil {
		return tradeWithHintParam{}, errors.Wrap(err, "cannot find method for correspond data")
	}
	var resp tradeWithHintParam
	if method.Name != "tradeWithHint" {
		return tradeWithHintParam{}, errors.New("try to decode data is not from tradeWithHint")
	}
	// unpack method inputs
	err = method.Inputs.Unpack(&resp, data[4:])
	if err != nil {
		return tradeWithHintParam{}, errors.Wrap(err, "unpack tradeWithHint param failed")
	}
	return resp, nil
}
func (crawler *Crawler) assembleTradeLogsV3(eventLogs []types.Log) ([]common.TradeLog, error) {
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
				return result, errors.New("could not update trade log basic info")
			}
			tradeLog.TransactionFee = big.NewInt(0).Mul(tradeLog.GasPrice, big.NewInt(int64(tradeLog.GasUsed)))
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
