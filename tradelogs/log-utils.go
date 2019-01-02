package tradelogs

import (
	"errors"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (crawler *Crawler) handleEventLogWithBlockNumber(log types.Log, tradeLog *common.TradeLog) (bool, error) {
	//this should not happen, but still...
	if len(log.Topics) == 0 {
		return false, errors.New("log item has no topic")
	}
	//if blockNumber < startingBlock/
	if log.BlockNumber < startingBlockV3 {
		switch log.Topics[0].Hex() {
		case etherReceivalEvent:
			amount, err := logDataToEtherReceivalParams(log.Data)
			if err != nil {
				return false, err
			}
			tradeLog.EtherReceivalSender = ethereum.BytesToAddress(log.Topics[1].Bytes())
			tradeLog.EtherReceivalAmount = amount.Big()
			return false, nil
		case tradeEvent:
			srcAddr, destAddr, srcAmount, destAmount, err := logDataToTradeParams(log.Data)
			if err != nil {
				return false, err
			}

			tradeLog.SrcAddress = srcAddr
			tradeLog.DestAddress = destAddr
			tradeLog.SrcAmount = srcAmount.Big()
			tradeLog.DestAmount = destAmount.Big()
			tradeLog.UserAddress = ethereum.BytesToAddress(log.Topics[1].Bytes())
			tradeLog.Index = log.Index

			tradeLog.BlockNumber = log.BlockNumber
			tradeLog.TransactionHash = log.TxHash

			if tradeLog.Timestamp, err = crawler.txTime.Resolve(log.BlockNumber); err != nil {
				return false, err
			}

			crawler.sugar.Infow("gathered new trade log", "trade_log", tradeLog)
			return true, nil
			// prepare to fulfill new TradeLog from next event logs
		default:
			crawler.sugar.Info("Unknown log topic.")
			return false, nil
		}
	}

	// We're not picking up reserver Addressese from KyberTrade event, they will be concerned at burnfee
	if log.Topics[0].Hex() == kyberTradeEvent {
		if len(log.Data) < 256 {
			return false, errors.New("log topic kyberTradeEvent came with wrong data len. Expect at least 256 bytes")
		}
		tradeLog.SrcAddress = ethereum.BytesToAddress(log.Data[0:32])
		tradeLog.DestAddress = ethereum.BytesToAddress(log.Data[32:64])
		tradeLog.SrcAmount = ethereum.BytesToHash(log.Data[64:96]).Big()
		tradeLog.DestAmount = ethereum.BytesToHash(log.Data[96:128]).Big()
		tradeLog.UserAddress = ethereum.BytesToAddress(log.Topics[1].Bytes())
		tradeLog.EtherReceivalAmount = ethereum.BytesToHash(log.Data[160:192]).Big()

		tradeLog.EtherReceivalSender = ethereum.BytesToAddress(log.Topics[1].Bytes())
		tradeLog.Index = log.Index
		tradeLog.BlockNumber = log.BlockNumber
		tradeLog.TransactionHash = log.TxHash
		var err error
		if tradeLog.Timestamp, err = crawler.txTime.Resolve(log.BlockNumber); err != nil {
			return false, err
		}
		crawler.sugar.Infow("gathered new trade log", "trade_log", tradeLog)
		return true, nil
	}
	return false, nil
}
