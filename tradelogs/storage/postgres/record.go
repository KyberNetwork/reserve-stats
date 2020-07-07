package postgres

import (
	"database/sql"
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

func (tldb *TradeLogDB) calculateDstAmountV4(log common.TradelogV4) (float64, error) {
	var (
		srcDecimals, dstDecimals int64
		err                      error
		dstAmount                float64
	)
	if len(log.E2TReserves) != 0 {
		srcDecimals, err = tldb.tokenAmountFormatter.GetDecimals(blockchain.ETHAddr)
		if err != nil {
			return dstAmount, err
		}
		dstDecimals, err = tldb.tokenAmountFormatter.GetDecimals(log.TokenInfo.DestAddress)
		if err != nil {
			return dstAmount, err
		}
	} else {
		srcDecimals, err = tldb.tokenAmountFormatter.GetDecimals(log.TokenInfo.SrcAddress)
		if err != nil {
			return dstAmount, err
		}
		dstDecimals, err = tldb.tokenAmountFormatter.GetDecimals(blockchain.ETHAddr)
		if err != nil {
			return dstAmount, err
		}
	}
	if dstDecimals >= srcDecimals {
		precision := new(big.Float).SetInt(new(big.Int).Exp(
			big.NewInt(10), big.NewInt(18), nil,
		))
		exp := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(dstDecimals-srcDecimals), nil)
		tmp := log.DestAmount.Mul(log.DestAmount, exp)
		dstAmountInt, _ := new(big.Float).Quo(new(big.Float).SetInt(tmp), precision).Int(nil)
		dstAmount, err = tldb.tokenAmountFormatter.FromWei(log.TokenInfo.DestAddress, dstAmountInt)
		if err != nil {
			return dstAmount, err
		}
	} else {
		precision := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil)
		exp := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(srcDecimals-dstDecimals), nil)
		tmp := big.NewInt(0).Mul(exp, precision)
		dstAmountInt, _ := new(big.Float).Quo(new(big.Float).SetInt(log.DestAmount), new(big.Float).SetInt(tmp)).Int(nil)
		dstAmount, err = tldb.tokenAmountFormatter.FromWei(log.TokenInfo.DestAddress, dstAmountInt)
		if err != nil {
			return dstAmount, err
		}
	}
	return dstAmount, nil
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
		dstAmount, err = tldb.calculateDstAmountV4(log)
		if err != nil {
			return nil, err
		}
	} else {
		dstAmount, err = tldb.tokenAmountFormatter.FromWei(log.TokenInfo.DestAddress, log.DestAmount)
		if err != nil {
			return nil, err
		}
	}

	transactionFee, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, log.TxDetail.TransactionFee)
	if err != nil {
		return nil, err
	}
	gasPrice, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, log.TxDetail.GasPrice)
	if err != nil {
		return nil, err
	}
	var (
		t2eSrcAmounts, e2tSrcAmounts, t2eRates, e2tRates []float64
	)
	for _, am := range log.T2ESrcAmount {
		amount, err := tldb.tokenAmountFormatter.FromWei(log.TokenInfo.SrcAddress, am)
		if err != nil {
			return nil, err
		}
		t2eSrcAmounts = append(t2eSrcAmounts, amount)
	}

	for _, am := range log.E2TSrcAmount {
		amount, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, am)
		if err != nil {
			return nil, err
		}
		e2tSrcAmounts = append(e2tSrcAmounts, amount)
	}

	for _, r := range log.T2ERates {
		rate, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, r)
		if err != nil {
			return nil, err
		}
		t2eRates = append(t2eRates, rate)
	}

	for _, r := range log.E2TRates {
		rate, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, r)
		if err != nil {
			return nil, err
		}
		e2tRates = append(e2tRates, rate)
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
		T2ESrcAmount:      t2eSrcAmounts,
		E2TSrcAmount:      e2tSrcAmounts,
		T2ERates:          t2eRates,
		E2TRates:          e2tRates,
		SrcAmount:         srcAmount,
		DestAmount:        dstAmount,
		WalletAddress:     log.WalletAddress.String(),
		WalletName:        log.WalletName,
		SrcReserveAddress: log.SrcReserveAddress.Hex(),
		DstReserveAddress: log.DstReserveAddress.Hex(),
		IntegrationApp:    log.IntegrationApp,
		IP:                sql.NullString{String: log.User.IP, Valid: log.User.IP != ""},
		Country:           sql.NullString{String: log.User.Country, Valid: log.User.Country != ""},
		ETHUSDRate:        log.ETHUSDRate,
		ETHUSDProvider:    log.ETHUSDProvider,
		Index:             strconv.FormatUint(uint64(log.Index), 10),
		Kyced:             log.User.UID != "",
		TxSender:          log.TxDetail.TxSender.Hex(),
		ReceiverAddress:   log.ReceiverAddress.Hex(),
		TransactionFee:    transactionFee,
		GasPrice:          gasPrice,
		GasUsed:           log.TxDetail.GasUsed,
		Version:           log.Version,
		Fee:               log.Fees,
	}, nil
}
