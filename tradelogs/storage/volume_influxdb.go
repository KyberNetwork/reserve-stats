package storage

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	influxModel "github.com/influxdata/influxdb/models"
)

const (
	tokenVolumeField = "token_volume"
	ethVolumeField   = "eth_volume"
	fiatVolumeField  = "usd_volume"
)

var (
	errCantConvert  = errors.New("cannot convert response from influxDB to pre-defined struct")
	measurementName = map[string]string{
		"h": "volume_hour",
		"d": "volume_day",
	}
)

// GetAssetVolume returns the volume of a specific assset(token) between a period and with desired frequency
func (is *InfluxStorage) GetAssetVolume(token core.Token, fromTime, toTime uint64, frequency string) (map[time.Time]*common.VolumeStats, error) {
	var (
		logger = is.sugar.With("asset Volume", token.Address, "from", fromTime, "to", toTime)
	)
	mName, ok := measurementName[strings.ToLower(frequency)]
	if !ok {
		return nil, fmt.Errorf("frequency %s is not supported", frequency)
	}
	var (
		tokenAddr  = ethereum.HexToAddress(token.Address).Hex()
		timeFilter = fmt.Sprintf("(time >=%d%s AND time <= %d%s)", fromTime, timePrecision, toTime, timePrecision)
		addrFilter = fmt.Sprintf("(dst_addr='%s' OR src_addr='%s')", tokenAddr, tokenAddr)
		cmd        = fmt.Sprintf("SELECT SUM(token_volume) as %s, SUM(eth_volume) as %s, sum(usd_volume) as %s FROM %s WHERE %s AND %s GROUP BY time(1%s) fill(0)", tokenVolumeField, ethVolumeField, fiatVolumeField, mName, timeFilter, addrFilter, frequency)
	)

	logger.Debugw("query CMD rendered is:", "query", cmd)
	response, err := is.queryDB(is.influxClient, cmd)
	logger.Debugw("response is :", "response", response)

	if err != nil {
		return nil, err
	}
	if len(response) == 0 || len(response[0].Series) == 0 {
		return nil, nil
	}
	return convertQueryResultToVolume(response[0].Series[0])
}

func convertQueryResultToVolume(row influxModel.Row) (map[time.Time]*common.VolumeStats, error) {
	result := make(map[time.Time]*common.VolumeStats)
	if len(row.Values) == 0 {
		return nil, nil
	}
	for _, v := range row.Values {
		ts, vol, err := convertRowValueToVolume(v)
		if err != nil {
			return nil, err
		}
		result[ts] = vol
	}
	return result, nil
}

func convertRowValueToVolume(v []interface{}) (time.Time, *common.VolumeStats, error) {
	if len(v) == 0 {
		return time.Time{}, nil, errors.New("value fields is empty")
	}
	timestampString, ok := v[0].(string)
	if !ok {
		return time.Time{}, nil, errCantConvert
	}
	ts, err := time.Parse(time.RFC3339, timestampString)
	if err != nil {
		return time.Time{}, nil, err
	}
	volume, err := influxdb.GetFloat64FromInterface(v[1])
	if err != nil {
		return time.Time{}, nil, err
	}
	ethVolume, err := influxdb.GetFloat64FromInterface(v[2])
	if err != nil {
		return time.Time{}, nil, err
	}
	usdVolume, err := influxdb.GetFloat64FromInterface(v[3])
	if err != nil {
		return time.Time{}, nil, err
	}
	return ts, &common.VolumeStats{
		Volume:    volume,
		ETHAmount: ethVolume,
		USDAmount: usdVolume,
	}, nil
}
