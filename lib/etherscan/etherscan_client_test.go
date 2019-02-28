package etherscanclient

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/nanmu42/etherscan-api"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func TestEtherScanClient(t *testing.T) {
	testutil.SkipExternal(t)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

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
	testutil.SkipExternal(t)
	var (
		rps = 6
		wg  = &sync.WaitGroup{}
	)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()
	etherscanAPIKey, ok := os.LookupEnv("ETHERSCAN_API_KEY")
	if !ok {
		sugar.Info("etherscan api is not provided, using public api")
	}

	etherscanClient := etherscan.New(etherscan.Network("api"), etherscanAPIKey)
	limiter := rate.NewLimiter(rate.Limit(float64(rps)), 1)
	etherscanClient.BeforeRequest = limitRate(limiter)

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			_, err := etherscanClient.EtherTotalSupply()
			sugar.Debugw("finshied a request", "request_number", index, "finish_time", time.Now(), "error", err)
		}(i)
	}
	wg.Wait()
}
