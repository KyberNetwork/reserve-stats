package marketdata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"go.uber.org/zap"
)

// Client to call to market data
type Client struct {
	client  *http.Client
	baseURL string
	sugar   *zap.SugaredLogger
}

// NewMarketDataClient return new client for market data
func NewMarketDataClient(baseURL string, sugar *zap.SugaredLogger) *Client {
	c := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &Client{
		client:  c,
		baseURL: baseURL,
		sugar:   sugar,
	}
}

// PairSupportedResponse ...
type PairSupportedResponse struct {
	Valid bool `json:"valid"`
}

// PairSupported return if pair is supported by binance
func (c *Client) PairSupported(source, symbol string) (bool, error) {
	var (
		result PairSupportedResponse
		logger = c.sugar.With("func", caller.GetCurrentFunctionName(),
			"source", source,
			"symbol", symbol,
		)
	)
	endpoint := fmt.Sprintf("%s/is-valid-symbol?source=%s&symbol=%s", c.baseURL, source, symbol)
	logger.Infow("pair support endpoint", "enpoint", endpoint)
	req, err := http.NewRequest(
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return result.Valid, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return result.Valid, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result.Valid, nil
		}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return result.Valid, err
		}
	default:
		return result.Valid, fmt.Errorf("market data return with error code: %d", resp.StatusCode)
	}
	return result.Valid, nil
}
