package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

//AccountingReserveRate is a struct to hold the rate as defined.
//it holds infomation as followed example
// {
//     "2018-10-15"{
//     "rates": {
//         "ETH": {
//             "KNC": 917.431192,
//             "ZIL": 5205.351102
//         },
//         "USD": {
//             "ETH": 0.009434
//         }
//     }
// }
//}
type AccountingReserveRates map[time.Time]map[string]map[string]float64

//MarshalJSON implement custom JSON marshaller for AccountingReserveRate to short form date format
func (acrr AccountingReserveRates) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	const shortForm = "2006-01-02"

	length := len(acrr)
	count := 0
	for k, v := range acrr {
		jsonVal, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("\"%s\":%s", k.Format(shortForm), string(jsonVal)))
		count++
		if count < length {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
