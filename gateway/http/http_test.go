package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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

//WrappedRecorded wrap the gin response from proxy server
//added closed chan to fulfilled assert function
//This type is exported to fulfilled golint requirement
type WrappedRecorded struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func NewWrappedRecorder() *WrappedRecorded {
	return &WrappedRecorded{
		ResponseRecorder: httptest.NewRecorder(),
		closed:           make(chan bool, 1),
	}
}

func (c *WrappedRecorded) CloseNotify() <-chan bool {
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
	time, err := strconv.ParseUint(fromTime, 10, 64)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{"from": time},
	)
}

func mockServer() error {
	r := gin.Default()
	r.GET("/trade-logs", getTradeLogs)
	return r.Run(tradeLogsURL)
}

func TestReverseProxy(t *testing.T) {
	go mockServer()

	time.Sleep(5 * time.Second)

	// assert.Nil(t, err, "mockserver should be start ok")
	tradeLogsAddr := fmt.Sprintf("http://%s", tradeLogsURL)
	testServer, err := NewServer(testAddr, tradeLogsAddr, "", "", "123108")
	assert.Nil(t, err, "reverse proxy server should initiate successfully")

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
					t.Errorf("Reverse proxy not working, %s", err.Error())
				}
				assert.Equal(t, result.FromTime, fromParams, "Reverse proxy should receive correct params")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { runHTTPTestCase(t, tc, testServer.r) })
	}
}
