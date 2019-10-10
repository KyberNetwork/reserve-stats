package influx

import (
	"fmt"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/tradelog"
	volSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/volume"
)

//GetUserVolume return volume of an address from time to time by a frequency
func (is *Storage) GetUserVolume(userAddress ethereum.Address, from, to time.Time, freq string) (map[uint64]common.UserVolume, error) {
	var (
		userAddr    = userAddress.Hex()
		measurement string
	)

	result := map[uint64]common.UserVolume{}

	logger := is.sugar.With("user address", userAddr, "freq", freq,
		"fromTime", from, "toTime", to)

	switch strings.ToLower(freq) {
	case day:
		measurement = common.UserVolumeDayMeasurementName
	case hour:
		measurement = common.UserVolumeHourMeasurementName
	}

	q := fmt.Sprintf(`
		SELECT %[1]s, %[2]s from "%[3]s"
		WHERE %[4]s = '%[5]s' AND time >= '%[6]s' AND time <= '%[7]s'`,
		volSchema.ETHVolume.String(),
		volSchema.USDVolume.String(),
		measurement,
		logSchema.UserAddr.String(),
		userAddr,
		from.UTC().Format(time.RFC3339),
		to.UTC().Format(time.RFC3339))

	logger.Debug(q)

	res, err := influxdb.QueryDB(is.influxClient, q, is.dbName)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot query user volume from influx: %s", err.Error()))
		return result, err
	}

	if len(res[0].Series) == 0 {
		return result, nil
	}
	idxs, err := volSchema.NewFieldsRegistrar(res[0].Series[0].Columns)
	if err != nil {
		return result, err
	}
	for _, row := range res[0].Series[0].Values {
		ts, ethAmount, usdAmount, err := is.rowToAggregatedUserVolume(row, idxs)
		if err != nil {
			return nil, err
		}
		key := timeutil.TimeToTimestampMs(ts)
		result[key] = common.UserVolume{
			ETHAmount: ethAmount,
			USDAmount: usdAmount,
		}
	}

	return result, err
}
