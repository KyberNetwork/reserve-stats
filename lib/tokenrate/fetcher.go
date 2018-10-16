package tokenrate

import (
	"fmt"
	"strconv"
	"time"

	"github.com/KyberNetwork/tokenrate"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

//ETHUSDRateFetcher represent fetcher for ETH-USD rate
type ETHUSDRateFetcher struct {
	sugar        *zap.SugaredLogger
	dbName       string
	influxClient client.Client
	rateProvider tokenrate.ETHUSDRateProvider
}

//NewETHUSDRateFetcher return new instance of ETHUSDFetcher
func NewETHUSDRateFetcher(sugar *zap.SugaredLogger, dbName string, client client.Client, rateProvider tokenrate.ETHUSDRateProvider) (*ETHUSDRateFetcher, error) {

	fetcher := &ETHUSDRateFetcher{
		sugar:        sugar,
		dbName:       dbName,
		influxClient: client,
		rateProvider: rateProvider,
	}
	err := fetcher.createDB()
	if err != nil {
		return nil, err
	}
	return fetcher, nil
}

//FetchRates fetch ETH-USD rate and save to db
func (ef *ETHUSDRateFetcher) FetchRates(blockNumber uint64, timestamp time.Time) (float64, error) {
	var (
		ethRate float64
		rate    ETHUSDRate
		err     error
	)
	if ethRate, err = ef.rateProvider.USDRate(timestamp); err != nil {
		ef.sugar.Errorw("failed to get ETH/USD rate", "timestamp", timestamp.String())
		return 0, err
	}
	if ethRate != 0 {
		ef.sugar.Debugw("got ETH/USD rate",
			"rate", ethRate,
			"timestamp", timestamp.String())

		// Save rates to db
		rate = ETHUSDRate{
			Timestamp:   timestamp,
			Rate:        ethRate,
			Provider:    ef.rateProvider.Name(),
			BlockNumber: blockNumber,
		}

		if err := ef.SaveTokenRate(rate); err != nil {
			return ethRate, err
		}
	}
	return ethRate, nil
}

// createDB creates the database will be used for storing trade logs measurements.
func (ef *ETHUSDRateFetcher) createDB() error {
	_, err := ef.queryDB(ef.influxClient, fmt.Sprintf("CREATE DATABASE %s", ef.dbName))
	return err
}

// queryDB convenience function to query the database
func (ef *ETHUSDRateFetcher) queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: ef.dbName,
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

//SaveTokenRate into influx
func (ef ETHUSDRateFetcher) SaveTokenRate(rate ETHUSDRate) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  ef.dbName,
		Precision: "ms",
	})
	if err != nil {
		return err
	}

	pt, err := ef.tokenRateToPoint(rate)
	if err != nil {
		return err
	}

	bp.AddPoint(pt)

	if err := ef.influxClient.Write(bp); err != nil {
		return err
	}

	return nil
}

func (ef *ETHUSDRateFetcher) tokenRateToPoint(rate ETHUSDRate) (*client.Point, error) {
	var (
		pt  *client.Point
		err error
	)
	tags := map[string]string{
		"block_number": strconv.FormatUint(rate.BlockNumber, 10),
	}
	fields := map[string]interface{}{
		"provider": rate.Provider,
		"rate":     rate.Rate,
	}
	pt, err = client.NewPoint("token_rate", tags, fields, rate.Timestamp)
	if err != nil {
		return nil, err
	}
	return pt, err
}
