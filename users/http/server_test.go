package http

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/users/cmc"
	"github.com/KyberNetwork/reserve-stats/users/http/httputil"
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

func connectToTestDB() *storage.UserDB {
	return storage.NewDB(
		postgresHost,
		postgresUser,
		postgresPassword,
		postgresDatabase,
	)
}

func tearDownDB(t *testing.T, storage *storage.UserDB) {
	if err := storage.CloseDBConnection(); err != nil {
		t.Fatalf("cannot teardown db connection")
	}
}

func TestUserHTTPServer(t *testing.T) {
	userStorage := connectToTestDB()
	cmc := cmc.NewCMCEthUSDRate()
	// sleep so cmc fetcher can get rate from cmc
	time.Sleep(1 * time.Second)

	userStats := stats.NewUserStats(cmc, userStorage)
	defer tearDownDB(t, userStorage)

	s := NewServer(userStats, host)
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

	var tests = []testCase{
		{
			msg:      "empty db",
			endpoint: fmt.Sprintf("%s/%s", requestEndpoint, queryAddress),
			method:   http.MethodGet,
			assert:   httputil.ExpectSuccess,
		},
		{
			msg:      "email is not valid",
			endpoint: requestEndpoint,
			method:   http.MethodPost,
			data: map[string]string{
				"user":       wrongUserEmail,
				"addresses":  userAddresses,
				"timestamps": userTimeStamp,
			},
			assert: httputil.ExpectBadRequest,
		},
		{
			msg:      "update user addresses with wrong number of addresses",
			endpoint: requestEndpoint,
			method:   http.MethodPost,
			data: map[string]string{
				"user":       userEmail,
				"addresses":  wrongNumberOfAddresses,
				"timestamps": userTimeStamp,
			},
			assert: httputil.ExpectFailure,
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
			msg:      "request malformed",
			endpoint: requestEndpoint,
			method:   http.MethodPost,
			assert:   httputil.ExpectFailure,
		},
		{
			msg:      "update correct user addresses",
			endpoint: requestEndpoint,
			method:   http.MethodPost,
			data: map[string]string{
				"user":       userEmail,
				"addresses":  userAddresses,
				"timestamps": userTimeStamp,
			},
			assert: httputil.ExpectSuccess,
		},
		{
			msg:      "user is not kyced",
			endpoint: fmt.Sprintf("%s/%s", requestEndpoint, nonKycAddress),
			method:   http.MethodGet,
			assert:   httputil.ExpectNonKYC,
		},
		{
			msg:      "user is kyced",
			endpoint: fmt.Sprintf("%s/%s", requestEndpoint, queryAddress),
			method:   http.MethodGet,
			assert:   httputil.ExpectKYC,
		},
	}

	for _, tc := range tests {
		t.Run(tc.msg, func(t *testing.T) { testHTTPRequest(t, tc, s.r) })
	}
}
