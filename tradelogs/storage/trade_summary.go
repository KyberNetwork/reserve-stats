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

func (is *InfluxStorage) GetTradeSummary(from, to uint64) (map[uint64]*common.TradeSummary, error) {
	var (
		logger = is.sugar.With(
			"func", "tradelogs/storage/InfluxStorage.GetTradeSummary",
			"from", from,
			"to", to,
		)
		timeFilter = fmt.Sprintf("(time >=%d%s AND time <= %d%s)", from, timePrecision, to, timePrecision)
	)
	cmd := fmt.Sprintf("SELECT time,eth_per_trade, new_unique_addresses, total_eth_volume, total_trade,total_usd_amount ,usd_per_trade, unique_addresses FROM trade_summary WHERE %s GROUP BY(1d) fill(0)", timeFilter)
	logger.Debugw("get trade summary", "query", cmd)

	response, err := is.queryDB(is.influxClient, cmd)
	if err != nil {
		return nil, err
	}

	logger.Debugw("got result for trade summary query", "response", response)

	if len(response) == 0 || len(response[0].Series) == 0 {
		return nil, nil
	}
	return convertQueryResultToSummary(response[0].Series[0])
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
	// number of fields in record result
	// - time
	// - token_volume
	// - eth_volume
	// - usd_volume
	if len(v) != 4 {
		return 0, nil, errors.New("value fields is empty")
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
	newUniqueAddr, err := influxdb.GetUint64FromInterface(v[2])
	if err != nil {
		return 0, nil, err
	}
	ethVolume, err := influxdb.GetFloat64FromInterface(v[3])
	if err != nil {
		return 0, nil, err
	}
	totalTrade, err := influxdb.GetUint64FromInterface(v[4])
	if err != nil {
		return 0, nil, err
	}
	usdVolume, err := influxdb.GetFloat64FromInterface(v[5])
	if err != nil {
		return 0, nil, err
	}
	usdPertrade, err := influxdb.GetFloat64FromInterface(v[6])
	uniqueAddr, err := influxdb.GetUint64FromInterface(v[7])

	return tsUint64, &common.TradeSummary{
		ETHVolume:     ethVolume,
		USDAmount:     usdVolume,
		TotalTrade:    totalTrade,
		UniqueAddress: uniqueAddr,
		NewUquAddress: newUniqueAddr,
		usdPertrade:   usdPertrade,
		ETHPerTrade:   ethPerTrade,
	}, nil
}
