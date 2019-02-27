package huobi

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

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
	require.True(t, ok, "huobi api key is not set")

	huobiSecretKey, ok := os.LookupEnv("HUOBI_SECRET_KEY")
	require.True(t, ok, "huobi secret key is not set")

	huobiClient := NewHuobiClient(huobiAPIKey, huobiSecretKey, sugar)
	accounts, err := huobiClient.GetAccounts()
	assert.NoError(t, err, fmt.Sprintf("get account fee error: %s", err))
	assert.NotEmpty(t, accounts, "get accounts nil")

	//fixed timestamp for test
	startDate := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC)

	_, err = huobiClient.GetTradeHistory("bixeth", startDate, endDate)
	assert.NoError(t, err, fmt.Sprintf("get history error: %v", err))

	_, err = huobiClient.GetWithdrawHistory("ETH", 0)
	assert.NoError(t, err)
}
