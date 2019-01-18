package tradelogs

import (
	"math/big"
	"os"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const defaultEthereumNode = "https://mainnet.infura.io"

type mockBroadCastClient struct{}

func newMockBroadCastClient() *mockBroadCastClient {
	return &mockBroadCastClient{}
}

func (*mockBroadCastClient) GetTxInfo(tx string) (string, string, error) {
	return "8.8.8.8", "US", nil
}

func assertTradeLog(t *testing.T, tradeLog common.TradeLog) {
	t.Helper()

	assert.NotZero(t, tradeLog.Timestamp)
	assert.NotZero(t, tradeLog.BlockNumber)
	assert.NotZero(t, tradeLog.TransactionHash)

	if tradeLog.SrcAddress != blockchain.ETHAddr {
		assert.NotZero(t, tradeLog.EtherReceivalAmount)
	}

	assert.NotZero(t, tradeLog.UserAddress)
	assert.NotZero(t, tradeLog.SrcAddress)
	assert.NotZero(t, tradeLog.DestAddress)
	assert.NotZero(t, tradeLog.SrcAmount)
	assert.NotZero(t, tradeLog.DestAmount)

	if blockchain.IsBurnable(tradeLog.SrcAddress) || blockchain.IsBurnable(tradeLog.DestAddress) {
		assert.NotZero(t, tradeLog.BurnFees)
	}

	assert.NotZero(t, tradeLog.ETHUSDRate)
	assert.NotZero(t, tradeLog.ETHUSDProvider)
	assert.NotZero(t, tradeLog.Index)
}

func TestCrawlerGetTradeLogs(t *testing.T) {
	t.Skip("disable as this test require external resource")

	node, ok := os.LookupEnv("ETHEREUM_NODE")
	if !ok {
		node = defaultEthereumNode
	}

	logger, err := zap.NewDevelopment()
	require.Nil(t, err, "logger should be initiated successfully")
	sugar := logger.Sugar()

	client, err := ethclient.Dial(node)
	require.NoError(t, err)

	// v3 contract addresses
	v3Addresses := []ethereum.Address{
		ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"), // network contract
		ethereum.HexToAddress("0x9ae49C0d7F8F9EF4B864e004FE86Ac8294E20950"), // internal network contract
		ethereum.HexToAddress("0x52166528FCC12681aF996e409Ee3a421a4e128A3"), // burner contract
	}
	c, err := NewCrawler(
		sugar,
		client,
		newMockBroadCastClient(),
		tokenrate.NewMock(),
		v3Addresses,
		deployment.StartingBlocks[deployment.Production])
	require.NoError(t, err)

	tradeLogs, err := c.GetTradeLogs(big.NewInt(7025000), big.NewInt(7025100), time.Minute)
	require.NoError(t, err)
	require.Len(t, tradeLogs, 10)
	for _, tradeLog := range tradeLogs {
		assertTradeLog(t, tradeLog)
	}

	// v2 contract addresses
	v2Addresses := []ethereum.Address{
		ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"), // network contract
		ethereum.HexToAddress("0x91a502C678605fbCe581eae053319747482276b9"), // internal network contract
		ethereum.HexToAddress("0xed4f53268bfdFF39B36E8786247bA3A02Cf34B04"), // burner contract
	}

	c, err = NewCrawler(
		sugar,
		client,
		newMockBroadCastClient(),
		tokenrate.NewMock(),
		v2Addresses,
		deployment.StartingBlocks[deployment.Production])
	require.NoError(t, err)

	tradeLogs, err = c.GetTradeLogs(big.NewInt(7000000), big.NewInt(7000100), time.Minute)
	require.NoError(t, err)
	require.Len(t, tradeLogs, 12)
	for _, tradeLog := range tradeLogs {
		assertTradeLog(t, tradeLog)
	}
}
