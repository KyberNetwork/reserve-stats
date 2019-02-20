package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	walletFeeVolumeSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/walletfee_volume"
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
	fromTime, toTime time.Time, timezone int8) (map[uint64]float64, error) {
	var (
		err         error
		measurement string
		result      = make(map[uint64]float64)
	)

	logger := is.sugar.With("reserveAddr", reserveAddr, "walletAddr", walletAddr, "freq", freq,
		"fromTime", fromTime, "toTime", toTime, "timezone", timezone)

	switch strings.ToLower(freq) {
	case day:
		measurement = common.WalletFeeVolumeMeasurementDay
	case hour:
		measurement = common.WalletFeeVolumeMeasurementHour
	default:
		return nil, fmt.Errorf("Frequency %s is not supported", freq)
	}

	// in cq we will add timezone as time offset interval
	q := fmt.Sprintf(`
		SELECT sum_amount FROM "%[1]s"
		WHERE (src_rsv_addr = '%[2]s' OR dst_rsv_addr= '%[2]s') AND wallet_addr = '%[3]s'
		AND time >= '%[4]s' AND time <= '%[5]s' 
	`, measurement, reserveAddr, walletAddr,
		fromTime.UTC().Format(time.RFC3339), toTime.UTC().Format(time.RFC3339), freq)

	logger.Debugw("GetAggregatedWalletFee", "query", q)
	res, err := influxdb.QueryDB(is.influxClient, q, is.dbName)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot query wallet fee from influx: %s", err.Error()))
		return result, err
	}

	if len(res[0].Series) == 0 {
		logger.Debug("the aggregated measurement is empty")
		return result, nil
	}

	for _, row := range res[0].Series[0].Values {
		idxs, err := walletFeeVolumeSchema.NewFieldsRegistrar(res[0].Series[0].Columns)
		if err != nil {
			return nil, err
		}
		ts, amount, err := is.rowToAggregatedFee(row, idxs)
		if err != nil {
			return nil, err
		}
		key := timeutil.TimeToTimestampMs(ts)
		//if the result is already there, that mean there was either src/dst wallet fee
		if _, avail := result[key]; avail {
			result[key] += amount
		} else {
			result[key] = amount
		}
	}

	return result, err
}
