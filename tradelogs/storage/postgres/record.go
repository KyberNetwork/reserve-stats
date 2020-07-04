package postgres

import (
	"database/sql"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

type record struct {
	Timestamp         time.Time            `db:"timestamp"`
	BlockNumber       uint64               `db:"block_number"`
	TransactionHash   string               `db:"tx_hash"`
	EthAmount         float64              `db:"eth_amount"`
	OriginalEthAmount float64              `db:"original_eth_amount"`
	UserAddress       string               `db:"user_address"`
	SrcAddress        string               `db:"src_address"`
	DestAddress       string               `db:"dst_address"`
	SrcReserveAddress string               `db:"src_reserve_address"`
	DstReserveAddress string               `db:"dst_reserve_address"`
	T2EReserves       [][32]byte           `db:"t2e_reserves"`
	E2TReserves       [][32]byte           `db:"e2t_reserves"`
	T2ESrcAmount      []float64            `db:"t2e_src_amount"`
	E2TSrcAmount      []float64            `db:"t2e_src_amount"`
	T2ERates          []float64            `db:"t2e_rates"`
	E2TRates          []float64            `db:"e2t_rates"`
	SrcAmount         float64              `db:"src_amount"`
	DestAmount        float64              `db:"dst_amount"`
	WalletAddress     string               `db:"wallet_address"`
	WalletName        string               `db:"wallet_name"`
	IntegrationApp    string               `db:"integration_app"`
	IP                sql.NullString       `db:"ip"`
	Country           sql.NullString       `db:"country"`
	ETHUSDRate        float64              `db:"eth_usd_rate"`
	ETHUSDProvider    string               `db:"eth_usd_provider"`
	Index             string               `db:"index"`
	Kyced             bool                 `db:"kyced"`
	IsFirstTrade      bool                 `db:"is_first_trade"`
	TxSender          string               `db:"tx_sender"`
	ReceiverAddress   string               `db:"receiver_address"`
	GasUsed           uint64               `db:"gas_used"`
	GasPrice          float64              `db:"gas_price"`
	TransactionFee    float64              `db:"transaction_fee"`
	Version           uint                 `db:"version"`
	Fee               []common.TradelogFee `db:"fee"`
}

func (tldb *TradeLogDB) recordFromTradeLog(log common.TradelogV4) (*record, error) {
	var dstAmount float64
	ethAmount, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, log.EthAmount)
	if err != nil {
		return nil, err
	}

	originalEthAmount, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, log.OriginalEthAmount)
	if err != nil {
		return nil, err
	}

	srcAmount, err := tldb.tokenAmountFormatter.FromWei(log.TokenInfo.SrcAddress, log.SrcAmount)
	if err != nil {
		return nil, err
	}
	if log.Version == 4 {
		dstAmount, err = tldb.tokenAmountFormatter.FromWei(log.TokenInfo.SrcAddress, log.DestAmount)
		if err != nil {
			return nil, err
		}
		floatDstAmount := big.NewFloat(0).SetFloat64(dstAmount)
		bigDstAmount := big.NewInt(0)
		floatDstAmount.Int(bigDstAmount)
		dstAmount, err = tldb.tokenAmountFormatter.FromWei(log.TokenInfo.DestAddress, bigDstAmount)
		if err != nil {
			return nil, err
		}
		if dstAmount == 0 {
			fmt.Println(bigDstAmount)
			fmt.Println(log.TransactionHash.Hex())
			panic(log.DestAmount)
		}
	} else {
		dstAmount, err = tldb.tokenAmountFormatter.FromWei(log.TokenInfo.DestAddress, log.DestAmount)
		if err != nil {
			return nil, err
		}
	}

	// srcBurnAmount, dstBurnAmount, err := utils.GetBurnAmount(tldb.sugar, tldb.tokenAmountFormatter, log, tldb.kncAddr)
	// if err != nil {
	// 	return nil, err
	// }

	// srcWalletFee, dstWalletFee, err := tldb.getWalletFeeAmount(log)
	// if err != nil {
	// 	return nil, err
	// }
	transactionFee, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, log.TxDetail.TransactionFee)
	if err != nil {
		return nil, err
	}
	gasPrice, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, log.TxDetail.GasPrice)
	if err != nil {
		return nil, err
	}

	return &record{
		Timestamp:         log.Timestamp.UTC(),
		BlockNumber:       log.BlockNumber,
		TransactionHash:   log.TransactionHash.String(),
		EthAmount:         ethAmount,
		OriginalEthAmount: originalEthAmount,
		UserAddress:       log.User.UserAddress.String(),
		SrcAddress:        log.TokenInfo.SrcAddress.String(),
		DestAddress:       log.TokenInfo.DestAddress.String(),
		T2EReserves:       log.T2EReserves,
		E2TReserves:       log.E2TReserves,
		SrcAmount:         srcAmount,
		DestAmount:        dstAmount,
		WalletAddress:     log.WalletAddress.String(),
		WalletName:        log.WalletName,
		// SrcBurnAmount:     srcBurnAmount,
		// DstBurnAmount:     dstBurnAmount,
		// TODO: fill wallet fee amount
		// SrcWalletFeeAmount: srcWalletFee,
		// DstWalletFeeAmount: dstWalletFee,
		IntegrationApp:  log.IntegrationApp,
		IP:              sql.NullString{String: log.User.IP, Valid: log.User.IP != ""},
		Country:         sql.NullString{String: log.User.Country, Valid: log.User.Country != ""},
		ETHUSDRate:      log.ETHUSDRate,
		ETHUSDProvider:  log.ETHUSDProvider,
		Index:           strconv.FormatUint(uint64(log.Index), 10),
		Kyced:           log.User.UID != "",
		TxSender:        log.TxDetail.TxSender.Hex(),
		ReceiverAddress: log.ReceiverAddress.Hex(),
		TransactionFee:  transactionFee,
		GasPrice:        gasPrice,
		GasUsed:         log.TxDetail.GasUsed,
		Version:         log.Version,
		Fee:             log.Fees,
	}, nil
}

// func (tldb *TradeLogDB) getWalletFeeAmount(log common.TradelogV4) (float64, float64, error) {
// 	var (
// 		logger = tldb.sugar.With(
// 			"func", caller.GetCurrentFunctionName(),
// 			"log", log,
// 		)
// 		dstAmount    float64
// 		srcAmount    float64
// 		srcAmountSet bool
// 	)
// 	for _, walletFee := range log.WalletFees {
// 		amount, err := tldb.tokenAmountFormatter.FromWei(blockchain.KNCAddr, walletFee.Amount)
// 		if err != nil {
// 			return dstAmount, srcAmount, err
// 		}

// 		switch {
// 		case walletFee.ReserveAddress == log.SrcReserveAddress && !srcAmountSet:
// 			srcAmount = amount
// 			srcAmountSet = true
// 		case walletFee.ReserveAddress == log.DstReserveAddress:
// 			dstAmount = amount
// 		default:
// 			logger.Warnw("unexpected wallet fees with unrecognized reserve address", "wallet fee", walletFee)
// 		}
// 	}
// 	return srcAmount, dstAmount, nil
// }
