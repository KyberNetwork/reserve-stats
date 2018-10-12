package crawler

import (
	"reflect"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	rsvRateCommon "github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/common"
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/storage/influx"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

const (
	testRsvAddress = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
	testInfluxURL  = "http://127.0.0.1:8086"
)

func newTestCrawler() (*ResreveRatesCrawler, error) {
	var (
		addrs = []ethereum.Address{ethereum.HexToAddress(testRsvAddress)}
	)
	sett := core.MockClient{}
	wrpContract := contracts.MockVersionedWrapper{}
	bltimeRsver := blockchain.MockBlockTimeResolve{}
	lger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: testInfluxURL,
	})
	dbInstance, err := influx.NewRateInfluxDBStorage(influxClient)
	if err != nil {
		return nil, err
	}
	return &ResreveRatesCrawler{
		wrapperContract: &wrpContract,
		Addresses:       addrs,
		tokenSetting:    &sett,
		sugar:           lger.Sugar(),
		blkTimeRsv:      &bltimeRsver,
		db:              dbInstance,
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
	crawler, err := newTestCrawler()
	if err != nil {
		t.Fatal(err)
	}
	rates, err := crawler.GetReserveRates(0)
	if err != nil {
		t.Error(err)
	}
	rate, ok := rates[testRsvAddress]
	if !ok {
		crawler.sugar.Errorf("result did not contain rate for reserve %s", testRsvAddress)
		t.Fail()
	}
	rateEntry, ok := rate.Data["ETH-KNC"]
	if !ok {
		crawler.sugar.Error("result did not contain rate for ETH-KNC pair")
		t.Fail()
	}
	if !reflect.DeepEqual(rateEntry, testRateEntry) {
		crawler.sugar.Error("RateEntry ETH-KNC did not match the expected result")
		t.Fail()
	}
}
