package http

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	timeutil "github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
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
)

func getTestServer(sugar *zap.SugaredLogger) (*Server, error) {
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: testInfluxURL,
	})
	if err != nil {
		return nil, err
	}
	rateStorage, err := influxRateStorage.NewRateInfluxDBStorage(sugar, influxClient)
	if err != nil {
		return nil, err
	}
	var testReserveRate common.ReserveRates
	if err = json.Unmarshal([]byte(testRsvRateJSON), &testReserveRate); err != nil {
		return nil, err
	}
	testRecords := map[string]common.ReserveRates{testRsvAddress: testReserveRate}
	if err = rateStorage.UpdateRatesRecords(testRecords); err != nil {
		return nil, err
	}
	return NewServer(host, rateStorage, sugar)
}

func TestHTTPRateServer(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	server, err := getTestServer(logger.Sugar())
	if err != nil {
		t.Fatal(err)
	}
	server.register()
	const (
		requestEndpoint = "reserve-rates"
	)
	var testReserveRate common.ReserveRates
	if err = json.Unmarshal([]byte(testRsvRateJSON), &testReserveRate); err != nil {
		t.Error(err)
	}
	fromTime := timeutil.TimeToTimestampMs(testReserveRate.Timestamp)
	var tests = []httputil.HTTPTestCase{
		{Msg: "success query",
			Endpoint: fmt.Sprintf("%s/%s?from=%d&to=%d&reserve=%s", host, requestEndpoint, fromTime, fromTime, testRsvAddress),
			Method:   "GET",
			Assert:   expectCorrectRate,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, server.r) })
	}
}
