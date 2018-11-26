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
	fromParams      uint64 = 12342082
	tradeLogsURL           = "127.0.0.1:7000"
	testAddr               = "127.0.0.1:7001"
	writeKeyID             = "writekeyID"
	writeSigningKey        = "fdsr122541"
	readKeyID              = "readKeyID"
	readSigningKey         = "xxx123xxx"
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

func runHTTPTestCase(t *testing.T, tc httputil.HTTPTestCase, handler http.Handler, key string, signingKey string) {
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

	req, err = httputil.Sign(req, key, signingKey)
	assert.Nil(t, err, "sign request should be success")

	resp := NewWrappedRecorder()
	handler.ServeHTTP(resp, req)
	tc.Assert(t, resp.ResponseRecorder)
}

// dummy response to check proxy
func getTradeLogs(c *gin.Context) {
	// check getting correct params
	fromTimeParam := c.Query("fromTime")
	fromTime, err := strconv.ParseUint(fromTimeParam, 10, 64)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{"from": fromTime},
	)
}

// this function is to test http signing
// if request come here then it pass signing
// return ok
func updateUsers(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{},
	)
}

func mockServer() error {
	r := gin.Default()
	r.GET("/trade-logs", getTradeLogs)
	r.POST("/users", updateUsers)
	return r.Run(tradeLogsURL)
}

func TestReverseProxy(t *testing.T) {
	go mockServer()

	// ping until mock server response
	ok := make(chan bool)
	go func() {
		for {
			_, err := http.Get("http://127.0.0.1:7000/trade-logs")
			if err == nil {
				ok <- true
			}
		}
	}()

	select {
	case <-ok:
		break
	case <-time.After(5 * time.Second):
		t.Fatal("mock server not running after 5 second")
	}

	// assert.Nil(t, err, "mockserver should be start ok")
	testURL := fmt.Sprintf("http://%s", tradeLogsURL)
	auth, err := NewAuthenticator(readKeyID, readSigningKey, writeKeyID, writeSigningKey)
	if err != nil {
		t.Fatal(err)
	}
	perm, err := NewPermissioner(readKeyID, writeKeyID)
	if err != nil {
		t.Fatal(err)
	}
	testServer, err := NewServer(testAddr, testURL, testURL, testURL, testURL, auth, perm)
	assert.Nil(t, err, "reverse proxy server should initiate successfully")

	var testCaseReadKey = []httputil.HTTPTestCase{
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
	var testCaseWriteKey = []httputil.HTTPTestCase{
		{
			Msg:      "test sign request body is empty",
			Endpoint: fmt.Sprintf("/users"),
			Method:   http.MethodPost,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			Msg:      "test sign request body is not empty",
			Endpoint: fmt.Sprintf("/users"),
			Method:   http.MethodPost,
			Body:     []byte(`{"user":"something"}`),
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
	}

	var testsFailedWithoutKey = []httputil.HTTPTestCase{
		{
			Msg:      "test sign request without body",
			Endpoint: fmt.Sprintf("/users"),
			Method:   http.MethodPost,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
			},
		},
	}

	var testFailedWhenWriteWithReadKey = []httputil.HTTPTestCase{
		{
			Msg:      "test sign write request with read Key",
			Endpoint: fmt.Sprintf("/users"),
			Method:   http.MethodPost,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
			},
		},
	}

	for _, tc := range testCaseReadKey {
		t.Run(tc.Msg, func(t *testing.T) { runHTTPTestCase(t, tc, testServer.r, readKeyID, readSigningKey) })
	}

	for _, tc := range testCaseWriteKey {
		t.Run(tc.Msg, func(t *testing.T) { runHTTPTestCase(t, tc, testServer.r, writeKeyID, writeSigningKey) })
	}

	for _, tc := range testsFailedWithoutKey {
		t.Run(tc.Msg, func(t *testing.T) { runHTTPTestCase(t, tc, testServer.r, "", "") })
	}

	for _, tc := range testFailedWhenWriteWithReadKey {
		t.Run(tc.Msg, func(t *testing.T) { runHTTPTestCase(t, tc, testServer.r, readKeyID, readSigningKey) })
	}
}
