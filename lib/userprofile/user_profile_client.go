package userprofile

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
)

// Client is the the implementation to query userprofile info.
type Client struct {
	sugar  *zap.SugaredLogger
	client *http.Client
	url    string
}

func (c *Client) newRequest(method, endpoint string, params map[string]string) (*http.Request, error) {
	logger := c.sugar.With(
		"func", "lib/userprofile/client.go/newRequest",
		"method", method,
		"endpoint", endpoint,
	)

	logger.Debug("creating new Core API HTTP request")

	url := fmt.Sprintf("%s/%s",
		strings.TrimRight(c.url, "/"),
		strings.Trim(endpoint, "/"),
	)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	logger = logger.With("raw_query", req.URL.RawQuery)

	return req, nil
}

// NewClient creates a new core client instance.
func NewClient(sugar *zap.SugaredLogger, url string) (*Client, error) {
	const timeout = time.Minute
	client := &http.Client{Timeout: timeout}
	return &Client{sugar: sugar, url: url, client: client}, nil
}

// LookUpUserProfile will look for the UserProfile of input addr from server
func (c *Client) LookUpUserProfile(addr ethereum.Address) (UserProfile, error) {
	const endpoint = "/api/wallet_info"
	var (
		params = make(map[string]string)
		result = UserProfile{}
	)
	params["wallet_address"] = addr.Hex()
	req, err := c.newRequest(http.MethodGet, endpoint, params)
	if err != nil {
		return result, err
	}

	rsp, err := c.client.Do(req)
	if err != nil {
		return result, err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected return code: %d", rsp.StatusCode)
	}

	if err = json.NewDecoder(rsp.Body).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}
