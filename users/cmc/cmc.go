package cmc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/utils"
	"go.uber.org/zap"
)

const (
	cmcEthereumPricingAPIEndpoint = "https://graphs2.coinmarketcap.com/currencies/ethereum/"
	cmcTopUSDPricingAPIEndpoint   = "https://api.coinmarketcap.com/v1/ticker/?convert=USD&limit=10"
)

//CoinCapRateResponse response from CMC
type CoinCapRateResponse []struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Rank     string `json:"rank"`
	PriceUSD string `json:"price_usd"`
}

//EthUSDRate represent the rate from cmc
type EthUSDRate struct {
	sugar             *zap.SugaredLogger
	mu                *sync.RWMutex
	cachedRates       [][]float64
	currentCacheMonth uint64
	realtimeTimepoint uint64
	realtimeRate      float64
}

//RateLogResponse rate response from CMC
type RateLogResponse struct {
	PriceUSD [][]float64 `json:"price_usd"`
}

//GetTimestamp return timestamp
func GetTimestamp(year int, month time.Month, day int, hour int, minute int, sec int, nanosec int, loc *time.Location) uint64 {
	return uint64(time.Date(year, month, day, hour, minute, sec, nanosec, loc).Unix() * 1000)
}

//GetMonthTimestamp return month timestamp
func GetMonthTimestamp(timepoint uint64) uint64 {
	t := time.Unix(int64(timepoint/1000), 0).UTC()
	month, year := t.Month(), t.Year()
	return GetTimestamp(year, month, 1, 0, 0, 0, 0, time.UTC)
}

//GetNextMonth return next month timestamp
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

//GetUSDRate return ETH-USD rate for provided timepoint
func (cmc *EthUSDRate) GetUSDRate(timepoint uint64) float64 {
	if timepoint >= cmc.realtimeTimepoint {
		return cmc.realtimeRate
	}
	return cmc.rateFromCache(timepoint)
}

func (cmc *EthUSDRate) rateFromCache(timepoint uint64) float64 {
	cmc.mu.Lock()
	defer cmc.mu.Unlock()
	monthTimeStamp := GetMonthTimestamp(timepoint)
	if monthTimeStamp != cmc.currentCacheMonth {
		ethRates, err := fetchRate(timepoint, cmc.sugar)
		if err != nil {
			cmc.sugar.Error("Cannot get rate from coinmarketcap")
			return cmc.realtimeRate
		}
		rate, err := findEthRate(ethRates, timepoint)
		if err != nil {
			log.Println(err)
			return cmc.realtimeRate
		}
		cmc.currentCacheMonth = monthTimeStamp
		cmc.cachedRates = ethRates
		return rate
	}
	rate, err := findEthRate(cmc.cachedRates, timepoint)
	if err != nil {
		return cmc.realtimeRate
	}
	return rate
}

func fetchRate(timepoint uint64, sugar *zap.SugaredLogger) ([][]float64, error) {
	t := time.Unix(int64(timepoint/1000), 0).UTC()
	month, year := t.Month(), t.Year()
	fromTime := GetTimestamp(year, month, 1, 0, 0, 0, 0, time.UTC)
	toMonth, toYear := GetNextMonth(int(month), year)
	toTime := GetTimestamp(toYear, time.Month(toMonth), 1, 0, 0, 0, 0, time.UTC)
	api := cmcEthereumPricingAPIEndpoint + strconv.FormatInt(int64(fromTime), 10) + "/" + strconv.FormatInt(int64(toTime), 10) + "/"
	resp, err := http.Get(api)
	if err != nil {
		return [][]float64{}, err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			sugar.Errorf("Response body close error: %s", cErr.Error())
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

//FetchEthRate fetcher for eth rate
func (cmc *EthUSDRate) FetchEthRate() (err error) {
	resp, err := http.Get(cmcTopUSDPricingAPIEndpoint)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			cmc.sugar.Errorf("Response body close error: %s", cErr.Error())
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	rateResponse := CoinCapRateResponse{}
	err = json.Unmarshal(body, &rateResponse)
	if err != nil {
		log.Printf("Getting eth-usd rate failed: %+v", err)
	} else {
		for _, rate := range rateResponse {
			if rate.Symbol == "ETH" {
				newrate, err := strconv.ParseFloat(rate.PriceUSD, 64)
				if err != nil {
					cmc.sugar.Errorf("Cannot get usd rate: %s", err.Error())
					return err
				}
				if cmc.realtimeRate == 0 {
					// set realtimeTimepoint to the timepoint that realtime rate is updated for the
					// first time
					cmc.realtimeTimepoint = utils.GetTimepoint()
				}
				cmc.realtimeRate = newrate
				return nil
			}
		}
	}
	return nil
}

//Run the fetcher
func (cmc *EthUSDRate) Run() {
	tick := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			err := cmc.FetchEthRate()
			if err != nil {
				cmc.sugar.Error(err)
			}
			<-tick.C
		}
	}()
}

//NewCMCEthUSDRate return new instance for cmc fetcher
func NewCMCEthUSDRate(sugar *zap.SugaredLogger) *EthUSDRate {
	result := &EthUSDRate{
		sugar: sugar,
		mu:    &sync.RWMutex{},
	}
	result.Run()
	return result
}
