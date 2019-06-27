package postgrestorage

import (
	"fmt"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"

	ethereum "github.com/ethereum/go-ethereum/common"
)

func (tldb *TradeLogDB) GetUserVolume(userAddress ethereum.Address, from, to time.Time, freq string) (map[uint64]common.UserVolume, error) {
	var (
		timeField string
		logger    = tldb.sugar.With("from", from, "to", to, "userAddress",
			userAddress, "freq", freq)
	)

	switch strings.ToLower(freq) {
	case "h":
		timeField = schema.BuildDateTruncField("hour", 0)
		from = schema.RoundTime(from, "hour", 0)
		to = schema.RoundTime(to, "hour", 0).Add(time.Hour)
	case "d":
		timeField = schema.BuildDateTruncField("day", 0)
		from = schema.RoundTime(from, "day", 0)
		to = schema.RoundTime(to, "day", 0).Add(time.Hour * 24)

	default:
		return nil, fmt.Errorf("frequency not supported: %v", freq)
	}

	query :=
		`SELECT ` + timeField + ` AS time, SUM(eth_amount) eth_volume,
			SUM(eth_amount * eth_usd_rate) usd_volume
		FROM "` + schema.TradeLogsTableName + `" a
		WHERE timestamp >= $1 AND timestamp < $2
		AND EXISTS (SELECT NULL FROM "` + schema.UserTableName + `" WHERE user_address_id = id and address = $3)
		GROUP BY time`

	logger.Debugw("prepare statement", "stmt", query)

	var records []struct {
		Time      time.Time `db:"time"`
		EthAmount float64   `db:"eth_volume"`
		UsdAmount float64   `db:"usd_volume"`
	}

	if err := tldb.db.Select(&records, query, from.UTC().Format(schema.DefaultDateFormat),
		to.UTC().Format(schema.DefaultDateFormat), userAddress.Hex()); err != nil {
		return nil, err
	}

	result := make(map[uint64]common.UserVolume)
	for _, r := range records {
		key := timeutil.TimeToTimestampMs(r.Time)
		result[key] = common.UserVolume{
			ETHAmount: r.EthAmount,
			USDAmount: r.UsdAmount,
		}
	}
	return result, nil
}
