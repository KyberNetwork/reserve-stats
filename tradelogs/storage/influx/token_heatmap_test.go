package influx

import (
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/stretchr/testify/require"

	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/cq"
	ethereum "github.com/ethereum/go-ethereum/common"
)

func aggregateHeatMap(is *Storage) error {
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

func TestInfluxStorage_GetTokenHeatmap(t *testing.T) {
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

	is, err := newTestInfluxStorage(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, is.tearDown())
	}()
	require.NoError(t, loadTestData(dbName))
	require.NoError(t, aggregateHeatMap(is))

	integrationVol, err := is.GetTokenHeatmap(ethereum.HexToAddress(accessAddress),
		fromTime, toTime, 0)
	require.NoError(t, err)
	require.Contains(t, integrationVol, country)
	t.Logf("%+v", integrationVol[country])

	require.Equal(t, ethExpectedValue, integrationVol[country].TotalETHValue)
	require.Equal(t, tokenExpectedValue, integrationVol[country].TotalTokenValue)
	require.Equal(t, usdExpectedValue, integrationVol[country].TotalFiatValue)

}
