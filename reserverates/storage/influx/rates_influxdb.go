package influx

import (
	"bytes"
	"errors"
	"strconv"
	"text/template"

	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage/influx/schema"
	ethereum "github.com/ethereum/go-ethereum/common"
	influxClient "github.com/influxdata/influxdb/client/v2"
	influxModel "github.com/influxdata/influxdb/models"
)

const (
	//RateTableName is the name of influx table storing reserveRate
	RateTableName = "reserve_rate"
	//timePrecision is the precision configured for influxDB
	timePrecision = "ms"
)

// RateStorage is the implementation of influxclient to serve as ReserveRate storage
type RateStorage struct {
	sugar  *zap.SugaredLogger
	client influxClient.Client
	dbName string
}

// NewRateInfluxDBStorage return an instance of influx client to store ReserveRate
func NewRateInfluxDBStorage(sugar *zap.SugaredLogger, client influxClient.Client, dbName string) (*RateStorage, error) {
	q := influxClient.NewQuery("CREATE DATABASE "+dbName, "", timePrecision)
	response, err := client.Query(q)
	if err != nil {
		return nil, err
	}
	if response.Error() != nil {
		return nil, response.Error()
	}
	return &RateStorage{sugar: sugar, client: client, dbName: dbName}, nil
}

// LastBlock returns last stored rate block number from database.
func (rs *RateStorage) LastBlock() (int64, error) {
	q := influxClient.NewQuery("SELECT block_number, buy_rate from reserve_rate ORDER BY time DESC limit 1", rs.dbName, timePrecision)

	res, err := rs.client.Query(q)
	if err != nil {
		return 0, err
	}

	if res.Error() != nil {
		return 0, res.Error()
	}

	if len(res.Results) != 1 || len(res.Results[0].Series) != 1 || len(res.Results[0].Series[0].Values[0]) != 3 {
		rs.sugar.Infow("no result returned for last block query", "res", res)
		return 0, nil
	}

	return influxdb.GetInt64FromTagValue(res.Results[0].Series[0].Values[0][1])
}

// UpdateRatesRecords update all the rate records from different reserve to influxDB in one go.
// It take a map[reserveAddress] ReserveRates and return error if occurs.
func (rs *RateStorage) UpdateRatesRecords(rateRecords map[string]common.ReserveRates) error {
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
		for pair, rate := range rateRecord.Data {
			tags := map[string]string{
				schema.Reserve.String():     rsvAddr,
				schema.Pair.String():        pair,
				schema.BlockNumber.String(): strconv.FormatUint(rateRecord.BlockNumber, 10),
			}
			fields := map[string]interface{}{
				schema.BuyRate.String():        rate.BuyReserveRate,
				schema.SellRate.String():       rate.SellReserveRate,
				schema.BuySanityRate.String():  rate.BuySanityRate,
				schema.SellSanityRate.String(): rate.SellSanityRate,
			}
			pt, err := influxClient.NewPoint(RateTableName, tags, fields, rateRecord.Timestamp)
			if err != nil {
				return err
			}
			bp.AddPoint(pt)
		}
	}
	return rs.client.Write(bp)
}

// GetRatesByTimePoint returns all the rate record in a period of time of a reserve
func (rs *RateStorage) GetRatesByTimePoint(addrs []ethereum.Address, fromTime, toTime uint64) (map[string]map[uint64]common.ReserveRates, error) {
	const queryTmpl = `SELECT * FROM "{{.TableName}}" WHERE {{.FromTime }}{{.TimePrecision}} <= time AND time <= {{.ToTime}}{{.TimePrecision}} ` +
		`{{if len .Addrs}}AND ({{range $index, $element := .Addrs}}"reserve" = '{{$element}}'{{if ne $index $.AddrsLastIndex}} OR {{end}}{{end}}){{end}}`
	var (
		logger = rs.sugar.With("reserves", len(addrs),
			"from", fromTime,
			"to", toTime,
		)
		addrsStrs []string
	)

	logger.Debugw("before rendering query statement from template", "query_tempalte", queryTmpl)
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
		TableName:      RateTableName,
		FromTime:       fromTime,
		ToTime:         toTime,
		TimePrecision:  timePrecision,
		Addrs:          addrsStrs,
		AddrsLastIndex: len(addrs) - 1,
	}); err != nil {
		return nil, err
	}

	logger.Debugw("rendered query statement", "rendered_tempalte", queryStmtBuf.String())
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

func convertRowValueToReserveRate(v []interface{}, idxs schema.FieldsRegistrar) (*common.ReserveRates, error) {
	rate := common.ReserveRates{
		Data: make(map[string]common.ReserveRateEntry),
	}
	// Get Time
	intNumber, err := influxdb.GetInt64FromInterface(v[idxs[schema.Time]])
	if err != nil {
		return nil, err
	}
	timeStamp := timeutil.TimestampMsToTime(uint64(intNumber))
	rate.Timestamp = timeStamp

	// get Block number
	blockNumberStr, ok := v[idxs[schema.BlockNumber]].(string)
	if !ok {
		return nil, errors.New("cannot convert influx interface to string")
	}
	blockNumber, err := strconv.ParseUint(blockNumberStr, 10, 64)
	if err != nil {
		return nil, err
	}
	rate.BlockNumber = blockNumber
	// get pair
	pairName, convertible := v[idxs[schema.Pair]].(string)
	if !convertible {
		return nil, errors.New("cannot convert influx interface to string")
	}
	buyRate, err := influxdb.GetFloat64FromInterface(v[idxs[schema.BuyRate]])
	if err != nil {
		return nil, err
	}
	sellRate, err := influxdb.GetFloat64FromInterface(v[idxs[schema.SellRate]])
	if err != nil {
		return nil, err
	}
	buySanityRate, err := influxdb.GetFloat64FromInterface(v[(idxs)[schema.BuySanityRate]])
	if err != nil {
		return nil, err
	}
	sellSanityRate, err := influxdb.GetFloat64FromInterface(v[idxs[schema.SellSanityRate]])
	if err != nil {
		return nil, err
	}
	reserve, ok := v[idxs[schema.Reserve]].(string)
	if !ok {
		return nil, errors.New("cannot conver influx interface to string")
	}
	rate.Reserve = reserve

	rate.Data[pairName] = common.ReserveRateEntry{
		BuyReserveRate:  buyRate,
		SellReserveRate: sellRate,
		BuySanityRate:   buySanityRate,
		SellSanityRate:  sellSanityRate,
	}
	return &rate, nil
}

func convertQueryResultToRate(row influxModel.Row) (map[string]map[uint64]common.ReserveRates, error) {
	var (
		result = make(map[string]map[uint64]common.ReserveRates)
	)
	if len(row.Values) == 0 {
		return nil, nil
	}
	idxs, err := schema.NewFieldsRegistrar(row.Columns)
	if err != nil {
		return nil, err
	}
	//rates := make(map[uint64]common.ReserveRates)
	for _, v := range row.Values {
		rate, err := convertRowValueToReserveRate(v, idxs)
		if err != nil {
			return nil, err
		}

		rates, ok := result[rate.Reserve]
		if !ok {
			result[rate.Reserve] = map[uint64]common.ReserveRates{rate.BlockNumber: *rate}
			continue
		}

		if _, ok = rates[rate.BlockNumber]; !ok {
			result[rate.Reserve][rate.BlockNumber] = *rate
			continue
		}

		//append this rate.Pair to the total record.
		for pair, rateEntry := range rate.Data {
			result[rate.Reserve][rate.BlockNumber].Data[pair] = rateEntry
		}

	}
	return result, nil
}
