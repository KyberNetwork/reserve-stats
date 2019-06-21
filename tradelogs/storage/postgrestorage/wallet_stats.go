package postgrestorage

import (
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"
	"time"
)

func (tldb *TradeLogDB) GetWalletStats(from, to time.Time, walletAddr string, timezone int8) (map[uint64]common.WalletStats, error) {
	var (
		err          error
		ethCondition string
	)
	logger := tldb.sugar.With("from", from, "to", to, "walletAddr", walletAddr,
		"timezone", timezone, "func", "tradelogs/storage/postgrestorage/TradeLogDB.GetWalletStats")

	from = schema.RoundTime(from, "day", timezone)
	to = schema.RoundTime(to, "day", timezone).Add(time.Hour * 24)
	timeField := schema.BuildDateTruncField("day", timezone)

	if ethCondition, err = schema.BuildEthWethExcludingCondition(); err != nil {
		return nil, err
	}

	walletStatsQuery := `
		SELECT ` + timeField + ` as time,
			COUNT(DISTINCT(user_address_id)) AS unique_address,
			SUM(src_burn_amount) + SUM(dst_burn_amount) AS total_burn_fee
		FROM "` + schema.TradeLogsTableName + `" 
		WHERE timestamp >= $1 AND timestamp < $2
		AND EXISTS (SELECT NULL FROM "` + schema.WalletTableName + `" WHERE address = $3 AND id=wallet_address_id)
		GROUP BY time
	`

	logger.Debugw("prepare statement", "stmt", walletStatsQuery)
	var records []struct {
		TimeStamp     time.Time `db:"time"`
		UniqueAddress int64     `db:"unique_address"`
		TotalBurnFee  float64   `db:"total_burn_fee"`
	}
	err = tldb.db.Select(&records, walletStatsQuery, from.UTC().Format(schema.DefaultDateFormat),
		to.UTC().Format(schema.DefaultDateFormat), walletAddr)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	}

	results := make(map[uint64]common.WalletStats)
	for _, record := range records {
		ts := timeutil.TimeToTimestampMs(record.TimeStamp)
		results[ts] = common.WalletStats{
			UniqueAddresses: record.UniqueAddress,
			BurnFee:         record.TotalBurnFee,
		}
	}

	walletStatsQuery = `
		SELECT ` + timeField + ` AS time, 
		SUM(eth_amount) total_eth_volume, 
		AVG(eth_amount) eth_per_trade,
		SUM(eth_amount*eth_usd_rate) as total_usd_volume, 
		AVG(eth_amount*eth_usd_rate) usd_per_trade, count(1) as total_trade
		FROM "` + schema.TradeLogsTableName + `" 
		WHERE timestamp >= $1 AND timestamp < $2
		AND EXISTS (SELECT NULL FROM "` + schema.WalletTableName + `" WHERE address = $3 AND id=wallet_address_id)
		AND ` + ethCondition + `
		GROUP BY time
	`
	logger.Debugw("prepare statement", "stmt", walletStatsQuery)
	var records2 []struct {
		Time           time.Time `db:"time"`
		TotalEthVolume float64   `db:"total_eth_volume"`
		EthPerTrade    float64   `db:"eth_per_trade"`
		TotalUsdVolume float64   `db:"total_usd_volume"`
		UsdPerTrade    float64   `db:"usd_per_trade"`
		TotalTrade     int64     `db:"total_trade"`
	}

	err = tldb.db.Select(&records2, walletStatsQuery, from.UTC().Format(schema.DefaultDateFormat),
		to.UTC().Format(schema.DefaultDateFormat), walletAddr)
	if err != nil {
		return nil, err
	}

	if len(records2) == 0 {
		logger.Debugw("no exclude eth weth trades")
	}

	for _, record := range records2 {
		ts := timeutil.TimeToTimestampMs(record.Time)
		walletStats, ok := results[ts]
		if !ok {
			logger.Warn("key not found", "ts", ts)
			continue
		}
		walletStats.USDVolume = record.TotalUsdVolume
		walletStats.ETHVolume = record.TotalEthVolume
		walletStats.ETHPerTrade = record.EthPerTrade
		walletStats.USDPerTrade = record.UsdPerTrade
		walletStats.TotalTrade = record.TotalTrade
		results[ts] = walletStats
	}
	return results, nil
}
