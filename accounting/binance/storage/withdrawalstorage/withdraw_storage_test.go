package withdrawalstorage

import (
	"testing"

	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestBinanceWithdrawStorage(t *testing.T) {
	logger := testutil.MustNewDevelopmentSugaredLogger()
	logger.Info("test binance trade storage")
	var (
		testData = []binance.WithdrawHistory{
			{
				ID:        "3c3bd6d1adb742f0bf8586bb7bb614cb",
				Amount:    "4.7",
				Address:   "0x93Dc33d2EAFcD212879d4833202F99eC453A6e18",
				Asset:     "KNC",
				TxID:      "0x102556d7ebb4e8aea93dca7c61c5926946312af98d3c38e48b062e06582 4b70f",
				ApplyTime: "2018-01-25 13:23:14",
				Status:    6,
				TxFee:     "1432",
			},
			{
				ID:        "53bb6b37ce61443f9d7fd99c21652baa",
				Amount:    "0.64",
				Address:   "0xe813dee553d09567D4873d9bd5A 4914796367082",
				Asset:     "ETH",
				TxID:      "0x679d514dafb4c8eee1fce3b00a984167bed02bf69ca278e49fa4c4a8fb2308ed",
				ApplyTime: "2018-03-26 04:09:12",
				Status:    6,
				TxFee:     "4208",
			},
			{
				ID:        "53bb6b37ce61443f9d7fd99c21652baaa",
				Amount:    "0.64",
				Address:   "0xe813dee553d09567D4873d9bd5A 4914796367082",
				Asset:     "ETH",
				TxID:      "0x679d514dafb4c8eee1fce3b00a984167bed02bf69ca278e49fa4c4a8fb2308ed",
				ApplyTime: "2018-03-26 04:09:12",
				Status:    6,
				TxFee:     "29742",
			},
		}
		expectedData = map[string][]binance.WithdrawHistory{
			"binance_1": {
				{
					ID:        "3c3bd6d1adb742f0bf8586bb7bb614cb",
					Amount:    "4.7",
					Address:   "0x93Dc33d2EAFcD212879d4833202F99eC453A6e18",
					Asset:     "KNC",
					TxID:      "0x102556d7ebb4e8aea93dca7c61c5926946312af98d3c38e48b062e06582 4b70f",
					ApplyTime: "2018-01-25 13:23:14",
					Status:    6,
					TxFee:     "1432",
				},
				{
					ID:        "53bb6b37ce61443f9d7fd99c21652baa",
					Amount:    "0.64",
					Address:   "0xe813dee553d09567D4873d9bd5A 4914796367082",
					Asset:     "ETH",
					TxID:      "0x679d514dafb4c8eee1fce3b00a984167bed02bf69ca278e49fa4c4a8fb2308ed",
					ApplyTime: "2018-03-26 04:09:12",
					Status:    6,
					TxFee:     "4208",
				},
				{
					ID:        "53bb6b37ce61443f9d7fd99c21652baaa",
					Amount:    "0.64",
					Address:   "0xe813dee553d09567D4873d9bd5A 4914796367082",
					Asset:     "ETH",
					TxID:      "0x679d514dafb4c8eee1fce3b00a984167bed02bf69ca278e49fa4c4a8fb2308ed",
					ApplyTime: "2018-03-26 04:09:12",
					Status:    6,
					TxFee:     "29742",
				},
			},
		}
	)

	db, teardown := testutil.MustNewDevelopmentDB()
	binanceStorage, err := NewDB(logger, db)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, teardown())
	}()

	_, err = binanceStorage.GetLastStoredTimestamp("binance_1")
	require.NoError(t, err)

	err = binanceStorage.UpdateWithdrawHistory(testData, "binance_1")
	assert.NoError(t, err)

	lastStoredTimestamp, err := binanceStorage.GetLastStoredTimestamp("binance_1")
	require.NoError(t, err)
	assert.Equal(t, uint64(1522037352000), timeutil.TimeToTimestampMs(lastStoredTimestamp))

	// test get trade history from database
	fromTime := timeutil.TimestampMsToTime(1516886594000)
	toTime := timeutil.TimestampMsToTime(1522037352000)

	withdrawHistory, err := binanceStorage.GetWithdrawHistory(fromTime, toTime)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, withdrawHistory)

	// test stored duplicate data
	err = binanceStorage.UpdateWithdrawHistory(testData, "binance_1")
	assert.NoError(t, err)

	withdrawHistory, err = binanceStorage.GetWithdrawHistory(fromTime, toTime)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, withdrawHistory)
}
