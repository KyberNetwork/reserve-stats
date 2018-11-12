package storage

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const measurement = "user_volume_day"

//GetUserList return list of user info
func (is *InfluxStorage) GetUserList(fromTime, toTime uint64) ([]common.UserInfo, error) {
	var (
		err error
	)
	result := []common.UserInfo{}

	logger := is.sugar.With("from time", fromTime, "to time", toTime)

	q := fmt.Sprintf(`
		SELECT sum(eth_volume) as eth_amount, sum(usd_volume) as usd_amount from "%s"
		WHERE time >= %d%s AND TIME <= %d%s GROUP BY user_addr
	`, measurement, fromTime, timePrecision, toTime, timePrecision)

	res, err := is.queryDB(is.influxClient, q)
	if err != nil {
		return result, err
	}
	jsonValue, _ := json.Marshal(res)
	log.Printf("influx result: %s", jsonValue)

	if len(res[0].Series) == 0 {
		return result, nil
	}

	logger.Debug(res)

	for _, serie := range res[0].Series {
		userAddr := serie.Tags["user_addr"]
		for _, row := range serie.Values {
			ethAmount, usdAmount, err := is.rowToUserInfo(row)
			if err != nil {
				return result, err
			}
			result = append(result, common.UserInfo{
				Addr:      userAddr,
				ETHVolume: ethAmount,
				USDVolume: usdAmount,
			})
		}
	}

	return result, err
}
