package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	postgresHost     = "127.0.0.1"
	postgresPort     = 5432
	postgresUser     = "reserve_stats"
	postgresPassword = "reserve_stats"
	postgresDatabase = "reserve_stats"
)

func newTestDB(sugar *zap.SugaredLogger) (*storage.UserDB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDatabase,
	)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return storage.NewDB(sugar, db)
}

func tearDown(t *testing.T, storage *storage.UserDB, influxClient client.Client) {
	assert.Nil(t, storage.DeleteAllTables(), "database should be deleted completely")
	_, err := influxClient.Query(client.Query{
		Command: fmt.Sprintf("DROP DATABASE %s", "test_db"),
	})
	assert.Nil(t, err, "influx test db should be tear down successfully")
}

func TestUserHTTPServer(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err, "logger should be initiated successfully")

	sugar := logger.Sugar()
	userStorage, err := newTestDB(sugar)
	assert.Nil(t, err, "user database should be initiated successfully")

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	assert.Nil(t, err, "influx client should be created successfully")

	defer tearDown(t, userStorage, influxClient)

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

	s := NewServer(sugar, tokenrate.NewMock(), userStorage, "", influxStorage)
	s.register()

	// test case
	const (
		requestEndpoint = "/users"
		queryAddress    = "0xc9a658f87d7432ff897f31dce318f0856f66acb7"
		nonKycAddress   = "0xb8df4cf4b7ad086cd5139a75033566164e41a0b4"
	)

	var (
		tests = []httputil.HTTPTestCase{
			{
				Msg:      "empty db",
				Endpoint: fmt.Sprintf("%s?address=%s", requestEndpoint, queryAddress),
				Method:   http.MethodGet,
				Assert:   expectSuccess,
			},
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
				// TODO: should not return 500
				Assert: expectInternalServerError,
			},
			{
				Msg:      "user is not kyced",
				Endpoint: fmt.Sprintf("%s?address=%s", requestEndpoint, nonKycAddress),
				Method:   http.MethodGet,
				Assert:   expectNonKYCed,
			},
			{
				Msg:      "user is kyced",
				Endpoint: fmt.Sprintf("%s?address=%s", requestEndpoint, queryAddress),
				Method:   http.MethodGet,
				Assert:   expectKYCed,
			},
			{
				Msg:      "user remove address",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Body: []byte(`
					{
						"email": "test@gmail.com",
						"user_info": [
							{
								"address": "0x2ea6200a999f4c6c982be525f8dc294f14f4cb08",
								"timestamp": 1538380682000
							}
						]
					}`),
				Assert: expectSuccess,
			},
			{
				Msg:      "address is removed",
				Endpoint: fmt.Sprintf("%s?address=%s", requestEndpoint, queryAddress),
				Method:   http.MethodGet,
				Assert:   expectNonKYCed,
			},
			{
				Msg:      "address have not trade, rich is false",
				Endpoint: fmt.Sprintf("%s?address=%s", requestEndpoint, queryAddress),
				Method:   http.MethodGet,
				Assert:   expectRichStatus,
			},
		}
	)

	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, s.r) })
	}
}
