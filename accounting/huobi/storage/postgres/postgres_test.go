package postgres

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestSaveAndGetAccountingRates(t *testing.T) {
	var (
		testData = []huobi.TradeHistory{
			{
				ID:              15584072551,
				Symbol:          "cmtetsh",
				AccountID:       3375841,
				Amount:          "6000.000",
				Price:           "0.00045",
				CreatedAt:       1540793585678,
				Type:            "buy-limit",
				FieldAmount:     "6000.000",
				FieldCashAmount: "2.73336",
				FieldFees:       "12.00000",
				FinishedAt:      1540796135588,
				UserID:          0,
				Source:          "web",
				State:           "filled",
				CanceledAt:      0,
			},
		}
	)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	sugar := logger.Sugar()

	db, err := sqlx.Connect("postgres", "host=127.0.0.1 port=5432 user=reserve_stats password=reserve_stats dbname=reserve_stats sslmode=disable")
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		td := testData[0]
		td.ID++
		td.CreatedAt += 100
		testData = append(testData, td)
	}

	hdb, err := NewDB(sugar, db, WithTradeTableName("test_huobi_trades"))
	require.NoError(t, err)

	defer func() {
		if err := hdb.TearDown(); err != nil {
			t.Fatalf("teardown error : %v", err)
		}
		hdb.Close()
	}()

	err = hdb.UpdateTradeHistory(testData)
	require.NoError(t, err)
	data, err := hdb.GetTradeHistory(timeutil.TimestampMsToTime(1540793585600), timeutil.TimestampMsToTime(1540793585699))
	require.NoError(t, err)
	assert.Equal(t, len(data), 1)
	assert.Equal(t, testData[0].FieldAmount, data[0].FieldAmount)
}
