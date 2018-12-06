package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KyberNetwork/reserve-stats/app-names/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
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

func newTestDB(sugar *zap.SugaredLogger) (*storage.AppNameDB, error) {
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
	return storage.NewAppNameDB(sugar, db)
}

func tearDown(t *testing.T, storage *storage.AppNameDB) {
	assert.Nil(t, storage.DeleteAllTables(), "database should be deleted completely")
}

func TestAppNameHTTPServer(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err, "logger should be initiated successfully")

	sugar := logger.Sugar()
	appNameStorage, err := newTestDB(sugar)
	assert.Nil(t, err, "user database should be initiated successfully")

	defer tearDown(t, appNameStorage)

	s, err := NewServer("", appNameStorage, sugar)
	assert.Nil(t, err, "server should be create successfully")
	s.register()

	// test case
	const (
		requestEndpoint = "/application-names"
		appID           = 1
		appIDNotExist   = 2
	)

	var (
		tests = []httputil.HTTPTestCase{
			{
				Msg:      "empty db",
				Endpoint: fmt.Sprintf("%s/%d", requestEndpoint, appID),
				Method:   http.MethodGet,
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusInternalServerError, resp.Code)
				},
			},
			{
				Msg:      "create success",
				Method:   http.MethodPost,
				Endpoint: fmt.Sprintf("%s", requestEndpoint),
				Body: []byte(`
				{
					"app_name": "first_app",
					"addresses": [
						"0x3baE9b9e1dca462Ad8827f62F4A8b5b3714d7700",
						"0x804aDa8c08A2E8ecff1a6535bf28DC4f1EfF4f8e"
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusOK, resp.Code)
				},
			},
			{
				Msg:      "get all apps",
				Endpoint: fmt.Sprintf("%s", requestEndpoint),
				Method:   http.MethodGet,
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusOK, resp.Code)
				},
			},
			{
				Msg:      "update app not exist",
				Method:   http.MethodPut,
				Endpoint: fmt.Sprintf("%s/%d", requestEndpoint, appIDNotExist),
				Body: []byte(`
				{
					"app_name": "first_app",
					"addresses": [
						"0x804aDa8c08A2E8ecff1a6535bf28DC4f1EfF4f8e"
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusInternalServerError, resp.Code)
				},
			},
			{
				Msg:      "update address success",
				Method:   http.MethodPut,
				Endpoint: fmt.Sprintf("%s/%d", requestEndpoint, appID),
				Body: []byte(`
				{
					"app_name": "first_app",
					"addresses": [
						"0xde6a6fb70b0375d9c761f67f2db3de97f21362dc"
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusOK, resp.Code)
				},
			},
			{
				Msg:      "update address success",
				Method:   http.MethodPut,
				Endpoint: fmt.Sprintf("%s/%d", requestEndpoint, appID),
				Body: []byte(`
				{
					"app_name": "first_app",
					"addresses": [
						"0xde6a6fb70b0375d9c761f67f2db3de97f21362dc"
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusOK, resp.Code)
				},
			},
			{
				Msg:      "delete app success",
				Method:   http.MethodDelete,
				Endpoint: fmt.Sprintf("%s/%d", requestEndpoint, appID),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusOK, resp.Code)
				},
			},
		}
	)

	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, s.r) })
	}
}
