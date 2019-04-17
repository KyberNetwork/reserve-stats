package http

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/influxdata/influxdb/client/v2"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/users/storage"
)

const (
	requestEndpoint = "/users"
)

var (
	userStorage   *storage.UserDB
	err           error
	influxClient  client.Client
	influxStorage *storage.InfluxStorage
	s             *Server
)

func tearDown(storage *storage.UserDB, influxClient client.Client) error {
	if err := storage.DeleteAllTables(); err != nil {
		return err
	}
	_, err := influxClient.Query(client.Query{
		Command: fmt.Sprintf("DROP DATABASE %s", "test_db"),
	})
	return err
}

func TestUserHTTPServer(t *testing.T) {
	// test case
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

// RunHTTPTestCase run http request test case
func runHTTPTestCase(tc httputil.HTTPTestCase, handler http.Handler) {
	req, err := http.NewRequest(tc.Method, tc.Endpoint, bytes.NewBuffer(tc.Body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	q := req.URL.Query()
	for k, v := range tc.Params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		log.Fatal(err)
	}
}

func BenchmarkUserAPI(b *testing.B) {
	var (
		tc = httputil.HTTPTestCase{
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
		}
	)
	for i := 0; i < b.N; i++ {
		runHTTPTestCase(tc, s.r)
	}
}

func TestMain(m *testing.M) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	_, db := testutil.MustNewDevelopmentDB()
	if userStorage, err = storage.NewDB(sugar, db); err != nil {
		log.Fatal("user database should be initiated successfully")
	}

	if influxClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	}); err != nil {
		log.Fatal("influx client should be created successfully")
	}

	// create test db
	if _, err = influxClient.Query(client.Query{
		Command: "CREATE DATABASE test_db",
	}); err != nil {
		log.Fatal("influx should create database successfully")
	}

	if influxStorage, err = storage.NewInfluxStorage(
		sugar,
		"test_db",
		influxClient,
	); err != nil {
		log.Fatal("influx storage should be created successfully")
	}

	s = NewServer(sugar, tokenrate.NewMock(), userStorage, "", influxStorage)
	s.register()
	ret := m.Run()
	if err = tearDown(userStorage, influxClient); err != nil {
		log.Fatal(err)
	}
	os.Exit(ret)
}
