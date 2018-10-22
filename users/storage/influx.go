package storage

import (
	"fmt"

	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

const (
	ethDecimals int64  = 18
	ethAddress  string = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
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
func (inf *InfluxStorage) IsExceedDailyLimit(address string, dailyLimit float64) (bool, error) {
	query := fmt.Sprintf(`SELECT SUM(amount) as daily_fiat_amount FROM (SELECT eth_amount*eth_usd_rate as amount 
FROM trades WHERE user_addr='%s' AND time <= now() AND time >= (now()-24h))`,
		address)
	res, err := inf.queryDB(inf.influxClient, query)

	inf.sugar.Debugw("result from query", "result", res)

	if err != nil {
		return false, err
	}
	var userTradeAmount float64
	if len(res[0].Series) > 0 {
		userTradeAmount = (res[0].Series[0].Values[0][1]).(float64)
	}
	return userTradeAmount >= dailyLimit, nil
}
