package postgrestorage

import (
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"
	"time"
)

//TODO implement this
func (tldb *TradeLogDB) GetTradeSummary(from, to time.Time, timezone int8) (map[uint64]*common.TradeSummary, error) {
	var (
		err           error
		ethCondition  string
		tradelogQuery string
		timeField     = schema.BuildDateTruncField("day", timezone)
	)

	from = schema.RoundTime(from, "day", timezone)
	to = schema.RoundTime(to, "day", timezone).Add(time.Hour * 24)

	logger := tldb.sugar.With("from", from, "to", to, "time_zone", timezone,
		"func", "tradelogs/storage/postgrestorage/TradeLogDB.GetTradeSummary")
	results := make(map[uint64]*common.TradeSummary)
	if ethCondition, err = schema.BuildEthWethExcludingCondition(); err != nil {
		return nil, err
	}

	tradelogQuery = `SELECT ` + timeField + ` AS time, 
		COUNT(DISTINCT(user_address_id)) AS unique_address,
		SUM(src_burn_amount) + SUM(dst_burn_amount) AS total_burn_fee
		FROM tradelogs
		WHERE timestamp >= $1 AND timestamp < $2
		GROUP BY time
	`
	logger.Debugw("prepare statement", "stmt", tradelogQuery)
	var datas2 []struct {
		Time          time.Time `db:"time"`
		UniqueAddress uint64    `db:"unique_address"`
		TotalBurnFee  float64   `db:"total_burn_fee"`
	}
	if err = tldb.db.Select(&datas2, tradelogQuery, from.UTC().Format(schema.DefaultDateFormat),
		to.UTC().Format(schema.DefaultDateFormat)); err != nil {
		return nil, err
	}

	if len(datas2) == 0 {
		return nil, nil
	}

	for _, data := range datas2 {
		ts := timeutil.TimeToTimestampMs(data.Time)
		results[ts] = &common.TradeSummary{
			UniqueAddresses: data.UniqueAddress,
			TotalBurnFee:    data.TotalBurnFee,
		}
	}

	tradelogQuery =
		`SELECT ` + timeField + ` AS time, 
		SUM(eth_amount) total_eth_volume, 
		AVG(eth_amount) eth_per_trade,
		SUM(eth_amount*eth_usd_rate) as total_usd_volume, 
		AVG(eth_amount*eth_usd_rate) usd_per_trade, count(1) as total_trade 
	FROM ` + schema.TradeLogsTableName + `
	WHERE ` + ethCondition + `
	AND timestamp >= $1 AND timestamp < $2
	GROUP BY time
	`
	logger.Debugw("prepare statement", "stmt", tradelogQuery)
	var datas []struct {
		Time           time.Time `db:"time"`
		TotalEthVolume float64   `db:"total_eth_volume"`
		EthPerTrade    float64   `db:"eth_per_trade"`
		TotalUsdVolume float64   `db:"total_usd_volume"`
		UsdPerTrade    float64   `db:"usd_per_trade"`
		TotalTrade     uint64    `db:"total_trade"`
	}

	err = tldb.db.Select(&datas, tradelogQuery, from.UTC().Format(schema.DefaultDateFormat),
		to.UTC().Format(schema.DefaultDateFormat))
	if err != nil {
		return nil, err
	}

	for _, data := range datas {
		ts := timeutil.TimeToTimestampMs(data.Time)
		summary, ok := results[ts]
		if !ok {
			logger.Warn("key not found", "ts", ts)
			continue
		}
		summary.USDAmount = data.TotalUsdVolume
		summary.ETHVolume = data.TotalEthVolume
		summary.ETHPerTrade = data.EthPerTrade
		summary.USDPerTrade = data.UsdPerTrade
		summary.TotalTrade = data.TotalTrade
	}

	return results, nil
}
