package storage

import (
	"fmt"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	kycedschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/kyced"
	logschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

const (
	//timePrecision is the precision configured for influxDB
	timePrecision           = "s"
	tradeLogMeasurementName = "trades"
)

// InfluxStorage represent a client to store trade data to influx DB
type InfluxStorage struct {
	sugar                *zap.SugaredLogger
	dbName               string
	influxClient         client.Client
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface
	kycChecker           KycChecker

	// traded stored traded addresses to use in a single SaveTradeLogs
	traded map[ethereum.Address]struct{}
}

// NewInfluxStorage init an instance of InfluxStorage
func NewInfluxStorage(sugar *zap.SugaredLogger, dbName string, influxClient client.Client,
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface, kycChecker KycChecker) (*InfluxStorage, error) {
	storage := &InfluxStorage{
		sugar:                sugar,
		dbName:               dbName,
		influxClient:         influxClient,
		tokenAmountFormatter: tokenAmountFormatter,
		kycChecker:           kycChecker,
		traded:               make(map[ethereum.Address]struct{}),
	}
	if err := storage.createDB(); err != nil {
		return nil, err
	}
	return storage, nil
}

// SaveTradeLogs persist trade logs to DB
func (is *InfluxStorage) SaveTradeLogs(logs []common.TradeLog) error {
	logger := is.sugar.With(
		"func", "tradelogs/storage/InfluxStorage.SaveTradeLogs",
	)
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
		logger.Errorw("saving error", "error", err)
		return err
	}

	if len(logs) > 0 {
		logger.Debugw("saved trade logs into influxdb",
			"first_block", logs[0].BlockNumber,
			"last_block", logs[len(logs)-1].BlockNumber,
			"trade_logs", len(logs))
	} else {
		logger.Debugw("no trade log to store")
	}

	// reset traded map to avoid ever growing size
	is.traded = make(map[ethereum.Address]struct{})
	return nil
}

// LastBlock returns last stored trade log block number from database.
func (is InfluxStorage) LastBlock() (int64, error) {
	q := fmt.Sprintf(`SELECT "block_number","eth_amount" from "trades" ORDER BY time DESC limit 1`)

	res, err := influxdb.QueryDB(is.influxClient, q, is.dbName)
	if err != nil {
		return 0, err
	}

	if len(res) != 1 || len(res[0].Series) != 1 || len(res[0].Series[0].Values[0]) != 3 {
		is.sugar.Info("no result returned for last block query")
		return 0, nil
	}

	return influxdb.GetInt64FromInterface(res[0].Series[0].Values[0][1])
}

func prepareTradeLogQuery() string {
	var (
		tradeLogQueryFields = []logschema.FieldName{
			logschema.Time,
			logschema.BlockNumber,
			logschema.EthAmount,
			logschema.UserAddr,
			logschema.SrcAddr,
			logschema.DstAddr,
			logschema.SrcAmount,
			logschema.DstAmount,
			logschema.IP,
			logschema.Country,
			logschema.IntegrationApp,
			logschema.SourceBurnAmount,
			logschema.DestBurnAmount,
			logschema.LogIndex,
			logschema.TxHash,
			logschema.SrcReserveAddr,
			logschema.DstReserveAddr,
			logschema.SourceWalletFeeAmount,
			logschema.DestWalletFeeAmount,
			logschema.WalletAddress,
		}
		tradeLogQuery string
	)
	for _, field := range tradeLogQueryFields {
		tradeLogQuery += field.String() + ", "
	}
	fiatAmount := fmt.Sprintf("(%s * %s) AS %s", logschema.EthAmount.String(), logschema.EthUSDRate.String(), logschema.FiatAmount.String())
	tradeLogQuery += fiatAmount
	return tradeLogQuery
}

// LoadTradeLogs return trade logs from DB
func (is *InfluxStorage) LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error) {
	var (
		result = make([]common.TradeLog, 0)
		q      = fmt.Sprintf(
			`
		SELECT %[1]s FROM %[2]s WHERE time >= '%[3]s' AND time <= '%[4]s';
		`,
			prepareTradeLogQuery(),
			tradeLogMeasurementName,
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

	res, err := influxdb.QueryDB(is.influxClient, q, is.dbName)
	if err != nil {
		return nil, err
	}

	// Get TradeLogs
	if len(res) == 0 || len(res[0].Series) == 0 || len(res[0].Series[0].Values) == 0 {
		is.sugar.Debug("empty trades in query result")
		return result, nil
	}
	idxs, err := logschema.NewFieldsRegistrar(res[0].Series[0].Columns)
	if err != nil {
		return nil, err
	}
	for _, row := range res[0].Series[0].Values {

		tradeLog, err := is.rowToTradeLog(row, idxs)
		if err != nil {
			return nil, err
		}
		result = append(result, tradeLog)
	}

	return result, nil
}

// createDB creates the database will be used for storing trade logs measurements.
func (is *InfluxStorage) createDB() error {
	_, err := influxdb.QueryDB(is.influxClient, fmt.Sprintf("CREATE DATABASE %s", is.dbName), is.dbName)
	return err
}

//getBurnAmount return the burn amount in float for src and
func (is *InfluxStorage) getBurnAmount(log common.TradeLog) (float64, float64, error) {
	var (
		logger = is.sugar.With(
			"func", "tradelogs/storage/getBurnAmount",
			"log", log,
		)
		srcAmount float64
		dstAmount float64
	)
	if blockchain.IsBurnable(log.SrcAddress) {
		if len(log.BurnFees) < 1 {
			logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "at least 1 burn fees (src)")
			return srcAmount, dstAmount, nil
		}
		srcAmount, err := is.tokenAmountFormatter.FromWei(blockchain.KNCAddr, log.BurnFees[0].Amount)
		if err != nil {
			return srcAmount, dstAmount, err
		}

		if blockchain.IsBurnable(log.DestAddress) {
			if len(log.BurnFees) < 2 {
				logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "2 burn fees (src-dst)")
				return srcAmount, dstAmount, nil
			}
			dstAmount, err = is.tokenAmountFormatter.FromWei(blockchain.KNCAddr, log.BurnFees[1].Amount)
			if err != nil {
				return srcAmount, dstAmount, err
			}
			return srcAmount, dstAmount, nil
		}
		return srcAmount, dstAmount, nil
	}

	if blockchain.IsBurnable(log.DestAddress) {
		if len(log.BurnFees) < 1 {
			logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "at least 1 burn fees (dst)")
			return srcAmount, dstAmount, nil
		}
		dstAmount, err := is.tokenAmountFormatter.FromWei(blockchain.KNCAddr, log.BurnFees[0].Amount)
		if err != nil {
			return srcAmount, dstAmount, err
		}
		return srcAmount, dstAmount, nil
	}

	return srcAmount, dstAmount, nil
}

func (is *InfluxStorage) getWalletFeeAmount(log common.TradeLog) (float64, float64, error) {
	var (
		logger = is.sugar.With(
			"func", "tradelogs/storage/getWalletFeeAmount",
			"log", log,
		)
		dstAmount    float64
		srcAmount    float64
		srcAmountSet bool
	)
	for _, walletFee := range log.WalletFees {
		amount, err := is.tokenAmountFormatter.FromWei(blockchain.KNCAddr, walletFee.Amount)
		if err != nil {
			return dstAmount, srcAmount, err
		}

		if walletFee.ReserveAddress == log.SrcReserveAddress && !srcAmountSet {
			srcAmount = amount
			srcAmountSet = true
		} else if walletFee.ReserveAddress == log.DstReserveAddress {
			dstAmount = amount
		} else {
			logger.Warnw("unexpected wallet fees with unrecognized reserve address", "wallet fee", walletFee)
		}
	}
	return srcAmount, dstAmount, nil
}

func (is *InfluxStorage) tradeLogToPoint(log common.TradeLog) ([]*client.Point, error) {
	var points []*client.Point
	var walletAddr ethereum.Address
	if len(log.WalletFees) > 0 {
		walletAddr = log.WalletFees[0].WalletAddress
	}

	tags := map[string]string{

		logschema.UserAddr.String(): log.UserAddress.String(),

		logschema.SrcAddr.String():        log.SrcAddress.String(),
		logschema.DstAddr.String():        log.DestAddress.String(),
		logschema.IntegrationApp.String(): log.IntegrationApp,
		logschema.LogIndex.String():       strconv.FormatUint(uint64(log.Index), 10),

		logschema.Country.String(): log.Country,

		logschema.LogIndex.String(): strconv.FormatUint(uint64(log.Index), 10),
	}

	if !blockchain.IsZeroAddress(log.SrcReserveAddress) {
		tags[logschema.SrcReserveAddr.String()] = log.SrcReserveAddress.String()
	}
	if !blockchain.IsZeroAddress(log.DstReserveAddress) {
		tags[logschema.DstReserveAddr.String()] = log.DstReserveAddress.String()
	}
	if !blockchain.IsZeroAddress(log.WalletAddress) {
		tags[logschema.WalletAddress.String()] = walletAddr.String()
	}

	ethAmount, err := is.tokenAmountFormatter.FromWei(blockchain.ETHAddr, log.EthAmount)
	if err != nil {
		return nil, err
	}

	srcAmount, err := is.tokenAmountFormatter.FromWei(log.SrcAddress, log.SrcAmount)
	if err != nil {
		return nil, err
	}

	dstAmount, err := is.tokenAmountFormatter.FromWei(log.DestAddress, log.DestAmount)
	if err != nil {
		return nil, err
	}

	srcBurnAmount, dstBurnAmount, err := is.getBurnAmount(log)
	if err != nil {
		return nil, err
	}

	srcWalletFee, dstWalletFee, err := is.getWalletFeeAmount(log)
	if err != nil {
		return nil, err
	}
	fields := map[string]interface{}{

		logschema.SrcAmount.String():        srcAmount,
		logschema.DstAmount.String():        dstAmount,
		logschema.EthUSDRate.String():       log.ETHUSDRate,
		logschema.SourceBurnAmount.String(): srcBurnAmount,
		logschema.DestBurnAmount.String():   dstBurnAmount,

		logschema.EthAmount.String():             ethAmount,
		logschema.BlockNumber.String():           int64(log.BlockNumber),
		logschema.TxHash.String():                log.TransactionHash.String(),
		logschema.IP.String():                    log.IP,
		logschema.EthUSDProvider.String():        log.ETHUSDProvider,
		logschema.SourceWalletFeeAmount.String(): srcWalletFee,
		logschema.DestWalletFeeAmount.String():   dstWalletFee,
	}
	tradePoint, err := client.NewPoint(tradeLogMeasurementName, tags, fields, log.Timestamp)
	if err != nil {
		return nil, err
	}

	points = append(points, tradePoint)

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
	response, err := influxdb.QueryDB(is.influxClient, q, is.dbName)
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

	kyced, err := is.kycChecker.IsKYCedAtTime(logItem.UserAddress, logItem.Timestamp)
	if err != nil {
		return nil, err
	}

	if !kyced {
		logger.Debugw("user has not been kyced yet")
		return nil, nil
	}

	logger.Debugw("user has been kyced")
	tags := map[string]string{
		kycedschema.UserAddress.String(): logItem.UserAddress.Hex(),
		kycedschema.Country.String():     logItem.Country,
	}

	for _, walletFee := range logItem.WalletFees {
		tags[kycedschema.WalletAddress.String()] = walletFee.WalletAddress.Hex()
	}

	fields := map[string]interface{}{
		kycedschema.KYCed.String(): true,
	}

	point, err := client.NewPoint("kyced", tags, fields, logItem.Timestamp)
	return point, err
}
