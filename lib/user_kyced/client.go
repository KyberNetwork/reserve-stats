package userkyced

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// Client is the the implementation to query user kyced status info.
type Client struct {
	sugar        *zap.SugaredLogger
	client       *http.Client
	url          string
	signingKey   string
	signingKeyID string
}

func (c *Client) newRequest(method, endpoint string, params map[string]string) (*http.Request, error) {
	logger := c.sugar.With(
		"func", "lib/user_kyced/client.go/newRequest",
		"method", method,
		"endpoint", endpoint,
	)
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

	req, err = httputil.Sign(req, c.signingKeyID, c.signingKey)
	if err != nil {
		return nil, err
	}
	logger.Debugw("User Kyced HTTP request created", "host", req.URL.Host, "raw query", req.URL.RawQuery)
	return req, nil
}

// NewClient creates a new user kyc lookup client instance.
func NewClient(sugar *zap.SugaredLogger, url, signingKey, signingKeyID string) (*Client, error) {
	const timeout = time.Minute
	client := &http.Client{Timeout: timeout}
	return &Client{sugar: sugar,
		url:          url,
		client:       client,
		signingKey:   signingKey,
		signingKeyID: signingKeyID,
	}, nil
}

// IsKYCed will look for the user kyc status of input addr from server
func (c *Client) IsKYCed(addr ethereum.Address, timePoint time.Time) (bool, error) {
	const endpoint = "/kyced"
	var (
		params = make(map[string]string)
		result = userKycedReply{}
	)
	params["address"] = addr.Hex()
	req, err := c.newRequest(http.MethodGet, endpoint, params)
	if err != nil {
		return result.Kyced, err
	}

	rsp, err := c.client.Do(req)
	if err != nil {
		return result.Kyced, err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return result.Kyced, fmt.Errorf("unexpected return code: %d", rsp.StatusCode)
	}
	if err = json.NewDecoder(rsp.Body).Decode(&result); err != nil {
		return result.Kyced, err
	}
	if !result.Success {
		return result.Kyced, fmt.Errorf("failed to get  user profile, reason : %v", result.Reason)
	}
	return result.Kyced, nil
}
