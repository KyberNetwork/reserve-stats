package binance

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

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
	if !ok {
		t.Skip("Binance API key is not available")
	}

	binanceSecretKey, ok := os.LookupEnv("BINANCE_SECRET_KEY")
	if !ok {
		t.Skip("Binance secret key is not available")
	}

	limiter := rate.NewLimiter(rate.Limit(1200.0), 1)

	binanceClient := NewBinance(binanceAPIKey, binanceSecretKey, sugar, limiter)

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
