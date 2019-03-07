package postgres

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
)

func TestSaveAndGetAccountingRates(t *testing.T) {
	var (
		testData = huobi.TradeHistory{
			ID:              15584072551,
			Symbol:          "cmtetsh",
			AccountID:       3375841,
			Amount:          "6000.000",
			Price:           "0.00045",
			CreateAt:        1540793585678,
			Type:            "buy-limit",
			FieldAmount:     "6000.000",
			FieldCashAmount: "2.73336",
			FieldFee:        "12.00000",
			FinishedAt:      1540796135588,
			UserID:          0,
			Source:          "web",
			State:           "filled",
			CanceledAt:      0,
		}
	)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	sugar := logger.Sugar()

	db, err := sqlx.Connect("postgres", "host=127.0.0.1 port=5432 user=reserve_stats password=reserve_stats dbname=reserve_stats sslmode=disable")
	require.NoError(t, err)

	hdb, err := NewDB(sugar, db, "test_huobi_trades")
	require.NoError(t, err)

	defer func() {
		hdb.Close()
	}()

	err = hdb.UpdateTradeHistory(testData)
	require.NoError(t, err)
}
