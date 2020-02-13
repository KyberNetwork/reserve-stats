package postgres

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
)

// GetTradeSummary returns trade summary group by day, eth_volume and usd volume are only accounted for not ETH-WETH trades
func (tldb *TradeLogDB) GetTradeSummary(from, to time.Time, timezone int8) (map[uint64]*common.TradeSummary, error) {
	var (
		err           error
		tradelogQuery string
		timeField     = schema.BuildDateTruncField("day", timezone)
	)
	logger := tldb.sugar.With("from", from, "to", to, "time_zone", timezone,
		"func", caller.GetCurrentFunctionName())
	from = schema.RoundTime(from, "day", timezone)
	to = schema.RoundTime(to, "day", timezone).Add(time.Hour * 24)
	results := make(map[uint64]*common.TradeSummary)

	tradelogQuery = `SELECT ` + timeField + ` AS time, 
		COUNT(DISTINCT(user_address_id)) AS unique_address,
		SUM(src_burn_amount) + SUM(dst_burn_amount) AS total_burn_fee,
		COUNT(CASE WHEN kyced THEN 1 END) AS kyced,
		COUNT(CASE WHEN is_first_trade THEN 1 END) AS count_new_trades
		FROM tradelogs
		WHERE timestamp >= $1 AND timestamp < $2
		GROUP BY time
	`
	logger.Debugw("prepare statement", "stmt", tradelogQuery)
	var countRecords []struct {
		Time           time.Time `db:"time"`
		UniqueAddress  uint64    `db:"unique_address"`
		TotalBurnFee   float64   `db:"total_burn_fee"`
		CountNewTrades uint64    `db:"count_new_trades"`
		Kyced          uint64    `db:"kyced"`
	}
	if err = tldb.db.Select(&countRecords, tradelogQuery, from, to); err != nil {
		return nil, err
	}

	if len(countRecords) == 0 {
		return nil, nil
	}
	for _, r := range countRecords {
		ts := timeutil.TimeToTimestampMs(r.Time)
		results[ts] = &common.TradeSummary{
			UniqueAddresses:    r.UniqueAddress,
			TotalBurnFee:       r.TotalBurnFee,
			NewUniqueAddresses: r.CountNewTrades,
			KYCEDAddresses:     r.Kyced,
		}
	}

	tradelogQuery = fmt.Sprintf(
		`SELECT %[1]s AS time, 
		SUM(eth_amount) total_eth_volume, 
		AVG(eth_amount) eth_per_trade,
		SUM(eth_amount*eth_usd_rate) as total_usd_volume, 
		AVG(eth_amount*eth_usd_rate) usd_per_trade, count(1) as total_trade 
	FROM "%[2]s"
	AND timestamp >= $1 AND timestamp < $2
	GROUP BY time
	`, timeField, schema.TradeLogsTableName)
	logger.Debugw("prepare statement", "stmt", tradelogQuery)
	var records []struct {
		Time           time.Time `db:"time"`
		TotalEthVolume float64   `db:"total_eth_volume"`
		EthPerTrade    float64   `db:"eth_per_trade"`
		TotalUsdVolume float64   `db:"total_usd_volume"`
		UsdPerTrade    float64   `db:"usd_per_trade"`
		TotalTrade     uint64    `db:"total_trade"`
	}
	err = tldb.db.Select(&records, tradelogQuery, from, to)
	if err != nil {
		return nil, err
	}
	for _, r := range records {
		ts := timeutil.TimeToTimestampMs(r.Time)
		summary, ok := results[ts]
		if !ok {
			logger.Warn("key not found", "ts", ts)
			continue
		}
		summary.USDAmount = r.TotalUsdVolume
		summary.ETHVolume = r.TotalEthVolume
		summary.ETHPerTrade = r.EthPerTrade
		summary.USDPerTrade = r.UsdPerTrade
		summary.TotalTrade = r.TotalTrade
	}
	return results, nil
}
