package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"

	"github.com/influxdata/influxdb/client/v2"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/users/storage"
)

func tearDown(t *testing.T, teardown func() error, influxClient client.Client) {
	assert.NoError(t, teardown(), "database should be deleted completely")

	_, err := influxClient.Query(client.Query{
		Command: fmt.Sprintf("DROP DATABASE %s", "test_db"),
	})
	assert.NoError(t, err, "influx test db should be tear down successfully")
}

func TestUserHTTPServer(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err, "logger should be initiated successfully")

	sugar := logger.Sugar()
	db, teardown := testutil.MustNewDevelopmentDB()
	userStorage, err := storage.NewDB(sugar, db)
	assert.Nil(t, err, "user database should be initiated successfully")

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	assert.Nil(t, err, "influx client should be created successfully")

	defer tearDown(t, teardown, influxClient)

	// create test db
	_, err = influxClient.Query(client.Query{
		Command: "CREATE DATABASE test_db",
	})
	assert.Nil(t, err, "influx should create database successfully")

	influxStorage, err := storage.NewInfluxStorage(
		sugar,
		"test_db",
		influxClient,
	)
	assert.Nil(t, err, "influx storage should be created successfully")

	s := NewServer(sugar, tokenrate.NewMock(), userStorage, "", influxStorage, nil)
	s.register()

	// test case
	const (
		requestEndpoint = "/users"
	)

	var (
		tests = []httputil.HTTPTestCase{
			{
				Msg:      "email is not valid",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Body: []byte(`
{
  "email": "test",
  "user_info": [
    {
      "address": "0xc9a658f87d7432ff897f31dce318f0856f66acb7",
      "timestamp": 1538380670000
    },
    {
      "address": "0x2ea6200a999f4c6c982be525f8dc294f14f4cb08",
      "timestamp": 1538380682000
    }
  ]
}`),
				Assert: expectBadRequest,
			},
			{
				Msg:      "user address is empty",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Body: []byte(`
{
  "email": "test@gmail.com"",
  "user_info": [
    {
      "address": "",
      "timestamp": 1538380670000
    },
    {
      "address": "0x2ea6200a999f4c6c982be525f8dc294f14f4cb08",
      "timestamp": 1538380682000
    }
  ]
}`),
				Assert: expectBadRequest,
			},
			{
				Msg:      "timestamp is empty",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Body: []byte(`
{
  "email": "test@gmail.com",
  "user_info": [
    {
      "address": "0xc9a658f87d7432ff897f31dce318f0856f66acb7",
    },
    {
      "address": "0x2ea6200a999f4c6c982be525f8dc294f14f4cb08",
      "timestamp": 1538380682000
    }
  ]
}`),
				Assert: expectBadRequest,
			},
			{
				Msg:      "invalid user address",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Body: []byte(`
			{
			 "email": "test@gmail.com",
			 "user_info": [
			   {
			     "address": "0x1497340a82",
			     "timestamp": 1538380670000
			   },
			   {
			     "address": "not a valid address",
			     "timestamp": 1538380682000
			   }
			 ]
			}`),
				Assert: expectBadRequest,
			},
			{
				Msg:      "update correct user addresses",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Body: []byte(`
{
  "email": "test@gmail.com",
  "user_info": [
    {
      "address": "0xc9a658f87d7432ff897f31dce318f0856f66acb7",
      "timestamp": 1538380670000
    },
    {
      "address": "0x2ea6200a999f4c6c982be525f8dc294f14f4cb08",
      "timestamp": 1538380682000
    }
  ]
}`),
				Assert: expectSuccess,
			},
			{
				Msg:      "address is not unique",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Body: []byte(`
{
  "email": "test2@gmail.com",
  "user_info": [
    {
      "address": "0xc9a658f87d7432ff897f31dce318f0856f66acb7",
      "timestamp": 1538380670000
    },
    {
      "address": "0x2ea6200a999f4c6c982be525f8dc294f14f4cb08",
      "timestamp": 1538380682000
    }
  ]
}`),
				Assert: expectSuccess,
			},
			{
				Msg:      "user remove address",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Body: []byte(`
					{
						"email": "test2@gmail.com",
						"user_info": [
							{
								"address": "0x2ea6200a999f4c6c982be525f8dc294f14f4cb08",
								"timestamp": 1538380682000
							}
						]
					}`),
				Assert: expectSuccess,
			},
		}
	)

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, s.r) })
	}
}
