package postgres

import (
	"fmt"
	"math/big"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

const (
	getBigTradesQuery = `
SELECT bt.tradelog_id, 
a.timestamp AS timestamp, 
a.block_number, 
eth_amount, 
original_eth_amount, 
eth_usd_rate*original_eth_amount as fiat_amount, 
e.address AS src_address, 
f.address AS dst_address,
e.symbol AS src_symbol, 
f.symbol AS dst_symbol, 
index, 
tx_hash, 
tx_sender, 
g.name as wallet_name
FROM big_tradelogs AS bt
INNER JOIN tradelogs as a ON a.id = bt.tradelog_id
INNER JOIN token AS e ON a.src_address_id = e.id
INNER JOIN token AS f ON a.dst_address_id = f.id
INNER JOIN wallet AS g on g.id = a.wallet_address_id
WHERE bt.twitted is false AND a.timestamp >= $1 AND a.timestamp <= $2;
`

	insertionBigTradelogsTemplate = `
INSERT INTO big_tradelogs (tradelog_id) (
	SELECT tradelog_id.id FROM tradelogs AS tradelog_id 
	INNER JOIN token AS src_token ON src_token.id = tradelog_id.src_address_id
	INNER JOIN token AS dst_token ON dst_token.id = tradelog_id.dst_address_id
	WHERE original_eth_amount > $1 
	AND tradelog_id.block_number >= $2 
	AND src_token.symbol != 'WETH' AND dst_token.symbol != 'WETH'
	AND tradelogs.timestamp >= now() - interval '1' hour
)
ON CONFLICT (tradelog_id) DO NOTHING;
`
	updateBigTradesQuery = `UPDATE big_tradelogs SET twitted = true WHERE tradelog_id = $1 RETURNING tradelog_id;`
)

type bigTradeLogDBData struct {
	TradelogID        uint64    `db:"tradelog_id"`
	SrcSymbol         string    `db:"src_symbol"`
	DstSymbol         string    `db:"dst_symbol"`
	WalletName        string    `db:"wallet_name"`
	Timestamp         time.Time `db:"timestamp"`
	TransactionHash   string    `db:"tx_hash"`
	EthAmount         float64   `db:"eth_amount"`
	OriginalETHAmount float64   `db:"original_eth_amount"`
	FiatAmount        float64   `db:"fiat_amount"`
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
			ethAmountInWei, originalEthAmountInWei *big.Int
		)
		if ethAmountInWei, err = tldb.tokenAmountFormatter.ToWei(blockchain.ETHAddr, r.EthAmount); err != nil {
			logger.Debugw("failed to parse eth amount", "error", err)
			return nil, err
		}

		if originalEthAmountInWei, err = tldb.tokenAmountFormatter.ToWei(blockchain.ETHAddr, r.OriginalETHAmount); err != nil {
			logger.Debugw("failed to parse original eth amount", "error", err)
			return nil, err
		}

		bigTradeLog := common.BigTradeLog{
			TradelogID:        r.TradelogID,
			Timestamp:         r.Timestamp,
			TransactionHash:   ethereum.HexToHash(r.TransactionHash),
			EthAmount:         ethAmountInWei,
			OriginalETHAmount: originalEthAmountInWei,
			SrcSymbol:         r.SrcSymbol,
			DestSymbol:        r.DstSymbol,
			FiatAmount:        r.FiatAmount,
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
