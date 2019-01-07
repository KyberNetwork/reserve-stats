package tokenrate

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/tokenrate"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

const (
	timePrecision = "s"
)

//RateFetcher represent fetcher for anyToken-USD rate
type RateFetcher struct {
	sugar        *zap.SugaredLogger
	dbName       string
	influxClient client.Client
	rateProvider tokenrate.Provider
}

//NewRateFetcher return a RateFetcher for any Token_USD rate
func NewRateFetcher(sugar *zap.SugaredLogger, dbName string, client client.Client, rp tokenrate.Provider) (*RateFetcher, error) {
	fetcher := &RateFetcher{
		sugar:        sugar,
		dbName:       dbName,
		influxClient: client,
		rateProvider: rp,
	}

	if err := fetcher.createDB(); err != nil {
		return nil, err
	}
	return fetcher, nil
}

func (rf *RateFetcher) createDB() error {
	_, err := rf.queryDB(rf.influxClient, fmt.Sprintf("CREATE DATABASE %s", rf.dbName))
	return err
}

// queryDB convenience function to query the database
func (rf *RateFetcher) queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: rf.dbName,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

//FetchRatesInRanges fetch and saves times into influxDB
func (rf *RateFetcher) FetchRatesInRanges(from, to time.Time, tokenID, currency string) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  rf.dbName,
		Precision: timePrecision,
	})
	if err != nil {
		return err
	}
	rf.sugar.Debugw("fetching rates", "from", from.String(), "to", to.String())

	for d := from.Truncate(time.Hour * 24); !d.After(to.Truncate(time.Hour * 24)); d = d.AddDate(0, 0, 1) {
		rate, err := rf.FetchRates(tokenID, currency, d)
		if err != nil {
			rf.sugar.Errorw("Failed to get rate", "error", err)
			return err
		}
		rf.sugar.Infow("Rate return", "rate", rate, "time", d.String())
		pt, err := rf.tokenRateToPoint(rate)
		if err != nil {
			rf.sugar.Errorw("Failed to get influx point from rate", "error", err)
			return err
		}
		bp.AddPoint(pt)
	}
	return rf.influxClient.Write(bp)
}

//FetchRates return the rate for pair token-currency at timestamp timeStamp
func (rf *RateFetcher) FetchRates(token, currency string, timeStamp time.Time) (TokenRate, error) {
	var (
		USDRate float64
		rate    TokenRate
		err     error
	)
	if USDRate, err = rf.rateProvider.Rate(token, currency, timeStamp); err != nil {
		rf.sugar.Errorw(fmt.Sprintf("failed to get %s/%s rate", token, currency), "timestamp", timeStamp.String())
		return rate, err
	}
	if USDRate != 0 {
		rf.sugar.Debugw(
			fmt.Sprintf("got %s/%s rate", token, currency),
			"rate", USDRate,
			"timestamp", timeStamp.String())

		rate = TokenRate{
			Timestamp: timeStamp,
			Rate:      USDRate,
			Provider:  rf.rateProvider.Name(),
			TokenID:   token,
			Currency:  currency,
		}
	}

	return rate, nil
}

func (rf *RateFetcher) tokenRateToPoint(rate TokenRate) (*client.Point, error) {
	var (
		pt  *client.Point
		err error
	)

	tags := map[string]string{
		"provider": rate.Provider,
	}

	fields := map[string]interface{}{
		"rate": rate.Rate,
	}

	pt, err = client.NewPoint(fmt.Sprintf("%s_%s", getTokenSymbolFromTokenID(rate.TokenID), rate.Currency), tags, fields, rate.Timestamp)
	if err != nil {
		return nil, err
	}
	return pt, nil
}
