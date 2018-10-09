package crawler

import (
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	rsvRateCommon "github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/common"
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/storage/influx"
	ethereum "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

const (
	testRsvAddress = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
	testInfluxURL  = "http://localhost:8086"
)

func getTestCrawler() (*ResreveRatesCrawler, error) {
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
	dbInstance, err := influx.NewRateInfluxDBStorage(testInfluxURL, "", "")
	if err != nil {
		return nil, err
	}
	return &ResreveRatesCrawler{
		wrapperContract: &wrpContract,
		Addresses:       addrs,
		tokenSetting:    &sett,
		logger:          lger.Sugar(),
		blkTimeRsv:      &bltimeRsver,
		db:              dbInstance,
	}, nil
}

func TestGetReserveRate(t *testing.T) {
	var testRateEntry = rsvRateCommon.ReserveRateEntry{
		BuyReserveRate:  1.0,
		SellReserveRate: 2.0,
		BuySanityRate:   3.0,
		SellSanityRate:  4.0,
	}
	crawler, err := getTestCrawler()
	if err != nil {
		t.Fatal(err)
	}
	rates, err := crawler.GetReserveRates(0)
	if err != nil {
		t.Error(err)
	}
	rate, ok := rates[testRsvAddress]
	if !ok {
		crawler.logger.Errorf("result did not contain rate for reserve %s", testRsvAddress)
		t.Fail()
	}
	rateEntry, ok := rate.Data["ETH-KNC"]
	if !ok {
		crawler.logger.Error("result did not contain rate for ETH-KNC pair")
		t.Fail()
	}
	if !reflect.DeepEqual(rateEntry, testRateEntry) {
		crawler.logger.Error("RateEntry ETH-KNC did not match the expected result")
		t.Fail()
	}
}
