package ethrate

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	cmcEthereumPricingAPIEndpoint = "https://graphs2.coinmarketcap.com/currencies/ethereum/"
	cmcTopUSDPricingAPIEndpoint   = "https://api.coinmarketcap.com/v1/ticker/?convert=USD&limit=10"
	historyRateTimeRange          = time.Hour * 24 * 30
)

// CMCRateResponse represents response from CoinMarketCap for digital
// currency pricing in USD
type CMCRateResponse []struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Rank     string `json:"rank"`
	PriceUSD string `json:"price_usd"`
}

// CMCRate is a fetcher which fetch ETH pricing in USD from CoinMarketCap.
type CMCRate struct {
	mu                    *sync.RWMutex
	cachedRates           [][]float64 // contain history rates within historyRateTimeRange started at currentCacheTimepoint,
	currentCacheTimepoint time.Time   // represents the start timepoint of cached history rates
	realtimeTimepoint     time.Time   // the first timepoint when it get recent rate from CMC
	realtimeRate          float64     // recent eth-usd rate from CMC
	sugar                 *zap.SugaredLogger
}

// RateLogResponse represents ETH price in USD at a timestamp
type RateLogResponse struct {
	PriceUSD [][]float64 `json:"price_usd"`
}

// GetUSDRate get the ETH price in USD at given timepoint
func (cmcRate *CMCRate) GetUSDRate(timepoint time.Time) float64 {
	if timepoint.After(cmcRate.realtimeTimepoint) {
		return cmcRate.realtimeRate
	}
	return cmcRate.rateFromCache(timepoint)
}

// rateFromCache return rate from cache, given a timepoint.
// The cache contain rates in a time range specify by historyRateTimeRange,
// started at cmcRate.currentCacheTimepoint.
// If the timepoint not in cache, we call CMC to update it.
func (cmcRate *CMCRate) rateFromCache(timepoint time.Time) float64 {
	cmcRate.mu.Lock()
	defer cmcRate.mu.Unlock()

	roundedTimepoint := timepoint.Round(historyRateTimeRange).UTC()
	if roundedTimepoint.After(timepoint) {
		roundedTimepoint = roundedTimepoint.Add(-historyRateTimeRange)
	}

	if !roundedTimepoint.Equal(cmcRate.currentCacheTimepoint) {
		ethRates, err := cmcRate.fetchHistoryRates(roundedTimepoint)
		if err != nil {
			cmcRate.sugar.Error("Cannot get rate from CoinMarketCap")
			return cmcRate.realtimeRate
		}
		rate, err := findEthRate(ethRates, timepoint)
		if err != nil {
			cmcRate.sugar.Error(err)
			return cmcRate.realtimeRate
		}
		cmcRate.currentCacheTimepoint = roundedTimepoint
		cmcRate.cachedRates = ethRates
		return rate
	}

	rate, err := findEthRate(cmcRate.cachedRates, timepoint)
	if err != nil {
		return cmcRate.realtimeRate
	}
	return rate
}

// fetchHistoryRates get rate from CMC in historyRateTimeRange, which contain the
// timepoint
func (cmcRate *CMCRate) fetchHistoryRates(timepoint time.Time) ([][]float64, error) {
	fromTime := timepoint.UnixNano() / int64(time.Millisecond)
	toTime := fromTime + int64(historyRateTimeRange/time.Millisecond)
	url := fmt.Sprintf("%s/%d/%d/", cmcEthereumPricingAPIEndpoint, fromTime, toTime)

	resp, err := http.Get(url)
	if err != nil {
		cmcRate.sugar.Errorf("Getting history eth-usd rate failed: %v. Url: %s", err, url)
		return [][]float64{}, err
	}

	defer resp.Body.Close()
	rateResponse := RateLogResponse{}

	err = json.NewDecoder(resp.Body).Decode(&rateResponse)
	if err != nil {
		cmcRate.sugar.Errorf("Got bad response for history eth-usd rate: %v", err)
		return [][]float64{}, err
	}

	return rateResponse.PriceUSD, nil
}

func findEthRate(ethRateLog [][]float64, timepoint time.Time) (float64, error) {
	timepointInMillis := timepoint.UnixNano() / int64(time.Millisecond)

	var ethRate float64
	for _, e := range ethRateLog {
		if e[0] >= float64(timepointInMillis) {
			ethRate = e[1]
			return ethRate, nil
		}
	}
	return 0, errors.New("Cannot find ether rate corresponding with the timepoint")
}

// RunGetEthRate schedule to get rate each 10 minutes.
func (cmcRate *CMCRate) RunGetEthRate() {
	tick := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			err := cmcRate.FetchEthRate()
			if err != nil {
				cmcRate.sugar.Error(err)
			}
			<-tick.C
		}
	}()
}

// FetchEthRate get current rate ETH vs USD from CoinMarketCap.
func (cmcRate *CMCRate) FetchEthRate() (err error) {
	resp, err := http.Get(cmcTopUSDPricingAPIEndpoint)
	if err != nil {
		cmcRate.sugar.Errorf("Getting eth-usd rate failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	rateResponse := CMCRateResponse{}
	err = json.NewDecoder(resp.Body).Decode(&rateResponse)
	if err != nil {
		cmcRate.sugar.Errorf("Got bad response for eth-usd rate: %v", err)
		return err
	}

	for _, rate := range rateResponse {
		if rate.Symbol == "ETH" {
			newrate, err := strconv.ParseFloat(rate.PriceUSD, 64)
			if err != nil {
				cmcRate.sugar.Infof("Cannot get USD rate: %s", err.Error())
				return err
			}

			if cmcRate.realtimeRate == 0 {
				// set realtimeTimepoint to the timepoint that realtime rate is updated for the
				// first time
				cmcRate.realtimeTimepoint = time.Now()
			}
			cmcRate.realtimeRate = newrate
			return nil
		}
	}
	return nil
}

// Run start a schedule to fetch data from CoinMarketCap
func (cmcRate *CMCRate) Run() {
	// run real time fetcher
	cmcRate.RunGetEthRate()
}

// NewCMCRate fetch rate of ETH vs USD from CoinMarketCap
func NewCMCRate(sugar *zap.SugaredLogger) *CMCRate {
	result := &CMCRate{
		mu:    &sync.RWMutex{},
		sugar: sugar,
	}
	result.Run()
	return result
}
