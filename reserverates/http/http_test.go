package http

import (
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage"
	influxRateStorage "github.com/KyberNetwork/reserve-stats/reserverates/storage/influx"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

const (
	host           = "http://127.0.0.1:9001"
	testInfluxURL  = "http://127.0.0.1:8086"
	testRsvAddress = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
	dbName         = "test_reserve_rate"
	testFromBlock  = 123
	testToBlock    = 124
	testTs         = 1539143833304
)

var (
	testRates = map[string]map[string]common.ReserveRateEntry{
		testRsvAddress: {
			"ETH-KNC": {
				BuyReserveRate:  1,
				BuySanityRate:   2,
				SellReserveRate: 3,
				SellSanityRate:  4,
			},
			"ETH-ZRX": {
				BuyReserveRate:  5,
				BuySanityRate:   6,
				SellReserveRate: 7,
				SellSanityRate:  8,
			},
		},
	}
)

func newTestServer(sugar *zap.SugaredLogger, dbInstance storage.ReserveRatesStorage) (*Server, error) {
	if err := dbInstance.UpdateRatesRecords(123, testRates); err != nil {
		return nil, err
	}
	return NewServer(host, dbInstance, sugar)
}

func tearDownTestDB(t *testing.T, influxClient client.Client) {
	cmd := fmt.Sprintf("DROP DATABASE %s", dbName)
	_, err := influxClient.Query(client.Query{
		Command: cmd,
	})
	assert.Nil(t, err, "rate storage test db should be teardown")
}

func TestHTTPRateServer(t *testing.T) {
	const requestEndpoint = "reserve-rates"

	logger, err := zap.NewDevelopment()
	assert.Nil(t, err, "logger should be created")

	defer logger.Sync()
	sugar := logger.Sugar()

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: testInfluxURL,
	})
	assert.Nil(t, err, "influx client should be created successfully")

	dbInstance, err := influxRateStorage.NewRateInfluxDBStorage(sugar, influxClient, dbName, blockchain.NewMockBlockTimeResolve(timeutil.TimestampMsToTime(testTs)))
	assert.Nil(t, err, "Rate storage should be created successfully")

	defer tearDownTestDB(t, influxClient)

	server, err := newTestServer(sugar, dbInstance)
	assert.Nil(t, err, "test server should be created succesfully")

	server.register()

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "success query",
			Endpoint: fmt.Sprintf("%s/%s?from=%d&to=%d&reserve=%s", host, requestEndpoint, testTs, testTs, testRsvAddress),
			Method:   http.MethodGet,
			Assert:   expectCorrectRate,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, server.r) })
	}
}
