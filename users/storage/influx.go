package storage

import (
	"errors"
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
	return storage, nil
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

	if err != nil {
		inf.sugar.Debugw("error from query", "error", err)
		return false, err
	}
	inf.sugar.Debugw("result from query", "result", res)

	var userTradeAmount float64
	var ok bool
	if len(res[0].Series) > 0 {
		if len(res[0].Series[0].Values) > 0 {
			userTradeAmount, ok = (res[0].Series[0].Values[0][1]).(float64)
			if !ok {
				inf.sugar.Debugw("values second should be float", "value", res[0].Series[0].Values[0][1])
				return false, errors.New("trade amount values is not a float")
			}
		} else {
			inf.sugar.Debugw("return values from influx", "values", res[0].Series[0].Values)
			return false, nil
		}
	} else {
		inf.sugar.Debugw("user address was not found from trade db", "user address", address)
		return false, nil
	}
	inf.sugar.Debugw("user rich is", "address", address, "user trade amount", userTradeAmount, "daily limit", dailyLimit)
	return userTradeAmount >= dailyLimit, nil
}
