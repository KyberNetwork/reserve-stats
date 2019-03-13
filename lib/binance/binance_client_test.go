package binance

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
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

	binanceClient := NewBinance(binanceAPIKey, binanceSecretKey, sugar)

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

func TestBinanceClientWithLimiter(t *testing.T) {
	//Uncomment the skip to run the test in dev mode.
	//Alter these number to test binance's behaviour
	t.Skip()
	var (
		rps     = 20.0
		limiter = rate.NewLimiter(rate.Limit(rps), 5)
		wg      = &sync.WaitGroup{}
	)
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

	binanceClient := NewBinance(binanceAPIKey, binanceSecretKey, sugar, WithRateLimiter(limiter))
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err = binanceClient.GetTradeHistory("KNCETH", 0)
			if err != nil {
				panic(err)
			}
			assert.NoError(t, err, "binance client get trade history error: %s", err)
		}(i)
	}
	wg.Wait()
}
