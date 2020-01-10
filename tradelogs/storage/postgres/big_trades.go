package postgres

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
)

const (
	getBigTradesQuery = `
SELECT bt.tradelog_id, a.timestamp AS timestamp, block_number, eth_amount, original_eth_amount, eth_usd_rate, d.address AS user_address,
e.address AS src_address, f.address AS dst_address,
src_amount, dst_amount, ip, country, integration_app, src_burn_amount, dst_burn_amount,
index, tx_hash, b.address AS src_rsv_address, c.address AS dst_rsv_address, src_wallet_fee_amount, dst_wallet_fee_amount,
g.address AS wallet_addr, tx_sender, receiver_address
FROM "` + schema.BigTradeLogsTableName + `" AS bt
INNER JOIN tradelogs as a ON a.id = bt.tradelog_id
INNER JOIN reserve AS b ON a.src_reserve_address_id = b.id
INNER JOIN reserve AS c ON a.dst_reserve_address_id = c.id
INNER JOIN users AS d ON a.user_address_id = d.id
INNER JOIN token AS e ON a.src_address_id = e.id
INNER JOIN token AS f ON a.dst_address_id = f.id
INNER JOIN wallet AS g ON a.wallet_address_id = g.id
WHERE bt.twitted is false;
`

	insertionBigTradelogsTemplate = `
INSERT INTO "` + schema.BigTradeLogsTableName + `" (tradelog_id) VALUES ($1) ON CONFLICT (tradelog_id) DO NOTHING;
`
	detectBigTradeQuery = `
	SELECT id from tradelogs WHERE eth_amount > $1 AND block_number >= $2;
`
	updateBigTradesQuery = `UPDATE "` + schema.BigTradeLogsTableName + `" SET twitted = true WHERE tradelog_id = $1;`
)

type bigTradeLogDBData struct {
	TradelogID uint64 `db:"tradelog_id"`
	tradeLogDBData
}

// GetNotTwittedTrades return big trades that is not twitted yet
func (tldb *TradeLogDB) GetNotTwittedTrades() ([]common.BigTradeLog, error) {
	var (
		logger      = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		queryResult = []bigTradeLogDBData{}
		result      = []common.BigTradeLog{}
	)
	err := tldb.db.Select(&queryResult, getBigTradesQuery)
	if err != nil {
		return nil, err
	}

	if len(queryResult) == 0 {
		logger.Debugw("empty result returned", "query", selectTradeLogsWithTxHashQuery)
		return result, nil
	}

	for _, r := range queryResult {
		tradeLog, err := tldb.tradeLogFromDBData(r.tradeLogDBData)
		if err != nil {
			logger.Errorw("cannot parse db data to trade log", "error", err)
			return nil, err
		}
		bigTradeLog := common.BigTradeLog{
			TradelogID: r.TradelogID,
			TradeLog:   tradeLog,
		}
		result = append(result, bigTradeLog)
	}
	return result, nil
}

// SaveBigTrades save trades into db
func (tldb *TradeLogDB) SaveBigTrades(bigVolume float32, fromBlock uint64) error {
	var (
		logger    = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		bigTrades = []uint64{}
	)
	if err := tldb.db.Select(&bigTrades, detectBigTradeQuery, bigVolume, fromBlock); err != nil {
		return fmt.Errorf("cannot get big trades: %s", err.Error())
	}
	logger.Infow("number of big trades", "number", len(bigTrades))
	for _, tradelogID := range bigTrades {
		if _, err := tldb.db.Exec(insertionBigTradelogsTemplate, tradelogID); err != nil {
			return fmt.Errorf("cannot insert big trades into db: %s", err.Error())
		}
	}
	return nil
}

// UpdateBigTradesTwitted update trades to twitted
func (tldb *TradeLogDB) UpdateBigTradesTwitted(trades []uint64) error {
	var (
		logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
	)
	logger.Infow("update big trade twitted", "len", len(trades))
	for _, tradelogID := range trades {
		if _, err := tldb.db.Exec(updateBigTradesQuery, tradelogID); err != nil {
			logger.Errorw("failed to update big trades twitted", "error", err)
			return err
		}
	}
	return nil
}
