package storage

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/tradelog"

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

// Last24hVolume returns last 24h eth volume of user with given uid.
func (inf *InfluxStorage) Last24hVolume(uid string) (float64, error) {
	var (
		logger = inf.sugar.With("func", "users/storage/InfluxStorage.Last24hVolume",
			"uid", uid)

		influxQueryFmt = `SELECT SUM(amount) FROM (SELECT %[1]s * %[2]s AS amount FROM trades WHERE time >= now() - 24h AND "%[3]s" = '%[4]s') WHERE time >= now() - 24h`
		query          string
	)

	query = fmt.Sprintf(influxQueryFmt,
		tradelog.EthAmount.String(),
		tradelog.EthUSDRate.String(),
		tradelog.UID.String(),
		uid,
	)

	logger = logger.With("query", query)

	logger.Debugw("querying InfluxDB for last 24h ETH volume")
	res, err := influxdb.QueryDB(inf.influxClient, query, inf.dbName)
	if err != nil {
		return 0, err
	}

	if len(res) == 0 || len(res[0].Series) == 0 || len(res[0].Series[0].Values) == 0 || len(res[0].Series[0].Values[0]) < 2 {
		logger.Debugw("no record found")
		return 0, nil
	}

	volume, err := influxdb.GetFloat64FromInterface(res[0].Series[0].Values[0][1])
	if err != nil {
		logger.Debugw("invalid volume returned by query",
			"value", res[0].Series[0].Values[0][1],
			"error", err.Error())
		return 0, nil
	}
	return volume, nil
}
