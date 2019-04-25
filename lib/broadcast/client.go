package broadcast

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KyberNetwork/httpsign-utils/sign"
	"go.uber.org/zap"
)

// errNotFoundMsg is the error returned by broadcast API
// if the given transaction address not found.
var errNotFoundMsg = "Can not find the transaction. Check Tx again"

// Client is the the real implementation of broadcast client interface
type Client struct {
	host   string
	sugar  *zap.SugaredLogger
	client *http.Client

	accessKeyID     string
	secretAccessKey string
}

type tradeLogGeoInfoResp struct {
	Success bool   `json:"success"`
	Err     string `json:"err"`
	Data    struct {
		IP      string `json:"IP"`
		Country string `json:"Country"`
	} `json:"data"`
}

const timeout = time.Minute * 5

// ClientOption option to Client constructor
type ClientOption func(*Client)

// WithAuth is option to create Client with auth keys
func WithAuth(accessKeyID, secretAccessKey string) ClientOption {
	return func(c *Client) {
		c.accessKeyID = accessKeyID
		c.secretAccessKey = secretAccessKey
	}
}

// NewClient creates a new broadcast client instance.
func NewClient(sugar *zap.SugaredLogger, host string, options ...ClientOption) *Client {
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

// GetTxInfo get ip, country info of a tx
func (c *Client) GetTxInfo(tx string) (ip string, country string, err error) {
	url := fmt.Sprintf("%s/get-tx-info/%s", c.host, tx)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}
	if c.accessKeyID != "" && c.secretAccessKey != "" {
		req, err = sign.Sign(req, c.accessKeyID, c.secretAccessKey)
		if err != nil {
			return "", "", err
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			c.sugar.Errorw("failed to close body", "err", cErr.Error())
		}
	}()
	response := tradeLogGeoInfoResp{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", "", err
	}
	if !response.Success {
		if response.Err == errNotFoundMsg {
			c.sugar.Debugw("transaction not found", "tx", tx, "err", response.Err)
			return "", "", nil
		}
		c.sugar.Errorw("server returns unknown error")
		return "", "", err
	}
	return response.Data.IP, response.Data.Country, nil
}
