package ipinfo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"go.uber.org/zap"
)

func TestIPLocatorHTTPServer(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Error("Get error while create logger", "error", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	s, err := NewHTTPServer(sugar, "testdata", int(httputil.IPLocatorPort))
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
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, s.r) })
	}
}

func validResult(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	if resp.Code != http.StatusOK {
		t.Error("wrong return code", "return code", resp.Code, "expected", http.StatusOK)
	}
}

func invalidResult(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	if resp.Code != http.StatusBadRequest {
		t.Error("wrong return code", "return code", resp.Code, "expected", http.StatusBadRequest)
	}
}
