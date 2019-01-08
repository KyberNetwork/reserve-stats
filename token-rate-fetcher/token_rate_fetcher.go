package tokenratefetcher

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/token-rate-fetcher/common"
	"github.com/KyberNetwork/reserve-stats/token-rate-fetcher/storage"
	"github.com/KyberNetwork/tokenrate"
	"go.uber.org/zap"
)

//RateFetcher represent fetcher for anyToken-USD rate
type RateFetcher struct {
	sugar         *zap.SugaredLogger
	influxStorage *storage.InfluxStorage
	rateProvider  tokenrate.Provider
}

//NewRateFetcher return a RateFetcher for any Token_USD rate
func NewRateFetcher(sugar *zap.SugaredLogger, is *storage.InfluxStorage, rp tokenrate.Provider) (*RateFetcher, error) {
	fetcher := &RateFetcher{
		sugar:         sugar,
		influxStorage: is,
		rateProvider:  rp,
	}

	return fetcher, nil
}

//FetchRatesInRanges fetch and saves times into influxDB
func (rf *RateFetcher) FetchRatesInRanges(from, to time.Time, tokenID, currency string) error {
	var rates []common.TokenRate

	rf.sugar.Debugw("fetching rates", "from", from.String(), "to", to.String())

	for d := from.Truncate(time.Hour * 24); !d.After(to.Truncate(time.Hour * 24)); d = d.AddDate(0, 0, 1) {
		rate, err := rf.fetchRates(tokenID, currency, d)
		if err != nil {
			rf.sugar.Errorw("Failed to get rate", "error", err)
			return err
		}
		rf.sugar.Infow("Rate return", "rate", rate, "time", d.String())
		rates = append(rates, rate)
	}
	return rf.influxStorage.SaveRates(rates)
}

//FetchRates return the rate for pair token-currency at timestamp timeStamp
func (rf *RateFetcher) fetchRates(token, currency string, timeStamp time.Time) (common.TokenRate, error) {
	var (
		USDRate float64
		rate    common.TokenRate
		err     error
	)
	if USDRate, err = rf.rateProvider.Rate(token, currency, timeStamp); err != nil {
		rf.sugar.Errorw(fmt.Sprintf("failed to get %s/%s rate", token, currency), "timestamp", timeStamp.String())
		return rate, err
	}
	rf.sugar.Debugw(
		fmt.Sprintf("got %s/%s rate", token, currency),
		"rate", USDRate,
		"timestamp", timeStamp.String())

	rate = common.TokenRate{
		Timestamp: timeStamp,
		Rate:      USDRate,
		Provider:  rf.rateProvider.Name(),
		TokenID:   token,
		Currency:  currency,
	}

	return rate, nil
}
