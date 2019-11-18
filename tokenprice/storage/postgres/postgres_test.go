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
		token     = "ETH"
		currency  = "USD"
		coinbase  = "coinbase"
		coingecko = "coingecko"
		timeS     = "2019-02-06"
		price     = 100.1
		newPrice  = 101.2
	)
	timestamp, err := time.Parse("2006-01-02", timeS)
	require.NoError(t, err)
	err = trdb.SaveTokenPrice(token, currency, coinbase, timestamp, price)
	require.NoError(t, err)

	priceDB, err := trdb.GetTokenPrice(token, currency, coinbase, timestamp)
	require.NoError(t, err)
	require.Equal(t, price, priceDB)

	err = trdb.SaveTokenPrice(token, currency, coinbase, timestamp, newPrice)
	require.NoError(t, err)

	newPriceDB, err := trdb.GetTokenPrice(token, currency, coinbase, timestamp)
	require.NoError(t, err)
	require.Equal(t, newPrice, newPriceDB)

	_, err = trdb.GetTokenPrice("KNC", currency, coinbase, timestamp)
	require.EqualError(t, err, ErrNotFound.Error())

	_, err = trdb.GetTokenPrice("KNC", currency, coingecko, timestamp)
	require.EqualError(t, err, ErrNotFound.Error())
}
