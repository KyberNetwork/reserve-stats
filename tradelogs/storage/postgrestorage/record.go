package postgrestorage

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

type record struct {
	Timestamp          time.Time      `db:"timestamp"`
	BlockNumber        uint64         `db:"block_number"`
	TransactionHash    string         `db:"tx_hash"`
	EthAmount          float64        `db:"eth_amount"`
	UserAddress        string         `db:"user_address"`
	SrcAddress         string         `db:"src_address"`
	DestAddress        string         `db:"dest_address"`
	SrcReserveAddress  string         `db:"src_reserveaddress"`
	DstReserveAddress  string         `db:"dst_reserveaddress"`
	SrcAmount          float64        `db:"src_amount"`
	DestAmount         float64        `db:"dest_amount"`
	WalletAddress      string         `db:"wallet_address"`
	SrcBurnAmount      float64        `db:"src_burn_amount"`
	DstBurnAmount      float64        `db:"dst_burn_amount"`
	SrcWalletFeeAmount float64        `db:"src_wallet_fee_amount"`
	DstWalletFeeAmount float64        `db:"dst_wallet_fee_amount"`
	IntegrationApp     string         `db:"integration_app"`
	IP                 sql.NullString `db:"ip"`
	Country            sql.NullString `db:"country"`
	ETHUSDRate         float64        `db:"ethusd_rate"`
	ETHUSDProvider     string         `db:"ethusd_provider"`
	Index              string         `db:"index"`
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

	srcAmount, err := tldb.tokenAmountFormatter.FromWei(log.SrcAddress, log.SrcAmount)
	if err != nil {
		return nil, err
	}

	dstAmount, err := tldb.tokenAmountFormatter.FromWei(log.DestAddress, log.DestAmount)
	if err != nil {
		return nil, err
	}

	srcBurnAmount, dstBurnAmount, err := tldb.getBurnAmount(log)
	if err != nil {
		return nil, err
	}

	srcWalletFee, dstWalletFee, err := tldb.getWalletFeeAmount(log)
	if err != nil {
		return nil, err
	}
	return &record{
		Timestamp:          log.Timestamp,
		BlockNumber:        log.BlockNumber,
		TransactionHash:    log.TransactionHash.String(),
		EthAmount:          ethAmount,
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
	}, nil
}

//getBurnAmount return the burn amount in float for src and
func (tldb *TradeLogDB) getBurnAmount(log common.TradeLog) (float64, float64, error) {
	var (
		logger = tldb.sugar.With(
			"func", "tradelogs/storage/postgrestorage/getBurnAmount",
			"log", log,
		)
		srcAmount float64
		dstAmount float64
	)
	if blockchain.IsBurnable(log.SrcAddress) {
		if len(log.BurnFees) < 1 {
			logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "at least 1 burn fees (src)")
			return srcAmount, dstAmount, nil
		}
		srcAmount, err := tldb.tokenAmountFormatter.FromWei(blockchain.KNCAddr, log.BurnFees[0].Amount)
		if err != nil {
			return srcAmount, dstAmount, err
		}

		if blockchain.IsBurnable(log.DestAddress) {
			if len(log.BurnFees) < 2 {
				logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "2 burn fees (src-dst)")
				return srcAmount, dstAmount, nil
			}
			dstAmount, err = tldb.tokenAmountFormatter.FromWei(blockchain.KNCAddr, log.BurnFees[1].Amount)
			if err != nil {
				return srcAmount, dstAmount, err
			}
			return srcAmount, dstAmount, nil
		}
		return srcAmount, dstAmount, nil
	}

	if blockchain.IsBurnable(log.DestAddress) {
		if len(log.BurnFees) < 1 {
			logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "at least 1 burn fees (dst)")
			return srcAmount, dstAmount, nil
		}
		dstAmount, err := tldb.tokenAmountFormatter.FromWei(blockchain.KNCAddr, log.BurnFees[0].Amount)
		if err != nil {
			return srcAmount, dstAmount, err
		}
		return srcAmount, dstAmount, nil
	}

	return srcAmount, dstAmount, nil
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

		if walletFee.ReserveAddress == log.SrcReserveAddress && !srcAmountSet {
			srcAmount = amount
			srcAmountSet = true
		} else if walletFee.ReserveAddress == log.DstReserveAddress {
			dstAmount = amount
		} else {
			logger.Warnw("unexpected wallet fees with unrecognized reserve address", "wallet fee", walletFee)
		}
	}
	return srcAmount, dstAmount, nil
}
