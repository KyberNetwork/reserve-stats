package influx

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"text/template"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	influxClient "github.com/influxdata/influxdb/client/v2"
	influxModel "github.com/influxdata/influxdb/models"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage/influx/schema"
)

const (
	//rateTableName is the name of influx table storing reserveRate
	rateTableName = "rates"
	//timePrecision is the precision configured for influxDB
	timePrecision   = "s"
	dbDuration      = "30d"
	dbShardDuration = "1d"
)

// RateStorage is the implementation of influxclient to serve as ReserveRate storage
type RateStorage struct {
	sugar      *zap.SugaredLogger
	client     influxClient.Client
	dbName     string
	blkTimeRsv blockchain.BlockTimeResolverInterface

	duration      time.Duration
	shardDuration time.Duration
}

// RateStorageOption configures the RateStorage constructor behaviour.
type RateStorageOption func(rs *RateStorage)

// RateStorageOptionWithRetentionPolicy configures the database with given retention policy.
func RateStorageOptionWithRetentionPolicy(duration, shardDuration time.Duration) RateStorageOption {
	return func(rs *RateStorage) {
		rs.duration = duration
		rs.shardDuration = shardDuration
	}
}

// NewRateInfluxDBStorage return an instance of influx client to store ReserveRate
func NewRateInfluxDBStorage(
	sugar *zap.SugaredLogger,
	client influxClient.Client,
	dbName string,
	blkTimeRsv blockchain.BlockTimeResolverInterface,
	options ...RateStorageOption,
) (*RateStorage, error) {
	var (
		rs   = &RateStorage{sugar: sugar, client: client, dbName: dbName, blkTimeRsv: blkTimeRsv}
		stmt = fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)
	)

	for _, option := range options {
		option(rs)
	}

	if rs.duration != 0 && rs.shardDuration != 0 {
		stmt = fmt.Sprintf(`"%s" WITH DURATION %s REPLICATION 1 SHARD DURATION %s NAME "rates"`,
			stmt,
			dbDuration,
			dbShardDuration,
		)
	}

	q := influxClient.NewQuery(stmt, "", timePrecision)
	response, err := client.Query(q)
	if err != nil {
		return nil, err
	}
	if response.Error() != nil {
		return nil, response.Error()
	}
	return rs, nil
}

// LastBlock returns last stored rate block number from database.
func (rs *RateStorage) LastBlock() (int64, error) {
	stmt := fmt.Sprintf(`SELECT "%s" from "%s" ORDER BY time DESC limit 1`,
		schema.ToBlock.String(),
		rateTableName)
	q := influxClient.NewQuery(stmt, rs.dbName, timePrecision)

	res, err := rs.client.Query(q)
	if err != nil {
		return 0, err
	}

	if res.Error() != nil {
		return 0, res.Error()
	}

	if len(res.Results) != 1 || len(res.Results[0].Series) != 1 || len(res.Results[0].Series[0].Values[0]) != 2 {
		rs.sugar.Infow("no result returned for last block query", "res", res)
		return 0, nil
	}

	return influxdb.GetInt64FromInterface(res.Results[0].Series[0].Values[0][1])
}

func (rs *RateStorage) lastRates(reserveAddr string) (map[string]common.LastRate, error) {
	var (
		logger = rs.sugar.With(
			"func", "reserverates/storage/influx/RateStorage.lastRates",
			"reserve_addr", reserveAddr,
		)
		lastRates = make(map[string]common.LastRate)
	)

	stmt := fmt.Sprintf(`SELECT "%s", "%s", "%s", "%s", "%s", "%s" from "%s" WHERE "%s" = '%s' GROUP BY "%s" ORDER BY time desc limit 1`,
		schema.BuyRate.String(),
		schema.SellRate.String(),
		schema.BuySanityRate.String(),
		schema.SellSanityRate.String(),
		schema.FromBlock.String(),
		schema.ToBlock.String(),
		rateTableName,
		schema.Reserve.String(),
		reserveAddr,
		schema.Pair.String(),
	)

	logger.Debugw("query statement", "query", stmt)

	q := influxClient.NewQuery(stmt,
		rs.dbName,
		timePrecision)

	res, err := rs.client.Query(q)
	if err != nil {
		return nil, err
	}

	if res.Error() != nil {
		return nil, res.Error()
	}

	if len(res.Results) == 0 || len(res.Results[0].Series) == 0 {
		logger.Infow("no rates record found")
		return nil, nil
	}

	for _, serie := range res.Results[0].Series {
		buyRate, err := influxdb.GetFloat64FromInterface(serie.Values[0][1])
		if err != nil {
			return nil, err
		}

		sellRate, err := influxdb.GetFloat64FromInterface(serie.Values[0][2])
		if err != nil {
			return nil, err
		}

		buySanityRate, err := influxdb.GetFloat64FromInterface(serie.Values[0][3])
		if err != nil {
			return nil, err
		}

		sellSanityRate, err := influxdb.GetFloat64FromInterface(serie.Values[0][4])
		if err != nil {
			return nil, err
		}

		fromBlock, err := influxdb.GetUint64FromTagValue(serie.Values[0][5])
		if err != nil {
			return nil, err
		}

		toBlock, err := influxdb.GetUint64FromInterface(serie.Values[0][6])
		if err != nil {
			return nil, err
		}

		lastRates[serie.Tags["pair"]] = common.LastRate{
			Rate: &common.ReserveRateEntry{
				BuyReserveRate:  buyRate,
				BuySanityRate:   buySanityRate,
				SellReserveRate: sellRate,
				SellSanityRate:  sellSanityRate,
			},
			FromBlock: fromBlock,
			ToBlock:   toBlock,
		}
	}

	return lastRates, nil
}

func (rs *RateStorage) constructDataPoint(rsvAddr, pair string, fromBlock, toBlock uint64, rate common.ReserveRateEntry) (*influxClient.Point, error) {
	tags := map[string]string{
		schema.Reserve.String(): rsvAddr,
		schema.Pair.String():    pair,
	}

	if toBlock == 0 {
		toBlock = fromBlock + 1
	}

	fields := map[string]interface{}{
		schema.BuyRate.String():        rate.BuyReserveRate,
		schema.SellRate.String():       rate.SellReserveRate,
		schema.BuySanityRate.String():  rate.BuySanityRate,
		schema.SellSanityRate.String(): rate.SellSanityRate,
		schema.ToBlock.String():        int64(toBlock),
		schema.FromBlock.String():      strconv.FormatUint(fromBlock, 10),
	}

	if rs.blkTimeRsv == nil {
		return nil, errors.New("block time resolver is not available")
	}

	ts, err := rs.blkTimeRsv.Resolve(fromBlock)
	if err != nil {
		return nil, err
	}

	return influxClient.NewPoint(rateTableName, tags, fields, ts)
}

// UpdateRatesRecords update all the rate records from different reserve to influxDB in one go.
// It take a map[reserveAddress] ReserveRates and return error if occurs.
func (rs *RateStorage) UpdateRatesRecords(blockNumber uint64, rateRecords map[string]map[string]common.ReserveRateEntry) error {
	var logger = rs.sugar.With(
		"func", "reserverates/storage/influx/RateStorage.UpdateRatesRecord",
		"block_number", blockNumber,
	)
	bp, err := influxClient.NewBatchPoints(
		influxClient.BatchPointsConfig{
			Database:  rs.dbName,
			Precision: timePrecision,
		},
	)
	if err != nil {
		return err
	}

	for rsvAddr, rateRecord := range rateRecords {
		lastRates, fErr := rs.lastRates(rsvAddr)
		if fErr != nil {
			return fErr
		}
		for pair, rate := range rateRecord {
			var (
				fromBlock uint64
				toBlock   uint64
			)
			lastRate := lastRates[pair].Rate
			lastFromBlock := lastRates[pair].FromBlock
			lastToBlock := lastRates[pair].ToBlock

			switch {
			case lastRate == nil:
				logger.Debugw("no last rate available",
					"reserve_addr", rsvAddr,
					"pair", pair)
				fromBlock = blockNumber
				toBlock = 0
			case *lastRate != rate:
				logger.Debugw("rate changed, starting new rate group",
					"reserve_addr", rsvAddr,
					"last_to_block", lastToBlock,
					"pair", pair,
					"last_rate", lastRate,
					"rate", rate,
				)
				fromBlock = blockNumber
				toBlock = 0
			default:
				logger.Debugw("rate is remain the same as last stored record",
					"reserve_addr", rsvAddr,
					"last_to_block", lastToBlock,
					"pair", pair)
				fromBlock = lastFromBlock
				toBlock = lastToBlock + 1
			}

			pt, fErr := rs.constructDataPoint(rsvAddr, pair, fromBlock, toBlock, rate)
			if fErr != nil {
				return fErr
			}
			bp.AddPoint(pt)
		}
	}
	return rs.client.Write(bp)
}

// GetRatesByTimePoint returns all the rate record in a period of time of a reserve
func (rs *RateStorage) GetRatesByTimePoint(addrs []ethereum.Address, fromTime, toTime uint64) (map[string]map[string][]common.ReserveRates, error) {
	const queryTmpl = `SELECT * FROM "{{.TableName}}" WHERE {{.FromTime }}{{.TimePrecision}} <= time AND time <= {{.ToTime}}{{.TimePrecision}} ` +
		`{{if len .Addrs}}AND ({{range $index, $element := .Addrs}}"reserve" = '{{$element}}'{{if ne $index $.AddrsLastIndex}} OR {{end}}{{end}}){{end}}`
	var (
		logger = rs.sugar.With("reserves", len(addrs),
			"from", fromTime,
			"to", toTime,
		)
		addrsStrs []string
	)

	logger.Debugw("before rendering query statement from template", "query_template", queryTmpl)
	tmpl, err := template.New("queryStmt").Parse(queryTmpl)
	if err != nil {
		return nil, err
	}

	var queryStmtBuf bytes.Buffer
	for _, rsvAddr := range addrs {
		addrsStrs = append(addrsStrs, rsvAddr.Hex())
	}
	if err = tmpl.Execute(&queryStmtBuf, struct {
		TableName      string
		FromTime       uint64
		ToTime         uint64
		TimePrecision  string
		Addrs          []string
		AddrsLastIndex int
	}{
		TableName: rateTableName,
		// convert from milliseconds to seconds
		FromTime:       fromTime / 1000,
		ToTime:         toTime / 1000,
		TimePrecision:  timePrecision,
		Addrs:          addrsStrs,
		AddrsLastIndex: len(addrs) - 1,
	}); err != nil {
		return nil, err
	}

	logger.Debugw("rendered query statement", "rendered_template", queryStmtBuf.String())
	q := influxClient.NewQuery(queryStmtBuf.String(), rs.dbName, timePrecision)
	response, err := rs.client.Query(q)
	if err != nil {
		return nil, err
	}

	if response.Error() != nil {
		return nil, response.Error()
	}

	if len(response.Results) == 0 || len(response.Results[0].Series) == 0 {
		return nil, nil
	}

	return convertQueryResultToRate(response.Results[0].Series[0])
}

func convertRowValueToReserveRate(v []interface{}, idxs schema.FieldsRegistrar) (
	string, // reserve address
	string, // token pair
	common.ReserveRates, // rates
	error) {
	var (
		ts          time.Time
		fromBlock   uint64
		toBlock     uint64
		reserveAddr string
		rate        common.ReserveRates
	)

	intNumber, err := influxdb.GetInt64FromInterface(v[idxs[schema.Time]])
	if err != nil {
		return "", "", common.ReserveRates{}, err
	}
	ts = timeutil.TimestampMsToTime(uint64(intNumber))

	fromBlockStr, ok := v[idxs[schema.FromBlock]].(string)
	if !ok {
		return "", "", common.ReserveRates{}, errors.New("cannot convert influx interface to string")
	}

	if fromBlock, err = strconv.ParseUint(fromBlockStr, 10, 64); err != nil {
		return "", "", common.ReserveRates{}, err
	}

	toBlock, err = influxdb.GetUint64FromInterface(v[idxs[schema.ToBlock]])
	if err != nil {
		return "", "", common.ReserveRates{}, errors.New("failed to convert to block")
	}

	// get pair
	pairName, convertible := v[idxs[schema.Pair]].(string)
	if !convertible {
		return "", "", common.ReserveRates{}, errors.New("cannot convert influx interface to string")
	}
	buyRate, err := influxdb.GetFloat64FromInterface(v[idxs[schema.BuyRate]])
	if err != nil {
		return "", "", common.ReserveRates{}, err
	}
	sellRate, err := influxdb.GetFloat64FromInterface(v[idxs[schema.SellRate]])
	if err != nil {
		return "", "", common.ReserveRates{}, err
	}
	buySanityRate, err := influxdb.GetFloat64FromInterface(v[(idxs)[schema.BuySanityRate]])
	if err != nil {
		return "", "", common.ReserveRates{}, err
	}
	sellSanityRate, err := influxdb.GetFloat64FromInterface(v[idxs[schema.SellSanityRate]])
	if err != nil {
		return "", "", common.ReserveRates{}, err
	}

	reserveAddr, ok = v[idxs[schema.Reserve]].(string)
	if !ok {
		return "", "", common.ReserveRates{}, errors.New("cannot convert influx interface to string")
	}

	rate = common.ReserveRates{
		Timestamp: ts,
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Rates: common.ReserveRateEntry{
			BuyReserveRate:  buyRate,
			SellReserveRate: sellRate,
			BuySanityRate:   buySanityRate,
			SellSanityRate:  sellSanityRate,
		}}
	return reserveAddr, pairName, rate, nil
}

func convertQueryResultToRate(row influxModel.Row) (map[string]map[string][]common.ReserveRates, error) {
	var result = make(map[string]map[string][]common.ReserveRates)

	if len(row.Values) == 0 {
		return nil, nil
	}

	idxs, err := schema.NewFieldsRegistrar(row.Columns)
	if err != nil {
		return nil, err
	}

	for _, v := range row.Values {
		reserveAddr, pairName, rate, err := convertRowValueToReserveRate(v, idxs)
		if err != nil {
			return nil, err
		}

		_, ok := result[reserveAddr]
		if !ok {
			result[reserveAddr] = map[string][]common.ReserveRates{
				pairName: {rate},
			}
			continue
		}

		if _, ok = result[reserveAddr][pairName]; !ok {
			result[reserveAddr][pairName] = []common.ReserveRates{rate}
			continue
		}

		result[reserveAddr][pairName] = append(result[reserveAddr][pairName], rate)

	}
	return result, nil
}
