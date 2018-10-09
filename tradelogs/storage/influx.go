package storage

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"strconv"

	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// InfluxStorage represent a client to store trade data to influx DB
type InfluxStorage struct {
	dbName       string
	influxClient client.Client
	amountFmt    tokenAmountFormatter
	sugar        *zap.SugaredLogger
}

// NewInfluxStorage init an instance of InfluxStorage
func NewInfluxStorage(sugar *zap.SugaredLogger, dbName string, influxClient client.Client, amountFmt tokenAmountFormatter) (*InfluxStorage, error) {
	storage := &InfluxStorage{
		dbName:       dbName,
		influxClient: influxClient,
		amountFmt:    amountFmt,
		sugar:        sugar,
	}
	err := storage.createDB()
	if err != nil {
		return nil, err
	}
	return storage, nil
}

// SaveTradeLogs persist trade logs to DB
func (is *InfluxStorage) SaveTradeLogs(logs []common.TradeLog) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  is.dbName,
		Precision: "ms",
	})
	if err != nil {
		return err
	}

	for _, log := range logs {
		pt, err := is.tradeLogToPoint(log)
		if err != nil {
			return err
		}

		bp.AddPoint(pt)
	}

	if err := is.influxClient.Write(bp); err != nil {
		return err
	}

	is.sugar.Debugw("saved trade logs into influxdb", "trade logs", logs)

	return nil
}

// createDB creates the database will be used for storing trade logs measurements.
func (is *InfluxStorage) createDB() error {
	_, err := is.queryDB(is.influxClient, fmt.Sprintf("CREATE DATABASE %s", is.dbName))
	return err
}

// queryDB convenience function to query the database
func (is *InfluxStorage) queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: is.dbName,
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

func (is *InfluxStorage) tradeLogToPoint(log common.TradeLog) (*client.Point, error) {
	tags := map[string]string{
		"block_number": strconv.FormatUint(log.BlockNumber, 10),
		"tx_hash":      log.TransactionHash.String(),

		"eth_receival_sender": log.EtherReceivalSender.String(),

		"user_addr": log.UserAddress.String(),

		"src_addr": log.SrcAddress.String(),
		"dst_addr": log.DestAddress.String(),

		"reserve_addr": log.ReserveAddress.String(),
		"wallet_addr":  log.WalletAddress.String(),

		"country": log.Country,
		"ip":      log.IP,
	}

	ethReceivalAmount, err := is.amountFmt.FormatAmount(blockchain.ETHAddr, log.EtherReceivalAmount)
	if err != nil {
		return nil, err
	}

	srcAmount, err := is.amountFmt.FormatAmount(log.SrcAddress, log.SrcAmount)
	if err != nil {
		return nil, err
	}

	dstAmount, err := is.amountFmt.FormatAmount(log.DestAddress, log.DestAmount)
	if err != nil {
		return nil, err
	}

	walletFee, err := is.amountFmt.FormatAmount(blockchain.KNCAddr, log.WalletFee)
	if err != nil {
		return nil, err
	}

	burnFee, err := is.amountFmt.FormatAmount(blockchain.KNCAddr, log.BurnFee)
	if err != nil {
		return nil, err
	}

	fields := map[string]interface{}{
		"eth_receival_amount": ethReceivalAmount,

		"src_amount":  srcAmount,
		"dst_amount":  dstAmount,
		"fiat_amount": log.FiatAmount,

		"wallet_fee": walletFee,
		"burn_fee":   burnFee,
	}

	pt, err := client.NewPoint("trade", tags, fields, log.Timestamp)
	if err != nil {
		return nil, err
	}

	return pt, nil
}
