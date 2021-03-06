package crawler

import (
	"reflect"
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	rsvRateCommon "github.com/KyberNetwork/reserve-stats/reserverates/common"
)

const (
	testRsvAddress = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
)

var (
	knc = ethereum.HexToAddress("0xdd974D5C2e2928deA5F71b9825b8b646686BD200") //KNC
	zrx = ethereum.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498") //ZRX
)

type mockSupportedTokens struct{}

func (mst *mockSupportedTokens) Tokens(ethereum.Address, uint64) ([]blockchain.TokenInfo, error) {
	return []blockchain.TokenInfo{
		{
			Address: knc,
			Symbol:  "KNC",
		},
		{
			Address: zrx,
			Symbol:  "ZRX",
		},
	}, nil
}

func newTestCrawler(sugar *zap.SugaredLogger) (*ReserveRatesCrawler, error) {
	var (
		wrpContract = contracts.MockVersionedWrapper{}
	)

	return &ReserveRatesCrawler{
		wrapperContract: &wrpContract,
		rtf:             &mockSupportedTokens{},
		sugar:           sugar,
	}, nil
}

// TestGetReserveRate query the mock blockchain for reserve rate result
// and ensure that the result is the rate configured
func TestGetReserveRate(t *testing.T) {
	var (
		addrs = []ethereum.Address{ethereum.HexToAddress(testRsvAddress)}

		testRateEntry = rsvRateCommon.ReserveRateEntry{
			BuyReserveRate:  1.0,
			SellReserveRate: 2.0,
			BuySanityRate:   3.0,
			SellSanityRate:  4.0,
		}
	)
	sugar := testutil.MustNewDevelopmentSugaredLogger()

	crawler, err := newTestCrawler(sugar)
	assert.Nil(t, err, "test crawler should be created")

	rates, err := crawler.GetReserveRatesWithAddresses(addrs, 0)
	assert.Nil(t, err, "reserve rate should be generate")

	rate, ok := rates[testRsvAddress]
	if !ok {
		sugar.Errorf("result did not contain rate for reserve %s", testRsvAddress)
		t.Fail()
	}

	_, ok = rate["ETH-KNC"]
	if !ok {
		t.Fatal("result did not contain rate for ETH-KNC pair")
	}

	if !reflect.DeepEqual(rate["ETH-KNC"], testRateEntry) {
		sugar.Error("RateEntry ETH-KNC did not match the expected result")
		t.Fail()
	}
}
