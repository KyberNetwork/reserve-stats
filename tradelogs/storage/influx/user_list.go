package influx

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

//GetUserList return list of user info
func (is *Storage) GetUserList(fromTime, toTime time.Time) ([]common.UserInfo, error) {
	var (
		err             error
		result          []common.UserInfo
		measurementName = "trades"
	)

	logger := is.sugar.With("from time", fromTime, "to time", toTime, "func", "/tradelogs/storage.GetUserList")

	q := fmt.Sprintf(`
		SELECT sum(eth_volume) as eth_amount, sum(usd_volume) as usd_amount
		FROM (SELECT eth_amount as eth_volume, eth_amount*eth_usd_rate as usd_volume FROM "%s")
		WHERE time >= '%s' AND TIME <= '%s' GROUP BY user_addr
	`, measurementName, fromTime.UTC().Format(time.RFC3339), toTime.UTC().Format(time.RFC3339))

	logger.Debug(q)

	res, err := influxdb.QueryDB(is.influxClient, q, is.dbName)
	if err != nil {
		return result, err
	}

	if len(res[0].Series) == 0 {
		return result, nil
	}

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
