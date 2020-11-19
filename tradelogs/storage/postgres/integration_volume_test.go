package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestTradeLogDB_GetIntegrationVolume(t *testing.T) {
	t.Skip()
	const (
		dbName = "test_integration_volume"
		// These params are expected to be change when export.dat changes.
		timeStamp = "2018-10-11T00:00:00Z"
	)

	fromTime := timeutil.TimestampMsToTime(1539216000000)
	toTime := timeutil.TimestampMsToTime(1539254666000)
	// set up load demo db
	tldb, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tldb.tearDown(dbName))
	}()
	require.NoError(t, loadTestData(tldb.db, testDataFile))

	integrationVol, err := tldb.GetIntegrationVolume(fromTime, toTime)
	require.NoError(t, err)
	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := integrationVol[timeUint]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}
	// change this number from 5.3909054905423455 because eth_amount is double for non-burnable trade
	assert.Equal(t, float64(5.392022969243365), result.KyberSwapVolume)
	assert.Equal(t, float64(12), result.NonKyberSwapVolume)
}
