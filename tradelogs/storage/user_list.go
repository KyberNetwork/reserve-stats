package storage

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const measurement = "trades"

//GetUserList return list of user info
func (is *InfluxStorage) GetUserList(fromTime, toTime uint64) ([]common.UserInfo, error) {
	var (
		err    error
		result []common.UserInfo
	)

	logger := is.sugar.With("from time", fromTime, "to time", toTime, "func", "/tradelogs/storage.GetUserList")

	from := timeutil.TimestampMsToTime(fromTime)
	to := timeutil.TimestampMsToTime(toTime)

	q := fmt.Sprintf(`
		SELECT sum(eth_volume) as eth_amount, sum(usd_volume) as usd_amount
		FROM (SELECT eth_amount as eth_volume, eth_amount*eth_usd_rate as usd_volume FROM "%s")
		WHERE time >= '%s' AND TIME <= '%s' GROUP BY user_addr
	`, measurement, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))

	logger.Debug(q)

	res, err := is.queryDB(is.influxClient, q)
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
