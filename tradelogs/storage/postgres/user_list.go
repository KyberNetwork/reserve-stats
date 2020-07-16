package postgres

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// GetUserList return list of user with his amount
func (tldb *TradeLogDB) GetUserList(fromTime, toTime time.Time) ([]common.UserInfo, error) {
	logger := tldb.sugar.With("from", fromTime, "to", toTime, "func", caller.GetCurrentFunctionName())

	userListQuery := `
		SELECT 
			b.address user_address,
			SUM(eth_amount) total_eth_volume,
			SUM(eth_amount * eth_usd_rate) total_usd_volume
		FROM "tradelogs" a
		INNER JOIN "users" b ON a.user_address_id =b.id
		WHERE a.timestamp >= $1 and a.timestamp <= $2
		GROUP BY user_address
	`
	logger.Debugw("prepare statement", "stmt", userListQuery)

	var result []common.UserInfo
	if err := tldb.db.Select(&result, userListQuery, fromTime,
		toTime); err != nil {
		return nil, err
	}
	if len(result) == 0 {
		logger.Debugw("empty user list result", "query", userListQuery)
		return nil, nil
	}
	return result, nil
}
