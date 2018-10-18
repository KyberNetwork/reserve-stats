package crawler

import (
	"reflect"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	rsvRateCommon "github.com/KyberNetwork/reserve-stats/reserverates/common"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage/influx"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

const (
	testRsvAddress = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
	testInfluxURL  = "http://127.0.0.1:8086"
	dbName         = "test_rate"
)

func newTestCrawler(sugar *zap.SugaredLogger, dbInstance storage.ReserveRatesStorage) (*ResreveRatesCrawler, error) {
	var (
		addrs       = []ethereum.Address{ethereum.HexToAddress(testRsvAddress)}
		sett        = core.MockClient{}
		wrpContract = contracts.MockVersionedWrapper{}
		bltimeRsver = blockchain.MockBlockTimeResolve{}
	)

	return &ResreveRatesCrawler{
		wrapperContract: &wrpContract,
		Addresses:       addrs,
		tokenSetting:    &sett,
		sugar:           sugar,
		blkTimeRsv:      &bltimeRsver,
		db:              dbInstance,
	}, nil
}

func newTestDB(sugar *zap.SugaredLogger) (*influx.RateStorage, error) {
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: testInfluxURL,
	})
	if err != nil {
		return nil, err
	}
	return influx.NewRateInfluxDBStorage(sugar, influxClient, dbName)
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
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	dbInstance, err := newTestDB(sugar)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		tErr := dbInstance.TearDown()
		if tErr != nil {
			t.Fatal(tErr)
		}
	}()
	crawler, err := newTestCrawler(sugar, dbInstance)
	if err != nil {
		t.Fatal(err)
	}
	rates, err := crawler.GetReserveRates(0)
	if err != nil {
		t.Error(err)
	}
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
