package ipinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

type responseNormal struct {
	Country string `json:"country"`
}

type responseError struct {
	Error string `json:"error"`
}

func TestIPLocatorHTTPServer(t *testing.T) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	s, err := NewHTTPServer(sugar, "testdata", fmt.Sprintf(":%d", httputil.IPLocatorPort))
	if err != nil {
		t.Error("Could not create HTTP server", "error", err.Error())
	}
	s.register()
	// test case
	const (
		requestEndpoint = "/ip"
		correctIP       = "81.2.69.142"
		wrongIPFormat   = "22"
	)
	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test valid IP",
			Endpoint: fmt.Sprintf("%s/%s", requestEndpoint, correctIP),
			Method:   http.MethodGet,
			Assert:   validResult,
		},
		{
			Msg:      "Test invalid IP",
			Endpoint: fmt.Sprintf("%s/%s", requestEndpoint, wrongIPFormat),
			Method:   http.MethodGet,
			Assert:   invalidResult,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, s.r) })
	}
}

func validResult(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	if resp.Code != http.StatusOK {
		t.Error("wrong return code", "return code", resp.Code, "expected", http.StatusOK)
	}
	var response responseNormal
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Error("Could not decode out put", "err", err)
	}
	if response.Country != "GB" {
		t.Error("Get location of ip was incorrect", "result", response.Country, "expected", "GB")
	}

}

func invalidResult(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	if resp.Code != http.StatusBadRequest {
		t.Error("wrong return code", "return code", resp.Code, "expected", http.StatusBadRequest)
	}
	var response responseError
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Error("Could not decode out put", "err", err)
	}
	if response.Error != ErrInvalidIP {
		t.Error("Wrong response error message", "result", response.Error, "expected", ErrInvalidIP)
	}
}
