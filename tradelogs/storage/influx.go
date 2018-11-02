package storage

import (
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"strconv"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	//timePrecision is the precision configured for influxDB
	timePrecision = "ms"
)

// InfluxStorage represent a client to store trade data to influx DB
type InfluxStorage struct {
	dbName       string
	influxClient client.Client
	coreClient   core.Interface
	sugar        *zap.SugaredLogger
}

// NewInfluxStorage init an instance of InfluxStorage
func NewInfluxStorage(sugar *zap.SugaredLogger, dbName string, influxClient client.Client, coreClient core.Interface) (*InfluxStorage, error) {
	storage := &InfluxStorage{
		dbName:       dbName,
		influxClient: influxClient,
		coreClient:   coreClient,
		sugar:        sugar,
	}
	if err := storage.createDB(); err != nil {
		return nil, err
	}
	return storage, nil
}

// SaveTradeLogs persist trade logs to DB
func (is *InfluxStorage) SaveTradeLogs(logs []common.TradeLog) error {
	defer is.influxClient.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  is.dbName,
		Precision: timePrecision,
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

	if len(logs) > 0 {
		is.sugar.Debugw("saved trade logs into influxdb",
			"first_block", logs[0].BlockNumber,
			"last_block", logs[len(logs)-1].BlockNumber,
			"trade_logs", len(logs))
	} else {
		is.sugar.Debugw("no trade log to store")
	}

	return is.influxClient.Close()
}

// LastBlock returns last stored trade log block number from database.
func (is InfluxStorage) LastBlock() (int64, error) {
	q := fmt.Sprintf(`SELECT "block_number","eth_amount" from "trades" ORDER BY time DESC limit 1`)

	res, err := is.queryDB(is.influxClient, q)
	if err != nil {
		return 0, err
	}

	if len(res) != 1 || len(res[0].Series) != 1 || len(res[0].Series[0].Values[0]) != 3 {
		is.sugar.Info("no result returned for last block query")
		return 0, nil
	}

	return influxdb.GetInt64FromTagValue(res[0].Series[0].Values[0][1])
}

// LoadTradeLogs return trade logs from DB
func (is *InfluxStorage) LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error) {
	q := fmt.Sprintf(
		`
		SELECT %[1]s FROM burn_fees WHERE time >= '%[4]s' AND time <= '%[5]s';
		SELECT %[2]s FROM wallet_fees WHERE time >= '%[4]s' AND time <= '%[5]s';
		SELECT %[3]s FROM trades WHERE time >= '%[4]s' AND time <= '%[5]s';
		`,
		"time, tx_hash, reserve_addr, amount",
		"time, tx_hash, reserve_addr, wallet_addr, amount",
		`
		time, block_number, tx_hash, 
		eth_receival_sender, eth_receival_amount, 
		user_addr, src_addr, dst_addr, src_amount, dst_amount, (eth_amount * eth_usd_rate) as fiat_amount, 		
		ip, country
		`,
		from.Format(time.RFC3339),
		to.Format(time.RFC3339),
	)

	logger := is.sugar.With("from", from, "to", to)
	logger.Debug("prepared query statement", "query", q)

	res, err := is.queryDB(is.influxClient, q)
	if err != nil {
		return nil, err
	}

	// Get BurnFees
	burnFeesByTxHash := make(map[ethereum.Hash][]common.BurnFee)

	if len(res[0].Series) == 0 {
		is.sugar.Debug("empty burn fee in query result")
		return nil, nil
	}

	for _, row := range res[0].Series[0].Values {
		txHash, burnFee, err := is.rowToBurnFee(row)
		if err != nil {
			return nil, err
		}
		burnFeesByTxHash[txHash] = append(burnFeesByTxHash[txHash], burnFee)
	}

	// Get WalletFees
	walletFeesByTxHash := make(map[ethereum.Hash][]common.WalletFee)

	if len(res[1].Series) == 0 {
		is.sugar.Debug("empty wallet fee in query result")
	} else {
		for _, row := range res[1].Series[0].Values {
			txHash, walletFee, err := is.rowToWalletFee(row)
			if err != nil {
				return nil, err
			}
			walletFeesByTxHash[txHash] = append(walletFeesByTxHash[txHash], walletFee)
		}
	}

	// Get TradeLogs
	var result []common.TradeLog

	if len(res[2].Series) == 0 {
		is.sugar.Debug("empty trades in query result")
		return nil, nil
	}

	for _, row := range res[2].Series[0].Values {
		tradeLog, err := is.rowToTradeLog(row, burnFeesByTxHash, walletFeesByTxHash)
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

		"eth_rate_provider": log.ETHUSDProvider,
	}

	logger := is.sugar.With(
		"func", "tradelogs/storage/tradeLogToPoint",
		"log", log,
	)

	if blockchain.IsBurnable(log.SrcAddress) {
		if blockchain.IsBurnable(log.DestAddress) {
			if len(log.BurnFees) == 2 {
				tags["src_rsv_addr"] = log.BurnFees[0].ReserveAddress.String()
				tags["dst_rsv_addr"] = log.BurnFees[1].ReserveAddress.String()
			} else {
				logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "2 burn fees (src-dst)")
			}
		} else {
			if len(log.BurnFees) == 1 {
				tags["src_rsv_addr"] = log.BurnFees[0].ReserveAddress.String()
			} else {
				logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "1 burn fees (src)")
			}
		}
	} else if blockchain.IsBurnable(log.DestAddress) {
		if len(log.BurnFees) == 1 {
			tags["dst_rsv_addr"] = log.BurnFees[0].ReserveAddress.String()
		} else {
			logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "1 burn fees (dst)")
		}
	}

	ethReceivalAmount, err := is.coreClient.FromWei(blockchain.ETHAddr, log.EtherReceivalAmount)
	if err != nil {
		return nil, err
	}

	srcAmount, err := is.coreClient.FromWei(log.SrcAddress, log.SrcAmount)
	if err != nil {
		return nil, err
	}

	dstAmount, err := is.coreClient.FromWei(log.DestAddress, log.DestAmount)
	if err != nil {
		return nil, err
	}

	var ethAmount float64

	if log.SrcAddress == blockchain.ETHAddr {
		ethAmount = srcAmount
	} else if log.DestAddress == blockchain.ETHAddr {
		ethAmount = dstAmount
	} else {
		ethAmount = ethReceivalAmount
	}

	fields := map[string]interface{}{
		"eth_receival_amount": ethReceivalAmount,

		"src_amount":   srcAmount,
		"dst_amount":   dstAmount,
		"eth_usd_rate": log.ETHUSDRate,

		"eth_amount": ethAmount,
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
			"ordinal":      strconv.Itoa(idx), // prevent overwrite by other event belong to same trade log
		}

		burnAmount, err := is.coreClient.FromWei(blockchain.KNCAddr, burn.Amount)
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
			"ordinal":      strconv.Itoa(idx), // prevent overwrite by other event belong to same trade log
		}

		amount, err := is.coreClient.FromWei(blockchain.KNCAddr, walletFee.Amount)
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
