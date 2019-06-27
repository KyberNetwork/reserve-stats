package influxstorage

import (
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/stretchr/testify/require"

	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influxstorage/cq"
)

func aggregateContryStats(is *InfluxStorage) error {
	cqs, err := tradelogcq.CreateCountryCqs(is.dbName)
	if err != nil {
		return err
	}
	for _, cq := range cqs {
		err = cq.Execute(is.influxClient, is.sugar)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestInfluxStorage_GetCountryStats(t *testing.T) {
	const (
		dbName  = "test_country_stats_volume"
		country = "DE"

		ethExpectedValue     = 0.47151091316103044
		usdExpectedValue     = 106.40210331359417
		burnFeeExpectedValue = 0.6542213920109298

		timeStamp = "2018-10-11T08:00:00Z"
	)

	var (
		fromTime = timeutil.TimestampMsToTime(1539216000000)
		toTime   = timeutil.TimestampMsToTime(1539254666000)
		timeZone = int8(-8)
	)

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	require.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)

	is, err := newTestInfluxStorage(dbName)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, is.tearDown())
	}()

	require.NoError(t, loadTestData(dbName))
	require.NoError(t, aggregateContryStats(is))
	integrationVol, err := is.GetCountryStats(country, fromTime, toTime, timeZone)
	require.NoError(t, err)
	require.Contains(t, integrationVol, timeUint)
	t.Logf("%+v", integrationVol[timeUint])

	require.Equal(t, ethExpectedValue, integrationVol[timeUint].TotalETHVolume)
	require.Equal(t, usdExpectedValue, integrationVol[timeUint].TotalUSDVolume)
	require.Equal(t, burnFeeExpectedValue, integrationVol[timeUint].TotalBurnFee)
}
