package postgres

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// GetStats return tradelogs stats in a time range
func (tldb *TradeLogDB) GetStats(from, to time.Time) (common.StatsResponse, error) {
	var (
		logger = tldb.sugar.With(
			"from", from,
			"to", to,
			"func", caller.GetCurrentFunctionName(),
		)
		query = `
		SELECT 
		COALESCE(SUM(split.eth_amount), 0) AS eth_volume,
		COALESCE(SUM(split.eth_amount*eth_usd_rate), 0) AS usd_volume,
		COALESCE(SUM(platform_fee+burn+rebate+reward), 0) as collected_fee,
		COUNT(DISTINCT(tx_hash, tradelogs.index)) as total_trades,
		COUNT(CASE WHEN is_first_trade THEN 1 END) AS new_users,
		COUNT(distinct(user_address_id)) AS unique_addresses,
		COALESCE(AVG(split.eth_amount*eth_usd_rate), 0) as average_trade_size
		FROM tradelogs
		LEFT JOIN fee ON fee.trade_id = tradelogs.id
		LEFT JOIN split ON split.trade_id = tradelogs.id
	  WHERE timestamp >= $1 AND timestamp <= $2
	`
		statsRecord struct {
			ETHVolume        float64 `db:"eth_volume"`
			USDVolume        float64 `db:"usd_volume"`
			CollectedFee     float64 `db:"collected_fee"`
			TotalTrades      uint64  `db:"total_trades"`
			NewUsers         uint64  `db:"new_users"`
			UniqueAddresses  uint64  `db:"unique_addresses"`
			AverageTradeSize float64 `db:"average_trade_size"`
		}
	)
	logger.Infow("query to get tradelogs stats", "query", query)
	if err := tldb.db.Get(&statsRecord, query, from, to); err != nil {
		return common.StatsResponse{}, err
	}
	return common.StatsResponse{
		ETHVolume:        statsRecord.ETHVolume,
		USDVolume:        statsRecord.USDVolume,
		FeeCollected:     statsRecord.CollectedFee,
		TotalTrades:      statsRecord.TotalTrades,
		NewAdresses:      statsRecord.NewUsers,
		UniqueAddresses:  statsRecord.UniqueAddresses,
		AverageTradeSize: statsRecord.AverageTradeSize,
	}, nil
}

// GetTopTokens return top tokens by volume in a time range
func (tldb *TradeLogDB) GetTopTokens(from, to time.Time, limit uint64) (common.TopTokens, error) {
	var (
		logger = tldb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"from", from,
			"to", to,
			"limit", limit,
		)
		query = `
	  SELECT
		address as token_address,
		symbol as token_symbol,
		sum(usd_amount) as usd_amount
	  FROM
	  (
	  SELECT
	    sum(tradelogs.original_eth_amount*tradelogs.eth_usd_rate) usd_amount,
		token.address,
		token.symbol,
		token.id
	  FROM tradelogs
	    left join token on tradelogs.src_address_id = token.id
	  WHERE
		timestamp >= $1 AND timestamp <= $2
	  GROUP BY token.id
	  UNION ALL
	  SELECT
	    sum(tradelogs.original_eth_amount*tradelogs.eth_usd_rate) usd_amount,
		token.address,
		token.symbol,
		token.id
	  FROM tradelogs
	    left join token on tradelogs.dst_address_id = token.id
	  WHERE
		timestamp >= $1 AND timestamp <= $2
	  GROUP BY token.id
	  ) a GROUP BY a.address, a.symbol ORDER BY usd_amount DESC
		`
		topTokens []struct {
			TokenAddress string  `db:"token_address"`
			TokenSymbol  string  `db:"token_symbol"`
			USDAmount    float64 `db:"usd_amount"`
		}
	)
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	logger.Infow("query to get top tokens", "query", query)
	if err := tldb.db.Select(&topTokens, query, from, to); err != nil {
		return common.TopTokens{}, err
	}
	var result = make(common.TopTokens)
	for _, token := range topTokens {
		if token.TokenSymbol == "" {
			result[token.TokenAddress] = token.USDAmount
		} else {
			result[token.TokenSymbol] = token.USDAmount
		}
	}

	return result, nil
}

// GetTopIntegrations return top integrations by volume
func (tldb *TradeLogDB) GetTopIntegrations(from, to time.Time, limit uint64) (common.TopIntegrations, error) {
	var (
		logger = tldb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"from", from,
			"to", to,
			"limit", limit,
		)
		query = `
		SELECT
	  wallet.address,
	  wallet.name AS wallet_name,
	  sum(tradelogs.eth_amount*tradelogs.eth_usd_rate) usd_amount
	FROM tradelogs
	  left join wallet on tradelogs.wallet_address_id = wallet.id
	WHERE
		timestamp >= $1 AND timestamp <= $2
	GROUP BY wallet.address, wallet.name ORDER BY usd_amount DESC
		`
		topIntegrations []struct {
			Address   string  `db:"address"`
			Name      string  `db:"wallet_name"`
			USDAmount float64 `db:"usd_amount"`
		}
	)
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	logger.Infow("get top integrations", "query", query)
	if err := tldb.db.Select(&topIntegrations, query, from, to); err != nil {
		return common.TopIntegrations{}, err
	}

	result := make(common.TopIntegrations)
	for _, integration := range topIntegrations {
		if integration.Name != "" {
			result[integration.Name] = integration.USDAmount
		} else {
			result[integration.Address] = integration.USDAmount
		}
	}

	return result, nil
}

// GetTopReserves return top reserves by volume
func (tldb *TradeLogDB) GetTopReserves(from, to time.Time, limit uint64) (common.TopReserves, error) {
	var (
		logger = tldb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"from", from,
			"to", to,
			"limit", limit,
		)
		query = `
	  SELECT 
		  reserve.address as reserve_address, 
		  SUM(
			CASE 
		  		WHEN split.src = '0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE'
				THEN split.src_amount*tradelogs.eth_usd_rate
				ELSE split.dst_amount*tradelogs.eth_usd_rate
			END
		  ) AS usd_amount
	  FROM split
	  JOIN tradelogs on tradelogs.id = split.trade_id
	  JOIN reserve on split.reserve_id = reserve.id
	  WHERE tradelogs.timestamp >= $1 AND tradelogs.timestamp <= $2
	  GROUP BY reserve.address ORDER BY usd_amount DESC
		`
		topReserves []struct {
			ReserveAddress string  `db:"reserve_address"`
			USDAmount      float64 `db:"usd_amount"`
		}
	)
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	logger.Infow("get top reserves", "query", query)
	if err := tldb.db.Select(&topReserves, query, from, to); err != nil {
		return common.TopReserves{}, err
	}
	var result = make(common.TopReserves)
	for _, reserve := range topReserves {
		reserveName, err := tradelogs.ReserveAddressToName(ethereum.HexToAddress(reserve.ReserveAddress))
		if err == nil {
			result[reserveName] = reserve.USDAmount
		} else {
			logger.Warnw("reserve address does not have name", "error", err)
			result[reserve.ReserveAddress] = reserve.USDAmount
		}
	}
	return result, nil
}
