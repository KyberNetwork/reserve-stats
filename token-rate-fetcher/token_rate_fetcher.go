package tokenratefetcher

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/token-rate-fetcher/common"
	"github.com/KyberNetwork/reserve-stats/token-rate-fetcher/storage"
	"github.com/KyberNetwork/tokenrate"
	"go.uber.org/zap"
)

//RateFetcher represent fetcher for anyToken-USD rate
type RateFetcher struct {
	sugar        *zap.SugaredLogger
	storage      storage.Interface
	rateProvider tokenrate.Provider
}

//NewRateFetcher return a RateFetcher for any Token_USD rate
func NewRateFetcher(sugar *zap.SugaredLogger, str storage.Interface, rp tokenrate.Provider) (*RateFetcher, error) {
	fetcher := &RateFetcher{
		sugar:        sugar,
		storage:      str,
		rateProvider: rp,
	}

	return fetcher, nil
}

//FetchRatesInRanges fetch and saves times into influxDB
func (rf *RateFetcher) FetchRatesInRanges(from, to time.Time, tokenID, currency string) error {
	var (
		logger = rf.sugar.With(
			"func", "tokenratefetcher/RateFetcher.FetchRatesInRanges",
			"from", from,
			"to", to,
			"token_id", tokenID,
			"currency", currency,
		)
	)

	// normalizes time range input
	from = timeutil.Midnight(from)
	to = timeutil.Midnight(to)

	logger.Debugw("normalized time range input",
		"normalized_from", from,
		"normalized_to", to,
	)

	logger.Debug("fetching rates")
	for d := from.UTC(); !d.After(to); d = d.AddDate(0, 0, 1) {
		logger.Infow("fetching rates at specific timestamp",
			"timestamp", d,
		)
		rate, err := rf.fetchRates(tokenID, currency, d)
		if err != nil {
			rf.sugar.Errorw("failed to get rate", "error", err)
			return err
		}
		rf.sugar.Infow("got token rate", "rate", rate, "time", d.String())
		if err = rf.storage.SaveRates([]common.TokenRate{rate}); err != nil {
			return err
		}
		logger.Debugw("successfully stored rate to database",
			"timestamp", d,
		)
	}
	return nil
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
