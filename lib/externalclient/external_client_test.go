package externalclient

import (
	"fmt"
	"os"
	"testing"

	"github.com/nanmu42/etherscan-api"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

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
	assert.True(t, ok, "binance api key is not set")

	binanceSecretKey, ok := os.LookupEnv("BINANCE_SECRET_KEY")
	assert.True(t, ok, "binance secret key is not set")

	binanceClient := NewBinanceClient(binanceAPIKey, binanceSecretKey, sugar)

	assetDetail, err := binanceClient.GetAssetDetail()
	assert.Nil(t, err, "binance client get asset detail error: %s", err)
	assert.NotNil(t, assetDetail, "asset detail should not be nil")

	tradeHistory, err := binanceClient.GetTradeHistory("KNCETH", 0)
	assert.Nil(t, err, "binance client get asset detail error: %s", err)
	assert.NotNil(t, tradeHistory, "asset detail should not be nil")

	withdrawHistory, err := binanceClient.GetWithdrawalHistory(0, currentTimePoint())
	assert.Nil(t, err, "binance client get asset detail error: %s", err)
	assert.NotNil(t, withdrawHistory, "asset detail should not be nil")
}

func TestHuobiClient(t *testing.T) {
	testutil.SkipExternal(t)
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
	startDate := "2018-01-01"
	endDate := "2019-03-08"

	tradeHistory, err := huobiClient.GetTradeHistory("bixeth", startDate, endDate)
	assert.Nil(t, err, fmt.Sprintf("get history error: %v", err))
	assert.NotNil(t, tradeHistory, "trade history is nil")

	withdrawHistory, err := huobiClient.GetWithdrawHistory()
	assert.Nil(t, err)
	assert.NotNil(t, withdrawHistory)
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
	assert.NotZero(t, totalSupply)

	//get list normal trade
	testAddress := "0x2c1ba59d6f58433fb1eaee7d20b26ed83bda51a3" //random address
	normalTxs, err := etherscanClient.NormalTxByAddress(testAddress, nil, nil, 1, 10, true)
	assert.Nil(t, err, fmt.Sprintf("get normal tx error: %v", err))
	assert.NotNil(t, normalTxs, "get normal txs failed")
}
