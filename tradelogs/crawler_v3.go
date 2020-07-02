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
func (crawler *Crawler) fetchTradeLogV3(fromBlock, toBlock *big.Int, timeout time.Duration) (*common.CrawlResult, error) {
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

	return crawler.assembleTradeLogsV3(typeLogs)
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

func decodeTradeInputParam(data []byte) (out tradeWithHintParam, err error) { // decode txInput method signature
	if len(data) < 4 {
		return tradeWithHintParam{}, errors.New("input data not valid")
	}
	// recover Method from signature and ABI
	method, err := networkABI.MethodById(data[0:4])
	if err != nil {
		return tradeWithHintParam{}, errors.Wrap(err, "cannot find method for correspond data")
	}
	switch method.Name {
	case "trade", "tradeWithHint":
		// unpack method inputs
		var out tradeWithHintParam
		err = method.Inputs.Unpack(&out, data[4:])
		if err != nil {
			return tradeWithHintParam{}, errors.Wrap(err, "unpack param failed")
		}
		return out, nil
	case "swapTokenToToken", "swapTokenToEther", "swapEtherToToken":
		// no wallet this trade, just return empty
		return tradeWithHintParam{}, nil
	default:
		return tradeWithHintParam{}, errors.Errorf("unexpected method %s", method.Name)
	}
}
func (crawler *Crawler) assembleTradeLogsV3(eventLogs []types.Log) (*common.CrawlResult, error) {
	var (
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
		case feeToWalletEvent:
			if tradeLog, err = fillWalletFees(tradeLog, log); err != nil {
				return nil, err
			}
		case burnFeeEvent:
			if tradeLog, err = fillBurnFees(tradeLog, log); err != nil {
				return nil, err
			}
		case kyberTradeEvent:
			tradeLog.Version = 3 // tradelog version 3
			if tradeLog, err = fillKyberTradeV3(tradeLog, log, crawler.volumeExludedReserves); err != nil {
				return nil, err
			}
			receipt, err := crawler.getTransactionReceipt(tradeLog.TransactionHash, defaultTimeout)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get transaction receipt tx: %v", tradeLog.TransactionHash)
			}
			tradeLog.TxDetail.GasUsed = receipt.GasUsed
			if tradeLog.Timestamp, err = crawler.txTime.Resolve(log.BlockNumber); err != nil {
				return nil, errors.Wrapf(err, "failed to resolve timestamp by block_number %v", log.BlockNumber)
			}
			tradeLog, err = crawler.updateBasicInfo(log, tradeLog, defaultTimeout)
			if err != nil {
				return &result, errors.Wrap(err, "could not update trade log basic info")
			}
			tradeLog.TxDetail.TransactionFee = big.NewInt(0).Mul(tradeLog.TxDetail.GasPrice, big.NewInt(int64(tradeLog.TxDetail.GasUsed)))
			crawler.sugar.Infow("gathered new trade log", "trade_log", tradeLog)
			result.Trades = append(result.Trades, tradeLog)
			tradeLog = common.TradelogV4{}
		default:
			return nil, errUnknownLogTopic
		}
	}

	return &result, nil
}
