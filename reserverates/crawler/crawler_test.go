package crawler

import (
	"reflect"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	rsvRateCommon "github.com/KyberNetwork/reserve-stats/reserverates/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	testRsvAddress = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
)

func newTestCrawler(sugar *zap.SugaredLogger) (*ResreveRatesCrawler, error) {
	var (
		addrs       = []ethereum.Address{ethereum.HexToAddress(testRsvAddress)}
		sett        = core.NewMockClient()
		wrpContract = contracts.MockVersionedWrapper{}
		bltimeRsver = blockchain.MockBlockTimeResolve{}
	)

	return &ResreveRatesCrawler{
		wrapperContract: &wrpContract,
		Addresses:       addrs,
		tokenSetting:    sett,
		sugar:           sugar,
		blkTimeRsv:      &bltimeRsver,
	}, nil
}

// TestGetReserveRate query the mock blockchain for reserve rate result
// and ensure that the result is the rate configured
func TestGetReserveRate(t *testing.T) {
	var testRateEntry = rsvRateCommon.ReserveRateEntry{
		BuyReserveRate:  1.0,
		SellReserveRate: 2.0,
		BuySanityRate:   3.0,
		SellSanityRate:  4.0,
	}
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err, "logger should be created")

	defer logger.Sync()
	sugar := logger.Sugar()

	crawler, err := newTestCrawler(sugar)
	assert.Nil(t, err, "test crawler should be created")

	rates, err := crawler.GetReserveRates(0)
	assert.Nil(t, err, "reserve rate should be generate")

	rate, ok := rates[testRsvAddress]
	if !ok {
		sugar.Errorf("result did not contain rate for reserve %s", testRsvAddress)
		t.Fail()
	}
	rateEntry, ok := rate.Data["ETH-KNC"]
	if !ok {
		sugar.Error("result did not contain rate for ETH-KNC pair")
		t.Fail()
	}
	if !reflect.DeepEqual(rateEntry, testRateEntry) {
		sugar.Error("RateEntry ETH-KNC did not match the expected result")
		t.Fail()
	}
}
