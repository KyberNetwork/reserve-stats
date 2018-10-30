package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	fromParams   uint64 = 12342082
	tradeLogsURL        = "127.0.0.1:7000"
	testAddr            = "127.0.0.1:7001"
)

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
	t.Skip()
	go mockServer()
	// assert.Nil(t, err, "mockserver should be start ok")
	tradeLogsAddr := fmt.Sprintf("http://%s", tradeLogsURL)
	testServer, err := NewServer(testAddr, tradeLogsAddr, "", "")
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
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, testServer.r) })
	}
}
