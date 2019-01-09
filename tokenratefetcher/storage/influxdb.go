package storage

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/tokenratefetcher/common"
	schema "github.com/KyberNetwork/reserve-stats/tokenratefetcher/storage/schema/tokenrate"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

const (
	timePrecision = "s"
)

//InfluxStorage implement DB interface in influx
type InfluxStorage struct {
	sugar        *zap.SugaredLogger
	dbName       string
	influxClient client.Client
}

//NewInfluxStorage return the influx storage instance
func NewInfluxStorage(cl client.Client, dbName string, sugar *zap.SugaredLogger) (*InfluxStorage, error) {
	is := InfluxStorage{
		sugar:        sugar,
		dbName:       dbName,
		influxClient: cl,
	}
	if err := is.createDB(); err != nil {
		return nil, err
	}
	return &is, nil
}

// queryDB convenience function to query the database
func (is *InfluxStorage) queryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: is.dbName,
	}
	if response, err := is.influxClient.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func (is *InfluxStorage) createDB() error {
	_, err := is.queryDB(fmt.Sprintf("CREATE DATABASE %s", is.dbName))
	return err
}

//SaveRates save a list of rates into db
func (is *InfluxStorage) SaveRates(rates []common.TokenRate) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  is.dbName,
		Precision: timePrecision,
	})
	if err != nil {
		return err
	}
	for _, rate := range rates {
		pt, err := is.tokenRateToPoint(rate)
		if err != nil {
			return err
		}
		bp.AddPoint(pt)
	}
	return is.influxClient.Write(bp)
}

func (is *InfluxStorage) tokenRateToPoint(rate common.TokenRate) (*client.Point, error) {
	var (
		pt  *client.Point
		err error
	)

	tags := map[string]string{
		schema.Provider.String(): rate.Provider,
	}

	fields := map[string]interface{}{
		schema.Rate.String(): rate.Rate,
	}

	pt, err = client.NewPoint(fmt.Sprintf("%s_%s", common.GetTokenSymbolFromProviderNameTokenID(rate.Provider, rate.TokenID), rate.Currency), tags, fields, rate.Timestamp)
	if err != nil {
		return nil, err
	}
	return pt, nil
}
