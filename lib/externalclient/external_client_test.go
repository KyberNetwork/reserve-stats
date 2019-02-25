package externalclient

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/adshao/go-binance"
	"github.com/nanmu42/etherscan-api"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestBinanceClient(t *testing.T) {
	testutil.SkipExternal(t)
	binanceAPIKey, ok := os.LookupEnv("BINANCE_API_KEY")
	assert.True(t, ok, "binance api key is not set")

	binanceSecretKey, ok := os.LookupEnv("BINANCE_SECRET_KEY")
	assert.True(t, ok, "binance secret key is not set")

	binanceClient := binance.NewClient(binanceAPIKey, binanceSecretKey)

	err := binanceClient.NewPingService().Do(context.Background())
	assert.Nil(t, err, fmt.Sprintf("get withdraw fee error: %s", err))
}

func TestHuobiClient(t *testing.T) {
	testutil.SkipExternal(t)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	huobiAPIKey, ok := os.LookupEnv("BINANCE_API_KEY")
	assert.True(t, ok, "huobi api key is not set")

	huobiSecretKey, ok := os.LookupEnv("BINANCE_SECRET_KEY")
	assert.True(t, ok, "huobi secret key is not set")

	huobiClient := NewHuobiClient(huobiAPIKey, huobiSecretKey)
	accounts, err := huobiClient.GetAccounts()
	assert.Nil(t, err, fmt.Sprintf("get withdraw fee error: %s", err))

	sugar.Infow("huobi accounts", "accounts", accounts)
}

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
	assert.Nil(t, err, fmt.Sprintf("get total supply of ether error: %s", err))

	sugar.Infow("ether total supply", "value", totalSupply)
}
