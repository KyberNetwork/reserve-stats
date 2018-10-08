package storage

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// InfluxStorage represent a client to store trade data to influx DB
type InfluxStorage struct {
	dbName       string
	influxClient client.Client
	tokenUtil    *blockchain.TokenUtil
}

// NewInfluxStorage init an instance of InfluxStorage
func NewInfluxStorage(dbName, endpoint, username, password string, tokenUtil *blockchain.TokenUtil) (*InfluxStorage, error) {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     endpoint,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return &InfluxStorage{dbName: dbName, influxClient: c, tokenUtil: tokenUtil}, nil
}

// SaveTradeLogs persist trade logs to DB
func (is *InfluxStorage) SaveTradeLogs(logs []common.TradeLog) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  is.dbName,
		Precision: "s",
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

	return nil
}

// CreateDB create database in influx db
func (is *InfluxStorage) CreateDB() error {
	_, err := is.queryDB(is.influxClient, fmt.Sprintf("CREATE DATABASE %s", is.dbName))
	if err != nil {
		return err
	}
	return nil
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
		"eth_sender": log.EtherReceivalSender.String(),

		"user_addr": log.UserAddress.String(),

		"src_addr": log.SrcAddress.String(),
		"dst_addr": log.DestAddress.String(),

		"reserve_addr": log.ReserveAddress.String(),
		"wallet_addr":  log.WalletAddress.String(),

		"country": log.Country,
	}

	ethAmount, err := is.tokenUtil.GetTokenAmount(blockchain.ETHAddr, log.EtherReceivalAmount)
	if err != nil {
		return nil, err
	}

	srcAmount, err := is.tokenUtil.GetTokenAmount(log.SrcAddress, log.SrcAmount)
	if err != nil {
		return nil, err
	}

	dstAmount, err := is.tokenUtil.GetTokenAmount(log.DestAddress, log.DestAmount)
	if err != nil {
		return nil, err
	}

	walletFee, err := is.tokenUtil.GetTokenAmount(blockchain.KNCAddr, log.WalletFee)
	if err != nil {
		return nil, err
	}

	burnFee, err := is.tokenUtil.GetTokenAmount(blockchain.KNCAddr, log.BurnFee)
	if err != nil {
		return nil, err
	}

	fields := map[string]interface{}{
		"block_number": int64(log.BlockNumber),
		"tx_hash":      log.TransactionHash.String(),
		"tx_index":     int64(log.Index),

		"eth_amount": ethAmount,

		"src_amount":  srcAmount,
		"dst_amount":  dstAmount,
		"fiat_amount": log.FiatAmount,

		"wallet_fee": walletFee,
		"burn_fee":   burnFee,

		"ip": log.IP,
	}

	pt, err := client.NewPoint("trade_logs", tags, fields, log.Timestamp)
	if err != nil {
		return nil, err
	}

	return pt, nil
}
