package postgres

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
)

// GetAggregatedBurnFee Get aggregated Burn fee by hour or day
func (tldb *TradeLogDB) GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddrs []ethereum.Address) (map[ethereum.Address]map[string]float64, error) {
	var (
		timeField string
		err       error
	)
	logger := tldb.sugar.With("from", from, "to", to,
		"func", caller.GetCurrentFunctionName())

	switch strings.ToLower(freq) {
	case "h":
		timeField = schema.BuildDateTruncField("hour", 0)
		from = schema.RoundTime(from, "hour", 0)
		to = schema.RoundTime(to, "hour", 0).Add(time.Hour)
	case "d":
		timeField = schema.BuildDateTruncField("day", 0)
		from = schema.RoundTime(from, "day", 0)
		to = schema.RoundTime(to, "day", 0).Add(time.Hour * 24)

	default:
		return nil, fmt.Errorf("frequency not supported: %v", freq)
	}
	args := []interface{}{
		from,
		to,
	}

	hexAddrs := make([]string, 0)
	for _, rsvAddr := range reserveAddrs {
		hexAddrs = append(hexAddrs, rsvAddr.Hex())
	}

	addrCondition := ""
	if len(hexAddrs) != 0 {
		addrCondition = " AND c.address = ANY($3)"
		args = append(args, pq.Array(hexAddrs))
	}

	integrationQuery := fmt.Sprintf(`
		SELECT time, address , SUM(amount) as amount
		FROM (
			SELECT %[1]s as time, burn AS amount, reserve_address AS address
			FROM "fee" b
			JOIN tradelogs on tradelogs.id = fee.trade_id
			WHERE tradelogs.timestamp >= $1 AND tradelogs.timestamp < $2 %[2]s
		) a GROUP BY time,address
	`, timeField, addrCondition)

	var records []struct {
		Amount  float64   `db:"amount"`
		Address string    `db:"address"`
		Time    time.Time `db:"time"`
	}

	logger.Debugw("prepare statement", "stmt", integrationQuery)
	err = tldb.db.Select(&records, integrationQuery, args...)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		logger.Debugw("no trade found")
		return nil, nil
	}

	result := make(map[ethereum.Address]map[string]float64)
	for _, burnFee := range records {
		reserve := ethereum.HexToAddress(burnFee.Address)
		key := strconv.FormatUint(timeutil.TimeToTimestampMs(burnFee.Time), 10)
		_, ok := result[reserve]
		if !ok {
			result[reserve] = make(map[string]float64)
		}
		result[reserve][key] += burnFee.Amount
	}
	return result, nil

}
