package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	fromParams   uint64 = 12342082
	tradeLogsURL        = "127.0.0.1:7000"
	testAddr            = "127.0.0.1:7001"
)

type wrappedRecorded struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func NewWrappedRecorder() *wrappedRecorded {
	return &wrappedRecorded{
		ResponseRecorder: httptest.NewRecorder(),
		closed:           make(chan bool, 1),
	}
}

func (c *wrappedRecorded) CloseNotify() <-chan bool {
	return c.closed
}

func runHTTPTestCase(t *testing.T, tc httputil.HTTPTestCase, handler http.Handler) {
	t.Helper()
	req, err := http.NewRequest(tc.Method, tc.Endpoint, bytes.NewBuffer(tc.Body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	q := req.URL.Query()
	for k, v := range tc.Params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp := NewWrappedRecorder()
	handler.ServeHTTP(resp, req)
	tc.Assert(t, resp.ResponseRecorder)
}

// dummy response to check proxy
func getTradeLogs(c *gin.Context) {
	// check getting correct params
	fromTime := c.Query("fromTime")
	c.JSON(
		http.StatusOK,
		fromTime,
	)
}

func mockServer() error {
	r := gin.Default()
	r.GET("/trade-logs", getTradeLogs)
	return r.Run(tradeLogsURL)
}

func TestReverseProxy(t *testing.T) {
	go mockServer()

	// TODO: call the mock server with small timeout until it response
	time.Sleep(time.Second)

	// assert.Nil(t, err, "mockserver should be start ok")
	tradeLogsAddr := fmt.Sprintf("http://%s", tradeLogsURL)
	testServer, err := NewServer(testAddr, tradeLogsAddr, "", "", "123108")
	assert.Nil(t, err, "reverse proxy server should initiate successfully")
	log.Printf("%+v", testServer)

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test reverse proxy",
			Endpoint: fmt.Sprintf("/trade-logs?fromTime=%d", fromParams),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				var result struct {
					FromTime uint64 `json:"from"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Reverse proxy not working")
				}

				assert.Equal(t, result.FromTime, fromParams, "Reverse proxy should receive correct params")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { runHTTPTestCase(t, tc, testServer.r) })
	}
}
