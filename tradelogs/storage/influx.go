package storage

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/influxdata/influxdb/client/v2"
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
		points, err := is.tradeLogToPoint(log)
		if err != nil {
			return err
		}

		for _, pt := range points {
			bp.AddPoint(pt)
		}
	}

	if err := is.influxClient.Write(bp); err != nil {
		return err
	}

	is.sugar.Debugw("saved trade logs into influxdb", "trade logs", logs)

	return nil
}

// LoadTradeLogs return trade logs from DB
func (is *InfluxStorage) LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error) {
	q := fmt.Sprintf(
		"SELECT %s FROM %s WHERE time >= %d AND time <= %d",
		`
		time, block_number, tx_hash, 
		eth_receival_sender, eth_receival_amount, 
		user_addr, src_addr, dst_addr, src_amount, dst_amount, fiat_amount, 
		reserve_addr, wallet_addr, wallet_fee, burn_fee,
		ip, country
		`,
		"trade",
		from.UnixNano(),
		to.UnixNano(),
	)

	is.sugar.Debug(from.UnixNano())
	is.sugar.Debug(to.UnixNano())

	res, err := is.queryDB(is.influxClient, q)
	if err != nil {
		return nil, err
	}

	var result []common.TradeLog

	if len(res[0].Series) == 0 {
		return nil, nil
	}

	for _, row := range res[0].Series[0].Values {
		is.sugar.Debug(row)
		tradeLog, err := is.rowToTradeLog(row)
		if err != nil {
			return nil, err
		}
		result = append(result, tradeLog)
	}

	return result, nil
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

func (is *InfluxStorage) tradeLogToPoint(log common.TradeLog) ([]*client.Point, error) {
	var points []*client.Point

	tags := map[string]string{
		"block_number": strconv.FormatUint(log.BlockNumber, 10),
		"tx_hash":      log.TransactionHash.String(),

		"eth_receival_sender": log.EtherReceivalSender.String(),

		"user_addr": log.UserAddress.String(),

		"src_addr": log.SrcAddress.String(),
		"dst_addr": log.DestAddress.String(),

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

	fields := map[string]interface{}{
		"eth_receival_amount": ethReceivalAmount,

		"src_amount":  srcAmount,
		"dst_amount":  dstAmount,
		"fiat_amount": log.FiatAmount,
	}

	tradePoint, err := client.NewPoint("trades", tags, fields, log.Timestamp)
	if err != nil {
		return nil, err
	}

	points = append(points, tradePoint)

	// build burnFeePoint
	for idx, burn := range log.BurnFees {
		tags := map[string]string{
			"tx_hash":      log.TransactionHash.String(),
			"reserve_addr": burn.ReserveAddress.String(),
			"order":        strconv.Itoa(idx),
		}

		burnAmount, err := is.amountFmt.FormatAmount(blockchain.KNCAddr, burn.Amount)
		if err != nil {
			return nil, err
		}

		fields := map[string]interface{}{
			"amount": burnAmount,
		}

		burnPoint, err := client.NewPoint("burn_fees", tags, fields, log.Timestamp)
		if err != nil {
			return nil, err
		}

		points = append(points, burnPoint)
	}

	// build walletFeePoint
	for idx, walletFee := range log.WalletFees {
		tags := map[string]string{
			"tx_hash":      log.TransactionHash.String(),
			"reserve_addr": walletFee.ReserveAddress.String(),
			"wallet_addr":  walletFee.WalletAddress.String(),
			"order":        strconv.Itoa(idx),
		}

		amount, err := is.amountFmt.FormatAmount(blockchain.KNCAddr, walletFee.Amount)
		if err != nil {
			return nil, err
		}

		fields := map[string]interface{}{
			"amount": amount,
		}

		walletFeePoint, err := client.NewPoint("wallet_fees", tags, fields, log.Timestamp)
		if err != nil {
			return nil, err
		}

		points = append(points, walletFeePoint)
	}

	return points, nil
}

func (is *InfluxStorage) rowToTradeLog(row []interface{}) (common.TradeLog, error) {
	var tradeLog common.TradeLog

	timestamp, err := time.Parse(time.RFC3339, row[0].(string))
	if err != nil {
		return tradeLog, err
	}
	blockNumber, err := strconv.ParseUint(row[1].(string), 10, 64)

	ethReceivalAmount, err := row[4].(json.Number).Float64()
	if err != nil {
		return tradeLog, err
	}
	ethReceivalAmountInWei, err := is.amountFmt.ToWei(blockchain.ETHAddr, ethReceivalAmount)
	if err != nil {
		return tradeLog, err
	}

	srcAddress := ethereum.HexToAddress(row[6].(string))
	srcAmount, err := row[8].(json.Number).Float64()
	if err != nil {
		return tradeLog, err
	}
	srcAmountInWei, err := is.amountFmt.ToWei(srcAddress, srcAmount)
	if err != nil {
		return tradeLog, err
	}

	dstAddress := ethereum.HexToAddress(row[7].(string))
	dstAmount, err := row[9].(json.Number).Float64()
	if err != nil {
		return tradeLog, err
	}
	dstAmountInWei, err := is.amountFmt.ToWei(dstAddress, dstAmount)
	if err != nil {
		return tradeLog, err
	}

	fiatAmount, err := row[10].(json.Number).Float64()
	if err != nil {
		return tradeLog, err
	}

	walletFee, err := row[13].(json.Number).Float64()
	if err != nil {
		return tradeLog, err
	}
	walletFeeInWei, err := is.amountFmt.ToWei(blockchain.KNCAddr, walletFee)
	if err != nil {
		return tradeLog, err
	}

	burnFee, err := row[14].(json.Number).Float64()
	if err != nil {
		return tradeLog, err
	}
	burnFeeInWei, err := is.amountFmt.ToWei(blockchain.KNCAddr, burnFee)
	if err != nil {
		return tradeLog, err
	}

	ip, ok := row[15].(string)
	if !ok {
		ip = ""
	}

	country, ok := row[16].(string)
	if err != nil {
		country = ""
	}

	tradeLog = common.TradeLog{
		Timestamp:       timestamp,
		BlockNumber:     blockNumber,
		TransactionHash: ethereum.HexToHash(row[2].(string)),

		EtherReceivalSender: ethereum.HexToAddress(row[3].(string)),
		EtherReceivalAmount: ethReceivalAmountInWei,

		UserAddress: ethereum.HexToAddress(row[5].(string)),
		SrcAddress:  srcAddress,
		DestAddress: dstAddress,
		SrcAmount:   srcAmountInWei,
		DestAmount:  dstAmountInWei,
		FiatAmount:  fiatAmount,

		ReserveAddress: ethereum.HexToAddress(row[11].(string)),
		WalletAddress:  ethereum.HexToAddress(row[12].(string)),
		WalletFee:      walletFeeInWei,
		BurnFee:        burnFeeInWei,

		IP:      ip,
		Country: country,
	}

	return tradeLog, nil
}
