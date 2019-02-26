package huobi

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHuobiClient(t *testing.T) {
	// testutil.SkipExternal(t)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	huobiAPIKey, ok := os.LookupEnv("HUOBI_API_KEY")
	assert.True(t, ok, "huobi api key is not set")

	huobiSecretKey, ok := os.LookupEnv("HUOBI_SECRET_KEY")
	assert.True(t, ok, "huobi secret key is not set")

	huobiClient := NewHuobiClient(huobiAPIKey, huobiSecretKey, sugar)
	accounts, err := huobiClient.GetAccounts()
	assert.Nil(t, err, fmt.Sprintf("get account fee error: %s", err))
	assert.NotNil(t, accounts, "get accounts nil")

	//fixed timestamp for test
	startDate := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC)

	tradeHistory, err := huobiClient.GetTradeHistory("bixeth", startDate, endDate)
	assert.Nil(t, err, fmt.Sprintf("get history error: %v", err))
	assert.NotNil(t, tradeHistory, "trade history is nil")
	sugar.Info(tradeHistory)

	withdrawHistory, err := huobiClient.GetWithdrawHistory()
	assert.Nil(t, err)
	assert.NotNil(t, withdrawHistory)
	sugar.Info(withdrawHistory)
}
