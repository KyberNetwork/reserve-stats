package postgrestorage

import (
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestTradeLogDB_GetAssetVolume(t *testing.T) {
	const (
		dbName = "test_volume"
		// These params are expected to be change when export.dat changes.
		fromTime    = 1539248043000
		toTime      = 1539248666000
		ethAmount   = 238.33849929550047
		totalVolume = 1.056174642648189277
		freq        = "h"
		timeStamp   = "2018-10-11T09:00:00Z"
		ethAddress  = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	)

	tldb, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	from := timeutil.TimestampMsToTime(fromTime)
	to := timeutil.TimestampMsToTime(toTime)
	defer func() {
		require.NoError(t, tldb.tearDown(dbName))
	}()
	require.NoError(t, loadTestData(tldb.db, testDataFile))

	volume, err := tldb.GetAssetVolume(ethereum.HexToAddress(ethAddress), from, to, freq)
	require.NoError(t, err)
	t.Logf("Volume result %v", volume)
	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := volume[timeUint]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}
	require.Equal(t, ethAmount, result.USDAmount)
	require.Equal(t, totalVolume, result.Volume)
}

func TestTradeLogDB_GetReserveVolume(t *testing.T) {
	const (
		dbName = "test_rsv_volume"

		// These params are expected to be change when export.dat changes.
		fromTime    = 1539248043000
		toTime      = 1539248666000
		ethAmount   = 227.05539848662738
		totalVolume = 1.006174642648189232
		freq        = "h"
		timeStamp   = "2018-10-11T09:00:00Z"
		rsvAddrStr  = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
		ethAddress  = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	)

	tldb, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)

	from := timeutil.TimestampMsToTime(fromTime)
	to := timeutil.TimestampMsToTime(toTime)
	defer func() {
		require.NoError(t, tldb.tearDown(dbName))
	}()
	require.NoError(t, loadTestData(tldb.db, testDataFile))

	volume, err := tldb.GetReserveVolume(ethereum.HexToAddress(rsvAddrStr), ethereum.HexToAddress(ethAddress), from, to, freq)
	t.Logf("Volume result %v", volume)
	if err != nil {
		t.Fatal(err)
	}
	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	if err != nil {
		t.Fatal(err)
	}
	result, ok := volume[timeutil.TimeToTimestampMs(timeUnix)]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}

	require.Equal(t, ethAmount, result.USDAmount)
	require.Equal(t, totalVolume, result.Volume)
}
