package postgres

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
)

// GetCountryStats return country stats aggregated data
func (tldb *TradeLogDB) GetCountryStats(countryCode string, from, to time.Time, timezone int8) (map[uint64]*common.CountryStats, error) {
	var (
		err            error
		tradelogsQuery string
		timeField      = schema.BuildDateTruncField("day", timezone)
		results        = make(map[uint64]*common.CountryStats)
	)
	logger := tldb.sugar.With("from", from, "to", to,
		"func", caller.GetCurrentFunctionName())
	from = schema.RoundTime(from, "day", timezone)
	to = schema.RoundTime(to, "day", timezone).Add(time.Hour * 24)
	// get unique_address, kyced, total_burn_fee, count_new_trades
	tradelogsQuery = fmt.Sprintf(`SELECT %[1]s AS time, 
		COUNT(DISTINCT(user_address_id)) AS unique_address,
		-- SUM(src_burn_amount) + SUM(dst_burn_amount) AS total_burn_fee,
		COUNT(CASE WHEN kyced THEN 1 END) AS kyced,
		COUNT(CASE WHEN is_first_trade THEN 1 END) AS count_new_trades
		FROM tradelogs
		WHERE timestamp >= $1 AND timestamp < $2 AND country = $3
		GROUP BY time
	`, timeField)
	logger.Debugw("prepare statement", "stmt", tradelogsQuery)
	var records []struct {
		Time          time.Time `db:"time"`
		UniqueAddress uint64    `db:"unique_address"`
		// TotalBurnFee   float64   `db:"total_burn_fee"`
		CountNewTrades uint64 `db:"count_new_trades"`
		Kyced          uint64 `db:"kyced"`
	}
	if err = tldb.db.Select(&records, tradelogsQuery, from, to, countryCode); err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	}

	for _, r := range records {
		ts := timeutil.TimeToTimestampMs(r.Time)
		results[ts] = &common.CountryStats{
			UniqueAddresses: r.UniqueAddress,
			// TotalBurnFee:       r.TotalBurnFee,
			KYCEDAddresses:     r.Kyced,
			NewUniqueAddresses: r.CountNewTrades,
		}
	}
	// get eth_volume, eth_per_trade, usd_volume, usd_per_trade for not eth-weth trades
	tradelogsQuery = fmt.Sprintf(
		`SELECT %[1]s AS time, 
		SUM(eth_amount) total_eth_volume, 
		AVG(eth_amount) eth_per_trade,
		SUM(eth_amount*eth_usd_rate) as total_usd_volume, 
		AVG(eth_amount*eth_usd_rate) usd_per_trade, count(1) as total_trade 
	FROM tradelogs
	WHERE timestamp >= $1 AND timestamp < $2 AND country = $3
	GROUP BY time
	`, timeField)
	logger.Debugw("prepare statement", "stmt", tradelogsQuery)
	var volumeRecords []struct {
		Time           time.Time `db:"time"`
		TotalEthVolume float64   `db:"total_eth_volume"`
		EthPerTrade    float64   `db:"eth_per_trade"`
		TotalUsdVolume float64   `db:"total_usd_volume"`
		UsdPerTrade    float64   `db:"usd_per_trade"`
		TotalTrade     uint64    `db:"total_trade"`
	}
	err = tldb.db.Select(&volumeRecords, tradelogsQuery, from, to, countryCode)
	if err != nil {
		return nil, err
	}
	for _, r := range volumeRecords {
		ts := timeutil.TimeToTimestampMs(r.Time)
		summary, ok := results[ts]
		if !ok {
			logger.Warn("key not found", "ts", ts)
			continue
		}
		summary.TotalUSDVolume = r.TotalUsdVolume
		summary.TotalETHVolume = r.TotalEthVolume
		summary.ETHPerTrade = r.EthPerTrade
		summary.USDPerTrade = r.UsdPerTrade
		summary.TotalTrade = r.TotalTrade
	}
	return results, nil
}
