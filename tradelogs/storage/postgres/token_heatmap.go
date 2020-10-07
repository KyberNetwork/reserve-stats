package postgres

import (
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
)

// GetTokenHeatmap returns map of country to eth_volume, token_volume, usd_volume filter by timestamp and asset
func (tldb *TradeLogDB) GetTokenHeatmap(asset ethereum.Address, from, to time.Time, timezone int8) (map[string]common.Heatmap, error) {
	var (
		err               error
		tokenHeatMapQuery string
	)

	logger := tldb.sugar.With("from", from, "to", to, "timezone", timezone,
		"func", caller.GetCurrentFunctionName())

	from = schema.RoundTime(from, "day", timezone)
	to = schema.RoundTime(to, "day", timezone).Add(time.Hour * 24)
	// nested query with filter by src_address_id and dst_address_id
	tokenHeatMapQuery = `
		SELECT country, 
			SUM(eth_amount) AS eth_volume,
			SUM(token_volume) AS token_volume,
			SUM(usd_volume) AS usd_volume 
		FROM (
			SELECT country, src_amount AS token_volume, 
				eth_amount, eth_amount * eth_usd_rate AS usd_volume
			FROM "tradelogs"
			WHERE timestamp >= $1 and timestamp < $2
			AND EXISTS (SELECT NULL FROM "token" WHERE address = $3 and id = src_address_id)
			AND country IS NOT NULL
		UNION ALL
			SELECT country, dst_amount AS token_volume, eth_amount, 
				eth_amount*eth_usd_rate AS usd_volume
			FROM "tradelogs"
			WHERE timestamp >= $1 and timestamp < $2
			AND EXISTS (SELECT NULL FROM "token" WHERE address = $3 and id = dst_address_id)
			AND country IS NOT NULL
		)a GROUP BY country
	`
	logger.Debugw("prepare statement", "stmt", tokenHeatMapQuery)

	var records []struct {
		Country     string  `db:"country"`
		TokenVolume float64 `db:"token_volume"`
		EthVolume   float64 `db:"eth_volume"`
		UsdVolume   float64 `db:"usd_volume"`
	}
	err = tldb.db.Select(&records, tokenHeatMapQuery, from, to, asset.Hex())
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		logger.Debugw("no trade found")
		return nil, nil
	}
	results := make(map[string]common.Heatmap)
	for _, record := range records {
		results[record.Country] = common.Heatmap{
			TotalETHValue:   record.EthVolume,
			TotalTokenValue: record.TokenVolume,
			TotalFiatValue:  record.UsdVolume,
		}
	}
	return results, nil
}
