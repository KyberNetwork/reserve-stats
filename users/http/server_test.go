package http

import (
	"fmt"
	"github.com/go-pg/pg"
	"net/http"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/stats"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"github.com/KyberNetwork/tokenrate/coingecko"
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
	cmc := coingecko.New()
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
		secondUser      = "second@gmail.com"
		wrongUserEmail  = "test"
		queryAddress    = "0xc9a658f87d7432ff897f31dce318f0856f66acb7"
		nonKycAddress   = "0xb8df4cf4b7ad086cd5139a75033566164e41a0b4"
	)
	var (
		correctUserInfo = []common.Info{
			{
				Address:   "0xc9a658f87d7432ff897f31dce318f0856f66acb7",
				Timestamp: 1538380670000,
			},
			{
				Address:   "0x2ea6200a999f4c6c982be525f8dc294f14f4cb08",
				Timestamp: 1538380682000,
			},
		}
		addressIsEmpty = []common.Info{
			{
				Address:   "",
				Timestamp: 1538380670000,
			},
			{
				Address:   "0x2ea6200a999f4c6c982be525f8dc294f14f4cb08",
				Timestamp: 1538380682000,
			},
		}
		timestampEmpty = []common.Info{
			{
				Address: "0xc9a658f87d7432ff897f31dce318f0856f66acb7",
			},
			{
				Address:   "0x2ea6200a999f4c6c982be525f8dc294f14f4cb08",
				Timestamp: 1538380682000,
			},
		}
		// wrongAddress = []common.Info{
		// 	common.Info{
		// 		Address: "0xc9a658f87d7432ff897f31dce318f0856f66acb7",
		// 		Timestamp:1538380670000,
		// 	},
		// 	common.Info{
		// 		Address:"0x2ea6200a999f4c6c982be525f8dc294f14f4cb08",
		// 		Timestamp:1538380682000,
		// 	},
		// }
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
				Data: map[string]interface{}{
					"email":     wrongUserEmail,
					"user_info": correctUserInfo,
				},
				Assert: expectBadRequest,
			},
			{
				Msg:      "user address is empty",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Data: map[string]interface{}{
					"email":     userEmail,
					"user_info": addressIsEmpty,
				},
				Assert: expectBadRequest,
			},
			{
				Msg:      "timestamp is empty",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Data: map[string]interface{}{
					"email":     userEmail,
					"user_info": timestampEmpty,
				},
				Assert: expectBadRequest,
			},
			// { msg:      "wrong user addresses",
			// 	endpoint: requestEndpoint,
			// 	method:   http.MethodPost,
			// 	data: map[string]string{
			// 		"user":       userEmail,
			// 		"addresses":  wrongUserAddresses,
			// 	},
			// 	assert: httputil.ExpectFailure,
			// },
			{
				Msg:      "update correct user addresses",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Data: map[string]interface{}{
					"email":     userEmail,
					"user_info": correctUserInfo,
				},
				Assert: expectSuccess,
			},
			{
				Msg:      "address is not unique",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Data: map[string]interface{}{
					"email":     secondUser,
					"user_info": correctUserInfo,
				},
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
