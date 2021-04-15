package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KyberNetwork/httpsign-utils/sign"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"go.uber.org/zap"
)

const timeout = time.Minute * 5

// Client is implementation of tradelog client
type Client struct {
	host   string
	sugar  *zap.SugaredLogger
	client *http.Client

	accessKeyID     string
	secretAccessKey string
}

// TradeLogClientOption option to Client constructor
type TradeLogClientOption func(*Client)

// WithAuth is option to create Client with auth keys
func WithAuth(accessKeyID, secretAccessKey string) TradeLogClientOption {
	return func(c *Client) {
		c.accessKeyID = accessKeyID
		c.secretAccessKey = secretAccessKey
	}
}

// NewTradeLogClient creates a new tradelog client instance.
func NewTradeLogClient(sugar *zap.SugaredLogger, host string, options ...TradeLogClientOption) *Client {
	c := &Client{
		host:   host,
		sugar:  sugar,
		client: &http.Client{Timeout: timeout},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// GetTradeLogs get trade logs from `fromTime` to `toTime`
func (c *Client) GetTradeLogs(fromTime, toTime uint64) ([]common.Tradelog, error) {
	var (
		tradeLogs []common.Tradelog
		url       = fmt.Sprintf("%s/trade-logs?from=%d&to=%d", c.host, fromTime, toTime)
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return tradeLogs, err
	}
	if c.accessKeyID != "" && c.secretAccessKey != "" {
		req, err = sign.Sign(req, c.accessKeyID, c.secretAccessKey)
		if err != nil {
			return tradeLogs, err
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return tradeLogs, err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			c.sugar.Errorw("failed to close body", "err", cErr.Error())
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return tradeLogs, fmt.Errorf("unexpcted status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&tradeLogs); err != nil {
		return tradeLogs, err
	}
	return tradeLogs, nil
}
