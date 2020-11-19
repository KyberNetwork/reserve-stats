package utils

import (
	ethereum "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// GetBurnAmount return the burn amount in float for src and
func GetBurnAmount(sugar *zap.SugaredLogger, tokenAmountFormatter blockchain.TokenAmountFormatterInterface, log common.TradelogV4, kncAddr ethereum.Address) (float64, float64, error) {
	var (
		logger = sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"log", log,
		)
		srcAmount float64
		dstAmount float64
	)

	// the tradelogs is fee free
	if common.LengthBurnFees(log) == 0 {
		return 0, 0, nil
	}

	// if len(log.BurnFees) > 0, check normal logic
	if blockchain.IsBurnable(log.TokenInfo.SrcAddress) {
		if common.LengthBurnFees(log) < 1 {
			logger.Warnw("unexpected burn fees", "got", common.LengthBurnFees(log), "want", "at least 1 burn fees (src)")
			return srcAmount, dstAmount, nil
		}
		//TODO:
		// srcAmount, err := tokenAmountFormatter.FromWei(kncAddr, log.BurnFees[0].Amount)
		// if err != nil {
		// 	return srcAmount, dstAmount, err
		// }

		// if blockchain.IsBurnable(log.DestAddress) {
		// 	if len(log.BurnFees) < 2 {
		// 		logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "2 burn fees (src-dst)")
		// 		return srcAmount, dstAmount, nil
		// 	}
		// 	dstAmount, err = tokenAmountFormatter.FromWei(kncAddr, log.BurnFees[1].Amount)
		// 	if err != nil {
		// 		return srcAmount, dstAmount, err
		// 	}
		// 	return srcAmount, dstAmount, nil
		// }
		return srcAmount, dstAmount, nil
	}

	if blockchain.IsBurnable(log.TokenInfo.DestAddress) {
		if common.LengthBurnFees(log) < 1 {
			logger.Warnw("unexpected burn fees", "got", common.LengthBurnFees(log), "want", "at least 1 burn fees (dst)")
			return srcAmount, dstAmount, nil
		}
		// TODO:
		// dstAmount, err := tokenAmountFormatter.FromWei(kncAddr, log.BurnFees[0].Amount)
		// if err != nil {
		// 	return srcAmount, dstAmount, err
		// }
		return srcAmount, dstAmount, nil
	}

	return srcAmount, dstAmount, nil
}
