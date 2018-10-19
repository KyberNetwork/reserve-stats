package storage

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	tradecommon "github.com/KyberNetwork/reserve-stats/tradelogs/common"
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

// calculateFiatAmount returns new TradeLog with fiat amount calculated.
// * For ETH-Token or Token-ETH conversions, the ETH amount is taken from ExecuteTrade event.
// * For Token-Token, the ETH amount is reading from EtherReceival event.
func calculateFiatAmount(tradeLog tradecommon.TradeLog, rate float64) tradecommon.TradeLog {
	ethAmount := new(big.Float)

	if strings.ToLower(ethAddress) == strings.ToLower(tradeLog.SrcAddress.String()) {
		// ETH-Token
		ethAmount.SetInt(tradeLog.SrcAmount)
	} else if strings.ToLower(ethAddress) == strings.ToLower(tradeLog.DestAddress.String()) {
		// Token-ETH
		ethAmount.SetInt(tradeLog.DestAmount)
	} else if tradeLog.EtherReceivalAmount != nil {
		// Token-Token
		ethAmount.SetInt(tradeLog.EtherReceivalAmount)
	}

	// fiat amount = ETH amount * rate
	ethAmount = ethAmount.Mul(ethAmount, new(big.Float).SetFloat64(rate))
	ethAmount.Quo(ethAmount, new(big.Float).SetFloat64(math.Pow10(int(ethDecimals))))
	tradeLog.FiatAmount, _ = ethAmount.Float64()

	return tradeLog
}

//IsExceedDailyLimit return if add address trade over daily limit or not
func (inf *InfluxStorage) IsExceedDailyLimit(address string, dailyLimit float64) (bool, error) {
	query := fmt.Sprintf(`SELECT SUM(eth_receival_amount*eth_usd_rate) as daily_fiat_amount
	FROM trades WHERE user_addr='%s' AND time <= now() AND time >= (now()-24h)`,
		address)
	// 	query := fmt.Sprintf(`SELECT SUM(eth_receival_amount) as daily_fiat_amount
	// FROM trades WHERE user_addr='%s'`,
	// 		address)
	res, err := inf.queryDB(inf.influxClient, query)

	inf.sugar.Debugw("result from query", "result", res)

	if err != nil {
		return false, err
	}
	userTradeAmount := (res[0].Series[0].Values[0][1]).(float64)
	return userTradeAmount >= dailyLimit, nil
}
