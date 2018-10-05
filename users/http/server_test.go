package http

import (
	"fmt"
	"github.com/go-pg/pg"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/ethrate"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/users/stats"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"go.uber.org/zap"
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

func tearDownDB(t *testing.T, storage *storage.UserDB) {
	if err := storage.CloseDBConnection(); err != nil {
		t.Fatalf("cannot teardown db connection")
	}
}

func TestUserHTTPServer(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	sugar := logger.Sugar()
	userStorage := connectToTestDB(sugar)
	cmc := ethrate.NewCMCRate(sugar)
	// sleep so cmc fetcher can get rate from cmc
	time.Sleep(1 * time.Second)

	userStats := stats.NewUserStats(cmc, userStorage)
	defer tearDownDB(t, userStorage)

	s := NewServer(sugar, userStats, host)
	s.register()

	zap.S().Infof("Server instance: %+v", s)

	// test case
	const (
		requestEndpoint = "/users"
		userEmail       = "test@gmail.com"
		wrongUserEmail  = "test"

		userAddresses          = "0xc9a658f87d7432ff897f31dce318f0856f66acb7-0x2ea6200a999f4c6c982be525f8dc294f14f4cb08"
		wrongUserAddresses     = "wrong-address_0x13197"
		wrongNumberOfAddresses = "0xc9a658f87d7432ff897f31dce318f0856f66acb7-0x2ea6200a999f4c6c982be525f8dc294f14f4cb08-0x4e012a6445ba2a590b8b1ee4e95d03e345a0c2e5"
		userTimeStamp          = "1538380670000-1538380682000"
		queryAddress           = "0xc9a658f87d7432ff897f31dce318f0856f66acb7"
		nonKycAddress          = "0xb8df4cf4b7ad086cd5139a75033566164e41a0b4"
	)

	var tests = []httputil.HTTPTestCase{
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
			Data: map[string]string{
				"user":       wrongUserEmail,
				"addresses":  userAddresses,
				"timestamps": userTimeStamp,
			},
			Assert: expectBadRequest,
		},
		{
			Msg:      "update user addresses with wrong number of addresses",
			Endpoint: requestEndpoint,
			Method:   http.MethodPost,
			Data: map[string]string{
				"user":       userEmail,
				"addresses":  wrongNumberOfAddresses,
				"timestamps": userTimeStamp,
			},
			Assert: expectBadRequest,
		},
		// { msg:      "wrong user addresses",
		// 	endpoint: requestEndpoint,
		// 	method:   http.MethodPost,
		// 	data: map[string]string{
		// 		"user":       userEmail,
		// 		"addresses":  wrongUserAddresses,
		// 		"timestamps": userTimeStamp,
		// 	},
		// 	assert: httputil.ExpectFailure,
		// },
		{
			Msg:      "request malformed",
			Endpoint: requestEndpoint,
			Method:   http.MethodPost,
			Assert:   expectBadRequest,
		},
		{
			Msg:      "update correct user addresses",
			Endpoint: requestEndpoint,
			Method:   http.MethodPost,
			Data: map[string]string{
				"user":       userEmail,
				"addresses":  userAddresses,
				"timestamps": userTimeStamp,
			},
			Assert: expectSuccess,
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

	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, s.r) })
	}
}
