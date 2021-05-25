package tradelogs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	ether "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nanmu42/etherscan-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	infuraRopsten = "https://ropsten.infura.io/v3/3243b06bf0334cff8e468bf90ce6e08c"
)

type mockBroadCastClient struct{}

var ec = etherscan.New(etherscan.Mainnet, "")

func newMockBroadCastClient() *mockBroadCastClient {
	return &mockBroadCastClient{}
}

func (*mockBroadCastClient) GetTxInfo(tx string) (string, string, string, error) {
	return "123", "8.8.8.8", "US", nil
}

func assertTradeLog(t *testing.T, tradeLog common.TradelogV4) {
	t.Helper()

	assert.NotZero(t, tradeLog.Timestamp)
	assert.NotZero(t, tradeLog.BlockNumber)
	assert.NotZero(t, tradeLog.TransactionHash)

	if tradeLog.TokenInfo.SrcAddress != blockchain.ETHAddr {
		assert.NotZero(t, tradeLog.EthAmount)
	}

	assert.NotZero(t, tradeLog.User.UserAddress)
	assert.NotZero(t, tradeLog.TokenInfo.SrcAddress)
	assert.NotZero(t, tradeLog.TokenInfo.DestAddress)
	assert.NotZero(t, tradeLog.SrcAmount)
	assert.NotZero(t, tradeLog.DestAmount)
	assert.NotZero(t, tradeLog.ReceiverAddress)

	assert.NotZero(t, tradeLog.ETHUSDRate)
	assert.NotZero(t, tradeLog.ETHUSDProvider)
	assert.NotZero(t, tradeLog.Index)
}

var (
	nwProxyAddr      = ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755")
	kyberStorageAddr = ethereum.HexToAddress("")
	feeHandlerAddr   = ethereum.HexToAddress("")
	feeHandlerV2Addr = ethereum.HexToAddress("")
	kyberNetwork     = ethereum.HexToAddress("")
)

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
	c, err := NewCrawler(sugar, client, newMockBroadCastClient(), tokenrate.NewMock(), v3Addresses,
		deployment.StartingBlocks[deployment.Production], ec, []ethereum.Address{}, nwProxyAddr, kyberStorageAddr, feeHandlerAddr, feeHandlerV2Addr, kyberNetwork)
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
		deployment.StartingBlocks[deployment.Production], ec, []ethereum.Address{}, nwProxyAddr, kyberStorageAddr, feeHandlerAddr, feeHandlerV2Addr, kyberNetwork)
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

			assert.Greater(t, tradeLog.T2EReserves, 0) // "WETH --> ETH trade log must have T2E reserves"

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
			assert.Greater(t, tradeLog.E2TReserves, 0) // 	"ETH --> WETH trade log must have E2T reserves > 0"

			assert.Equal(t, ethereum.HexToAddress("0xf214dde57f32f3f34492ba3148641693058d4a9e"), tradeLog.ReceiverAddress,
				"Tradelog must have receiver address")
		}
	}
	assert.True(t, found, "transaction %s not found", sampleTxHash)

	// block: 7000184
	// tx: 0xbda96c208fee7812f463f1fff515a1c70d9148ffe8b40a91db419a10074d4cc1
	// conversion : ETH-GTO
	// ethAmount must equal to : 749378067533693720
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
				tradeLog.EthAmount,
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
		deployment.StartingBlocks[deployment.Production], ec, []ethereum.Address{}, nwProxyAddr, kyberStorageAddr, feeHandlerAddr, feeHandlerV2Addr, kyberNetwork)
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
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	client := testutil.MustNewDevelopmentwEthereumClient()
	c, err := NewCrawler(sugar, client, newMockBroadCastClient(), tokenrate.NewMock(), addresses,
		deployment.StartingBlocks[deployment.Production], ec, []ethereum.Address{}, nwProxyAddr, kyberStorageAddr, feeHandlerAddr, feeHandlerV2Addr, kyberNetwork)
	require.NoError(t, err)
	return c
}

// test function for get eth amount (only run locally)
func TestCrawler_GetEthAmount(t *testing.T) {
	testutil.SkipExternal(t)
	// test v3 token to token
	c := newTestCrawler(t, "v3")
	result, err := c.GetTradeLogs(big.NewInt(8166246), big.NewInt(8166247), time.Minute)
	require.NoError(t, err)
	require.Len(t, result.Trades, 1)
	require.Equal(t, big.NewInt(7543875834785386865), result.Trades[0].OriginalEthAmount)
	require.Equal(t, big.NewInt(0).Mul(big.NewInt(7543875834785386865), big.NewInt(2)), result.Trades[0].EthAmount)
	for _, tradeLog := range result.Trades {
		assertTradeLog(t, tradeLog)
	}
	// test v3 token to eth
	c = newTestCrawler(t, "v3")
	result, err = c.GetTradeLogs(big.NewInt(8180001), big.NewInt(8180002), time.Minute)
	require.NoError(t, err)
	require.Len(t, result.Trades, 2)
	require.Equal(t, big.NewInt(682000000000000000), result.Trades[0].EthAmount)
	// eth to weth
	require.Equal(t, big.NewInt(500000000000000000), result.Trades[1].OriginalEthAmount)
	require.Equal(t, int64(0), result.Trades[1].EthAmount.Int64())
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
	require.Equal(t, big.NewInt(int64(478695176421724747)), tradeLogs[0].OriginalEthAmount)
	require.Equal(t, big.NewInt(int64(478695176421724747)*2), tradeLogs[0].EthAmount)
	require.Equal(t, int64(0), tradeLogs[1].EthAmount.Int64())
	require.Equal(t, big.NewInt(10000000000000000), tradeLogs[1].OriginalEthAmount)
	require.Equal(t, big.NewInt(int64(1249340978082777639)), tradeLogs[2].EthAmount)

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

func TestDecodeTradeWithHintTx(t *testing.T) {
	// example of transaction input data
	txInput := "0x29589f61000000000000000000000000eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee00000000000000000000000000000000000000000000000002e0f43384a3591f0000000000000000000000000d8775f648430679a709e98d2b0cb6250d2887ef0000000000000000000000005eee96fa064a571dabcbfe43799d46e5de2f51f98000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000031c454332beb5e2eff000000000000000000000000440bbd6a888a36de6e2f6a25f65bc4e16874faa9000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000045045524d00000000000000000000000000000000000000000000000000000000"
	var data tradeWithHintParam
	data, err := decodeTradeInputParamV3(hexutil.MustDecode(txInput))
	require.NoError(t, err)

	t.Log("src", data.Src.String())
	t.Log("srcAmount", data.SrcAmount)
	t.Log("dest", data.Dest.String())
	t.Log("destAddr", data.DestAddress.String())
	t.Log("maxDestAmount", data.MaxDestAmount)
	t.Log("walletID", data.WalletID.String())
	t.Log("minConversionRate", data.MinConversionRate.String())

	assert.Equal(t, ethereum.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"), data.Src)
	assert.Equal(t, ethereum.HexToAddress("0x0D8775F648430679A709E98d2b0Cb6250d2887EF"), data.Dest)
	assert.Equal(t, ethereum.HexToAddress("0x440bBd6a888a36DE6e2F6A25f65bc4e16874faa9"), data.WalletID)
}

func TestDecodeTrade(t *testing.T) {
	txInput := "0xcb3c28c70000000000000000000000006b175474e89094c44da98b954eedeac495271d0f000000000000000000000000000000000000000000000002b5e3af16b1880000000000000000000000000000eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee000000000000000000000000f76cb72ecfa276d8d64a24f54a94554e2da8f712800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000147c9c7c10adb2000000000000000000000000f1aa99c69715f423086008eb9d06dc1e35cc504d"
	var data tradeWithHintParam
	data, err := decodeTradeInputParamV3(hexutil.MustDecode(txInput))
	require.NoError(t, err)
	t.Log("src", data.Src.String())
	t.Log("srcAmount", data.SrcAmount)
	t.Log("dest", data.Dest.String())
	t.Log("destAddr", data.DestAddress.String())
	t.Log("maxDestAmount", data.MaxDestAmount)
	t.Log("walletID", data.WalletID.String())
	t.Log("minConversionRate", data.MinConversionRate.String())

	assert.Equal(t, ethereum.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"), data.Src)
	assert.Equal(t, ethereum.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"), data.Dest)
	assert.Equal(t, ethereum.HexToAddress("0xF1AA99C69715F423086008eB9D06Dc1E35Cc504d"), data.WalletID)
}

func TestDecodeTradeWithHintAndFeeTxV4(t *testing.T) {
	// example of transaction input data
	txInput := "0xae591d54000000000000000000000000eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee0000000000000000000000000000000000000000000000001bc16d674ec80000000000000000000000000000d9ec3ff1f8be459bb9369b4e79e9ebcf7141c093000000000000000000000000017cf1489eb1ff78ef8797dfc09662c845e187b9ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000004c083087eefb8b28280000000000000000000000000440bbd6a888a36de6e2f6a25f65bc4e16874faa9000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000000000"
	var data tradeWithHintParamV4
	data, err := decodeTradeInputParamV4(hexutil.MustDecode(txInput))
	require.NoError(t, err)

	t.Log("src", data.Src.String())
	t.Log("srcAmount", data.SrcAmount)
	t.Log("dest", data.Dest.String())
	t.Log("destAddr", data.DestAddress.String())
	t.Log("maxDestAmount", data.MaxDestAmount)
	t.Log("walletID", data.PlatformWallet.String())
	t.Log("minConversionRate", data.MinConversionRate.String())

	assert.Equal(t, ethereum.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"), data.Src)
	assert.Equal(t, ethereum.HexToAddress("0xD9Ec3ff1f8be459Bb9369b4E79e9Ebcf7141C093"), data.Dest)
	assert.Equal(t, ethereum.HexToAddress("0x440bBd6a888a36DE6e2F6A25f65bc4e16874faa9"), data.PlatformWallet)
}

func TestTraceAddReserve(t *testing.T) {
	c, err := ethclient.Dial(infuraRopsten)
	require.NoError(t, err)
	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(addReserveToStorageEvent),
		},
	}
	query := ether.FilterQuery{
		FromBlock: big.NewInt(8008389),
		ToBlock:   big.NewInt(8008389),
		Addresses: []ethereum.Address{ethereum.HexToAddress("0xa4ead31a6c8e047e01ce1128e268c101ad391959")},
		Topics:    topics,
	}
	logs, err := c.FilterLogs(context.Background(), query)
	require.NoError(t, err)
	ab, err := abi.JSON(bytes.NewReader([]byte(contracts.KyberStorageABI)))
	require.NoError(t, err)

	type AddReserveToStorageEvent struct {
		ReserveType uint8
		Add         bool
	}

	for _, log := range logs {
		t.Log("hash", log.TxHash.String())
		var d AddReserveToStorageEvent
		for _, d := range log.Topics {
			t.Log("topic", d.String())
		}
		err = ab.Unpack(&d, "AddReserveToStorage", log.Data)
		require.NoError(t, err)
		t.Log("rawdata", log.Data)
	}
}

func TestCrawlFeeDistributed(t *testing.T) {
	c, err := ethclient.Dial(infuraRopsten)
	require.NoError(t, err)
	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(feeDistributedEvent),
		},
	}
	query := ether.FilterQuery{
		FromBlock: big.NewInt(8151682),
		ToBlock:   big.NewInt(8152370),
		Addresses: []ethereum.Address{ethereum.HexToAddress("0xe57B2c3b4E44730805358131a6Fc244C57178Da7")},
		Topics:    topics,
	}
	logs, err := c.FilterLogs(context.Background(), query)
	require.NoError(t, err)
	kyberFeeHandlerContract, err := contracts.NewKyberFeeHandler(ethereum.HexToAddress("0xe57B2c3b4E44730805358131a6Fc244C57178Da7"), c)
	require.NoError(t, err)

	for _, log := range logs {
		d, err := kyberFeeHandlerContract.ParseFeeDistributed(log)
		require.NoError(t, err)
		dataByte, _ := json.Marshal(d)
		fmt.Printf("%s\n", dataByte)
	}
}

func TestCrawlKyberTrade(t *testing.T) {
	c, err := ethclient.Dial(infuraRopsten)
	require.NoError(t, err)
	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(kyberTradeEventV4),
		},
	}
	query := ether.FilterQuery{
		FromBlock: big.NewInt(8111308),
		ToBlock:   big.NewInt(8111408),
		Addresses: []ethereum.Address{ethereum.HexToAddress("0x9EC49C41Fdc4C79fDb042AF37659f2E3220ad0a4")},
		Topics:    topics,
	}
	logs, err := c.FilterLogs(context.Background(), query)
	require.NoError(t, err)
	kyberNetworkContract, err := contracts.NewKyberNetwork(ethereum.HexToAddress("0x9EC49C41Fdc4C79fDb042AF37659f2E3220ad0a4"), c)
	require.NoError(t, err)

	for _, log := range logs {
		d, err := kyberNetworkContract.ParseKyberTrade(log)
		require.NoError(t, err)
		dataByte, _ := json.Marshal(d)
		fmt.Printf("%s\n", dataByte)
	}
}
