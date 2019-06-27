package postgrestorage

import (
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/stretchr/testify/require"
)

func TestTradeLogDB_GetCountryStats(t *testing.T) {
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

	tldb, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, tldb.tearDown(dbName))
	}()

	require.NoError(t, loadTestData(tldb.db, testDataFile))
	integrationVol, err := tldb.GetCountryStats(country,
		fromTime, toTime, timeZone)
	require.NoError(t, err)
	require.Contains(t, integrationVol, timeUint)
	t.Logf("%+v", integrationVol[timeUint])

	require.Equal(t, ethExpectedValue, integrationVol[timeUint].TotalETHVolume)
	require.Equal(t, usdExpectedValue, integrationVol[timeUint].TotalUSDVolume)
	require.Equal(t, burnFeeExpectedValue, integrationVol[timeUint].TotalBurnFee)
}
