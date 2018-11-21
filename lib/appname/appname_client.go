package appname

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

// Client is the the real implementation of addr to integration app name  interface.
type Client struct {
	sugar  *zap.SugaredLogger
	client *http.Client
	url    string
}

func (c *Client) newRequest(method, endpoint string) (*http.Request, error) {
	var (
		logger = c.sugar.With("func", "lib/appname/appname_client.go/newRequest()", "method", method, "endpoint", endpoint)
	)
	logger.Debug("creating new Intergration app name HTTP request")

	url := fmt.Sprintf("%s/%s",
		strings.TrimRight(c.url, "/"),
		strings.Trim(endpoint, "/"),
	)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	return req, nil
}

//GetAddrToAppName return a map of ethereum.Address to app name or error if occur
func (c *Client) GetAddrToAppName() (map[ethereum.Address]string, error) {
	const endpoint = "/addr-to-appname"
	req, err := c.newRequest(http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}
	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	c.sugar.Debugw("Get Addr To App Name", "body", rsp.Body)

	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected return code: %d", rsp.StatusCode)
	}
	var (
		tmp    = make(map[string]string)
		result = make(map[ethereum.Address]string)
	)
	if err := json.NewDecoder(rsp.Body).Decode(&tmp); err != nil {
		return nil, err
	}
	for k, v := range tmp {
		result[ethereum.HexToAddress(k)] = v
	}
	return result, nil
}

// NewClient creates a new client to addr to app name interfaces
func NewClient(sugar *zap.SugaredLogger, url string) (*Client, error) {
	const timeout = time.Minute
	client := &http.Client{Timeout: timeout}
	return &Client{
		sugar:  sugar,
		url:    url,
		client: client,
	}, nil
}
