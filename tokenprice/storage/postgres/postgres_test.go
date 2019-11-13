package postgres

import (
	"testing"
	"time"

	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestSaveNewTokenRate(t *testing.T) {
	db, teardown := testutil.MustNewDevelopmentDB()
	defer func() {
		require.NoError(t, teardown())
	}()
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	trdb, err := NewTokenRateDB(sugar, db)
	require.NoError(t, err)
	var (
		token    = "ETH"
		currency = "USD"
		source   = "coinbase"
		timeS    = "2019-02-06"
		rate     = 100.1
	)
	timestamp, err := time.Parse("2006-01-02", timeS)
	require.NoError(t, err)
	err = trdb.SaveTokenRate(token, currency, source, timestamp, rate)
	require.NoError(t, err)

	err = trdb.SaveTokenRate("KNC", currency, source, timestamp, rate)
	require.EqualError(t, err, ErrExists.Error())

	rateDB, err := trdb.GetTokenRate(token, currency, timestamp)
	require.NoError(t, err)
	require.Equal(t, rate, rateDB)

	_, err = trdb.GetTokenRate("KNC", currency, timestamp)
	require.EqualError(t, err, ErrNotFound.Error())
}
