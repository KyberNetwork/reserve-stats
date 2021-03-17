package postgres

import (
	"fmt"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
)

// GetAssetVolume returns eth_amount, usd_amount, volume filter by token addr in a time range group by day or hour
func (tldb *TradeLogDB) GetAssetVolume(token ethereum.Address, fromTime, toTime time.Time,
	frequency string) (map[uint64]*common.VolumeStats, error) {
	var (
		err       error
		timeField string
		logger    = tldb.sugar.With("from", fromTime, "to", toTime, "frequency", frequency, "token", token.Hex(),
			"func", caller.GetCurrentFunctionName())
	)

	switch strings.ToLower(frequency) {
	case "h":
		timeField = schema.BuildDateTruncField("hour", 0)
		fromTime = schema.RoundTime(fromTime, "hour", 0)
		toTime = schema.RoundTime(toTime, "hour", 0).Add(time.Hour)
	case "d":
		timeField = schema.BuildDateTruncField("day", 0)
		fromTime = schema.RoundTime(fromTime, "day", 0)
		toTime = schema.RoundTime(toTime, "day", 0).Add(time.Hour * 24)

	default:
		return nil, fmt.Errorf("frequency not supported: %v", frequency)
	}

	queryStmt := fmt.Sprintf(`
		SELECT time, 
		SUM(token_volume) token_volume, 
		SUM(usdt_amount) usd_volume
		FROM (
			SELECT %[1]s AS time, src_amount token_volume, usdt_amount
			FROM tradelogs
			WHERE EXISTS (SELECT NULL FROM "token" WHERE address = $3 AND id=src_address_id)
				AND timestamp >= $1 AND timestamp < $2
			UNION ALL
			SELECT %[1]s AS time, dst_amount token_volume, usdt_amount
			FROM "tradelogs" 
			WHERE EXISTS (SELECT NULL FROM "token" WHERE address = $3 AND id=dst_address_id)
				AND timestamp >= $1 AND timestamp < $2
		) a GROUP BY time;
	`, timeField)
	logger.Debugw("prepare statement", "stmt", queryStmt)

	var records []struct {
		TokenVolume float64   `db:"token_volume"`
		EthVolume   float64   `db:"eth_volume"`
		USDVolume   float64   `db:"usd_volume"`
		Time        time.Time `db:"time"`
	}
	err = tldb.db.Select(&records, queryStmt, fromTime, toTime, token.Hex())
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		logger.Debugw("return empty result", "prepare statement", queryStmt)
		return nil, nil
	}
	result := make(map[uint64]*common.VolumeStats)
	for _, data := range records {
		fmt.Println(data.TokenVolume, " ", data.USDVolume)
		result[timeutil.TimeToTimestampMs(data.Time)] = &common.VolumeStats{
			Volume:    data.TokenVolume,
			USDAmount: data.USDVolume,
		}
	}
	return result, nil
}

// GetReserveVolume returns eth_amount, usd_amount, volume filter by reserve addr and token addr in a time range group by day or hour
func (tldb *TradeLogDB) GetReserveVolume(rsvAddr ethereum.Address, token ethereum.Address,
	fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error) {
	var (
		err       error
		timeField string
		logger    = tldb.sugar.With("from", fromTime, "to", toTime, "frequency", frequency,
			"func", caller.GetCurrentFunctionName())
	)

	switch strings.ToLower(frequency) {
	case "h":
		timeField = schema.BuildDateTruncField("hour", 0)
		fromTime = schema.RoundTime(fromTime, "hour", 0)
		toTime = schema.RoundTime(toTime, "hour", 0).Add(time.Hour)
	case "d":
		timeField = schema.BuildDateTruncField("day", 0)
		fromTime = schema.RoundTime(fromTime, "day", 0)
		toTime = schema.RoundTime(toTime, "day", 0).Add(time.Hour * 24)

	default:
		return nil, fmt.Errorf("frequency not supported: %v", frequency)
	}

	reserveQuery := fmt.Sprintf(`
		SELECT 
			time, 
			SUM(token_volume) token_volume, 
			SUM(usdt_amount) usd_volume
		FROM (
			SELECT %[1]s AS time, usdt_amount,
				src_amount token_volume, 
			FROM "tradelogs" 
			WHERE EXISTS (SELECT NULL FROM "token" WHERE address = $1 AND id=src_address_id)
				AND EXISTS (SELECT NULL FROM "reserve" WHERE address = $2 AND (id= src_reserve_address_id OR id = dst_reserve_address_id))
				AND timestamp >= $3 AND timestamp < $4 
			UNION ALL
			SELECT %[1]s AS time, usdt_amount,
				dst_amount token_volume, 
			FROM "tradelogs"
			WHERE EXISTS (SELECT NULL FROM "token" WHERE address = $1 AND id=dst_address_id)
				AND EXISTS (SELECT NULL FROM "reserve" WHERE address = $2 AND (id= src_reserve_address_id OR id = dst_reserve_address_id))
				AND timestamp >= $3 AND timestamp < $4
			) a GROUP BY time
	`, timeField)
	logger.Debugw("prepare statement", "stmt", reserveQuery)
	var records []struct {
		TokenVolume float64   `db:"token_volume"`
		USDVolume   float64   `db:"usd_volume"`
		Time        time.Time `db:"time"`
	}

	err = tldb.db.Select(&records, reserveQuery, token.Hex(), rsvAddr.Hex(), fromTime, toTime)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		logger.Debugw("return empty result", "prepare statement", reserveQuery)
		return nil, nil
	}
	result := make(map[uint64]*common.VolumeStats)
	for _, r := range records {
		result[timeutil.TimeToTimestampMs(r.Time)] = &common.VolumeStats{
			Volume:    r.TokenVolume,
			USDAmount: r.USDVolume,
		}
	}
	return result, nil
}

// GetMonthlyVolume returns monthly volume
func (tldb *TradeLogDB) GetMonthlyVolume(rsvAddr ethereum.Address, from, to time.Time) (map[uint64]*common.VolumeStats, error) {
	return nil, nil
}
