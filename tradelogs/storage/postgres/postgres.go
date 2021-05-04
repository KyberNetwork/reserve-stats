package postgres

import (
	"database/sql"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
)

// TradeLogDB is storage of tradelog data
type TradeLogDB struct {
	sugar                *zap.SugaredLogger
	db                   *sqlx.DB
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface

	// used for calculate burn amount
	// as different environment have different knc address
	kncAddr ethereum.Address
}

//NewTradeLogDB create a new instance of TradeLogDB
func NewTradeLogDB(sugar *zap.SugaredLogger, db *sqlx.DB, tokenAmountFormatter blockchain.TokenAmountFormatterInterface, kncAddr ethereum.Address) (*TradeLogDB, error) {
	var logger = sugar.With("func", caller.GetCurrentFunctionName())
	var err error
	logger.Debug("initializing database schema")
	if _, err = db.Exec(schema.TradeLogsSchema); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")

	return &TradeLogDB{
		sugar:                sugar,
		db:                   db,
		tokenAmountFormatter: tokenAmountFormatter,
		kncAddr:              kncAddr,
	}, err
}

// LastBlock returns last stored trade log block number from database.
func (tldb *TradeLogDB) LastBlock() (int64, error) {
	var (
		logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		result sql.NullInt64
	)
	stmt := `SELECT MAX("block_number") FROM "tradelogs"`
	logger = logger.With("query", stmt)
	logger.Debug("Start query")
	err := tldb.db.Get(&result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		logger.Errorw("Get error ", "error", err)
		return 0, err
	}
	return result.Int64, nil
}

type tradeLogDBData struct {
	ID                  uint64         `db:"id"`
	Timestamp           time.Time      `db:"timestamp"`
	BlockNumber         uint64         `db:"block_number"`
	QuoteAmount         float64        `db:"usdt_amount"`
	OriginalQuoteAmount float64        `db:"original_usdt_amount"`
	UserAddress         pq.StringArray `db:"user_address"`
	SrcAddress          pq.StringArray `db:"src_address"`
	DstAddress          pq.StringArray `db:"dst_address"`
	ReserveAddress      string         `db:"reserve_address"`
	SrcAmount           float64        `db:"src_amount"`
	DstAmount           float64        `db:"dst_amount"`
	LogIndex            uint           `db:"index"`
	TxHash              string         `db:"tx_hash"`
	TxSender            string         `db:"tx_sender"`
	ReceiverAddr        string         `db:"receiver_address"`
	GasUsed             uint64         `db:"gas_used"`
	GasPrice            float64        `db:"gas_price"`
	TransactionFee      float64        `db:"transaction_fee"`
	Version             uint           `db:"version"`
}

func (tldb *TradeLogDB) tradeLogFromDBData(r tradeLogDBData) (common.Tradelog, error) {
	var (
		tradeLog common.Tradelog
		err      error

		usdtAmountInWei                    *big.Int
		srcAmountInWei                     *big.Int
		dstAmountInWei                     *big.Int
		originalUSDTAmountInWei            *big.Int
		gasPriceInWei, transactionFeeInWei *big.Int

		logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
	)

	if usdtAmountInWei, err = tldb.tokenAmountFormatter.ToWei(blockchain.USDTAddr, r.QuoteAmount); err != nil {
		logger.Debugw("failed to parse usdt amount", "error", err)
		return tradeLog, err
	}

	if originalUSDTAmountInWei, err = tldb.tokenAmountFormatter.ToWei(blockchain.USDTAddr, r.OriginalQuoteAmount); err != nil {
		logger.Debugw("failed to parse original usdt amount", "error", err)
		return tradeLog, err
	}
	SrcAddress := ethereum.HexToAddress(r.SrcAddress[0])
	if srcAmountInWei, err = tldb.tokenAmountFormatter.ToWei(SrcAddress, r.SrcAmount); err != nil {
		logger.Debugw("failed to parse src amount", "error", err)
		return tradeLog, err
	}
	DstAddress := ethereum.HexToAddress(r.DstAddress[0])
	if dstAmountInWei, err = tldb.tokenAmountFormatter.ToWei(DstAddress, r.DstAmount); err != nil {
		logger.Debugw("failed to parse dst amount", "error", err)
		return tradeLog, err
	}

	// these conversion below is from Gwei to wei which is used ^18 method, so I used ToWei function with ETHAddr - which have decimals of 18
	if gasPriceInWei, err = tldb.tokenAmountFormatter.ToWei(blockchain.USDTAddr, r.GasPrice); err != nil {
		logger.Debugw("failed to parse gas price", "error", err)
		return tradeLog, err
	}

	if transactionFeeInWei, err = tldb.tokenAmountFormatter.ToWei(blockchain.USDTAddr, r.TransactionFee); err != nil {
		logger.Debugw("failed to parse transaction fee", "error", err)
		return tradeLog, err
	}
	tradeLog = common.Tradelog{
		TransactionHash:     ethereum.HexToHash(r.TxHash),
		Index:               r.LogIndex,
		Timestamp:           r.Timestamp,
		BlockNumber:         r.BlockNumber,
		QuoteAmount:         usdtAmountInWei,
		OriginalQuoteAmount: originalUSDTAmountInWei,
		User: common.KyberUserInfo{
			UserAddress: ethereum.HexToAddress(r.UserAddress[0]),
		},
		TokenInfo: common.TradeTokenInfo{
			SrcAddress:  SrcAddress,
			DestAddress: DstAddress,
		},
		ReserveAddress:  ethereum.HexToAddress(r.ReserveAddress),
		SrcAmount:       srcAmountInWei,
		DestAmount:      dstAmountInWei,
		ReceiverAddress: ethereum.HexToAddress(r.ReceiverAddr),
		TxDetail: common.TxDetail{
			GasUsed:        r.GasUsed,
			GasPrice:       gasPriceInWei,
			TransactionFee: transactionFeeInWei,
			TxSender:       ethereum.HexToAddress(r.TxSender),
		},
		Version: r.Version,
	}
	return tradeLog, nil
}

// LoadTradeLogsByTxHash get list of tradelogs by tx hash
func (tldb *TradeLogDB) LoadTradeLogsByTxHash(tx ethereum.Hash) ([]common.Tradelog, error) {
	var (
		logger      = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		queryResult []tradeLogDBData
		result      = make([]common.Tradelog, 0)
	)
	err := tldb.db.Select(&queryResult, selectTradeLogsWithTxHashQuery, tx.Hex())
	if err != nil {
		logger.Errorw("failed to get tradelog from database", "error", err)
		return nil, err
	}

	if len(queryResult) == 0 {
		logger.Debugw("empty result returned", "query", selectTradeLogsWithTxHashQuery)
		return result, nil
	}

	for _, r := range queryResult {
		tradeLog, err := tldb.tradeLogFromDBData(r)
		if err != nil {
			logger.Errorw("cannot parse db data to trade log", "error", err)
			return nil, err
		}
		result = append(result, tradeLog)
	}
	return result, nil
}

// LoadTradeLogs get list of tradelogs by timestamp from time to time
func (tldb *TradeLogDB) LoadTradeLogs(from, to time.Time) ([]common.Tradelog, error) {
	var (
		logger      = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		queryResult []tradeLogDBData
		result      = make([]common.Tradelog, 0)
	)
	err := tldb.db.Select(&queryResult, selectTradeLogsQuery, from, to)
	if err != nil {
		return nil, err
	}

	if len(queryResult) == 0 {
		logger.Debugw("empty result returned", "query", selectTradeLogsQuery)
		return result, nil
	}

	for _, r := range queryResult {
		tradeLog, err := tldb.tradeLogFromDBData(r)
		if err != nil {
			logger.Errorw("cannot parse db data to trade log", "error", err)
			return nil, err
		}
		result = append(result, tradeLog)
	}
	return result, nil
}

const insertionAddressTemplate = `INSERT INTO token(
	address
) VALUES(
	unnest($1::TEXT[])
)
ON CONFLICT ON CONSTRAINT token_address_key DO NOTHING`

const insertionUserTemplate string = `
INSERT INTO users(
	address,
	timestamp
) VALUES (
	:user_address,
	:timestamp
)
ON CONFLICT (address) 
DO NOTHING;`

const selectTradeLogsQuery = `
SELECT a.id, a.timestamp AS timestamp, a.block_number, a.usdt_amount, original_usdt_amount, 
ARRAY_AGG(d.address) AS user_address,
ARRAY_AGG(e.address) AS src_address, 
ARRAY_AGG(f.address) AS dst_address,
a.src_amount, 
a.dst_amount, 
a.index, tx_hash, tx_sender, receiver_address, 
COALESCE(gas_used, 0) as gas_used, COALESCE(gas_price, 0) as gas_price, 
COALESCE(transaction_fee, 0) as transaction_fee, 
r.address as reserve_address,
version
FROM tradelogs AS a
INNER JOIN users AS d ON a.user_address_id = d.id
INNER JOIN token AS e ON a.src_address_id = e.id
INNER JOIN token AS f ON a.dst_address_id = f.id
INNER JOIN reserve AS r on a.reserve_address_id = r.id
WHERE a.timestamp >= $1 and a.timestamp <= $2
GROUP BY a.id, r.address;
`

const selectTradeLogsWithTxHashQuery = `
SELECT 
a.timestamp AS timestamp, 
a.block_number, 
a.usdt_amount, 
original_usdt_amount, 
ARRAY_AGG(d.address) AS user_address,
ARRAY_AGG(e.address) AS src_address, 
ARRAY_AGG(f.address) AS dst_address,
a.src_amount, 
a.dst_amount, 
a.index, 
tx_hash, 
tx_sender, 
receiver_address, 
COALESCE(gas_used, 0) as gas_used, 
COALESCE(gas_price, 0) as gas_price, 
COALESCE(transaction_fee, 0) as transaction_fee,
sr.address as reserve_address,
version
FROM tradelogs AS a
INNER JOIN users AS d ON a.user_address_id = d.id
INNER JOIN token AS e ON a.src_address_id = e.id
INNER JOIN token AS f ON a.dst_address_id = f.id
INNER JOIN reserve sr ON a.reserve_address_id= sr.id
WHERE a.tx_hash=$1
GROUP BY a.id, sr.address;
`
