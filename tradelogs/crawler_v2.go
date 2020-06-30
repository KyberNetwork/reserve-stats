package tradelogs

import (
	"context"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/mathutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	// KyberTrade (address srcAddress, address srcToken, uint256 srcAmount, address destAddress, address destToken, uint256 destAmount)
	kyberTradeEventV2 = "0x1c8399ecc5c956b9cb18c820248b10b634cca4af308755e07cd467655e8ec3c7"
)

func (crawler *Crawler) fetchTradeLogV2(fromBlock, toBlock *big.Int, timeout time.Duration) (*common.CrawlResult, error) {
	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(kyberTradeEventV2),
			ethereum.HexToHash(burnFeeEvent),
			ethereum.HexToHash(feeToWalletEvent),
			ethereum.HexToHash(etherReceivalEvent),
		},
	}

	typeLogs, err := crawler.fetchLogsWithTopics(fromBlock, toBlock, timeout, topics)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch log by topic")
	}

	return crawler.assembleTradeLogsV2(typeLogs)
}

func (crawler *Crawler) getTransactionReceipt(txHash ethereum.Hash, timeout time.Duration) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	receipt, err := crawler.ethClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func getReserveFromReceipt(receipt *types.Receipt, logIndex uint) ethereum.Address {
	var (
		reserveAddr ethereum.Address
	)
	for index := mathutil.MinUint64(uint64(len(receipt.Logs)-1), uint64(logIndex)); ; index-- {
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
	return reserveAddr
}

func (crawler *Crawler) assembleTradeLogsV2(eventLogs []types.Log) (*common.CrawlResult, error) {
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
		case etherReceivalEvent:
			if tradeLog, err = fillEtherReceival(tradeLog, log); err != nil {
				return nil, err
			}
		case kyberTradeEventV2:
			if tradeLog, err = fillKyberTradeV2(tradeLog, log); err != nil {
				return nil, err
			}
			if tradeLog.Timestamp, err = crawler.txTime.Resolve(log.BlockNumber); err != nil {
				return nil, errors.Wrapf(err, "failed to resolve timestamp by block_number %v", log.BlockNumber)
			}

			// tradeLog = assembleTradeLogsReserveAddr(tradeLog, crawler.sugar)

			receipt, err := crawler.getTransactionReceipt(tradeLog.TransactionHash, defaultTimeout)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get transaction receipt tx: %v", tradeLog.TransactionHash)
			}
			tradeLog.TxDetail.GasUsed = receipt.GasUsed

			// when the tradelog does not contain burnfee and etherReceival event
			// get tx receipt to get reserve address
			if common.LengthBurnFees(tradeLog) == 0 && len(tradeLog.T2EReserves)+len(tradeLog.E2TReserves) == 0 {
				crawler.sugar.Debug("trade logs has no burn fee, no ethReceival event, no wallet fee, getting reserve address from tx receipt")
				srcReserveAddress := getReserveFromReceipt(receipt, log.Index)
				tradeLog.SrcReserveAddress = srcReserveAddress
				// tradeLog.T2EReserves = append(tradeLog.T2EReserves, srcReserveAddress)
			}
			// set tradeLog.EthAmount
			if tradeLog.TokenInfo.SrcAddress == blockchain.ETHAddr {
				tradeLog.EthAmount = tradeLog.SrcAmount
			} else if tradeLog.TokenInfo.DestAddress == blockchain.ETHAddr {
				tradeLog.EthAmount = tradeLog.DestAmount
			}
			// keep OriginalEthAmount as origin amount of EthAmount
			tradeLog.OriginalEthAmount = big.NewInt(0).Set(tradeLog.EthAmount)
			tradeLog.EthAmount = tradeLog.EthAmount.Mul(tradeLog.EthAmount, big.NewInt(int64(common.LengthBurnFees(tradeLog))))

			tradeLog, err = crawler.updateBasicInfo(log, tradeLog, defaultTimeout)
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
