package huobi

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func TestHuobiClient(t *testing.T) {
	testutil.SkipExternal(t)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	huobiAPIKey, ok := os.LookupEnv("HUOBI_API_KEY")
	if !ok {
		t.Skip("Huobi API key is not available")
	}

	huobiSecretKey, ok := os.LookupEnv("HUOBI_SECRET_KEY")
	if !ok {
		t.Skip("Huobi secret key is not available")
	}

	huobiClient := NewClient(huobiAPIKey, huobiSecretKey, sugar)
	_, err = huobiClient.GetAccounts()
	assert.NoError(t, err, fmt.Sprintf("get account fee error: %s", err))

	//fixed timestamp for test
	startDate := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC)

	_, err = huobiClient.GetTradeHistory("bixeth", startDate, endDate)
	assert.NoError(t, err, fmt.Sprintf("get history error: %v", err))

	_, err = huobiClient.GetWithdrawHistory("ETH", 0)
	assert.NoError(t, err)
}

func TestHuobiClientWithLimiter(t *testing.T) {
	//Uncomment the skip to run the test in dev mode.
	//Alter these number to test huobi's behaviour
	t.Skip()
	var (
		rps       = 8.0
		limiter   = rate.NewLimiter(rate.Limit(rps), 1)
		wg        = &sync.WaitGroup{}
		startDate = time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
		endDate   = time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC)
	)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	huobiAPIKey, ok := os.LookupEnv("HUOBI_API_KEY")
	if !ok {
		t.Skip("Huobi API key is not available")
	}

	huobiSecretKey, ok := os.LookupEnv("HUOBI_SECRET_KEY")
	if !ok {
		t.Skip("Huobi secret key is not available")
	}

	huobiClient := NewClient(huobiAPIKey, huobiSecretKey, sugar, WithRateLimiter(limiter))

	assert.NoError(t, err, fmt.Sprintf("get history error: %v", err))
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err = huobiClient.GetTradeHistory("bixeth", startDate, endDate)
			sugar.Debugw("done request", "number", i, "at time", time.Now(), "Err ", err)
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()
}
