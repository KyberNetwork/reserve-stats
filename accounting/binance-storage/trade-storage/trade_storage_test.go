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
	return NewDB(sugar, db, tradeTableTest)
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

func TestBinanceWithdrawStorage(t *testing.T) {
	logger := testutil.MustNewDevelopmentSugaredLogger()
	logger.Info("test binance trade storage")
	var testData = []binance.WithdrawHistory{
		{
			ID:        "3c3bd6d1ad b742f0bf8586bb7bb614cb",
			Amount:    4.7,
			Address:   "0x93Dc33d2EAFcD212879d4833202F99eC453A6e18",
			Asset:     "KNC",
			TxID:      "0x102556d7ebb4e8aea93dca7c61c5926946312af98d3c38e48b062e06582 4b70f",
			ApplyTime: 1516886594000,
			Status:    6,
		},
		{
			ID:        "53bb6b37ce61443f9d7fd99c21652baa",
			Amount:    0.64,
			Address:   "0xe813dee553d09567D4873d9bd5A 4914796367082",
			Asset:     "ETH",
			TxID:      "0x679d514dafb4c8eee1fce3b00a984167bed02bf69ca278e49fa4c4a8fb2308ed",
			ApplyTime: 1522037352000,
			Status:    6,
		},
	}

	binanceStorage, err := newTestDB(logger)
	assert.NoError(t, err)

	defer teardown(t, binanceStorage)

	err = binanceStorage.UpdateWithdrawHistory(testData)
	assert.NoError(t, err)

	// test get trade history from database
	// fromTime := timeutil.TimestampMsToTime(1516439513145)
	// toTime := timeutil.TimestampMsToTime(1524570040118)

	// tradeHistories, err := binanceStorage.GetWithdrawHistory(fromTime, toTime)
	// assert.NoError(t, err)
	// assert.Equal(t, testData, tradeHistories)
}
