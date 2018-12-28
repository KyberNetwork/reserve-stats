package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	integrationVolumeSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/integrationvolume"
	influxModel "github.com/influxdata/influxdb/models"
)

// GetIntegrationVolume returns the volume of non-kyberSwap and kyberSwap volume in ETH
// between a period
func (is *InfluxStorage) GetIntegrationVolume(fromTime, toTime time.Time) (map[uint64]*common.IntegrationVolume, error) {
	var (
		logger = is.sugar.With("func", "tradelogs/storage/InfluxStorage.GetIntegrationVolume", "from", fromTime, "to", toTime)
	)
	timeFilter := fmt.Sprintf("(time >='%s' AND time <= '%s')", fromTime.UTC().Format(time.RFC3339), toTime.UTC().Format(time.RFC3339))
	cmd := fmt.Sprintf(
		"SELECT %[1]s, %[2]s FROM %[3]s WHERE %[4]s",
		integrationVolumeSchema.KyberSwapVolume.String(),
		integrationVolumeSchema.NonKyberSwapVolume.String(),
		common.IntegrationVolumeMeasurement,
		timeFilter)
	logger.Debugw("query rendered", "query", cmd)

	response, err := influxdb.QueryDB(is.influxClient, cmd, is.dbName)
	if err != nil {
		return nil, err
	}
	logger.Debugw("data resp", "response", response)

	if len(response) == 0 || len(response[0].Series) == 0 {
		return nil, nil
	}
	return convertQueryResultToIntegrationVolume(response[0].Series[0])
}

func convertQueryResultToIntegrationVolume(row influxModel.Row) (map[uint64]*common.IntegrationVolume, error) {
	result := make(map[uint64]*common.IntegrationVolume)
	if len(row.Values) == 0 {
		return nil, nil
	}
	idxs, err := integrationVolumeSchema.NewFieldsRegistrar(row.Columns)
	if err != nil {
		return result, err
	}
	for _, v := range row.Values {
		ts, vol, err := convertRowValueToIntegrationVolume(v, idxs)
		if err != nil {
			return nil, err
		}
		result[ts] = vol
	}
	return result, nil
}

func convertRowValueToIntegrationVolume(v []interface{}, idxs map[integrationVolumeSchema.FieldName]int) (uint64, *common.IntegrationVolume, error) {
	// number of fields in record result
	// - time
	// - kyber_swap_volume
	// - non_kyber_swap_volume
	if len(v) != 3 {
		return 0, nil, errors.New("value fields is empty")
	}

	timestampString, ok := v[idxs[integrationVolumeSchema.Time]].(string)
	if !ok {
		return 0, nil, errCantConvert
	}
	ts, err := time.Parse(time.RFC3339, timestampString)
	if err != nil {
		return 0, nil, err
	}
	tsUint64 := timeutil.TimeToTimestampMs(ts)

	ksvolume, err := influxdb.GetFloat64FromInterface(v[idxs[integrationVolumeSchema.KyberSwapVolume]])
	if err != nil {
		return 0, nil, err
	}
	nonksVolume, err := influxdb.GetFloat64FromInterface(v[idxs[integrationVolumeSchema.NonKyberSwapVolume]])
	if err != nil {
		return 0, nil, err
	}
	return tsUint64, &common.IntegrationVolume{
		KyberSwapVolume:    ksvolume,
		NonKyberSwapVolume: nonksVolume,
	}, nil
}
