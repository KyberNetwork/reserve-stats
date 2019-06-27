package postgrestorage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"
)

func (tldb *TradeLogDB) GetUserList(fromTime, toTime time.Time, timezone int8) ([]common.UserInfo, error) {
	logger := tldb.sugar.With("from", fromTime, "to", toTime, "timezone", timezone)
	fromTime = fromTime.UTC().Add(time.Duration(-timezone) * time.Hour)
	toTime = toTime.UTC().Add(time.Duration(-timezone) * time.Hour)

	var result []common.UserInfo
	if err := tldb.db.Select(&result, userListQuery, fromTime.Format(schema.DefaultDateFormat),
		toTime.Format(schema.DefaultDateFormat)); err != nil {
		return nil, err
	}
	if len(result) == 0 {
		logger.Debugw("empty user list result", "query", userListQuery)
		return nil, nil
	}
	return result, nil

}

const userListQuery = `
	SELECT b.address user_address,sum(eth_amount) total_eth_volume,
		sum(eth_amount * eth_usd_rate) total_usd_volume
	FROM "` + schema.TradeLogsTableName + `" a
	INNER JOIN "` + schema.UserTableName + `" b ON a.user_address_id =b.id
	WHERE a.timestamp >= $1 and a.timestamp <= $2
	GROUP BY user_address
`
