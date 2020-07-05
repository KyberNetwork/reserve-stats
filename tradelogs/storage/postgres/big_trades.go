package postgres

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
)

const (
	getBigTradesQuery = `
SELECT bt.tradelog_id, a.timestamp AS timestamp, a.block_number, eth_amount, original_eth_amount, eth_usd_rate, d.address AS user_address,
e.address AS src_address, f.address AS dst_address,
e.symbol AS src_symbol, f.symbol AS dst_symbol, src_amount, dst_amount,
index, tx_hash, tx_sender, receiver_address,
g.name as wallet_name
FROM "` + schema.BigTradeLogsTableName + `" AS bt
INNER JOIN tradelogs as a ON a.id = bt.tradelog_id
INNER JOIN users AS d ON a.user_address_id = d.id
INNER JOIN token AS e ON a.src_address_id = e.id
INNER JOIN token AS f ON a.dst_address_id = f.id
WHERE bt.twitted is false AND a.timestamp >= $1 AND a.timestamp <= $2;
`

	insertionBigTradelogsTemplate = `
INSERT INTO big_tradelogs (tradelog_id) (
	SELECT tradelog_id.id FROM tradelogs AS tradelog_id 
	INNER JOIN token AS src_token ON src_token.id = tradelog_id.src_address_id
	INNER JOIN token AS dst_token ON dst_token.id = tradelog_id.dst_address_id
	WHERE original_eth_amount > $1 AND tradelog_id.block_number >= $2 AND src_token.symbol != 'WETH' and dst_token.symbol != 'WETH'
)
ON CONFLICT (tradelog_id) DO NOTHING;
`
	updateBigTradesQuery = `UPDATE "` + schema.BigTradeLogsTableName + `" SET twitted = true WHERE tradelog_id = $1 RETURNING tradelog_id;`
)

type bigTradeLogDBData struct {
	TradelogID uint64 `db:"tradelog_id"`
	SrcSymbol  string `db:"src_symbol"`
	DstSymbol  string `db:"dst_symbol"`
	WalletName string `db:"wallet_name"`
	tradeLogDBData
}

// GetNotTwittedTrades return big trades that is not twitted yet
func (tldb *TradeLogDB) GetNotTwittedTrades(from, to time.Time) ([]common.BigTradeLog, error) {
	var (
		logger      = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		queryResult = []bigTradeLogDBData{}
		result      = []common.BigTradeLog{}
	)
	err := tldb.db.Select(&queryResult, getBigTradesQuery, from, to)
	if err != nil {
		return nil, err
	}

	if len(queryResult) == 0 {
		logger.Debugw("empty result returned", "query", selectTradeLogsWithTxHashQuery)
		return result, nil
	}

	for _, r := range queryResult {
		var (
			feeResult []feeRecord
		)
		tradeLog, err := tldb.tradeLogFromDBData(r.tradeLogDBData, feeResult)
		if err != nil {
			logger.Errorw("cannot parse db data to trade log", "error", err)
			return nil, err
		}
		bigTradeLog := common.BigTradeLog{
			TradelogID:        r.TradelogID,
			Timestamp:         tradeLog.Timestamp,
			TransactionHash:   tradeLog.TransactionHash,
			EthAmount:         tradeLog.EthAmount,
			OriginalETHAmount: tradeLog.OriginalEthAmount,
			SrcSymbol:         r.SrcSymbol,
			DestSymbol:        r.DstSymbol,
			FiatAmount:        tradeLog.FiatAmount,
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
	logger.Infow("query save big trades", "query", insertionBigTradelogsTemplate)
	if _, err := tldb.db.Exec(insertionBigTradelogsTemplate, bigVolume, fromBlock); err != nil {
		return fmt.Errorf("cannot update big trades: %s", err.Error())
	}
	logger.Infow("number of big trades", "number", len(bigTrades))
	return nil
}

// UpdateBigTradesTwitted update trades to twitted
func (tldb *TradeLogDB) UpdateBigTradesTwitted(trades []uint64) error {
	var (
		logger     = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		tradeLogID uint64
	)
	logger.Infow("update big trade twitted", "len", len(trades))
	for _, tradelogID := range trades {
		if err := tldb.db.Get(&tradeLogID, updateBigTradesQuery, tradelogID); err != nil {
			logger.Errorw("failed to update big trades twitted", "error", err)
			return err
		}
		logger.Infow("tradelog updated", "id", tradeLogID)
	}
	return nil
}
