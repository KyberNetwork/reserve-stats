package postgrestorage

import (
	"database/sql"
	"strconv"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/utils"
)

type record struct {
	Timestamp          time.Time      `db:"timestamp"`
	BlockNumber        uint64         `db:"block_number"`
	TransactionHash    string         `db:"tx_hash"`
	EthAmount          float64        `db:"eth_amount"`
	OriginalEthAmount  float64        `db:"original_eth_amount"`
	UserAddress        string         `db:"user_address"`
	SrcAddress         string         `db:"src_address"`
	DestAddress        string         `db:"dst_address"`
	SrcReserveAddress  string         `db:"src_reserve_address"`
	DstReserveAddress  string         `db:"dst_reserve_address"`
	SrcAmount          float64        `db:"src_amount"`
	DestAmount         float64        `db:"dst_amount"`
	WalletAddress      string         `db:"wallet_address"`
	SrcBurnAmount      float64        `db:"src_burn_amount"`
	DstBurnAmount      float64        `db:"dst_burn_amount"`
	SrcWalletFeeAmount float64        `db:"src_wallet_fee_amount"`
	DstWalletFeeAmount float64        `db:"dst_wallet_fee_amount"`
	IntegrationApp     string         `db:"integration_app"`
	IP                 sql.NullString `db:"ip"`
	Country            sql.NullString `db:"country"`
	ETHUSDRate         float64        `db:"eth_usd_rate"`
	ETHUSDProvider     string         `db:"eth_usd_provider"`
	Index              string         `db:"index"`
	Kyced              bool           `db:"kyced"`
	IsFirstTrade       bool           `db:"is_first_trade"`
	TxSender           string         `db:"tx_sender"`
	ReceiverAddress    string         `db:"receiver_address"`
}

func (tldb *TradeLogDB) recordFromTradeLog(log common.TradeLog) (*record, error) {
	var walletAddr ethereum.Address
	if len(log.WalletFees) > 0 {
		walletAddr = log.WalletFees[0].WalletAddress
	}
	ethAmount, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, log.EthAmount)
	if err != nil {
		return nil, err
	}

	originalEthAmount, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, log.OriginalEthAmount)
	if err != nil {
		return nil, err
	}

	srcAmount, err := tldb.tokenAmountFormatter.FromWei(log.SrcAddress, log.SrcAmount)
	if err != nil {
		return nil, err
	}

	dstAmount, err := tldb.tokenAmountFormatter.FromWei(log.DestAddress, log.DestAmount)
	if err != nil {
		return nil, err
	}

	srcBurnAmount, dstBurnAmount, err := utils.GetBurnAmount(tldb.sugar, tldb.tokenAmountFormatter, log)
	if err != nil {
		return nil, err
	}

	srcWalletFee, dstWalletFee, err := tldb.getWalletFeeAmount(log)
	if err != nil {
		return nil, err
	}
	return &record{
		Timestamp:          log.Timestamp.UTC(),
		BlockNumber:        log.BlockNumber,
		TransactionHash:    log.TransactionHash.String(),
		EthAmount:          ethAmount,
		OriginalEthAmount:  originalEthAmount,
		UserAddress:        log.UserAddress.String(),
		SrcAddress:         log.SrcAddress.String(),
		DestAddress:        log.DestAddress.String(),
		SrcReserveAddress:  log.SrcReserveAddress.String(),
		DstReserveAddress:  log.DstReserveAddress.String(),
		SrcAmount:          srcAmount,
		DestAmount:         dstAmount,
		WalletAddress:      walletAddr.String(),
		SrcBurnAmount:      srcBurnAmount,
		DstBurnAmount:      dstBurnAmount,
		SrcWalletFeeAmount: srcWalletFee,
		DstWalletFeeAmount: dstWalletFee,
		IntegrationApp:     log.IntegrationApp,
		IP:                 sql.NullString{String: log.IP, Valid: log.IP != ""},
		Country:            sql.NullString{String: log.Country, Valid: log.Country != ""},
		ETHUSDRate:         log.ETHUSDRate,
		ETHUSDProvider:     log.ETHUSDProvider,
		Index:              strconv.FormatUint(uint64(log.Index), 10),
		Kyced:              log.UID != "",
		TxSender:           log.TxSender.Hex(),
		ReceiverAddress:    log.ReceiverAddress.Hex(),
	}, nil
}

func (tldb *TradeLogDB) getWalletFeeAmount(log common.TradeLog) (float64, float64, error) {
	var (
		logger = tldb.sugar.With(
			"func", "tradelogs/storage/postgrestorage/getWalletFeeAmount",
			"log", log,
		)
		dstAmount    float64
		srcAmount    float64
		srcAmountSet bool
	)
	for _, walletFee := range log.WalletFees {
		amount, err := tldb.tokenAmountFormatter.FromWei(blockchain.KNCAddr, walletFee.Amount)
		if err != nil {
			return dstAmount, srcAmount, err
		}

		switch {
		case walletFee.ReserveAddress == log.SrcReserveAddress && !srcAmountSet:
			srcAmount = amount
			srcAmountSet = true
		case walletFee.ReserveAddress == log.DstReserveAddress:
			dstAmount = amount
		default:
			logger.Warnw("unexpected wallet fees with unrecognized reserve address", "wallet fee", walletFee)
		}
	}
	return srcAmount, dstAmount, nil
}
