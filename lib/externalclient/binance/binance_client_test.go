package binance

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.True(t, ok, "binance api key is not set")

	binanceSecretKey, ok := os.LookupEnv("BINANCE_SECRET_KEY")
	assert.True(t, ok, "binance secret key is not set")

	binanceClient := NewBinanceClient(binanceAPIKey, binanceSecretKey, sugar)

	assetDetail, err := binanceClient.GetAssetDetail()
	assert.Nil(t, err, "binance client get asset detail error: %s", err)
	assert.NotNil(t, assetDetail, "asset detail should not be nil")
	sugar.Info(assetDetail)

	tradeHistory, err := binanceClient.GetTradeHistory("KNCETH", 0)
	assert.Nil(t, err, "binance client get asset detail error: %s", err)
	assert.NotNil(t, tradeHistory, "asset detail should not be nil")

	fromTime := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	toTime := time.Now()
	withdrawHistory, err := binanceClient.GetWithdrawalHistory(fromTime, toTime)
	assert.Nil(t, err, "binance client get asset detail error: %s", err)
	assert.NotNil(t, withdrawHistory, "asset detail should not be nil")
}
