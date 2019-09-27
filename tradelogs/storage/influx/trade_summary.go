package influx

import (
	"errors"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	tradeSumSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/trade_summary"
	influxModel "github.com/influxdata/influxdb/models"
)

const (
	measurementName = "trade_summary"
)

// GetTradeSummary return an incompleted tradeSummary for the specified time periods
func (is *Storage) GetTradeSummary(from, to time.Time, timezone int8) (map[uint64]*common.TradeSummary, error) {
	var (
		logger = is.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"from", from,
			"to", to,
		)
		timeFilter = fmt.Sprintf("(time >='%s' AND time <= '%s')",
			from.UTC().Format(time.RFC3339),
			to.UTC().Format(time.RFC3339))
		results map[uint64]*common.TradeSummary
	)
	measurement := getMeasurementName(measurementName, timezone)

	tradeLogQuery := fmt.Sprintf(
		"SELECT %[1]s,%[2]s,%[3]s,%[4]s,%[5]s,%[6]s,%[7]s,%[8]s,%[9]s FROM %[10]s WHERE %[11]s",
		tradeSumSchema.Time.String(),
		tradeSumSchema.ETHPerTrade.String(),
		tradeSumSchema.TotalETHVolume.String(),
		tradeSumSchema.TotalTrade.String(),
		tradeSumSchema.TotalUSDAmount.String(),
		tradeSumSchema.USDPerTrade.String(),
		tradeSumSchema.UniqueAddresses.String(),
		tradeSumSchema.NewUniqueAddresses.String(),
		tradeSumSchema.KYCedAddresses.String(),
		measurement,
		timeFilter)
	logger.Debugw("getting trade summary", "query", tradeLogQuery)

	response, err := influxdb.QueryDB(is.influxClient, tradeLogQuery, is.dbName)
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

	burnFeeMName := getMeasurementName("burn_fee_summary", timezone)

	burnFeeQuery := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s ",
		tradeSumSchema.TotalBurnFee.String(),
		burnFeeMName,
		timeFilter)
	logger.Debugw("getting total burn fee", "query", burnFeeQuery)

	if response, err = influxdb.QueryDB(is.influxClient, burnFeeQuery, is.dbName); err != nil {
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
	idxs, err := tradeSumSchema.NewFieldsRegistrar(row.Columns)
	if err != nil {
		return nil, err
	}
	for _, v := range row.Values {
		ts, vol, err := convertRowValueToSummary(v, idxs)
		if err != nil {
			return nil, err
		}
		result[ts] = vol
	}
	return result, nil
}

func convertRowValueToSummary(v []interface{}, idxs map[tradeSumSchema.FieldName]int) (uint64, *common.TradeSummary, error) {
	if len(v) != 9 {
		return 0, nil, errors.New("value fields is invalid in len")
	}
	timestampString, ok := v[idxs[tradeSumSchema.Time]].(string)
	if !ok {
		return 0, nil, errCantConvert
	}
	ts, err := time.Parse(time.RFC3339, timestampString)
	if err != nil {
		return 0, nil, err
	}
	tsUint64 := timeutil.TimeToTimestampMs(ts)
	ethPerTrade, err := influxdb.GetFloat64FromInterface(v[idxs[tradeSumSchema.ETHPerTrade]])
	if err != nil {
		return 0, nil, err
	}
	ethVolume, err := influxdb.GetFloat64FromInterface(v[idxs[tradeSumSchema.TotalETHVolume]])
	if err != nil {
		return 0, nil, err
	}
	totalTrade, err := influxdb.GetUint64FromInterface(v[idxs[tradeSumSchema.TotalTrade]])
	if err != nil {
		return 0, nil, err
	}
	usdVolume, err := influxdb.GetFloat64FromInterface(v[idxs[tradeSumSchema.TotalUSDAmount]])
	if err != nil {
		return 0, nil, err
	}
	usdPerTrade, err := influxdb.GetFloat64FromInterface(v[idxs[tradeSumSchema.USDPerTrade]])
	if err != nil {
		return 0, nil, err
	}
	uniqueAddr, err := influxdb.GetUint64FromInterface(v[idxs[tradeSumSchema.UniqueAddresses]])
	if err != nil {
		return 0, nil, err
	}
	newUniqueAddress, err := influxdb.GetUint64FromInterface(v[idxs[tradeSumSchema.NewUniqueAddresses]])
	if err != nil {
		return 0, nil, err
	}
	kyced, err := influxdb.GetUint64FromInterface(v[idxs[tradeSumSchema.KYCedAddresses]])
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
