package tradelogs

import (
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
	// Transfer (index_topic_1 address from, index_topic_2 address to, uint256 value)
	transferEvent = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

	internalNetworkAddrV1 = "0x964F35fAe36d75B1e72770e244F6595B68508CF5"
)

func (crawler *Crawler) fetchTradeLogV1(fromBlock, toBlock *big.Int, timeout time.Duration) (*common.CrawlResult, error) {
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
		return nil, errors.Wrap(err, "failed to fetch log by topic")
	}

	return crawler.assembleTradeLogsV1(typeLogs)
}

func (crawler *Crawler) getTransactionReceiptV1(tradeLog common.TradelogV4, receipt *types.Receipt, logIndex uint,
	shouldGetReserveAddr bool) (common.TradelogV4, error) {

	for index := mathutil.MinInt64(int64(len(receipt.Logs)-1), int64(logIndex)); index >= 0; index-- {
		log := receipt.Logs[index]
		if len(log.Topics) == 0 {
			continue
		}
		topic := log.Topics[0].Hex()
		switch {
		case topic == transferEvent && blockchain.IsZeroAddress(tradeLog.ReceiverAddress):
			if len(log.Topics) != 3 {
				return tradeLog, errors.New("invalid transfer event topics data")
			}
			tradeLog.ReceiverAddress = ethereum.BytesToAddress(log.Topics[2].Bytes())
			if !shouldGetReserveAddr {
				break
			}
		case shouldGetReserveAddr && topic == tradeExecuteEvent:
			tradeLog.SrcReserveAddress = log.Address
			if !blockchain.IsZeroAddress(tradeLog.ReceiverAddress) {
				break
			}
		}
	}
	return tradeLog, nil
}

func (crawler *Crawler) getInternalTransaction(tradeLog common.TradelogV4) (common.TradelogV4, error) {
	blockInt := int(tradeLog.BlockNumber)
	internalTxs, err := crawler.etherscanClient.InternalTxByAddress(internalNetworkAddrV1, &blockInt, &blockInt, 0, 0, false)
	if err != nil {
		switch {
		case blockchain.IsEtherscanNotransactionFound(err):
			crawler.sugar.Warnw("internal transaction not found on etherscan", "err", err,
				"tx_hash", tradeLog.TransactionHash, "block_number", tradeLog.BlockNumber)
			return tradeLog, nil
		case blockchain.IsEtherscanRateLimit(err):
			crawler.sugar.Warnw("failed to get internal transaction", "err", err,
				"tx_hash", tradeLog.TransactionHash, "block_number", tradeLog.BlockNumber)
			return tradeLog, err
		default:
			return tradeLog, err
		}
	}
	for _, tx := range internalTxs {
		txHash := ethereum.HexToHash(tx.Hash)
		fromAddr := ethereum.HexToAddress(tx.From)
		if txHash == tradeLog.TransactionHash && fromAddr == ethereum.HexToAddress(internalNetworkAddrV1) {
			tradeLog.ReceiverAddress = ethereum.HexToAddress(tx.To)
			return tradeLog, nil
		}
	}
	return tradeLog, nil
}

func (crawler *Crawler) assembleTradeLogsV1(eventLogs []types.Log) (*common.CrawlResult, error) {
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
		case executeTradeEvent:
			tradeLog.Version = 1 // tradelog version 1
			if tradeLog, err = fillExecuteTrade(tradeLog, log); err != nil {
				return nil, err
			}
			if tradeLog.Timestamp, err = crawler.txTime.Resolve(log.BlockNumber); err != nil {
				return nil, errors.Wrapf(err, "failed to resolve timestamp by block_number %v", log.BlockNumber)
			}

			receipt, err := crawler.getTransactionReceipt(tradeLog.TransactionHash, defaultTimeout)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get transaction receipt tx: %v", tradeLog.TransactionHash)
			}
			tradeLog.TxDetail.GasUsed = receipt.GasUsed

			// when the tradelog does not contain burnfee and etherReceival event
			// or trasaction is not from token to ether
			// get tx receipt to get reserve address, receiver address
			shouldGetReserveAddr := common.LengthBurnFees(tradeLog) == 0 && len(tradeLog.T2EReserves)+len(tradeLog.E2TReserves) == 0
			if shouldGetReserveAddr || tradeLog.TokenInfo.DestAddress != blockchain.ETHAddr {
				crawler.sugar.Debug("trade logs  has no burn fee, no ethReceival event, no wallet fee, getting reserve address from tx receipt")
				tradeLog, err = crawler.getTransactionReceiptV1(tradeLog, receipt, log.Index, shouldGetReserveAddr)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get info from transaction receipt")
				}
			}

			if tradeLog.TokenInfo.DestAddress == blockchain.ETHAddr {
				// if transaction is from token to eth
				// get internal transaction to detect receiver address
				crawler.sugar.Debugw("get internal transaction", "tx_hash", tradeLog.TransactionHash.Hex(), "block_number", tradeLog.BlockNumber)
				tradeLog, err = crawler.getInternalTransaction(tradeLog)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get internal transaction")
				}
			}

			// set tradeLog.EthAmount
			if tradeLog.TokenInfo.SrcAddress == blockchain.ETHAddr {
				tradeLog.EthAmount = tradeLog.SrcAmount
			} else if tradeLog.TokenInfo.DestAddress == blockchain.ETHAddr {
				tradeLog.EthAmount = tradeLog.DestAmount
			}
			tradeLog.OriginalEthAmount = tradeLog.EthAmount // some case EthAmount
			// will be multiple so we keep OriginalEthAmount as a copy of original amount.

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
