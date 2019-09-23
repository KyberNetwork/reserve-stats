package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestTradeLogDB_GetWalletStats(t *testing.T) {
	const (
		dbName = "test_wallet_stats_volume"
		// These params are expected to be change when export.dat changes.

		timeStamp     = "2018-10-11T00:00:00Z"
		walletAddress = "0xDECAF9CD2367cdbb726E904cD6397eDFcAe6068D"
	)

	var (
		fromTime = timeutil.TimestampMsToTime(1539216000000)
		toTime   = timeutil.TimestampMsToTime(1539254666000)
	)

	tldb, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, tldb.tearDown(dbName))
	}()

	require.NoError(t, loadTestData(tldb.db, testDataFile))
	integrationVol, err := tldb.GetWalletStats(fromTime, toTime, walletAddress, 0)

	require.NoError(t, err)

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := integrationVol[timeUint]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there tldb none", timeUnix.Format(time.RFC3339))
	}

	t.Logf("%+v", result)
}
