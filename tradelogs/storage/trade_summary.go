package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	influxModel "github.com/influxdata/influxdb/models"
)

const (
	measurementName = "trade_summary"
)

// GetTradeSummary return an incompleted tradeSummary for the specified time periods
func (is *InfluxStorage) GetTradeSummary(from, to time.Time, timezone int64) (map[uint64]*common.TradeSummary, error) {
	var (
		logger = is.sugar.With(
			"func", "tradelogs/storage/InfluxStorage.GetTradeSummary",
			"from", from,
			"to", to,
		)
		timeFilter = fmt.Sprintf("(time >='%s' AND time <= '%s')",
			from.UTC().Format(time.RFC3339),
			to.UTC().Format(time.RFC3339))
		results = make(map[uint64]*common.TradeSummary)
	)

	measurement := getMeasurementName(measurementName, timezone)
	tradeLogQuery := fmt.Sprintf(`SELECT time,eth_per_trade,total_eth_volume,total_trade,total_usd_amount,
		usd_per_trade,unique_addresses,new_unique_addresses,kyced FROM %s WHERE %s`, measurement, timeFilter)
	logger.Debugw("getting trade summary", "query", tradeLogQuery)

	response, err := is.queryDB(is.influxClient, tradeLogQuery)
	if err != nil {
		return nil, err
	}

	logger.Debugw("got result for trade summary query", "response", response)
	if len(response) == 0 || len(response[0].Series) == 0 {
		result := make(map[uint64]*common.TradeSummary)
		return result, nil
	}

	logger.Debugw("Number of records returned", "nRecord", len(response[0].Series[0].Values))
	if results, err = convertQueryResultToSummary(response[0].Series[0]); err != nil {
		return nil, err
	}

	burnFeeQuery := fmt.Sprintf("SELECT total_burn_fee FROM burn_fee_summary WHERE %s ", timeFilter)
	logger.Debugw("getting total burn fee", "query", burnFeeQuery)

	if response, err = is.queryDB(is.influxClient, burnFeeQuery); err != nil {
		return nil, err
	}

	logger.Debugw("got result for total burn fee query", "response", response)
	if len(response) == 0 || len(response[0].Series) == 0 {
		result := make(map[uint64]*common.TradeSummary)
		return result, nil
	}

	for _, val := range response[0].Series[0].Values {
		ts, totalBurnFee, fErr := convertQueryValueToTotalBurnFee(val)
		if fErr != nil {
			return nil, fErr
		}

		summary, ok := results[ts]
		if !ok {
			logger.Warnw("no summary result found", "ts", ts)
			break
		}

		logger.Debugw("insert total burn fee for trade summary", "ts", ts)
		summary.TotalBurnFee = totalBurnFee
		results[ts] = summary
	}

	return results, nil
}

func convertQueryValueToTotalBurnFee(val []interface{}) (uint64, float64, error) {
	if len(val) != 2 {
		return 0, 0, fmt.Errorf("wrong number of val in total burn fee result: %d", len(val))
	}
	timestampString, ok := val[0].(string)
	if !ok {
		return 0, 0, errCantConvert
	}

	ts, err := time.Parse(time.RFC3339, timestampString)
	if err != nil {
		return 0, 0, err
	}

	totalBurnFee, err := influxdb.GetFloat64FromInterface(val[1])
	if err != nil {
		return 0, 0, err
	}

	return timeutil.TimeToTimestampMs(ts), totalBurnFee, nil
}

func convertQueryResultToSummary(row influxModel.Row) (map[uint64]*common.TradeSummary, error) {
	result := make(map[uint64]*common.TradeSummary)
	if len(row.Values) == 0 {
		return nil, nil
	}
	for _, v := range row.Values {
		ts, vol, err := convertRowValueToSummary(v)
		if err != nil {
			return nil, err
		}
		result[ts] = vol
	}
	return result, nil
}

func convertRowValueToSummary(v []interface{}) (uint64, *common.TradeSummary, error) {
	if len(v) != 9 {
		return 0, nil, errors.New("value fields is invalid in len")
	}
	timestampString, ok := v[0].(string)
	if !ok {
		return 0, nil, errCantConvert
	}
	ts, err := time.Parse(time.RFC3339, timestampString)
	if err != nil {
		return 0, nil, err
	}
	tsUint64 := timeutil.TimeToTimestampMs(ts)
	ethPerTrade, err := influxdb.GetFloat64FromInterface(v[1])
	if err != nil {
		return 0, nil, err
	}
	ethVolume, err := influxdb.GetFloat64FromInterface(v[2])
	if err != nil {
		return 0, nil, err
	}
	totalTrade, err := influxdb.GetUint64FromInterface(v[3])
	if err != nil {
		return 0, nil, err
	}
	usdVolume, err := influxdb.GetFloat64FromInterface(v[4])
	if err != nil {
		return 0, nil, err
	}
	usdPerTrade, err := influxdb.GetFloat64FromInterface(v[5])
	if err != nil {
		return 0, nil, err
	}
	uniqueAddr, err := influxdb.GetUint64FromInterface(v[6])
	if err != nil {
		return 0, nil, err
	}
	newUniqueAddress, err := influxdb.GetUint64FromInterface(v[7])
	if err != nil {
		return 0, nil, err
	}
	kyced, err := influxdb.GetUint64FromInterface(v[8])
	if err != nil {
		return 0, nil, err
	}
	return tsUint64, &common.TradeSummary{
		ETHVolume:          ethVolume,
		USDAmount:          usdVolume,
		TotalTrade:         totalTrade,
		UniqueAddresses:    uniqueAddr,
		USDPerTrade:        usdPerTrade,
		ETHPerTrade:        ethPerTrade,
		NewUniqueAddresses: newUniqueAddress,
		KYCEDAddresses:     kyced,
	}, nil
}
