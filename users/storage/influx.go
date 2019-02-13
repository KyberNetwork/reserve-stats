package storage

import (
	"errors"
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"

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
	var (
		query = fmt.Sprintf(`SELECT SUM(amount) as daily_fiat_amount FROM (SELECT eth_amount*eth_usd_rate as amount 
FROM trades WHERE user_addr='%s' AND time <= now() AND time >= (now()-24h))`,
			address)
		userTradeAmount float64
		err             error
	)

	res, err := inf.queryDB(inf.influxClient, query)
	if err != nil {
		inf.sugar.Debugw("error from query", "error", err)
		return false, err
	}
	inf.sugar.Debugw("result from query", "result", res)

	if len(res) == 0 || len(res[0].Series) == 0 || len(res[0].Series[0].Values) == 0 || len(res[0].Series[0].Values[0]) < 2 {
		inf.sugar.Debugw("user address was not found from trade db", "user address", address)
		return false, nil
	}

	userTradeAmount, err = influxdb.GetFloat64FromInterface(res[0].Series[0].Values[0][1])
	if err != nil {
		inf.sugar.Debugw("values second should be float", "value", res[0].Series[0].Values[0][1], "error", err.Error())
		return false, errors.New("trade amount values is not a float")
	}

	inf.sugar.Debugw("got last 24h total transaction",
		"address", address,
		"user trade amount", userTradeAmount,
		"daily limit", dailyLimit)
	return userTradeAmount >= dailyLimit, nil
}
