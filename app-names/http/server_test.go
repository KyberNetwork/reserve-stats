package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/app-names/common"
	"github.com/KyberNetwork/reserve-stats/app-names/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestAppNameHTTPServer(t *testing.T) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	db, fn := testutil.MustNewDevelopmentDB()
	defer func() { assert.NoError(t, fn()) }()

	appNameStorage, err := storage.NewAppNameDB(sugar, db)
	assert.Nil(t, err, "database should be initiated successfully")
	logger := sugar.Desugar()
	s, err := NewServer("", appNameStorage, logger)
	assert.Nil(t, err, "server should be create successfully")
	s.register()

	// test case
	const (
		requestEndpoint = "/applications"
		appID           = 1
	)

	var (
		tests = []httputil.HTTPTestCase{
			{
				Msg:      "get non existing app",
				Endpoint: fmt.Sprintf("%s/%d", requestEndpoint, appID),
				Method:   http.MethodGet,
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusNotFound, resp.Code)
				},
			},
			{
				Msg:      "fail to create with empty app name",
				Method:   http.MethodPost,
				Endpoint: requestEndpoint,
				Body: []byte(`
				{
					"addresses": [
						"0x3baE9b9e1dca462Ad8827f62F4A8b5b3714d7700",
						"0x804aDa8c08A2E8ecff1a6535bf28DC4f1EfF4f8e"
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusBadRequest, resp.Code)
				},
			},
			{
				Msg:      "fail to create with invalid address",
				Method:   http.MethodPost,
				Endpoint: requestEndpoint,
				Body: []byte(`
				{
					"name": "first_app",
					"addresses": [
						"WTF-invalid-address",
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusBadRequest, resp.Code)
				},
			},
			{
				Msg:      "create success",
				Method:   http.MethodPost,
				Endpoint: requestEndpoint,
				Body: []byte(`
				{
					"name": "first_app",
					"addresses": [
						"0x3baE9b9e1dca462Ad8827f62F4A8b5b3714d7700",
						"0x804aDa8c08A2E8ecff1a6535bf28DC4f1EfF4f8e"
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					var result common.Application
					assert.Equal(t, http.StatusCreated, resp.Code)
					assert.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
					assert.Equal(t, int64(1), result.ID)
					assert.Equal(t,
						[]ethereum.Address{
							ethereum.HexToAddress("0x3baE9b9e1dca462Ad8827f62F4A8b5b3714d7700"),
							ethereum.HexToAddress("0x804aDa8c08A2E8ecff1a6535bf28DC4f1EfF4f8e"),
						},
						result.Addresses,
					)
					assert.Equal(t, "first_app", result.Name)
				},
			},
			{
				Msg:      "get existing app",
				Endpoint: fmt.Sprintf("%s/%d", requestEndpoint, appID),
				Method:   http.MethodGet,
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					var result common.Application
					assert.Equal(t, http.StatusOK, resp.Code)
					assert.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
					assert.Equal(t, int64(1), result.ID)
					assert.Equal(t,
						[]ethereum.Address{
							ethereum.HexToAddress("0x3baE9b9e1dca462Ad8827f62F4A8b5b3714d7700"),
							ethereum.HexToAddress("0x804aDa8c08A2E8ecff1a6535bf28DC4f1EfF4f8e"),
						},
						result.Addresses,
					)
					assert.Equal(t, "first_app", result.Name)
				},
			},
			{
				Msg:      "create app with conflict address",
				Method:   http.MethodPost,
				Endpoint: requestEndpoint,
				Body: []byte(`
				{
					"name": "first_app_conflict_address",
					"addresses": [
						"0x3baE9b9e1dca462Ad8827f62F4A8b5b3714d7700"
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusCreated, resp.Code)
				},
			},
			{
				Msg:      "update addresses",
				Method:   http.MethodPost,
				Endpoint: requestEndpoint,
				Body: []byte(`
				{
					"id": 1,
					"name": "first_app",
					"addresses": [
						"0x587ecf600d304f831201c30ea0845118dd57516e"
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					var result common.Application
					assert.Equal(t, http.StatusOK, resp.Code)
					assert.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
					assert.Equal(t, int64(1), result.ID)
					assert.Equal(t,
						[]ethereum.Address{
							ethereum.HexToAddress("0x587ecf600d304f831201c30ea0845118dd57516e"),
						},
						result.Addresses,
					)
					assert.Equal(t, "first_app", result.Name)
				},
			},
			{
				Msg:      "update application name",
				Method:   http.MethodPost,
				Endpoint: requestEndpoint,
				Body: []byte(`
				{
					"id": 1,
					"name": "first_app_new_edition",
					"addresses": ["0x587ecf600d304f831201c30ea0845118dd57516e"]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					var result common.Application
					log.Printf("%+v", resp)
					assert.Equal(t, http.StatusOK, resp.Code)
					assert.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
					assert.Equal(t, int64(1), result.ID)
					assert.Equal(t,
						[]ethereum.Address{
							ethereum.HexToAddress("0x587ecf600d304f831201c30ea0845118dd57516e"),
						},
						result.Addresses,
					)
					assert.Equal(t, "first_app_new_edition", result.Name)
				},
			},
			{
				Msg:      "get all apps",
				Endpoint: requestEndpoint,
				Method:   http.MethodGet,
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusOK, resp.Code)
				},
			},
			{
				Msg:      "get application with name filter",
				Endpoint: fmt.Sprintf("%s?name=first_app_new_edition", requestEndpoint),
				Method:   http.MethodGet,
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					var result []common.Application
					assert.Equal(t, http.StatusOK, resp.Code)
					assert.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
					assert.Len(t, result, 1)
					if len(result) > 0 {
						app := result[0]
						assert.Equal(t, int64(1), app.ID)
						assert.Equal(t,
							[]ethereum.Address{
								ethereum.HexToAddress("0x587ecf600d304f831201c30ea0845118dd57516e"),
							},
							app.Addresses,
						)
						assert.Equal(t, "first_app_new_edition", app.Name)
					}
				},
			},
			{
				Msg:      "get application with address filter",
				Endpoint: fmt.Sprintf("%s?address=0x587ecf600d304f831201c30ea0845118dd57516e", requestEndpoint),
				Method:   http.MethodGet,
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					var result []common.Application
					assert.Equal(t, http.StatusOK, resp.Code)
					assert.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
					assert.Len(t, result, 1)
					if len(result) > 0 {
						app := result[0]
						assert.Equal(t, int64(1), app.ID)
						assert.Equal(t,
							[]ethereum.Address{
								ethereum.HexToAddress("0x587ecf600d304f831201c30ea0845118dd57516e"),
							},
							app.Addresses,
						)
						assert.Equal(t, "first_app_new_edition", app.Name)
					}
				},
			},
			{
				Msg:      "update app not exist",
				Method:   http.MethodPut,
				Endpoint: fmt.Sprintf("%s/%d", requestEndpoint, 100),
				Body: []byte(`
				{
					"name": "app_100",
					"addresses": [
						"0xd8c67d024db85b271b6f6eeac5234e29c4d6bbb5"
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusNotFound, resp.Code)
				},
			},
			{
				Msg:      "update address with invalid address",
				Method:   http.MethodPut,
				Endpoint: fmt.Sprintf("%s/%d", requestEndpoint, appID),
				Body: []byte(`
				{
					"addresses": [
						"OMG-INVALID-ADDRESS",
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusBadRequest, resp.Code)
				},
			},
			{
				Msg:      "update address success",
				Method:   http.MethodPut,
				Endpoint: fmt.Sprintf("%s/%d", requestEndpoint, appID),
				Body: []byte(`
				{
					"name": "first_app_updated",
					"addresses": [
						"0xde6a6fb70b0375d9c761f67f2db3de97f21362dc"
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					var result common.Application
					assert.Equal(t, http.StatusOK, resp.Code)
					assert.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
					assert.Equal(t, int64(1), result.ID)
					assert.Equal(t,
						[]ethereum.Address{
							ethereum.HexToAddress("0xde6a6fb70b0375d9c761f67f2db3de97f21362dc"),
						},
						result.Addresses,
					)
					assert.Equal(t, "first_app_updated", result.Name)
				},
			},
			{
				Msg:      "delete non existing application",
				Method:   http.MethodDelete,
				Endpoint: fmt.Sprintf("%s/%d", requestEndpoint, 101),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusNotFound, resp.Code)
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
			{
				Msg:      "get non existing app",
				Endpoint: fmt.Sprintf("%s/%d", requestEndpoint, appID),
				Method:   http.MethodGet,
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusNotFound, resp.Code)
				},
			},
			{
				Msg:      "get inactive apps",
				Endpoint: fmt.Sprintf("%s?active=false", requestEndpoint),
				Method:   http.MethodGet,
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					var result []common.Application
					assert.Equal(t, http.StatusOK, resp.Code)
					log.Printf("%+v", resp)
					assert.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
					if len(result) > 0 {
						assert.Equal(t, int64(1), result[0].ID)
						assert.Equal(t,
							[]ethereum.Address{
								ethereum.HexToAddress("0xde6a6fb70b0375d9c761f67f2db3de97f21362dc"),
							},
							result[0].Addresses,
						)
					}
				},
			},
			{
				Msg:      "re-active delete app",
				Endpoint: requestEndpoint,
				Method:   http.MethodPost,
				Body: []byte(`
				{
					"name": "first_app_updated",
					"addresses": [
						"0xde6a6fb70b0375d9c761f67f2db3de97f21362dc"
					]
				}
				`),
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusOK, resp.Code)
				},
			},
		}
	)

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, s.r) })
	}
}
