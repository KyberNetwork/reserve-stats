package tradestorage

import (
	"testing"

	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	tradeTableTest = "binance_trade_test"
)

func newTestDB(sugar *zap.SugaredLogger) (*BinanceStorage, error) {
	_, db := testutil.MustNewDevelopmentDB()
	return NewDB(sugar, db, WithTableName(tradeTableTest))
}

func teardown(t *testing.T, storage *BinanceStorage) {
	err := storage.DeleteTable()
	assert.NoError(t, err)
	err = storage.Close()
	assert.NoError(t, err)
}

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

	binanceStorage, err := newTestDB(logger)
	assert.NoError(t, err)

	defer teardown(t, binanceStorage)

	err = binanceStorage.UpdateTradeHistory(testData)
	assert.NoError(t, err)

	// test get trade history from database
	fromTime := timeutil.TimestampMsToTime(1516439513145)
	toTime := timeutil.TimestampMsToTime(1524570040118)

	tradeHistories, err := binanceStorage.GetTradeHistory(fromTime, toTime)
	assert.NoError(t, err)
	assert.Equal(t, testData, tradeHistories)
}
