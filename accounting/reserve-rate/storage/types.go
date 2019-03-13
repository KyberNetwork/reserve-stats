package storage

import (
	"encoding/json"
	"time"
)

//AccountingReserveRates is a struct to hold the rate as defined.
//it holds infomation as followed example
//	{
//     "2019-03-13": {
//       "ETH": {
//         "KNC": 917.431192,
//         "ZIL": 5205.351102
//       }
//     }
//   }
type AccountingReserveRates map[time.Time]map[string]map[string]float64

//MarshalJSON implement custom JSON marshaller for AccountingReserveRate to short form date format
func (acrr AccountingReserveRates) MarshalJSON() ([]byte, error) {
	var mapResult = make(map[string]map[string]map[string]float64)
	const shortForm = "2006-01-02"
	for k, v := range acrr {
		mapResult[k.Format(shortForm)] = v
	}
	return json.Marshal(mapResult)
}

//AccountingRatesReply is the wrapper object for ReserveRate data
type AccountingRatesReply struct {
	ReserveRates map[string]AccountingReserveRates `json:"reserves-rates"`
	EthUsdRates  AccountingReserveRates            `json:"eth-usd-rates"`
}
