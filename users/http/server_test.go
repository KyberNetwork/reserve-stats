package http

import (
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/users/stats"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

const (
	host             = ":9000"
	postgresHost     = "127.0.0.1:5432"
	postgresUser     = "reserve_stats"
	postgresPassword = "reserve_stats"
	postgresDatabase = "reserve_stats"
)

func connectToTestDB(sugar *zap.SugaredLogger) *storage.UserDB {
	return storage.NewDB(
		sugar,
		pg.Connect(&pg.Options{
			Addr:     postgresHost,
			User:     postgresUser,
			Password: postgresPassword,
			Database: postgresDatabase,
		},
		),
	)
}

func tearDown(t *testing.T, storage *storage.UserDB) {
	if err := storage.DeleteAllTables(); err != nil {
		t.Fatal("Cannot clear db after test")
	}
}

func TestUserHTTPServer(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	sugar := logger.Sugar()
	userStorage := connectToTestDB(sugar)
	userStats := stats.NewUserStats(tokenrate.NewMock(), userStorage)
	defer tearDown(t, userStorage)

	s := NewServer(sugar, userStats, host)
	s.register()

	zap.S().Infof("Server instance: %+v", s)

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
				Endpoint: fmt.Sprintf("%s/%s", requestEndpoint, queryAddress),
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
			//			{
			//				Msg:      "invalid user address",
			//				Endpoint: requestEndpoint,
			//				Method:   http.MethodPost,
			//				Body: `
			//{
			//  "email": "test",
			//  "user_info": [
			//    {
			//      "address": "0x001122",
			//      "timestamp": 1538380670000
			//    },
			//    {
			//      "address": "0x2ea6200a999f4c6c982be525f8dc294f14f4cb08",
			//      "timestamp": 1538380682000
			//    }
			//  ]
			//}`,
			//				Assert: expectBadRequest,
			//			},
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
				Endpoint: fmt.Sprintf("%s/%s", requestEndpoint, nonKycAddress),
				Method:   http.MethodGet,
				Assert:   expectNonKYCed,
			},
			{
				Msg:      "user is kyced",
				Endpoint: fmt.Sprintf("%s/%s", requestEndpoint, queryAddress),
				Method:   http.MethodGet,
				Assert:   expectKYCed,
			},
		}
	)

	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, s.r) })
	}
}
