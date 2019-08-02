package etherscan

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/nanmu42/etherscan-api"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestEtherScanClient(t *testing.T) {
	testutil.SkipExternal(t)
	sugar := testutil.MustNewDevelopmentSugaredLogger()

	etherscanAPIKey, ok := os.LookupEnv("ETHERSCAN_API_KEY")
	if !ok {
		sugar.Info("etherscan api is not provided, using public api")
	}

	etherscanClient := etherscan.New(etherscan.Network("api"), etherscanAPIKey)
	totalSupply, err := etherscanClient.EtherTotalSupply()
	assert.NoError(t, err, fmt.Sprintf("get total supply of ether error: %s", err))
	assert.NotZero(t, totalSupply)
	sugar.Info(totalSupply)

	//get list normal trade
	testAddress := "0x2c1ba59d6f58433fb1eaee7d20b26ed83bda51a3" //random address
	normalTxs, err := etherscanClient.NormalTxByAddress(testAddress, nil, nil, 1, 10, true)
	assert.NoError(t, err, fmt.Sprintf("get normal tx error: %v", err))
	assert.NotEmpty(t, normalTxs, "get normal txs failed")
}

func TestEtherScanClientWithRateLimiter(t *testing.T) {
	//commented out this and run with different rps to see its behaviour
	//at rps =5 everything goes through
	//at rps =6 etherscan will return 403 for some requests.
	t.Skip()
	var (
		rps = 5
		wg  = &sync.WaitGroup{}
	)
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	etherscanAPIKey, ok := os.LookupEnv("ETHERSCAN_API_KEY")
	if !ok {
		sugar.Info("etherscan api is not provided, using public api")
	}

	etherscanClient := etherscan.New(etherscan.Network("api"), etherscanAPIKey)
	limiter := rate.NewLimiter(rate.Limit(float64(rps)), 1)

	// semaphore by buffer channel
	semaphore := make(chan struct{}, rps)
	etherscanClient.BeforeRequest = limitRate(limiter, semaphore)
	etherscanClient.AfterRequest = func(module, action string, param map[string]interface{}, outcome interface{}, requestErr error) {
		<-semaphore
	}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			_, err := etherscanClient.EtherTotalSupply()
			sugar.Debugw("finshed a request", "request_number", index, "finish_time", time.Now(), "error", err)
			assert.NoError(t, err, fmt.Sprintf("got error: %v", err))
		}(i)
	}
	wg.Wait()
}
