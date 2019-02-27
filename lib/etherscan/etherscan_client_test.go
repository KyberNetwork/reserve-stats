package etherscanclient

import (
	"fmt"
	"os"
	"testing"

	"github.com/nanmu42/etherscan-api"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestEtherScanClient(t *testing.T) {
	// testutil.SkipExternal(t)
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
