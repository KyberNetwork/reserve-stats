package postgres

import (
	"testing"

	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestSaveAndGetAccountingRates(t *testing.T) {
	var (
		testData = []huobi.WithdrawHistory{
			{
				ID:         2272335,
				CreatedAt:  1525754125590,
				UpdatedAt:  1525754753403,
				Currency:   "ETH",
				Type:       "withdraw",
				Amount:     0.48957444,
				State:      "confirmed",
				Fee:        0.01,
				Address:    "f6a605cdd9b2471ffdff706f8b7665a12b862158",
				AddressTag: "",
				TxHash:     "cdef3adad017d9564e62282f5e0f0d87d72b995759f1f7f4e473137cc1b96e56",
			},
			{
				ID:         2272334,
				CreatedAt:  1525754125590,
				UpdatedAt:  1525754753403,
				Currency:   "ETH",
				Type:       "withdraw",
				Amount:     0.48957444,
				State:      "confirmed",
				Fee:        0.01,
				Address:    "f6a605cdd9b2471ffdff706f8b7665a12b862158",
				AddressTag: "",
				TxHash:     "cdef3adad017d9564e62282f5e0f0d87d72b995759f1f7f4e473137cc1b96e56",
			},
		}
	)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	sugar := logger.Sugar()

	db, teardown := testutil.MustNewDevelopmentDB()
	hdb, err := NewDB(sugar, db)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, teardown())
	}()

	err = hdb.UpdateWithdrawHistory(testData, "huobi_v1_main")
	require.NoError(t, err)

	lastID, err := hdb.GetLastIDStored("huobi_v1_main")
	require.NoError(t, err)
	assert.Equal(t, uint64(2272335), lastID)

	withdrawals, err := hdb.GetWithdrawHistory(timeutil.TimestampMsToTime(1525754125500), timeutil.TimestampMsToTime(1525754125600))
	require.NoError(t, err)
	assert.Equal(t, testData, withdrawals)

	//test does not stored duplicate records
	err = hdb.UpdateWithdrawHistory(testData, "huobi_v1_main")
	require.NoError(t, err)

	withdrawals, err = hdb.GetWithdrawHistory(timeutil.TimestampMsToTime(1525754125500), timeutil.TimestampMsToTime(1525754125600))
	require.NoError(t, err)
	assert.Equal(t, testData, withdrawals)
}
