package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	timeutil "github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage"
	influxRateStorage "github.com/KyberNetwork/reserve-stats/reserverates/storage/influx"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

const (
	host            = "http://127.0.0.1:9001"
	testInfluxURL   = "http://127.0.0.1:8086"
	testRsvAddress  = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
	testRsvRateJSON = `{
		"Timestamp": "2018-10-10T15:57:13.304+07:00",
		"BlockNumber": 123,
		"Data": {
		  "ETH-KNC": {
			"BuyReserveRate": 1,
			"BuySanityRate": 3,
			"SellReserveRate": 2,
			"SellSanityRate": 4
		  },
		  "ETH-ZRX": {
			"BuyReserveRate": 1,
			"BuySanityRate": 3,
			"SellReserveRate": 2,
			"SellSanityRate": 4
		  }
		}
	  }`
	dbName = "test_reserve_rate"
)

func newTestServer(sugar *zap.SugaredLogger, dbInstance storage.ReserveRatesStorage) (*Server, error) {
	var testReserveRate common.ReserveRates
	if err := json.Unmarshal([]byte(testRsvRateJSON), &testReserveRate); err != nil {
		return nil, err
	}
	testRecords := map[string]common.ReserveRates{testRsvAddress: testReserveRate}
	if err := dbInstance.UpdateRatesRecords(testRecords); err != nil {
		return nil, err
	}
	return NewServer(host, dbInstance, sugar)
}

func newTestDB(sugar *zap.SugaredLogger) (*influxRateStorage.RateStorage, error) {
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: testInfluxURL,
	})
	if err != nil {
		return nil, err
	}
	return influxRateStorage.NewRateInfluxDBStorage(sugar, influxClient, dbName)
}

func TestHTTPRateServer(t *testing.T) {
	const requestEndpoint = "reserve-rates"

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
	server, err := newTestServer(sugar, dbInstance)
	if err != nil {
		t.Fatal(err)
	}
	server.register()

	var testReserveRate common.ReserveRates
	if err = json.Unmarshal([]byte(testRsvRateJSON), &testReserveRate); err != nil {
		t.Error(err)
	}
	fromTime := timeutil.TimeToTimestampMs(testReserveRate.Timestamp)
	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "success query",
			Endpoint: fmt.Sprintf("%s/%s?from=%d&to=%d&reserve=%s", host, requestEndpoint, fromTime, fromTime, testRsvAddress),
			Method:   http.MethodGet,
			Assert:   expectCorrectRate,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, server.r) })
	}
}
