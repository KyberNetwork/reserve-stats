package userprofile

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
)

// Client is the the implementation to query userprofile info.
type Client struct {
	sugar      *zap.SugaredLogger
	client     *http.Client
	url        string
	signingKey string
}

func (c *Client) newRequest(method, endpoint string, params map[string]string) (*http.Request, error) {
	logger := c.sugar.With(
		"func", "lib/userprofile/client.go/newRequest",
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

	_, ok := params["nonce"]
	if ok {
		signed, err := c.sign(q.Encode())
		if err != nil {
			return nil, err
		}
		req.Header.Add("signed", signed)
		logger = logger.With("signed", signed)
	}

	logger.Debugw("User profile API HTTP request created", "host", req.URL.Host, "raw query", req.URL.RawQuery)
	return req, nil
}

// NewClient creates a new core client instance.
func NewClient(sugar *zap.SugaredLogger, url, signingKey string) (*Client, error) {
	const timeout = time.Minute
	client := &http.Client{Timeout: timeout}
	return &Client{sugar: sugar, url: url, client: client}, nil
}

// LookUpUserProfile will look for the UserProfile of input addr from server
func (c *Client) LookUpUserProfile(addr ethereum.Address) (UserProfile, error) {
	const endpoint = "/api/wallet_info"
	var (
		params = make(map[string]string)
		result = userClientReply{
			Data: UserProfile{},
		}
	)
	params["nonce"] = generateNonce()
	params["wallet_address"] = addr.Hex()
	req, err := c.newRequest(http.MethodGet, endpoint, params)
	if err != nil {
		return result.Data, err
	}

	rsp, err := c.client.Do(req)
	if err != nil {
		return result.Data, err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return result.Data, fmt.Errorf("unexpected return code: %d", rsp.StatusCode)
	}
	if err = json.NewDecoder(rsp.Body).Decode(&result); err != nil {
		return result.Data, err
	}
	if !result.Success {
		return result.Data, fmt.Errorf("failed to get  user profile, reason : %v", result.Reason)
	}
	return result.Data, nil
}

func (c *Client) sign(msg string) (string, error) {
	mac := hmac.New(sha512.New, []byte(c.signingKey))
	if _, err := mac.Write([]byte(msg)); err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// generateNonce returns nonce header required to use Core API,
// which is current timestamp in milliseconds.
func generateNonce() string {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(now, 10)
}
