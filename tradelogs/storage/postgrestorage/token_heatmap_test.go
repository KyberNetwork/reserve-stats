package postgrestorage

import (
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestTradeLogDB_GetTokenHeatmap(t *testing.T) {
	const (
		dbName = "test_heat_map_volume"
		// These params are expected to be change when export.dat changes.
		accessAddress = "0x514910771AF9Ca656af840dff83E8264EcF986CA"
		country       = "DE"

		ethExpectedValue   = 0.47151091316103044
		tokenExpectedValue = 295.71295676047947
		usdExpectedValue   = 106.40210331359417
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

	require.Equal(t, ethExpectedValue, integrationVol[country].TotalETHValue)
	require.Equal(t, tokenExpectedValue, integrationVol[country].TotalTokenValue)
	require.Equal(t, usdExpectedValue, integrationVol[country].TotalFiatValue)
}
