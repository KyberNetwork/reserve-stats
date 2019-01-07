package influxdb

import (
	"fmt"
	"strconv"

	"github.com/KyberNetwork/reserve-stats/burnedfees/common"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

const (
	dbName          = "burned_fees"
	measurementName = "burned_fees"
	timePrecision   = "s"
)

// BurnedFeesStorage is the implementation of burnedfees.storage.Interface
// that stores data in InfluxDB.
type BurnedFeesStorage struct {
	sugar                *zap.SugaredLogger
	client               client.Client
	blkTimeRsv           blockchain.BlockTimeResolverInterface
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface
}

// NewBurnedFeesStorage creates new instance of BurnedFeesStorage.
func NewBurnedFeesStorage(
	sugar *zap.SugaredLogger,
	influxClient client.Client,
	blkTimeRsv blockchain.BlockTimeResolverInterface,
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface) (*BurnedFeesStorage, error) {
	st := &BurnedFeesStorage{
		sugar:                sugar,
		client:               influxClient,
		blkTimeRsv:           blkTimeRsv,
		tokenAmountFormatter: tokenAmountFormatter,
	}

	q := client.NewQuery(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName), "", timePrecision)
	response, err := st.client.Query(q)
	if err != nil {
		return nil, err
	}
	if response.Error() != nil {
		return nil, response.Error()
	}
	return st, nil
}

// Store stores the given events to database.
func (bfs *BurnedFeesStorage) Store(events []common.BurnAssignedFeesEvent) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbName,
		Precision: timePrecision,
	})
	if err != nil {
		return err
	}

	for _, event := range events {
		tags := map[string]string{
			"block_number": strconv.FormatUint(event.BlockNumber, 10),
			"tx_hash":      event.TxHash.Hex(),
			"reserve":      event.Reserve.Hex(),
			"sender":       event.Sender.Hex(),
		}

		qty, err := bfs.tokenAmountFormatter.FromWei(blockchain.KNCAddr, event.Quantity)
		if err != nil {
			return err
		}

		fields := map[string]interface{}{
			"qty": qty,
		}

		ts, err := bfs.blkTimeRsv.Resolve(event.BlockNumber)
		if err != nil {
			return err
		}

		pt, err := client.NewPoint(measurementName, tags, fields, ts)
		if err != nil {
			return err
		}

		bp.AddPoint(pt)
	}

	return bfs.client.Write(bp)
}

func (bfs *BurnedFeesStorage) LastBlock() (int64, error) {
	stmt := fmt.Sprintf(`SELECT "block_number", "qty" from "%s" ORDER BY time DESC limit 1`,
		measurementName)
	q := client.NewQuery(stmt, dbName, timePrecision)

	resp, err := bfs.client.Query(q)
	if err != nil {
		return 0, err
	}

	if resp.Error() != nil {
		return 0, resp.Error()
	}

	results := resp.Results

	if len(results) != 1 || len(results[0].Series) != 1 || len(results[0].Series[0].Values[0]) != 3 {
		bfs.sugar.Info("no result returned for last block query")
		return 0, nil
	}

	return influxdb.GetInt64FromTagValue(results[0].Series[0].Values[0][1])
}
