package coinbase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const providerName = "coinbase"

const (
	timeLayout         = "2006-01-02"
	currentEndpoint    = "%s/prices/%s-%s/spot"
	historicalEndpoint = "%s/prices/%s-%s/spot?date=%s"
)

// CoinBase is the CoinBase implementation of Provider. The
// precision of CoinBase provider is up to day.
type CoinBase struct {
	client  *http.Client
	baseURL string
}

// New creates a new CoinBase instance.
func New() *CoinBase {
	const (
		defaultTimeout = time.Second * 10
		baseURL        = "https://api.coinbase.com/v2"
	)
	client := &http.Client{
		Timeout: defaultTimeout,
	}
	return &CoinBase{
		client:  client,
		baseURL: baseURL,
	}
}

type responseData struct {
	Data struct {
		Amount string `json:"amount"`
	}
}

// Rate returns the rate of given token in real world currency at given timestamp.
func (cb *CoinBase) Rate(token, currency string, timestamp time.Time) (float64, error) {
	var url string
	currentDate := time.Now().UTC().Format(timeLayout)
	queryDate := timestamp.UTC().Format(timeLayout)

	if currentDate == queryDate {
		url = fmt.Sprintf(currentEndpoint, cb.baseURL, token, currency)
	} else {
		url = fmt.Sprintf(historicalEndpoint, cb.baseURL, token, currency, queryDate)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Accept", "application/json")
	rsp, err := cb.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %s", rsp.Status)
	}

	var resp = &responseData{}
	if err = json.NewDecoder(rsp.Body).Decode(resp); err != nil {
		return 0, err
	}

	rate, err := strconv.ParseFloat(resp.Data.Amount, 64)
	if err != nil {
		return 0, err
	}
	return rate, nil
}

// USDRate returns the historical price of ETH.
func (cb *CoinBase) USDRate(timestamp time.Time) (float64, error) {
	const (
		ethereumID = "ETH"
		usdID      = "USD"
	)
	return cb.Rate(ethereumID, usdID, timestamp)
}

//Name return name of CoinBase provider name
func (cb *CoinBase) Name() string {
	return providerName
}
