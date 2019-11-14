package postgres

import (
	"testing"
	"time"

	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestSaveNewTokenPrice(t *testing.T) {
	db, teardown := testutil.MustNewDevelopmentDB()
	defer func() {
		require.NoError(t, teardown())
	}()
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	trdb, err := NewTokenPriceDB(sugar, db)
	require.NoError(t, err)
	var (
		token    = "ETH"
		currency = "USD"
		source   = "coinbase"
		timeS    = "2019-02-06"
		price     = 100.1
	)
	timestamp, err := time.Parse("2006-01-02", timeS)
	require.NoError(t, err)
	err = trdb.SaveTokenPrice(token, currency, source, timestamp, price)
	require.NoError(t, err)

	err = trdb.SaveTokenPrice(token, currency, source, timestamp, 1000)
	require.EqualError(t, err, ErrExists.Error())

	priceDB, err := trdb.GetTokenPrice(token, currency, timestamp)
	require.NoError(t, err)
	require.Equal(t, price, priceDB)

	_, err = trdb.GetTokenPrice("KNC", currency, timestamp)
	require.EqualError(t, err, ErrNotFound.Error())
}
