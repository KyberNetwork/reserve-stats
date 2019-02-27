package binance

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestBinanceClient(t *testing.T) {
	testutil.SkipExternal(t)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	binanceAPIKey, ok := os.LookupEnv("BINANCE_API_KEY")
	require.True(t, ok, "binance api key is not set")

	binanceSecretKey, ok := os.LookupEnv("BINANCE_SECRET_KEY")
	require.True(t, ok, "binance secret key is not set")

	binanceClient := NewBinanceClient(binanceAPIKey, binanceSecretKey, sugar)

	assetDetail, err := binanceClient.GetAssetDetail()
	assert.NoError(t, err, "binance client get asset detail error: %s", err)
	assert.NotEmpty(t, assetDetail, "asset detail should not be nil")
	// sugar.Info(assetDetail)

	_, err = binanceClient.GetTradeHistory("KNCETH", 0)
	assert.NoError(t, err, "binance client get trade history error: %s", err)

	fromTime := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	toTime := time.Now()
	_, err = binanceClient.GetWithdrawalHistory(fromTime, toTime)
	assert.NoError(t, err, "binance client get withdraw history error: %s", err)
}
