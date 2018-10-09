package coingecko

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CoinGecko is the CoinGecko implementation of Provider. The
// precision of CoinGecko provider is up to day.
type CoinGecko struct {
	client  *http.Client
	baseURL string
}

// New creates a new CoinGecko instance.
func New() *CoinGecko {
	const (
		defaultTimeout = time.Second * 10
		baseURL        = "https://api.coingecko.com/api/v3"
	)
	client := &http.Client{
		Timeout: defaultTimeout,
	}
	return &CoinGecko{
		client:  client,
		baseURL: baseURL,
	}
}

type historyResponse struct {
	MarketData marketData `json:"market_data"`
}

type marketData struct {
	CurrentPrice map[string]float64 `json:"current_price"`
}

// Rate returns the rate of given token in real world currency at given timestamp.
func (cg *CoinGecko) Rate(token, currency string, timestamp time.Time) (float64, error) {
	const timeLayout = "02-01-2006"

	url := fmt.Sprintf("%s/coins/%s/history", cg.baseURL, token)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("date", timestamp.UTC().Format(timeLayout))
	req.URL.RawQuery = q.Encode()
	rsp, err := cg.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %s", rsp.Status)
	}

	var history = &historyResponse{}
	if err = json.NewDecoder(rsp.Body).Decode(history); err != nil {
		return 0, err
	}
	rate, ok := history.MarketData.CurrentPrice[currency]
	if !ok {
		return 0, fmt.Errorf("currency %q not found in market data", currency)
	}
	return rate, nil
}

// USDRate returns the historical price of ETH.
func (cg *CoinGecko) USDRate(timestamp time.Time) (float64, error) {
	const (
		ethereumID = "ethereum"
		usdID      = "usd"
	)
	return cg.Rate(ethereumID, usdID, timestamp)
}
