package postgres

import (
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
	 SELECT SUM(eth_amount) as eth_volume,
	  SUM(eth_amount*eth_usd_rate) as usd_volume,
	  SUM(src_burn_amount+dst_burn_amount+src_wallet_fee_amount+dst_wallet_fee_amount) as collected_fee,
	  COUNT(*) as total_trades,
	  COUNT(CASE WHEN is_first_trade THEN 1 END) AS new_users,
	  COUNT(distinct(user_address_id)) AS unique_addresses,
	  AVG(eth_amount*eth_usd_rate) as average_trade_size
	  from tradelogs
	  WHERE timestamp >= $1 and timestamp <= $2;
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
func (tldb *TradeLogDB) GetTopTokens(from, to time.Time) (common.TopTokens, error) {
	var (
		logger = tldb.sugar.With(
			"from", from,
			"to", to,
			"func", caller.GetCurrentFunctionName(),
		)
		query = `
	  SELECT
	    address as token_address,
		sum(usd_amount) as usd_amount
	  FROM
	  (
	  SELECT
	    sum(tradelogs.eth_amount*tradelogs.eth_usd_rate) usd_amount,
	    token.address,
		token.id
	  FROM tradelogs
	    left join token on tradelogs.src_address_id = token.id
	  WHERE
		timestamp >= $1 AND timestamp <= $2
	    AND (src_burn_amount + dst_amount) > 0
	  GROUP BY token.id
	  UNION ALL
	  SELECT
	    sum(tradelogs.eth_amount*tradelogs.eth_usd_rate) usd_amount,
	    token.address,
		token.id
	  FROM tradelogs
	    left join token on tradelogs.dst_address_id = token.id
	  WHERE
		timestamp >= $1 AND timestamp <= $2
	    AND (src_burn_amount + dst_amount) > 0
	  GROUP BY token.id
	  ) a GROUP BY a.address;
		`
		topTokens []struct {
			TokenAddress string  `db:"token_address"`
			USDAmount    float64 `db:"usd_amount"`
		}
	)
	logger.Infow("query to get top tokens", "query", query)

	if err := tldb.db.Select(&topTokens, query, from, to); err != nil {
		return common.TopTokens{}, err
	}
	var result = make(common.TopTokens)
	for _, token := range topTokens {
		result[token.TokenAddress] = token.USDAmount
	}

	return result, nil
}

// GetTopIntegrations return top integrations by volume
func (tldb *TradeLogDB) GetTopIntegrations(from, to time.Time) (common.TopIntegrations, error) {
	var (
		logger = tldb.sugar.With(
			"from", from,
			"to", to,
			"func", caller.GetCurrentFunctionName(),
		)
		query = `
		SELECT
	  wallet.address,
	  wallet.name AS wallet_name,
	  sum(tradelogs.eth_amount*tradelogs.eth_usd_rate) usd_amount,
	FROM tradelogs
	  left join wallet on tradelogs.wallet_address_id = wallet.id
	WHERE
		timestamp >= $1 AND timestamp <= $2
	  AND (src_burn_amount + dst_amount) > 0
	GROUP BY wallet.address, wallet.name;
		`
	)
	logger.Infow("get top integrations", "query", query)
	return common.TopIntegrations{}, nil
}

// GetTopReserves return top reserves by volume
func (tldb *TradeLogDB) GetTopReserves(from, to time.Time) (common.TopReserves, error) {
	var (
		logger = tldb.sugar.With(
			"from", from,
			"to", to,
			"func", caller.GetCurrentFunctionName(),
		)
		query = `
	  SELECT
	    address as reserve_address,
		sum(usd_amount) as usd_amount 
	  FROM
	  (
	  SELECT
	    sum(tradelogs.eth_amount*tradelogs.eth_usd_rate) usd_amount,
	    reserve.address,
		reserve.id
	  FROM tradelogs
	    left join reserve on tradelogs.src_reserve_address_id = reserve.id
	  WHERE
		timestamp >= $1 AND timestamp <= $2
	    AND (src_burn_amount + dst_amount) > 0
	  GROUP BY reserve.id
	  UNION ALL
	  SELECT
	    sum(tradelogs.eth_amount*tradelogs.eth_usd_rate) usd_amount,
	    reserve.address,
		reserve.id
	  FROM tradelogs
	    left join reserve on tradelogs.dst_reserve_address_id = reserve.id
	  WHERE
		timestamp >= $1 AND timestamp <= $2
	    AND (src_burn_amount + dst_amount) > 0
	  GROUP BY reserve.id
	  ) a GROUP BY a.address;
		`
		topReserves []struct {
			ReserveAddress string  `db:"reserve_address"`
			USDAmount      float64 `db:"usd_amount"`
		}
	)
	logger.Infow("get top reserves", "query", query)
	if err := tldb.db.Select(&topReserves, query, from, to); err != nil {
		return common.TopReserves{}, err
	}
	var result = make(common.TopReserves)
	for _, reserve := range topReserves {
		reserveName := tradelogs.ReserveAddressToName(ethereum.HexToAddress(reserve.ReserveAddress))
		if reserveName != "" {
			result[reserveName] = reserve.USDAmount
		} else {
			result[reserve.ReserveAddress] = reserve.USDAmount
		}
	}
	return common.TopReserves{}, nil
}
