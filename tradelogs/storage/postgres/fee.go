package postgres

import (
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// GetFeeByTradelogID return fee by tradelog id
func (tldb *TradeLogDB) GetFeeByTradelogID(tradelogID uint64) (common.TradelogFee, error) {
	var (
		fee common.TradelogFee
	)
	return fee, nil
}
