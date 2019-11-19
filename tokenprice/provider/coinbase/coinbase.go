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
	client         *http.Client
	baseURL        string
	reqWaitingTime time.Duration
}

// New creates a new CoinBase instance.
func New(reqWaitingTime time.Duration) *CoinBase {
	const (
		defaultTimeout = time.Second * 10
		baseURL        = "https://api.coinbase.com/v2"
	)
	client := &http.Client{
		Timeout: defaultTimeout,
	}
	return &CoinBase{
		client:         client,
		baseURL:        baseURL,
		reqWaitingTime: reqWaitingTime,
	}
}

// Option option when init coinbase instance
type Option func(cb *CoinBase)

// WithReqWaitingTime set waiting time to avoid rate limit
func WithReqWaitingTime(reqWaitingTime time.Duration) Option {
	return func(cb *CoinBase) {
		cb.reqWaitingTime = reqWaitingTime
	}
}

type responseData struct {
	Data struct {
		Amount string `json:"amount"`
	}
}

// Price returns the price of given token in real world currency at given timestamp.
func (cb *CoinBase) Price(token, currency string, timestamp time.Time) (float64, error) {
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

	price, err := strconv.ParseFloat(resp.Data.Amount, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

// ETHPrice returns the historical price of ETH in USD.
func (cb *CoinBase) ETHPrice(timestamp time.Time) (float64, error) {
	const (
		ethereumID = "ETH"
		usdID      = "USD"
	)
	return cb.Price(ethereumID, usdID, timestamp)
}

// Wait sleep in reqWaitingTime to avoid rate limit
func (cb *CoinBase) Wait() {
	time.Sleep(cb.reqWaitingTime)
}

//Name return name of CoinBase provider name
func (cb *CoinBase) Name() string {
	return providerName
}