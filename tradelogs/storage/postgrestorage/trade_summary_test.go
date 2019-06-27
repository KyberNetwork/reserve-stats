package postgrestorage

import (
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTradeLogDB_GetTradeSummary(t *testing.T) {
	const (
		dbName = "test_trade_summary"
		// These params are expected to be change when export.dat changes.

		ethAmount        = 17.390905490542348
		usdAmount        = 3924.466796147473
		totalBurnFee     = 16.745037801749728
		totalTrade       = uint64(11)
		uniqueAddress    = uint64(6)
		kycedAddress     = uint64(0)
		newUniqueAddress = uint64(6)
		usdPerTrade      = 356.76970874067933
		ethPerTrade      = 1.5809914082311225
		timeStamp        = "2018-10-11T00:00:00Z"
	)

	var (
		fromTime = timeutil.TimestampMsToTime(1539216000000)
		toTime   = timeutil.TimestampMsToTime(1539254666000)
		timezone int8
	)

	tldb, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tldb.tearDown(dbName))
	}()
	err = loadTestData(tldb.db, testDataFile)
	require.NoError(t, err)

	summary, err := tldb.GetTradeSummary(fromTime, toTime, timezone)
	require.NoError(t, err)

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := summary[timeUint]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}

	require.Equal(t, ethAmount, result.ETHVolume)
	require.Equal(t, ethPerTrade, result.ETHPerTrade)
	require.Equal(t, usdAmount, result.USDAmount)
	require.Equal(t, usdPerTrade, result.USDPerTrade)
	require.Equal(t, totalTrade, result.TotalTrade)
	require.Equal(t, uniqueAddress, result.UniqueAddresses)
	require.Equal(t, totalBurnFee, result.TotalBurnFee)
	require.Equal(t, kycedAddress, result.KYCEDAddresses)
	require.Equal(t, newUniqueAddress, result.NewUniqueAddresses)

	t.Logf("%+v", result)

}
