package influx

import (
	"encoding/json"
	"errors"
	"fmt"

	utils "github.com/KyberNetwork/reserve-stats/lib/common"
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	influxClient "github.com/influxdata/influxdb/client/v2"
	influxModel "github.com/influxdata/influxdb/models"
)

const (
	//RateDBName is the name of influx database storing reserveRate
	RateDBName = "ReserveRate"
	//RateTableName is the name of influx table storing reserveRate
	RateTableName = "reserve_rate"
	//TimePrecision is the precision configured for influxDB
	TimePrecision = "ms"
)

var errCantConvert error = errors.New("cannot convert response from influxDB to pre-defined struct")

// RateStorage is the implementation of influxclient to serve as ReserveRate storage
type RateStorage struct {
	client influxClient.Client
}

// NewRateInfluxDBStorage return an instance of influx client to store ReserveRate
func NewRateInfluxDBStorage(client influxClient.Client) (*RateStorage, error) {
	q := influxClient.NewQuery("CREATE DATABASE "+RateDBName, "", TimePrecision)
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
			Precision: TimePrecision,
		},
	)
	if err != nil {
		return err
	}

	for rsvAddr, rateRecord := range rateRecords {

		for pair, rate := range rateRecord.Data {
			// InfluxDB get parsing error if the input is uint64. Must use int64
			tags := map[string]string{
				Reserve.String(): rsvAddr,
				Pair.String():    pair,
			}
			fields := map[string]interface{}{
				BuyRate.String():        rate.BuyReserveRate,
				SellRate.String():       rate.SellReserveRate,
				BuySanityRate.String():  rate.BuySanityRate,
				SellSanityRate.String(): rate.SellSanityRate,
				BlockNumber.String():    int64(rateRecord.BlockNumber),
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
func (rs *RateStorage) GetRatesByTimePoint(rsvAddr ethereum.Address, fromTime, toTime uint64) ([]common.ReserveRates, error) {
	result := []common.ReserveRates{}
	command := fmt.Sprintf("SELECT * FROM %s WHERE time >= %d%s AND \"reserve\"='%s' AND time<= %d%s Order By time", RateTableName, fromTime, TimePrecision, rsvAddr.Hex(), toTime, TimePrecision)
	q := influxClient.NewQuery(command, RateDBName, TimePrecision)
	response, err := rs.client.Query(q)
	if err != nil {
		return result, err
	}
	if response.Error() != nil {
		return result, response.Error()
	}
	if len(response.Results) == 0 || len(response.Results[0].Series) == 0 {
		return []common.ReserveRates{}, nil
	}
	return convertQueryResultTorRate(response.Results[0].Series[0])
}

func getIndexOfFieldS(fieldNames []string) map[RateSchemaFieldName]int {
	result := make(map[RateSchemaFieldName]int)
	for idx, fieldNameStr := range fieldNames {
		fieldName, ok := RateSchemaFields[fieldNameStr]
		if ok {
			result[fieldName] = idx
		}
	}
	return result
}

func getint64FromInterface(v interface{}) (int64, error) {
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

func convertQueryResultTorRate(row influxModel.Row) ([]common.ReserveRates, error) {
	if len(row.Values) == 0 {
		return []common.ReserveRates{}, nil
	}
	idxs := getIndexOfFieldS(row.Columns)
	rateEntry := make(map[string]common.ReserveRateEntry)
	rate := common.ReserveRates{
		Data: rateEntry,
	}
	rates := []common.ReserveRates{rate}
	firstRecordProcessed := false
	nRate := 0
	for _, v := range row.Values {
		// Get Time
		intNumber, err := getint64FromInterface(v[idxs[Time]])
		if err != nil {
			return nil, err
		}

		timeStamp := utils.TimestampMsToTime(uint64(intNumber))
		// New record with new Timestamp
		if rate.Timestamp != timeStamp && firstRecordProcessed {
			rates = append(rates, rate)
			rate = common.ReserveRates{}
			rateEntry = make(map[string]common.ReserveRateEntry)
			nRate++
		} else {
			rate = rates[nRate]
			rateEntry = rate.Data
		}
		rate.Timestamp = timeStamp
		// get Block number
		intNumber, err = getint64FromInterface(v[idxs[BlockNumber]])
		if err != nil {
			return nil, err
		}
		rate.BlockNumber = uint64(intNumber)
		// get pair
		pairName, convertible := v[idxs[Pair]].(string)
		if !convertible {
			return nil, errCantConvert
		}
		buyRate, err := getFloat64FromInterface(v[idxs[BuyRate]])
		if err != nil {
			return nil, err
		}
		sellRate, err := getFloat64FromInterface(v[idxs[SellRate]])
		if err != nil {
			return nil, err
		}
		buySanityRate, err := getFloat64FromInterface(v[idxs[BuySanityRate]])
		if err != nil {
			return nil, err
		}
		sellSanityRate, err := getFloat64FromInterface(v[idxs[SellSanityRate]])
		if !convertible {
			return nil, errCantConvert
		}
		rateEntry[pairName] = common.ReserveRateEntry{
			BuyReserveRate:  buyRate,
			SellReserveRate: sellRate,
			BuySanityRate:   buySanityRate,
			SellSanityRate:  sellSanityRate,
		}
		rate.Data = rateEntry
		rates[nRate] = rate
		firstRecordProcessed = true
	}
	return rates, nil
}
