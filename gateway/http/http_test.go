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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/KyberNetwork/httpsign-utils/authenticator"
	"github.com/KyberNetwork/httpsign-utils/sign"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

const (
	fromParams      = 12342082
	tradeLogsURL    = "127.0.0.1:7000"
	testAddr        = "127.0.0.1:7001"
	writeKeyID      = "writekeyID"
	writeSigningKey = "fdsr122541"
	readKeyID       = "readKeyID"
	readSigningKey  = "xxx123xxx"
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

	req, err = sign.Sign(req, key, signingKey)
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
	return r.Run(tradeLogsURL)
}

func TestReverseProxy(t *testing.T) {
	go func(t *testing.T) {
		assert.NoError(t, mockServer())
	}(t)

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
	keyPairs := []authenticator.KeyPair{
		{
			AccessKeyID:     readKeyID,
			SecretAccessKey: readSigningKey,
		},
		{
			AccessKeyID:     writeKeyID,
			SecretAccessKey: writeSigningKey,
		},
	}
	auth, err := authenticator.NewAuthenticator(keyPairs...)
	if err != nil {
		t.Fatal(err)
	}
	perm, err := NewPermissioner(readKeyID, writeKeyID)
	if err != nil {
		t.Fatal(err)
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	testServer, err := NewServer(testAddr, auth, perm, logger,
		WithTradeLogURL(testURL),
		WithReserveRatesURL(testURL),
		WithUserURL(testURL),
		WithPriceAnalyticURL(testURL),
	)
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
				assert.Equal(t, result.FromTime, uint64(fromParams), "Reverse proxy should receive correct params")
			},
		},
	}

	for _, tc := range testCaseReadKey {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { runHTTPTestCase(t, tc, testServer.r, readKeyID, readSigningKey) })
	}

}
