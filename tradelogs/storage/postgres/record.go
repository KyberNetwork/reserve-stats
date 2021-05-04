package postgres

import (
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

type record struct {
	Timestamp           time.Time `db:"timestamp"`
	BlockNumber         uint64    `db:"block_number"`
	TransactionHash     string    `db:"tx_hash"`
	QuoteAmount         float64   `db:"usdt_amount"`
	OriginalQuoteAmount float64   `db:"original_usdt_amount"`
	ReserveAddress      string    `db:"reserve_address"`
	UserAddress         string    `db:"user_address"`
	SrcAddress          string    `db:"src_address"`
	DestAddress         string    `db:"dst_address"`
	SrcReserveAddress   string    `db:"src_reserve_address"`
	DstReserveAddress   string    `db:"dst_reserve_address"`
	SrcAmount           float64   `db:"src_amount"`
	DestAmount          float64   `db:"dst_amount"`
	Index               string    `db:"index"`
	IsFirstTrade        bool      `db:"is_first_trade"`
	TxSender            string    `db:"tx_sender"`
	ReceiverAddress     string    `db:"receiver_address"`
	GasUsed             uint64    `db:"gas_used"`
	GasPrice            float64   `db:"gas_price"`
	TransactionFee      float64   `db:"transaction_fee"`
	Version             uint      `db:"version"`
}

func (tldb *TradeLogDB) recordFromTradeLog(log common.Tradelog) (*record, error) {
	var dstAmount float64
	usdtAmount, err := tldb.tokenAmountFormatter.FromWei(blockchain.USDTAddr, log.QuoteAmount)
	if err != nil {
		return nil, err
	}

	originalUSDTAmount, err := tldb.tokenAmountFormatter.FromWei(blockchain.USDTAddr, log.OriginalQuoteAmount)
	if err != nil {
		return nil, err
	}

	srcAmount, err := tldb.tokenAmountFormatter.FromWei(log.TokenInfo.SrcAddress, log.SrcAmount)
	if err != nil {
		return nil, err
	}

	dstAmount, err = tldb.tokenAmountFormatter.FromWei(log.TokenInfo.DestAddress, log.DestAmount)
	if err != nil {
		return nil, err
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
		Timestamp:           log.Timestamp.UTC(),
		BlockNumber:         log.BlockNumber,
		TransactionHash:     log.TransactionHash.String(),
		QuoteAmount:         usdtAmount,
		OriginalQuoteAmount: originalUSDTAmount,
		ReserveAddress:      log.ReserveAddress.Hex(),
		UserAddress:         log.User.UserAddress.String(),
		SrcAddress:          log.TokenInfo.SrcAddress.String(),
		DestAddress:         log.TokenInfo.DestAddress.String(),
		SrcAmount:           srcAmount,
		DestAmount:          dstAmount,
		SrcReserveAddress:   log.SrcReserveAddress.Hex(),
		DstReserveAddress:   log.DstReserveAddress.Hex(),
		Index:               strconv.FormatUint(uint64(log.Index), 10),
		TxSender:            log.TxDetail.TxSender.Hex(),
		ReceiverAddress:     log.ReceiverAddress.Hex(),
		TransactionFee:      transactionFee,
		GasPrice:            gasPrice,
		GasUsed:             log.TxDetail.GasUsed,
		Version:             log.Version,
	}, nil
}
