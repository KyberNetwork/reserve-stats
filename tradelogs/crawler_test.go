package tradelogs

import (
	"math/big"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/nanmu42/etherscan-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// const (
// 	alchemyRopsten = "https://eth-ropsten.alchemyapi.io/v2/GvXu_IIrL0U10ZpgTXKjOGIA06KcwzEK"
// )

type mockBroadCastClient struct{}

var ec = etherscan.New(etherscan.Mainnet, "")

func newMockBroadCastClient() *mockBroadCastClient {
	return &mockBroadCastClient{}
}

func (*mockBroadCastClient) GetTxInfo(tx string) (string, string, string, error) {
	return "123", "8.8.8.8", "US", nil
}

func assertTradeLog(t *testing.T, tradeLog common.Tradelog) {
	t.Helper()

	assert.NotZero(t, tradeLog.Timestamp)
	assert.NotZero(t, tradeLog.BlockNumber)
	assert.NotZero(t, tradeLog.TransactionHash)

	if tradeLog.TokenInfo.SrcAddress != blockchain.USDTAddr {
		assert.NotZero(t, tradeLog.QuoteAmount)
	}

	assert.NotZero(t, tradeLog.User.UserAddress)
	assert.NotZero(t, tradeLog.TokenInfo.SrcAddress)
	assert.NotZero(t, tradeLog.TokenInfo.DestAddress)
	assert.NotZero(t, tradeLog.SrcAmount)
	assert.NotZero(t, tradeLog.DestAmount)
	assert.NotZero(t, tradeLog.ReceiverAddress)

	assert.NotZero(t, tradeLog.Index)
}

func TestCrawlerGetTradeLogs(t *testing.T) {
	testutil.SkipExternal(t)

	sugar := testutil.MustNewDevelopmentSugaredLogger()
	client := testutil.MustNewDevelopmentwEthereumClient()

	// v3 contract addresses
	v3Addresses := []ethereum.Address{
		ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"), // network contract
		ethereum.HexToAddress("0x9ae49C0d7F8F9EF4B864e004FE86Ac8294E20950"), // internal network contract
		ethereum.HexToAddress("0x52166528FCC12681aF996e409Ee3a421a4e128A3"), // burner contract
	}
	reserveAddress := ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755")

	c, err := NewCrawler(sugar, client, newMockBroadCastClient(), tokenrate.NewMock(), v3Addresses,
		deployment.StartingBlocks[deployment.Production], reserveAddress, ec)
	require.NoError(t, err)

	result, err := c.GetTradeLogs(big.NewInt(7025000), big.NewInt(7025100), time.Minute)
	require.NoError(t, err)
	require.Len(t, result.Trades, 10)
	for _, tradeLog := range result.Trades {
		assertTradeLog(t, tradeLog)
	}

	// v2 contract addresses
	v2Addresses := []ethereum.Address{
		ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"), // network contract
		ethereum.HexToAddress("0x91a502C678605fbCe581eae053319747482276b9"), // internal network contract
		ethereum.HexToAddress("0xed4f53268bfdFF39B36E8786247bA3A02Cf34B04"), // burner contract
	}

	c, err = NewCrawler(sugar, client, newMockBroadCastClient(), tokenrate.NewMock(), v2Addresses,
		deployment.StartingBlocks[deployment.Production], reserveAddress, ec)
	require.NoError(t, err)

	result, err = c.GetTradeLogs(big.NewInt(6343120), big.NewInt(6343220), time.Minute)
	require.NoError(t, err)
	require.Len(t, result.Trades, 8)
	for _, tradeLog := range result.Trades {
		assertTradeLog(t, tradeLog)
	}

	result, err = c.GetTradeLogs(big.NewInt(6343120), big.NewInt(6343220), time.Minute)
	require.NoError(t, err)
	require.Len(t, result.Trades, 8)
	for _, tradeLog := range result.Trades {
		assertTradeLog(t, tradeLog)
	}

	// block: 6343124
	// tx: 0x4959968da434aa9b1da585479a19da0ce71c47b6e896aab85c6a285400d6de18
	// conversion: WETH --> ETH
	var sampleTxHash = "0x4959968da434aa9b1da585479a19da0ce71c47b6e896aab85c6a285400d6de18"
	var found = false
	for _, tradeLog := range result.Trades {
		if tradeLog.TransactionHash == ethereum.HexToHash(sampleTxHash) {
			found = true
			// assert.Equal(t, ethereum.HexToAddress("0x57f8160e1c59d16c01bbe181fd94db4e56b60495"), tradeLog.SrcReserveAddress,
			// 	"WETH --> ETH trade log must have source reserve address")

			assert.Equal(t, ethereum.HexToAddress("0xd064e4c8f55cff3f45f2e5af5d24bdf0107fe40e"), tradeLog.ReceiverAddress,
				"Tradelog must have receiver address")
		}
	}
	assert.True(t, found, "transaction %s not found", sampleTxHash)

	result, err = c.GetTradeLogs(big.NewInt(6325136), big.NewInt(6325137), time.Minute)
	require.NoError(t, err)
	require.Len(t, result.Trades, 3)
	for _, tradeLog := range result.Trades {
		assertTradeLog(t, tradeLog)
	}

	// block: 6325136
	// tx: 0xa1a0cec06413c5466f46d64bc8d6aa2606e82e2da466ec7b266331c056e20133
	// conversion: ETH --> WETH
	sampleTxHash = "0xa1a0cec06413c5466f46d64bc8d6aa2606e82e2da466ec7b266331c056e20133"
	found = false
	for _, tradeLog := range result.Trades {
		if tradeLog.TransactionHash == ethereum.HexToHash(sampleTxHash) {
			found = true
			// assert.Equal(t,
			// 	ethereum.HexToAddress("0x57f8160e1c59d16c01bbe181fd94db4e56b60495"),
			// 	tradeLog.SrcReserveAddress,
			// 	"ETH --> WETH trade log must have dest reserve address")

			assert.Equal(t, ethereum.HexToAddress("0xf214dde57f32f3f34492ba3148641693058d4a9e"), tradeLog.ReceiverAddress,
				"Tradelog must have receiver address")
		}
	}
	assert.True(t, found, "transaction %s not found", sampleTxHash)

	// block: 7000184
	// tx: 0xbda96c208fee7812f463f1fff515a1c70d9148ffe8b40a91db419a10074d4cc1
	// conversion : ETH-GTO
	// usdtAmount must equal to : 749378067533693720
	result, err = c.GetTradeLogs(big.NewInt(7000184), big.NewInt(7000184), time.Minute)
	require.NoError(t, err)
	require.Len(t, result.Trades, 1)
	for _, tradeLog := range result.Trades {
		assertTradeLog(t, tradeLog)
	}
	sampleTxHash = "0xbda96c208fee7812f463f1fff515a1c70d9148ffe8b40a91db419a10074d4cc1"
	found = false
	for _, tradeLog := range result.Trades {
		if tradeLog.TransactionHash == ethereum.HexToHash(sampleTxHash) {
			found = true
			assert.Equal(t,
				big.NewInt(749378067533693720),
				tradeLog.QuoteAmount,
				"trade log's ETH amount must equal to the 749378067533693720")

			assert.Equal(t, ethereum.HexToAddress("0x85c5c26dc2af5546341fc1988b9d178148b4838b"), tradeLog.ReceiverAddress,
				"Tradelog must have receiver address")
		}
	}
	assert.True(t, found, "transaction %s not found", sampleTxHash)

	// v1 contract addresses
	v1Addresses := []ethereum.Address{
		ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"), // network contract
		ethereum.HexToAddress("0x964F35fAe36d75B1e72770e244F6595B68508CF5"), // internal network contract
		ethereum.HexToAddress("0x07f6e905f2a1559cd9fd43cb92f8a1062a3ca706"), // burner contract
	}

	c, err = NewCrawler(sugar, client, newMockBroadCastClient(), tokenrate.NewMock(), v1Addresses,
		deployment.StartingBlocks[deployment.Production], reserveAddress, ec)
	require.NoError(t, err)

	result, err = c.GetTradeLogs(big.NewInt(5877442), big.NewInt(5877500), time.Minute)
	require.NoError(t, err)
	require.Len(t, result.Trades, 7)
	for _, tradeLog := range result.Trades {
		assertTradeLog(t, tradeLog)
	}

	// block: 5877442
	// tx: 0xe9316215b35fd6c7ba9e8e770a90f67ce6b9332c04bf67d272c04cb26fb43ec0
	// conversion: DAI --> ETH
	sampleTxHash = "0xe9316215b35fd6c7ba9e8e770a90f67ce6b9332c04bf67d272c04cb26fb43ec0"
	found = false
	for _, tradeLog := range result.Trades {
		if tradeLog.TransactionHash == ethereum.HexToHash(sampleTxHash) {
			found = true
			assert.Equal(t, ethereum.HexToAddress("0xfb0f16663c71a2f92bf009c7dc7b401ad372b6de"), tradeLog.ReceiverAddress,
				"Tradelog must have receiver address")
		}
	}
	assert.True(t, found, "transaction %s not found", sampleTxHash)
}

func newTestCrawler(t *testing.T, version string) *Crawler {
	var addresses []ethereum.Address
	switch version {
	case "v1":
		addresses = []ethereum.Address{
			ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"), // network contract
			ethereum.HexToAddress("0x964F35fAe36d75B1e72770e244F6595B68508CF5"), // internal network contract
			ethereum.HexToAddress("0x07f6e905f2a1559cd9fd43cb92f8a1062a3ca706"), // burner contract
		}
	case "v2":
		addresses = []ethereum.Address{
			ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"), // network contract
			ethereum.HexToAddress("0x91a502C678605fbCe581eae053319747482276b9"), // internal network contract
			ethereum.HexToAddress("0xed4f53268bfdFF39B36E8786247bA3A02Cf34B04"), // burner contract
		}
	case "v3":
		addresses = []ethereum.Address{
			ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"), // network contract
			ethereum.HexToAddress("0x9ae49C0d7F8F9EF4B864e004FE86Ac8294E20950"), // internal network contract
			ethereum.HexToAddress("0x52166528FCC12681aF996e409Ee3a421a4e128A3"), // burner contract
		}
	default:
		t.Fatal("not found crawler version")

	}
	reserveAddress := ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755")

	sugar := testutil.MustNewDevelopmentSugaredLogger()
	client := testutil.MustNewDevelopmentwEthereumClient()
	c, err := NewCrawler(sugar, client, newMockBroadCastClient(), tokenrate.NewMock(), addresses,
		deployment.StartingBlocks[deployment.Production], reserveAddress, ec)
	require.NoError(t, err)
	return c
}

// test function for get eth amount (only run locally)
func TestCrawler_GetUSDTAmount(t *testing.T) {
	testutil.SkipExternal(t)
	// test v3 token to token
	c := newTestCrawler(t, "v3")
	result, err := c.GetTradeLogs(big.NewInt(8166246), big.NewInt(8166247), time.Minute)
	require.NoError(t, err)
	require.Len(t, result.Trades, 1)
	require.Equal(t, big.NewInt(7543875834785386865), result.Trades[0].OriginalQuoteAmount)
	require.Equal(t, big.NewInt(0).Mul(big.NewInt(7543875834785386865), big.NewInt(2)), result.Trades[0].OriginalQuoteAmount)
	for _, tradeLog := range result.Trades {
		assertTradeLog(t, tradeLog)
	}
	// test v3 token to eth
	c = newTestCrawler(t, "v3")
	result, err = c.GetTradeLogs(big.NewInt(8180001), big.NewInt(8180002), time.Minute)
	require.NoError(t, err)
	require.Len(t, result.Trades, 2)
	require.Equal(t, big.NewInt(682000000000000000), result.Trades[0].QuoteAmount)
	// eth to weth
	require.Equal(t, big.NewInt(500000000000000000), result.Trades[1].OriginalQuoteAmount)
	require.Equal(t, int64(0), result.Trades[1].QuoteAmount.Int64())
	for _, tradeLog := range result.Trades {
		assertTradeLog(t, tradeLog)
	}

	// test v2
	c = newTestCrawler(t, "v2")
	result, err = c.GetTradeLogs(big.NewInt(6325136), big.NewInt(6325137), time.Minute)
	require.NoError(t, err)
	tradeLogs := result.Trades
	require.Len(t, tradeLogs, 3)
	for _, tradeLog := range tradeLogs {
		assertTradeLog(t, tradeLog)
	}
	require.Equal(t, big.NewInt(int64(478695176421724747)), tradeLogs[0].OriginalQuoteAmount)
	require.Equal(t, big.NewInt(int64(478695176421724747)*2), tradeLogs[0].QuoteAmount)
	require.Equal(t, int64(0), tradeLogs[1].QuoteAmount.Int64())
	require.Equal(t, big.NewInt(10000000000000000), tradeLogs[1].OriginalQuoteAmount)
	require.Equal(t, big.NewInt(int64(1249340978082777639)), tradeLogs[2].QuoteAmount)

	// test v1
	c = newTestCrawler(t, "v1")
	result, err = c.GetTradeLogs(big.NewInt(5877442), big.NewInt(5877500), time.Minute)
	tradeLogs = result.Trades
	require.NoError(t, err)
	require.Len(t, tradeLogs, 7)
	for _, tradeLog := range tradeLogs {
		assertTradeLog(t, tradeLog)
	}
}
