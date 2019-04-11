package huobi

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestHuobiClient(t *testing.T) {
	//t.Skip("skip huobi trade history test")
	sugar := testutil.MustNewDevelopmentSugaredLogger()

	huobiAPIKey, ok := os.LookupEnv("HUOBI_API_KEY")
	if !ok {
		t.Skip("Huobi API key is not available")
	}

	huobiSecretKey, ok := os.LookupEnv("HUOBI_SECRET_KEY")
	if !ok {
		t.Skip("Huobi secret key is not available")
	}

	huobiClient, err := NewClient(huobiAPIKey, huobiSecretKey, sugar)
	assert.NoError(t, err)

	_, err = huobiClient.GetAccounts()
	assert.NoError(t, err, fmt.Sprintf("get account fee error: %s", err))

	//fixed timestamp for test
	startDate := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC)

	_, err = huobiClient.GetTradeHistory("bixeth", startDate, endDate)
	assert.NoError(t, err, fmt.Sprintf("get history error: %v", err))

	withdrawHistory, err := huobiClient.GetWithdrawHistory("ETH", 0)
	assert.NoError(t, err)
	sugar.Infow("withdraw history", "value", withdrawHistory)
}

func TestHuobiClientWithLimiter(t *testing.T) {
	//Uncomment the skip to run the test in dev mode.
	//Alter these number to test huobi's behaviour
	t.Skip("skip huobi trade history test")
	var (
		rps       = 8.0
		limiter   = rate.NewLimiter(rate.Limit(rps), 1)
		wg        = &sync.WaitGroup{}
		startDate = time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
		endDate   = time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC)
	)
	sugar := testutil.MustNewDevelopmentSugaredLogger()

	huobiAPIKey, ok := os.LookupEnv("HUOBI_API_KEY")
	if !ok {
		t.Skip("Huobi API key is not available")
	}

	huobiSecretKey, ok := os.LookupEnv("HUOBI_SECRET_KEY")
	if !ok {
		t.Skip("Huobi secret key is not available")
	}

	huobiClient, err := NewClient(huobiAPIKey, huobiSecretKey, sugar, WithRateLimiter(limiter))
	assert.NoError(t, err)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := huobiClient.GetTradeHistory("bixeth", startDate, endDate)
			sugar.Debugw("done request", "number", i, "at time", time.Now(), "Err ", err)
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()
}
