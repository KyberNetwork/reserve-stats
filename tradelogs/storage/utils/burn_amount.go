package utils

import (
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"go.uber.org/zap"
)

// GetBurnAmount return the burn amount in float for src and
func GetBurnAmount(sugar *zap.SugaredLogger, tokenAmountFormatter blockchain.TokenAmountFormatterInterface, log common.TradeLog) (float64, float64, error) {
	var (
		logger = sugar.With(
			"func", "tradelogs/storage/getBurnAmount",
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
		srcAmount, err := tokenAmountFormatter.FromWei(blockchain.KNCAddr, log.BurnFees[0].Amount)
		if err != nil {
			return srcAmount, dstAmount, err
		}

		if blockchain.IsBurnable(log.DestAddress) {
			if len(log.BurnFees) < 2 {
				logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "2 burn fees (src-dst)")
				return srcAmount, dstAmount, nil
			}
			dstAmount, err = tokenAmountFormatter.FromWei(blockchain.KNCAddr, log.BurnFees[1].Amount)
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
		dstAmount, err := tokenAmountFormatter.FromWei(blockchain.KNCAddr, log.BurnFees[0].Amount)
		if err != nil {
			return srcAmount, dstAmount, err
		}
		return srcAmount, dstAmount, nil
	}

	return srcAmount, dstAmount, nil
}
