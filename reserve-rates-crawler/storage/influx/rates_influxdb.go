package influx

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/common"
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/storage/influx/schema"
	ethereum "github.com/ethereum/go-ethereum/common"
	influxClient "github.com/influxdata/influxdb/client/v2"
	influxModel "github.com/influxdata/influxdb/models"
)

const (
	//RateDBName is the name of influx database storing reserveRate
	RateDBName = "ReserveRate"
	//RateTableName is the name of influx table storing reserveRate
	RateTableName = "reserve_rate"
	//timePrecision is the precision configured for influxDB
	timePrecision = "ms"
)

var errCantConvert = errors.New("cannot convert response from influxDB to pre-defined struct")

// RateStorage is the implementation of influxclient to serve as ReserveRate storage
type RateStorage struct {
	client influxClient.Client
}

// NewRateInfluxDBStorage return an instance of influx client to store ReserveRate
func NewRateInfluxDBStorage(client influxClient.Client) (*RateStorage, error) {
	q := influxClient.NewQuery("CREATE DATABASE "+RateDBName, "", timePrecision)
	response, err := client.Query(q)
	if err != nil {
		return nil, err
	}
	if response.Error() != nil {
		return nil, response.Error()
	}
	return &RateStorage{client: client}, nil
}

// UpdateRatesRecords update all the rate records from different reserve to influxDB in one go.
// It take a map[reserveAddress] ReserveRates and return error if occurs.
func (rs *RateStorage) UpdateRatesRecords(rateRecords map[string]common.ReserveRates) error {
	bp, err := influxClient.NewBatchPoints(
		influxClient.BatchPointsConfig{
			Database:  RateDBName,
			Precision: timePrecision,
		},
	)
	if err != nil {
		return err
	}

	for rsvAddr, rateRecord := range rateRecords {

		for pair, rate := range rateRecord.Data {
			tags := map[string]string{
				schema.Reserve.String(): rsvAddr,
				schema.Pair.String():    pair,
			}
			fields := map[string]interface{}{
				schema.BuyRate.String():        rate.BuyReserveRate,
				schema.SellRate.String():       rate.SellReserveRate,
				schema.BuySanityRate.String():  rate.BuySanityRate,
				schema.SellSanityRate.String(): rate.SellSanityRate,
				// InfluxDB get parsing error if the input is uint64. Must use int64
				schema.BlockNumber.String(): int64(rateRecord.BlockNumber),
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
func (rs *RateStorage) GetRatesByTimePoint(rsvAddr ethereum.Address, fromTime, toTime uint64) (map[uint64]common.ReserveRates, error) {
	result := make(map[uint64]common.ReserveRates)
	command := fmt.Sprintf("SELECT * FROM %s WHERE time >= %d%s AND \"reserve\"='%s' AND time<= %d%s Order By time", RateTableName, fromTime, timePrecision, rsvAddr.Hex(), toTime, timePrecision)
	q := influxClient.NewQuery(command, RateDBName, timePrecision)
	response, err := rs.client.Query(q)
	if err != nil {
		return result, err
	}
	if response.Error() != nil {
		return result, response.Error()
	}
	if len(response.Results) == 0 || len(response.Results[0].Series) == 0 {
		return result, nil
	}
	return convertQueryResultToRate(response.Results[0].Series[0])
}

func getInt64FromInterface(v interface{}) (int64, error) {
	number, convertible := v.(json.Number)
	if !convertible {
		return 0, errCantConvert
	}
	return number.Int64()
}

func getFloat64FromInterface(v interface{}) (float64, error) {
	number, convertible := v.(json.Number)
	if !convertible {
		return 0, errCantConvert
	}
	return number.Float64()
}

func convertRowValueToReserveRate(v []interface{}, idxs *schema.FieldsRegistrar) (*common.ReserveRates, error) {
	rate := common.ReserveRates{
		Data: make(map[string]common.ReserveRateEntry),
	}
	// Get Time
	intNumber, err := getInt64FromInterface(v[(*idxs)[schema.Time]])
	if err != nil {
		return nil, err
	}
	timeStamp := timeutil.TimestampMsToTime(uint64(intNumber))
	// get Block number
	intNumber, err = getInt64FromInterface(v[(*idxs)[schema.BlockNumber]])
	if err != nil {
		return nil, err
	}
	rate.Timestamp = timeStamp
	rate.BlockNumber = uint64(intNumber)
	// get pair
	pairName, convertible := v[(*idxs)[schema.Pair]].(string)
	if !convertible {
		return nil, errCantConvert
	}
	buyRate, err := getFloat64FromInterface(v[(*idxs)[schema.BuyRate]])
	if err != nil {
		return nil, err
	}
	sellRate, err := getFloat64FromInterface(v[(*idxs)[schema.SellRate]])
	if err != nil {
		return nil, err
	}
	buySanityRate, err := getFloat64FromInterface(v[(*idxs)[schema.BuySanityRate]])
	if err != nil {
		return nil, err
	}
	sellSanityRate, err := getFloat64FromInterface(v[(*idxs)[schema.SellSanityRate]])
	if !convertible {
		return nil, errCantConvert
	}

	rate.Data[pairName] = common.ReserveRateEntry{
		BuyReserveRate:  buyRate,
		SellReserveRate: sellRate,
		BuySanityRate:   buySanityRate,
		SellSanityRate:  sellSanityRate,
	}
	return &rate, nil
}

func convertQueryResultToRate(row influxModel.Row) (map[uint64]common.ReserveRates, error) {
	if len(row.Values) == 0 {
		return nil, nil
	}
	idxs, err := schema.NewFieldsRegistrar(row.Columns)
	if err != nil {
		return nil, err
	}
	rates := make(map[uint64]common.ReserveRates)
	for _, v := range row.Values {
		rate, err := convertRowValueToReserveRate(v, idxs)
		if err != nil {
			return rates, err
		}
		curRate, ok := rates[rate.BlockNumber]
		if !ok {
			rates[rate.BlockNumber] = *rate
		} else {
			//append this rate.Pair to the total record.
			for pair, rateEntry := range rate.Data {
				curRate.Data[pair] = rateEntry
			}
			rates[rate.BlockNumber] = curRate
		}
	}
	return rates, nil
}
