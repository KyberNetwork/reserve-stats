package postgres

import (
	"math/big"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

type record struct {
	Timestamp          time.Time `db:"timestamp"`
	BlockNumber        uint64    `db:"block_number"`
	TransactionHash    string    `db:"tx_hash"`
	USDTAmount         float64   `db:"usdt_amount"`
	OriginalUSDTAmount float64   `db:"original_usdt_amount"`
	UserAddress        string    `db:"user_address"`
	SrcAddress         string    `db:"src_address"`
	DestAddress        string    `db:"dst_address"`
	SrcReserveAddress  string    `db:"src_reserve_address"`
	DstReserveAddress  string    `db:"dst_reserve_address"`
	SrcAmount          float64   `db:"src_amount"`
	DestAmount         float64   `db:"dst_amount"`
	Index              string    `db:"index"`
	IsFirstTrade       bool      `db:"is_first_trade"`
	TxSender           string    `db:"tx_sender"`
	ReceiverAddress    string    `db:"receiver_address"`
	GasUsed            uint64    `db:"gas_used"`
	GasPrice           float64   `db:"gas_price"`
	TransactionFee     float64   `db:"transaction_fee"`
	Version            uint      `db:"version"`
}

func (tldb *TradeLogDB) calculateDstAmount(log common.Tradelog) (float64, error) {
	// this formula is base on https://github.com/KyberNetwork/smart-contracts/blob/Katalyst/contracts/sol6/utils/Utils5.sol#L88
	var (
		srcDecimals, dstDecimals int64
		err                      error
		dstAmount                float64
	)
	srcDecimals, err = tldb.tokenAmountFormatter.GetDecimals(log.TokenInfo.SrcAddress)
	if err != nil {
		return dstAmount, err
	}
	dstDecimals, err = tldb.tokenAmountFormatter.GetDecimals(blockchain.USDTAddr)
	if err != nil {
		return dstAmount, err
	}
	if dstDecimals >= srcDecimals {
		precision := new(big.Float).SetInt(new(big.Int).Exp(
			big.NewInt(10), big.NewInt(18), nil,
		))
		exp := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(dstDecimals-srcDecimals), nil)
		tmp := log.DestAmount.Mul(log.DestAmount, exp) // log.DestAmount is is equal srcAmount*rate when fillKyberTradeEvent
		dstAmountInt, _ := new(big.Float).Quo(new(big.Float).SetInt(tmp), precision).Int(nil)
		dstAmount, err = tldb.tokenAmountFormatter.FromWei(log.TokenInfo.DestAddress, dstAmountInt)
		if err != nil {
			return dstAmount, err
		}
	} else {
		precision := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil)
		exp := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(srcDecimals-dstDecimals), nil)
		tmp := big.NewInt(0).Mul(exp, precision)
		dstAmountInt, _ := new(big.Float).Quo(new(big.Float).SetInt(log.DestAmount), new(big.Float).SetInt(tmp)).Int(nil) // log.DestAmount is is equal srcAmount*rate when fillKyberTradeEvent
		dstAmount, err = tldb.tokenAmountFormatter.FromWei(log.TokenInfo.DestAddress, dstAmountInt)
		if err != nil {
			return dstAmount, err
		}
	}
	return dstAmount, nil
}

func (tldb *TradeLogDB) recordFromTradeLog(log common.Tradelog) (*record, error) {
	var dstAmount float64
	usdtAmount, err := tldb.tokenAmountFormatter.FromWei(blockchain.USDTAddr, log.USDTAmount)
	if err != nil {
		return nil, err
	}

	originalUSDTAmount, err := tldb.tokenAmountFormatter.FromWei(blockchain.USDTAddr, log.OriginalUSDTAmount)
	if err != nil {
		return nil, err
	}

	srcAmount, err := tldb.tokenAmountFormatter.FromWei(log.TokenInfo.SrcAddress, log.SrcAmount)
	if err != nil {
		return nil, err
	}
	if log.Version == 1 {
		dstAmount, err = tldb.calculateDstAmount(log)
		if err != nil {
			return nil, err
		}
	}

	transactionFee, err := tldb.tokenAmountFormatter.FromWei(blockchain.USDTAddr, log.TxDetail.TransactionFee)
	if err != nil {
		return nil, err
	}
	gasPrice, err := tldb.tokenAmountFormatter.FromWei(blockchain.USDTAddr, log.TxDetail.GasPrice)
	if err != nil {
		return nil, err
	}

	return &record{
		Timestamp:          log.Timestamp.UTC(),
		BlockNumber:        log.BlockNumber,
		TransactionHash:    log.TransactionHash.String(),
		USDTAmount:         usdtAmount,
		OriginalUSDTAmount: originalUSDTAmount,
		UserAddress:        log.User.UserAddress.String(),
		SrcAddress:         log.TokenInfo.SrcAddress.String(),
		DestAddress:        log.TokenInfo.DestAddress.String(),
		SrcAmount:          srcAmount,
		DestAmount:         dstAmount,
		SrcReserveAddress:  log.SrcReserveAddress.Hex(),
		DstReserveAddress:  log.DstReserveAddress.Hex(),
		Index:              strconv.FormatUint(uint64(log.Index), 10),
		TxSender:           log.TxDetail.TxSender.Hex(),
		ReceiverAddress:    log.ReceiverAddress.Hex(),
		TransactionFee:     transactionFee,
		GasPrice:           gasPrice,
		GasUsed:            log.TxDetail.GasUsed,
		Version:            log.Version,
	}, nil
}
