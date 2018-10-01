package tradelogs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/common"
)

const (
	cmcEthereumPricingAPIEndpoint = "https://graphs2.coinmarketcap.com/currencies/ethereum/"
	cmcTopUSDPricingAPIEndpoint   = "https://api.coinmarketcap.com/v1/ticker/?convert=USD&limit=10"
)

// CoinCapRateResponse represents response from CoinMarketCap for digital
// currency pricing in USD
type CoinCapRateResponse []struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Rank     string `json:"rank"`
	PriceUSD string `json:"price_usd"`
}

// CMCEthUSDRate is a fetcher which fetch ETH pricing in USD from CoinMarketCap.
type CMCEthUSDRate struct {
	mu                *sync.RWMutex
	cachedRates       [][]float64
	currentCacheMonth uint64
	realtimeTimepoint uint64
	realtimeRate      float64
	sugar             *zap.SugaredLogger
}

// RateLogResponse represents ETH price in USD at a timestamp
type RateLogResponse struct {
	PriceUSD [][]float64 `json:"price_usd"`
}

// GetTimeStamp get timestamp in miliseconds
func GetTimeStamp(year int, month time.Month, day int, hour int, minute int, sec int, nanosec int, loc *time.Location) uint64 {
	return uint64(time.Date(year, month, day, hour, minute, sec, nanosec, loc).Unix() * 1000)
}

// GetMonthTimeStamp get timestamp at the beginning of given month
func GetMonthTimeStamp(timepoint uint64) uint64 {
	t := time.Unix(int64(timepoint/1000), 0).UTC()
	month, year := t.Month(), t.Year()
	return GetTimeStamp(year, month, 1, 0, 0, 0, 0, time.UTC)
}

// GetNextMonth get the next month, year of given month, year
func GetNextMonth(month, year int) (int, int) {
	var toMonth, toYear int
	if int(month) == 12 {
		toMonth = 1
		toYear = year + 1
	} else {
		toMonth = int(month) + 1
		toYear = year
	}
	return toMonth, toYear
}

// GetUSDRate get the ETH price in USD at given timepoint
func (cmcRate *CMCEthUSDRate) GetUSDRate(timepoint uint64) float64 {
	if timepoint >= cmcRate.realtimeTimepoint {
		return cmcRate.realtimeRate
	}
	return cmcRate.rateFromCache(timepoint)
}

func (cmcRate *CMCEthUSDRate) rateFromCache(timepoint uint64) float64 {
	cmcRate.mu.Lock()
	defer cmcRate.mu.Unlock()

	monthTimeStamp := GetMonthTimeStamp(timepoint)
	if monthTimeStamp != cmcRate.currentCacheMonth {
		ethRates, err := cmcRate.fetchRate(timepoint)
		if err != nil {
			cmcRate.sugar.Info("Cannot get rate from CoinMarketCap")
			return cmcRate.realtimeRate
		}
		rate, err := findEthRate(ethRates, timepoint)
		if err != nil {
			cmcRate.sugar.Info(err)
			return cmcRate.realtimeRate
		}
		cmcRate.currentCacheMonth = monthTimeStamp
		cmcRate.cachedRates = ethRates
		return rate
	}

	rate, err := findEthRate(cmcRate.cachedRates, timepoint)
	if err != nil {
		return cmcRate.realtimeRate
	}
	return rate
}

func (cmcRate *CMCEthUSDRate) fetchRate(timepoint uint64) ([][]float64, error) {
	t := time.Unix(int64(timepoint/1000), 0).UTC()
	month, year := t.Month(), t.Year()
	fromTime := GetTimeStamp(year, month, 1, 0, 0, 0, 0, time.UTC)
	toMonth, toYear := GetNextMonth(int(month), year)
	toTime := GetTimeStamp(toYear, time.Month(toMonth), 1, 0, 0, 0, 0, time.UTC)
	api := cmcEthereumPricingAPIEndpoint + strconv.FormatInt(int64(fromTime), 10) + "/" + strconv.FormatInt(int64(toTime), 10) + "/"
	resp, err := http.Get(api)
	if err != nil {
		return [][]float64{}, err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			cmcRate.sugar.Infof("Response body close error: %s", cErr.Error())
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return [][]float64{}, err
	}
	rateResponse := RateLogResponse{}
	err = json.Unmarshal(body, &rateResponse)
	if err != nil {
		return [][]float64{}, err
	}
	ethRates := rateResponse.PriceUSD
	return ethRates, nil
}

func findEthRate(ethRateLog [][]float64, timepoint uint64) (float64, error) {
	var ethRate float64
	for _, e := range ethRateLog {
		if uint64(e[0]) >= timepoint {
			ethRate = e[1]
			return ethRate, nil
		}
	}
	return 0, errors.New("Cannot find ether rate corresponding with the timepoint")
}

// RunGetEthRate schedule to get rate each 10 minutes.
func (cmcRate *CMCEthUSDRate) RunGetEthRate() {
	tick := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			err := cmcRate.FetchEthRate()
			if err != nil {
				cmcRate.sugar.Info(err)
			}
			<-tick.C
		}
	}()
}

// FetchEthRate get current rate ETH vs USD from CoinMarketCap.
func (cmcRate *CMCEthUSDRate) FetchEthRate() (err error) {
	resp, err := http.Get(cmcTopUSDPricingAPIEndpoint)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			cmcRate.sugar.Infof("Response body close error: %s", cErr.Error())
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	rateResponse := CoinCapRateResponse{}
	err = json.Unmarshal(body, &rateResponse)
	if err != nil {
		cmcRate.sugar.Infof("Getting eth-usd rate failed: %+v", err)
	} else {
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
					cmcRate.realtimeTimepoint = common.GetTimepoint()
				}
				cmcRate.realtimeRate = newrate
				return nil
			}
		}
	}
	return nil
}

// Run start a schedule to fetch data from CoinMarketCap
func (cmcRate *CMCEthUSDRate) Run() {
	// run real time fetcher
	cmcRate.RunGetEthRate()
}

// NewCMCEthUSDRate fetch rate of ETH vs USD from CoinMarketCap
func NewCMCEthUSDRate() *CMCEthUSDRate {
	result := &CMCEthUSDRate{
		mu: &sync.RWMutex{},
	}
	result.Run()
	return result
}
