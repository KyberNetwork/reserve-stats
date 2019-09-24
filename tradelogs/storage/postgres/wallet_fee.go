package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
)

// GetAggregatedWalletFee returns fee_amount filter by time, reserve addr, wallet addr group by hour or day
func (tldb *TradeLogDB) GetAggregatedWalletFee(reserveAddr, walletAddr, freq string,
	fromTime, toTime time.Time, timezone int8) (map[uint64]float64, error) {
	var (
		timeField string
		err       error
	)
	logger := tldb.sugar.With("from", fromTime, "to", toTime,
		"reserve", reserveAddr, "wallet", walletAddr,
		"func", caller.GetCurrentFunctionName())

	switch strings.ToLower(freq) {
	case "h":
		timeField = schema.BuildDateTruncField("hour", timezone)
		fromTime = schema.RoundTime(fromTime, "hour", timezone)
		toTime = schema.RoundTime(toTime, "hour", timezone).Add(time.Hour)
	case "d":
		timeField = schema.BuildDateTruncField("day", timezone)
		fromTime = schema.RoundTime(fromTime, "day", timezone)
		toTime = schema.RoundTime(toTime, "day", timezone).Add(time.Hour * 24)

	default:
		return nil, fmt.Errorf("frequency not supported: %v", freq)
	}

	integrationQuery := fmt.Sprintf(`
		SELECT time, SUM(fee_amount) as fee_amount 
		FROM (
			SELECT %[1]s as time, src_wallet_fee_amount AS fee_amount 
			FROM "%[2]s"
			WHERE timestamp >= $1 and timestamp < $2
			AND EXISTS (SELECT NULL FROM "%[3]s"
				WHERE address = $3 and id = wallet_address_id)
			AND EXISTS (SELECT NULL FROM "%[4]s"
				WHERE address = $4 and id = src_reserve_address_id)
		UNION ALL
			SELECT %[1]s as time, dst_wallet_fee_amount AS fee_amount 
			FROM "%[2]s"
			WHERE timestamp >= $1 and timestamp < $2
			AND EXISTS (SELECT NULL FROM "%[3]s"
				WHERE address = $3 and id = wallet_address_id)
			AND EXISTS (SELECT NULL FROM "%[4]s"
				WHERE address = $4 and id = dst_reserve_address_id)
		) a GROUP BY time
	`, timeField, schema.TradeLogsTableName, schema.WalletTableName, schema.ReserveTableName)

	var records []struct {
		Timestamp time.Time `db:"time"`
		FeeAmount float64   `db:"fee_amount"`
	}

	logger.Debugw("prepare statement", "stmt", integrationQuery)
	err = tldb.db.Select(&records, integrationQuery, fromTime, toTime, walletAddr, reserveAddr)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		logger.Debugw("no trade found")
		return nil, nil
	}
	results := make(map[uint64]float64)
	for _, record := range records {
		ts := timeutil.TimeToTimestampMs(record.Timestamp)
		results[ts] = record.FeeAmount
	}
	return results, nil
}
