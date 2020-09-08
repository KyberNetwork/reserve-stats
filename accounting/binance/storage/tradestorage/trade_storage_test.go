package tradestorage

import (
	"testing"

	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestBinanceTradeStorage(t *testing.T) {
	logger := testutil.MustNewDevelopmentSugaredLogger()
	logger.Info("test binance trade storage")
	var testData = []binance.TradeHistory{
		{
			Symbol:          "KNCETH",
			ID:              574401,
			OrderID:         5883434,
			Price:           "0.00404000",
			Quantity:        "50.00000000",
			Commission:      "0.05000000",
			CommissionAsset: "KNC",
			Time:            1516439513145,
			IsBuyer:         true,
			IsMaker:         true,
			IsBestMatch:     true,
		},
		{
			Symbol:          "KNCETH",
			ID:              961633,
			OrderID:         11488279,
			Price:           "0.00319130",
			Quantity:        "49.00000000",
			Commission:      "0.00015637",
			CommissionAsset: "ETH",
			Time:            1524570040118,
			IsBuyer:         false,
			IsMaker:         true,
			IsBestMatch:     true,
		},
	}

	db, teardown := testutil.MustNewDevelopmentDB()
	binanceStorage, err := NewDB(logger, db)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, teardown())
	}()

	ts, err := binanceStorage.GetLastStoredID("KNCETH", "binance_1")
	require.NoError(t, err)
	assert.Zero(t, ts)

	err = binanceStorage.UpdateTradeHistory(testData, "binance_1")
	assert.NoError(t, err)

	lastStoredID, err := binanceStorage.GetLastStoredID("KNCETH", "binance_1")
	require.NoError(t, err)
	assert.Equal(t, uint64(961633), lastStoredID)

	// test get trade history from database
	fromTime := timeutil.TimestampMsToTime(1516439513145)
	toTime := timeutil.TimestampMsToTime(1524570040118)

	tradeHistories, err := binanceStorage.GetTradeHistory(fromTime, toTime)
	require.NoError(t, err)
	assert.Equal(t, testData, tradeHistories)

	ts, err = binanceStorage.GetLastStoredID("KNCETH", "binance_1")
	require.NoError(t, err)
	assert.NotZero(t, ts)

	// test stored duplicate data
	err = binanceStorage.UpdateTradeHistory(testData, "binance_1")
	assert.NoError(t, err)

	tradeHistories, err = binanceStorage.GetTradeHistory(fromTime, toTime)
	require.NoError(t, err)
	assert.Equal(t, testData, tradeHistories)
}
