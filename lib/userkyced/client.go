package userkyced

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// Client is the the implementation to query user kyced status info.
type Client struct {
	sugar  *zap.SugaredLogger
	client *http.Client
	url    string
}

func (c *Client) newRequest(method, endpoint string, params map[string]string) (*http.Request, error) {
	logger := c.sugar.With(
		"func", "lib/userkyced/client.go/newRequest",
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

	logger.Debugw("User Kyced HTTP request created", "host", req.URL.Host, "raw query", req.URL.RawQuery)
	return req, nil
}

// NewClient creates a new user kyc lookup client instance.
func NewClient(sugar *zap.SugaredLogger, url string) (*Client, error) {
	const timeout = time.Minute
	client := &http.Client{Timeout: timeout}
	return &Client{sugar: sugar,
		url:    url,
		client: client,
	}, nil
}

// IsKYCedAtTime will look for the user kyc status of input addr from server
func (c *Client) IsKYCedAtTime(addr ethereum.Address, timePoint time.Time) (bool, error) {
	const endpoint = "/kyced"
	var (
		params = make(map[string]string)
		result = userKycedReply{}
	)
	params["address"] = addr.Hex()
	params["time"] = strconv.FormatUint(timeutil.TimeToTimestampMs(timePoint), 10)
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
		return result.Kyced, fmt.Errorf("unexpected return code: %d, reason : %v", rsp.StatusCode, result.Reason)
	}
	if err = json.NewDecoder(rsp.Body).Decode(&result); err != nil {
		return result.Kyced, err
	}
	return result.Kyced, nil
}
