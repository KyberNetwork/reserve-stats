package core

import (
	"crypto/hmac"
	"crypto/sha512"
	"go.uber.org/zap"
	"net/http"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Client is the the real implementation of core client interface.
type Client struct {
	sugar      *zap.SugaredLogger
	client     *http.Client
	url        string
	signingKey string
}

type commonResponse struct {
	Reason  string `json:"reason"`
	Success bool   `json:"success"`
}

// sortByKey sort all the params by key in string order
// This is required for the request to be signed correctly
func sortByKey(params map[string]string) map[string]string {
	newParams := make(map[string]string, len(params))
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		newParams[key] = params[key]
	}
	return newParams
}

func (c *Client) sign(msg string) (string, error) {
	mac := hmac.New(sha512.New, []byte(c.signingKey))
	if _, err := mac.Write([]byte(msg)); err != nil {
		return "", err
	}
	return common.Bytes2Hex(mac.Sum(nil)), nil
}

func (c *Client) newRequest(method, endpoint string, params map[string]string) (*http.Request, error) {
	logger := c.sugar.With(
		"method", method,
		"endpoint", endpoint,
	)

	logger.Debug("creating new Core API HTTP request")

	req, err := http.NewRequest(method, c.url+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	params = sortByKey(params)
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	logger = logger.With("raw_query", req.URL.RawQuery)

	_, ok := params["nonce"]
	if ok {
		logger.Debug("nonce is available, signing message")
		signed, err := c.sign(q.Encode())
		if err != nil {
			return nil, err
		}
		req.Header.Add("signed", signed)
		logger = logger.With("signed", signed)
	}

	logger.Debug("Core API HTTP request created")
	return req, nil
}

// NewClient creates a new core client instance.
func NewClient(sugar *zap.SugaredLogger, url, signingKey string) (*Client, error) {
	const timeout = time.Minute
	client := &http.Client{Timeout: timeout}
	return &Client{sugar: sugar, url: url, signingKey: signingKey, client: client}, nil
}
