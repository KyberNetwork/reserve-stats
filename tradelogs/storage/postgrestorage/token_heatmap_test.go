package postgrestorage

import (
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTradeLogDB_GetTokenHeatmap(t *testing.T) {
	const (
		dbName = "test_wallet_stats_volume"
		// These params are expected to be change when export.dat changes.
		accessAddress = "0x23Ccc43365D9dD3882eab88F43d515208f832430"
		country       = "NULL"
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
	integrationVol, err := tldb.GetTokenHeatmap(ethereum.HexToAddress(accessAddress),
		fromTime, toTime, 0)
	require.NoError(t, err)
	require.Contains(t, integrationVol, country)
	t.Logf("%+v", integrationVol[country])
}
