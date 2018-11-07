package storage

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	day  = "d"
	hour = "h"
)

//GetAggregatedWalletFee return wallet fee follow by
//provided reserve address, wallet address, from time, to time
//frequency (hour, day) and timezone
//daily_wallet_fee and hourly_wallet_fee measurement is calculate by CQ
func (is *InfluxStorage) GetAggregatedWalletFee(reserveAddr, walletAddr, freq string,
	fromTime, toTime time.Time, timezone int64) (map[uint64]float64, error) {
	var (
		result      map[uint64]float64
		err         error
		measurement string
	)

	logger := is.sugar.With("reserveAddr", reserveAddr, "walletAddr", walletAddr, "freq", freq,
		"fromTime", fromTime, "toTime", toTime, "timezone", timezone)

	switch freq {
	case day:
		measurement = "wallet_fee_day"
	case hour:
		measurement = "wallet_fee_hour"
	}

	// in cq we will add timezone as time offset interval
	q := fmt.Sprintf(`
		SELECT sum_amount from %s
		WHERE reserve_addr = '%s' AND wallet_addr = %s
		AND time >= '%s' AND time <= '%s' 
	`, measurement, reserveAddr, walletAddr,
		fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339))

	res, err := is.queryDB(is.influxClient, q)
	if err != nil {
		logger.Error("cannot query wallet fee from influx")
		return result, err
	}

	if len(res[0].Series) == 0 {
		logger.Debug("the aggregated measurement is empty")
		return result, nil
	}

	for _, row := range res[0].Series[0].Values {
		ts, amount, err := is.rowToAggregatedFee(row)
		if err != nil {
			return nil, err
		}
		key := timeutil.TimeToTimestampMs(ts)
		result[key] = amount
	}

	return result, err
}
