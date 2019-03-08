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
		testData = huobi.WithdrawHistory{
			ID:                2272335,
			TransactionID:     0,
			CreatedAt:         1525754125590,
			UpdatedAt:         1525754753403,
			CandidateCurrency: "",
			Currency:          "ETH",
			Type:              "withdraw",
			Direction:         "",
			Amount:            0.48957444,
			State:             "confirmed",
			Fees:              0.01,
			ErrorCode:         "",
			ErrorMsg:          "",
			ToAddress:         "f6a605cdd9b2471ffdff706f8b7665a12b862158",
			ToAddrTag:         "",
			TxHash:            "cdef3adad017d9564e62282f5e0f0d87d72b995759f1f7f4e473137cc1b96e56",
			Chain:             "ETH",
			Extra:             "",
		}
	)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	sugar := logger.Sugar()

	db, err := sqlx.Connect("postgres", "host=127.0.0.1 port=5432 user=reserve_stats password=reserve_stats dbname=reserve_stats sslmode=disable")
	require.NoError(t, err)

	hdb, err := NewDB(sugar, db, "test_huobi_withdraw")
	require.NoError(t, err)

	defer func() {
		hdb.Close()
	}()

	err = hdb.UpdateWithdrawHistory(testData)
	require.NoError(t, err)
}
