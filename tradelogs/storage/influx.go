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
	sugar        *zap.SugaredLogger
	dbName       string
	influxClient client.Client
	coreClient   core.Interface
	kycChecker   kycChecker

	// traded stored traded addresses to use in a single SaveTradeLogs
	traded map[ethereum.Address]struct{}
}

// NewInfluxStorage init an instance of InfluxStorage
func NewInfluxStorage(sugar *zap.SugaredLogger, dbName string, influxClient client.Client,
	coreClient core.Interface, kycChecker kycChecker) (*InfluxStorage, error) {
	storage := &InfluxStorage{
		sugar:        sugar,
		dbName:       dbName,
		influxClient: influxClient,
		coreClient:   coreClient,
		kycChecker:   kycChecker,
		traded:       make(map[ethereum.Address]struct{}),
	}
	if err := storage.createDB(); err != nil {
		return nil, err
	}
	return storage, nil
}

// SaveTradeLogs persist trade logs to DB
func (is *InfluxStorage) SaveTradeLogs(logs []common.TradeLog) error {
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

	// reset traded map to avoid ever growing size
	is.traded = make(map[ethereum.Address]struct{})
	return nil
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
	var (
		result = make([]common.TradeLog, 0)

		q = fmt.Sprintf(
			`
		SELECT %[1]s FROM burn_fees WHERE time >= '%[4]s' AND time <= '%[5]s' GROUP BY tx_hash, trade_log_index;
		SELECT %[2]s FROM wallet_fees WHERE time >= '%[4]s' AND time <= '%[5]s' GROUP BY tx_hash, trade_log_index;
		SELECT %[3]s FROM trades WHERE time >= '%[4]s' AND time <= '%[5]s' GROUP BY tx_hash, log_index;
		`,
			"time, reserve_addr, amount, log_index",
			"time, reserve_addr, wallet_addr, amount, log_index",
			`
		time, block_number, 
		eth_receival_sender, eth_receival_amount, 
		user_addr, src_addr, dst_addr, src_amount, dst_amount, (eth_amount * eth_usd_rate) as fiat_amount, 		
		ip, country
		`,
			from.Format(time.RFC3339),
			to.Format(time.RFC3339),
		)

		logger = is.sugar.With(
			"func", "tradelogs/storage/InfluxStorage.LoadTradLogs",
			"from", from,
			"to", to,
		)
	)

	logger.Debug("prepared query statement", "query", q)

	res, err := is.queryDB(is.influxClient, q)
	if err != nil {
		return nil, err
	}

	// Get BurnFees
	if len(res[0].Series) == 0 {
		is.sugar.Debug("empty burn fee in query result")
		return nil, nil
	}

	// map [tx_hash][trade_log_index][]common.BurnFee
	burnFeesByTxHash := make(map[ethereum.Hash]map[uint][]common.BurnFee)
	for _, row := range res[0].Series {
		txHash, tradeLogIndex, burnFees, err := is.rowToBurnFees(row)
		if err != nil {
			return nil, err
		}
		_, exist := burnFeesByTxHash[txHash]
		if !exist {
			burnFeesByTxHash[txHash] = make(map[uint][]common.BurnFee)
		}
		burnFeesByTxHash[txHash][uint(tradeLogIndex)] = burnFees
	}

	// Get WalletFees
	// map [tx_hash][trade_log_index][]common.WalletFee
	walletFeesByTxHash := make(map[ethereum.Hash]map[uint][]common.WalletFee)

	if len(res[1].Series) == 0 {
		is.sugar.Debug("empty wallet fee in query result")
	} else {
		for _, row := range res[1].Series {
			txHash, tradeLogIndex, walletFees, err := is.rowToWalletFees(row)
			if err != nil {
				return nil, err
			}
			_, exist := walletFeesByTxHash[txHash]
			if !exist {
				walletFeesByTxHash[txHash] = make(map[uint][]common.WalletFee)
			}
			walletFeesByTxHash[txHash][uint(tradeLogIndex)] = walletFees
		}
	}

	// Get TradeLogs
	if len(res[2].Series) == 0 {
		is.sugar.Debug("empty trades in query result")
		return result, nil
	}

	for _, row := range res[2].Series {
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
	var walletAddr ethereum.Address
	if len(log.WalletFees) > 0 {
		walletAddr = log.WalletFees[0].WalletAddress
	}

	tags := map[string]string{
		"block_number": strconv.FormatUint(log.BlockNumber, 10),
		"tx_hash":      log.TransactionHash.String(),

		"eth_receival_sender": log.EtherReceivalSender.String(),

		"user_addr": log.UserAddress.String(),

		"src_addr": log.SrcAddress.String(),
		"dst_addr": log.DestAddress.String(),

		"wallet_addr": walletAddr.String(),

		"country": log.Country,
		"ip":      log.IP,

		"eth_rate_provider": log.ETHUSDProvider,

		"log_index": strconv.FormatUint(uint64(log.Index), 10),
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
	for _, burn := range log.BurnFees {
		tags := map[string]string{
			"tx_hash":         log.TransactionHash.String(),
			"reserve_addr":    burn.ReserveAddress.String(),
			"log_index":       strconv.FormatUint(uint64(burn.Index), 10),
			"trade_log_index": strconv.FormatUint(uint64(log.Index), 10),
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
	for _, walletFee := range log.WalletFees {
		tags := map[string]string{
			"tx_hash":         log.TransactionHash.String(),
			"reserve_addr":    walletFee.ReserveAddress.String(),
			"wallet_addr":     walletFee.WalletAddress.String(),
			"log_index":       strconv.FormatUint(uint64(walletFee.Index), 10),
			"trade_log_index": strconv.FormatUint(uint64(log.Index), 10),
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

	firstTradePoint, err := is.assembleFirstTradePoint(log)
	if err != nil {
		return nil, err
	}
	if firstTradePoint != nil {
		points = append(points, firstTradePoint)
	}

	kycedPoint, err := is.assembleKYCPoint(log)
	if err != nil {
		return nil, err
	}

	if kycedPoint != nil {
		points = append(points, kycedPoint)
	}

	return points, nil
}

func (is *InfluxStorage) assembleFirstTradePoint(logItem common.TradeLog) (*client.Point, error) {
	var logger = is.sugar.With(
		"func", "tradelogs/storage/InfluxStorage.assembleFirstTradePoint",
		"timestamp", logItem.Timestamp.String(),
		"user_addr", logItem.UserAddress.Hex(),
		"country", logItem.Country,
	)

	if _, ok := is.traded[logItem.UserAddress]; ok {
		logger.Debug("user has already traded, ignoring")
		return nil, nil
	}

	traded, err := is.userTraded(logItem.UserAddress)
	if err != nil {
		return nil, err
	}

	if traded {
		return nil, nil
	}

	logger.Debugw("user first trade")
	tags := map[string]string{
		"user_addr": logItem.UserAddress.Hex(),
		"country":   logItem.Country,
	}

	for _, walletFee := range logItem.WalletFees {
		tags["wallet_addr"] = walletFee.WalletAddress.Hex()
	}

	fields := map[string]interface{}{
		"traded": true,
	}

	point, err := client.NewPoint("first_trades", tags, fields, logItem.Timestamp)
	if err != nil {
		return nil, err
	}

	is.traded[logItem.UserAddress] = struct{}{}
	return point, nil
}

func (is *InfluxStorage) userTraded(addr ethereum.Address) (bool, error) {
	q := fmt.Sprintf("SELECT traded FROM first_trades WHERE user_addr='%s'", addr.String())
	response, err := is.queryDB(is.influxClient, q)
	if err != nil {
		return false, err
	}
	// if there is no record, this mean the address has not traded yet
	if (len(response) == 0) || (len(response[0].Series) == 0) || (len(response[0].Series[0].Values) == 0) {
		return false, nil
	}
	return true, nil
}

func (is *InfluxStorage) assembleKYCPoint(logItem common.TradeLog) (*client.Point, error) {
	var logger = is.sugar.With(
		"func", "tradelogs/storage/InfluxStorage.assembleKYCPoint",
		"timestamp", logItem.Timestamp.String(),
		"user_addr", logItem.UserAddress.Hex(),
		"country", logItem.Country,
	)

	kyced, err := is.kycChecker.IsKYCed(logItem.UserAddress, logItem.Timestamp)
	if err != nil {
		return nil, err
	}

	if !kyced {
		logger.Debugw("user has not been kyced yet")
		return nil, nil
	}

	logger.Debugw("user has been kyced")
	tags := map[string]string{
		"user_addr": logItem.UserAddress.Hex(),
		"country":   logItem.Country,
	}

	for _, walletFee := range logItem.WalletFees {
		tags["wallet_addr"] = walletFee.WalletAddress.Hex()
	}

	fields := map[string]interface{}{
		"kyced": true,
	}

	point, err := client.NewPoint("kyced", tags, fields, logItem.Timestamp)
	return point, err
}
