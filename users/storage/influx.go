package storage

import (
	"fmt"

	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

// InfluxStorage represent a client to store trade data to influx DB
type InfluxStorage struct {
	dbName       string
	influxClient client.Client
	sugar        *zap.SugaredLogger
}

// NewInfluxStorage init an instance of InfluxStorage
func NewInfluxStorage(sugar *zap.SugaredLogger, dbName string, influxClient client.Client) (*InfluxStorage, error) {
	storage := &InfluxStorage{
		dbName:       dbName,
		influxClient: influxClient,
		sugar:        sugar,
	}
	err := storage.createDB()
	if err != nil {
		return nil, err
	}
	return storage, nil
}

// createDB creates the database will be used for storing trade logs measurements.
func (inf *InfluxStorage) createDB() error {
	_, err := inf.queryDB(inf.influxClient, fmt.Sprintf("CREATE DATABASE %s", inf.dbName))
	return err
}

// queryDB convenience function to query the database
func (inf *InfluxStorage) queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: inf.dbName,
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

//IsExceedDailyLimit return if add address trade over daily limit or not
func (inf *InfluxStorage) IsExceedDailyLimit(address string) (bool, error) {
	// 	query := fmt.Sprintf(`SELECT SUM(eth_receival_amount) as daily_fiat_amount
	// FROM trades WHERE user_addr='%s' AND time <= now() AND time >= (now()-24h)`,
	// 		address)
	query := fmt.Sprintf(`SELECT SUM(eth_receival_amount) as daily_fiat_amount 
FROM trades WHERE user_addr='%s'`,
		address)
	res, err := inf.queryDB(inf.influxClient, query)
	inf.sugar.Debugw("result from query", "result", res)
	if err != nil {
		return false, err
	}
	return false, err
}
