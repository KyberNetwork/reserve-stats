package huobi

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
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
	limiter := rate.NewLimiter(rate.Limit(5.0), 1)

	huobiClient := NewClient(huobiAPIKey, huobiSecretKey, sugar, limiter)
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
