package storage

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
	volSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/volume"
	ethereum "github.com/ethereum/go-ethereum/common"
	influxModel "github.com/influxdata/influxdb/models"
)

const (
	tokenVolumeField = "token_volume"
	ethVolumeField   = "eth_volume"
	fiatVolumeField  = "usd_volume"
)

var (
	errCantConvert       = errors.New("cannot convert response from influxDB to pre-defined struct")
	assetMeasurementName = map[string]string{
		"h": "volume_hour",
		"d": "volume_day",
	}
	rsvMeasurementName = map[string]string{
		"h": "rsv_volume_hour",
		"d": "rsv_volume_day",
	}
)

func getMeasurementName(baseMeasurement string, timezone int8) string {
	if timezone < 0 {
		baseMeasurement = fmt.Sprintf("%s_minus%dh", baseMeasurement, (-1 * timezone))
	} else if timezone > 0 {
		baseMeasurement = fmt.Sprintf("%s_%dh", baseMeasurement, timezone)
	}
	return baseMeasurement
}

// GetReserveVolume returns the volume of a specific asset(token) from a reserve
// between a period and with desired frequency
func (is *InfluxStorage) GetReserveVolume(rsvAddr ethereum.Address, token ethereum.Address,
	fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error) {
	var (
		rsvAddrHex   = rsvAddr.Hex()
		tokenAddrHex = token.Hex()
		logger       = is.sugar.With("reserve Address", rsvAddr.Hex(), "func", "tradelogs/storage/InfluxStorage.GetReserveVolume", "token Address", token.Hex(), "from", fromTime, "to", toTime)
	)
	mName, ok := rsvMeasurementName[strings.ToLower(frequency)]
	if !ok {
		return nil, fmt.Errorf("frequency %s is not supported", frequency)
	}

	addrFilter := fmt.Sprintf("((%[1]s='%[2]s' OR %[3]s='%[4]s') AND (%[5]s='%[6]s' OR %[7]s='%[8]s'))",
		logSchema.DstAddr.String(),
		tokenAddrHex,
		logSchema.SrcAddr.String(),
		tokenAddrHex,
		logSchema.DstReserveAddr.String(),
		rsvAddrHex,
		logSchema.SrcReserveAddr.String(),
		rsvAddrHex)
	timeFilter := fmt.Sprintf("(time >='%s' AND time <= '%s')", fromTime.UTC().Format(time.RFC3339), toTime.UTC().Format(time.RFC3339))
	cmd := fmt.Sprintf("SELECT SUM(%[1]s) AS %[1]s, SUM(%[2]s) AS %[2]s, SUM(%[3]s) AS %[3]s FROM %[4]s WHERE %[5]s AND %[6]s GROUP BY time(1%[7]s) FILL(none)",
		volSchema.TokenVolume.String(),
		volSchema.ETHVolume.String(),
		volSchema.USDVolume.String(),
		mName, timeFilter, addrFilter, frequency)

	logger.Debugw("query rendered", "query", cmd)

	response, err := influxdb.QueryDB(is.influxClient, cmd, is.dbName)
	if err != nil {
		return nil, err
	}
	if len(response) == 0 || len(response[0].Series) == 0 {
		return nil, nil
	}
	return convertQueryResultToVolume(response[0].Series[0])
}

// GetAssetVolume returns the volume of a specific assset(token) between a period and with desired frequency
func (is *InfluxStorage) GetAssetVolume(token ethereum.Address, fromTime, toTime time.Time,
	frequency string) (map[uint64]*common.VolumeStats, error) {
	var (
		logger = is.sugar.With(
			"func", "tradelogs/storage/InfluxStorage.GetAssetVolume",
			"token", token.Hex(),
			"from", fromTime,
			"to", toTime,
		)
		result = make(map[uint64]*common.VolumeStats)
	)
	mName, ok := assetMeasurementName[strings.ToLower(frequency)]
	if !ok {
		return nil, fmt.Errorf("frequency %s is not supported", frequency)
	}

	var (
		tokenAddr  = token.Hex()
		timeFilter = fmt.Sprintf("(time >='%s' AND time <= '%s')", fromTime.UTC().Format(time.RFC3339), toTime.UTC().Format(time.RFC3339))
		addrFilter = fmt.Sprintf("(dst_addr='%s' OR src_addr='%s')", tokenAddr, tokenAddr)
		cmd        = fmt.Sprintf("SELECT SUM(%[1]s) AS %[1]s, SUM(%[2]s) AS %[2]s, SUM(%[3]s) AS %[3]s FROM %[4]s WHERE %[5]s AND %[6]s GROUP BY time(1%[7]s) FILL(none)",
			volSchema.TokenVolume.String(),
			volSchema.ETHVolume.String(),
			volSchema.USDVolume.String(),
			mName, timeFilter, addrFilter, frequency)
	)

	logger.Debugw("get asset volume query rendered", "query", cmd)
	response, err := influxdb.QueryDB(is.influxClient, cmd, is.dbName)

	if err != nil {
		return result, err
	}

	if len(response) == 0 || len(response[0].Series) == 0 {
		return result, nil
	}
	return convertQueryResultToVolume(response[0].Series[0])
}

func convertQueryResultToVolume(row influxModel.Row) (map[uint64]*common.VolumeStats, error) {
	result := make(map[uint64]*common.VolumeStats)
	if len(row.Values) == 0 {
		return nil, nil
	}
	idxs, err := volSchema.NewFieldsRegistrar(row.Columns)
	if err != nil {
		return result, err
	}
	for _, v := range row.Values {
		ts, vol, err := convertRowValueToVolume(v, idxs)
		if err != nil {
			return nil, err
		}
		result[ts] = vol
	}
	return result, nil
}

func convertRowValueToVolume(v []interface{}, idxs volSchema.FieldsRegistrar) (uint64, *common.VolumeStats, error) {
	// number of fields in record result
	// - time
	// - token_volume
	// - eth_volume
	// - usd_volume
	if len(v) != 4 {
		return 0, nil, errors.New("value fields is empty")
	}

	timestampString, ok := v[idxs[volSchema.Time]].(string)
	if !ok {
		return 0, nil, errCantConvert
	}
	ts, err := time.Parse(time.RFC3339, timestampString)
	if err != nil {
		return 0, nil, err
	}
	tsUint64 := timeutil.TimeToTimestampMs(ts)
	volume, err := influxdb.GetFloat64FromInterface(v[idxs[volSchema.TokenVolume]])
	if err != nil {
		return 0, nil, err
	}
	ethVolume, err := influxdb.GetFloat64FromInterface(v[idxs[volSchema.ETHVolume]])
	if err != nil {
		return 0, nil, err
	}
	usdVolume, err := influxdb.GetFloat64FromInterface(v[idxs[volSchema.USDVolume]])
	if err != nil {
		return 0, nil, err
	}
	return tsUint64, &common.VolumeStats{
		Volume:    volume,
		ETHAmount: ethVolume,
		USDAmount: usdVolume,
	}, nil
}
