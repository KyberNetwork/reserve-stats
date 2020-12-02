package tradestorage

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

//BinanceStorage is storage for binance fetcher including trade history and withdraw history
type BinanceStorage struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//NewDB return a new instance of binance storage
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*BinanceStorage, error) {
	var (
		logger = sugar.With("func", caller.GetCurrentFunctionName())
	)

	const schemaFmt = `CREATE TABLE IF NOT EXISTS "binance_trades"
	(
	  id   bigint NOT NULL,
	  symbol TEXT NOT NULL,
	  data JSONB,
	  timestamp TIMESTAMP NOT NULL,
	  CONSTRAINT binance_trades_pk PRIMARY KEY(id, symbol)
	);
	CREATE INDEX IF NOT EXISTS binance_trades_time_idx ON binance_trades (timestamp);
	CREATE INDEX IF NOT EXISTS binance_trades_symbol_idx ON binance_trades (symbol);
	ALTER TABLE binance_trades ADD COLUMN IF NOT EXISTS account TEXT;

	CREATE TABLE IF NOT EXISTS "binance_margin_trades"
	(
		id bigint NOT NULL,
		symbol TEXT NOT NULL,
		data JSONB,
		timestamp TIMESTAMP NOT NULL,
		account TEXT,
		CONSTRAINT binance_margin_trades_pk PRIMARY KEY(id, symbol)
	);
	CREATE INDEX IF NOT EXISTS binance_margin_trades_time_idx ON binance_trades (timestamp);
	CREATE INDEX IF NOT EXISTS binance_margin_trades_symbol_idx ON binance_trades (symbol);

	CREATE TABLE IF NOT EXISTS "binance_convert_to_eth_price"
	(
		original_symbol TEXT NOT NULL,
		symbol TEXT NOT NULL,
		price FLOAT NOT NULL,
		timestamp BIGINT NOT NULL,
		original_trade JSONB,
		trade JSONB,
		CONSTRAINT binance_convert_to_eth_price_pk PRIMARY KEY(symbol, price, timestamp)
	);
	CREATE INDEX IF NOT EXISTS binance_convert_to_eth_price_time_idx ON binance_convert_to_eth_price(timestamp);
	`

	s := &BinanceStorage{
		sugar: sugar,
		db:    db,
	}

	logger.Debugw("create table query", "query", schemaFmt)

	if _, err := db.Exec(schemaFmt); err != nil {
		return nil, err
	}

	logger.Info("binance table init successfully")

	return s, nil
}

//Close database connection
func (bd *BinanceStorage) Close() error {
	if bd.db != nil {
		return bd.db.Close()
	}
	return nil
}

//UpdateTradeHistory save trade history into a postgres db
func (bd *BinanceStorage) UpdateTradeHistory(trades []binance.TradeHistory, account string) (err error) {
	var (
		logger     = bd.sugar.With("func", caller.GetCurrentFunctionName())
		tradeJSON  []byte
		dataJSON   [][]byte
		ids        []uint64
		timestamps []time.Time
		symbols    []string
	)
	const updateQuery = `INSERT INTO binance_trades (id, data, timestamp, symbol, account)
	VALUES(
		unnest($1::BIGINT[]),
		unnest($2::JSONB[]),
		unnest($3::TIMESTAMP[]),
		unnest($4::TEXT[]),
		$5
	) ON CONFLICT ON CONSTRAINT binance_trades_pk DO NOTHING;
	`

	tx, err := bd.db.Beginx()
	if err != nil {
		return err
	}

	defer pgsql.CommitOrRollback(tx, bd.sugar, &err)

	// prepare data for insert into db
	for _, trade := range trades {
		tradeJSON, err = json.Marshal(trade)
		if err != nil {
			logger.Errorw("failed to marshal trade", "error", err)
			return
		}
		time := timeutil.TimestampMsToTime(trade.Time)
		ids = append(ids, trade.ID)
		dataJSON = append(dataJSON, tradeJSON)
		timestamps = append(timestamps, time)
		symbols = append(symbols, trade.Symbol)
	}

	if _, err = tx.Exec(updateQuery, pq.Array(ids), pq.Array(dataJSON), pq.Array(timestamps), pq.Array(symbols), account); err != nil {
		return
	}

	return err
}

// TradeHistoryDB response from db
type TradeHistoryDB struct {
	Account string        `db:"account"`
	Data    pq.ByteaArray `db:"data"`
}

//GetTradeHistory return trade history from binance storage
func (bd *BinanceStorage) GetTradeHistory(fromTime, toTime time.Time) (map[string][]binance.TradeHistory, error) {
	var (
		logger   = bd.sugar.With("func", caller.GetCurrentFunctionName())
		result   = make(map[string][]binance.TradeHistory)
		dbResult []TradeHistoryDB
		tmp      binance.TradeHistory
	)
	const selectStmt = `SELECT account, ARRAY_AGG(data) as data FROM binance_trades WHERE timestamp >=$1::TIMESTAMP AND timestamp <=$2::TIMESTAMP GROUP BY account`

	logger.Debugw("querying trade history...", "query", selectStmt)

	if err := bd.db.Select(&dbResult, selectStmt, fromTime.UTC(), toTime.UTC()); err != nil {
		return result, err
	}
	for _, record := range dbResult {
		arrResult := []binance.TradeHistory{}
		for _, data := range record.Data {
			if err := json.Unmarshal(data, &tmp); err != nil {
				return result, err
			}
			arrResult = append(arrResult, tmp)
		}
		result[record.Account] = arrResult
	}

	return result, nil
}

//GetLastStoredID return last stored id
func (bd *BinanceStorage) GetLastStoredID(symbol, account string) (uint64, error) {
	var (
		logger = bd.sugar.With("func", caller.GetCurrentFunctionName())
		result uint64
	)
	const selectStmt = `SELECT COALESCE(MAX(id), 0) FROM binance_trades WHERE symbol=$1 AND account=$2`

	logger.Debugw("querying last stored id", "query", selectStmt)

	if err := bd.db.Get(&result, selectStmt, symbol, account); err != nil {
		return 0, err
	}

	return result, nil
}

// UpdateMarginTradeHistory update margin trades into db
func (bd *BinanceStorage) UpdateMarginTradeHistory(marginTrades []binance.TradeHistory, account string) error {
	var (
		logger     = bd.sugar.With("func", caller.GetCurrentFunctionName())
		tradeJSON  []byte
		dataJSON   [][]byte
		ids        []uint64
		timestamps []time.Time
		symbols    []string
	)
	const updateQuery = `INSERT INTO binance_margin_trades (id, data, timestamp, symbol, account)
	VALUES(
		unnest($1::BIGINT[]),
		unnest($2::JSONB[]),
		unnest($3::TIMESTAMP[]),
		unnest($4::TEXT[]),
		$5
	) ON CONFLICT ON CONSTRAINT binance_margin_trades_pk DO NOTHING;
	`

	tx, err := bd.db.Beginx()
	if err != nil {
		return err
	}

	defer pgsql.CommitOrRollback(tx, bd.sugar, &err)
	logger.Debugw("query update trade history", "query", updateQuery)

	// prepare data for insert into db
	for _, trade := range marginTrades {
		tradeJSON, err = json.Marshal(trade)
		if err != nil {
			return err
		}
		time := timeutil.TimestampMsToTime(trade.Time)
		ids = append(ids, trade.ID)
		dataJSON = append(dataJSON, tradeJSON)
		timestamps = append(timestamps, time)
		symbols = append(symbols, trade.Symbol)
	}

	if _, err = tx.Exec(updateQuery, pq.Array(ids), pq.Array(dataJSON), pq.Array(timestamps), pq.Array(symbols), account); err != nil {
		return err
	}

	return err
}

// GetMarginTradeHistory return list of trade history from time to time
func (bd *BinanceStorage) GetMarginTradeHistory(fromTime, toTime time.Time) (map[string][]binance.TradeHistory, error) {
	var (
		logger   = bd.sugar.With("func", caller.GetCurrentFunctionName())
		result   = make(map[string][]binance.TradeHistory)
		dbResult []TradeHistoryDB
		tmp      binance.TradeHistory
	)
	const selectStmt = `SELECT account, ARRAY_AGG(data) as data FROM binance_margin_trades WHERE timestamp >=$1::TIMESTAMP AND timestamp <=$2::TIMESTAMP GROUP BY account;`

	logger.Debugw("querying margin trade history...", "query", selectStmt)

	if err := bd.db.Select(&dbResult, selectStmt, fromTime.UTC(), toTime.UTC()); err != nil {
		return result, err
	}
	for _, record := range dbResult {
		arrResult := []binance.TradeHistory{}
		for _, data := range record.Data {
			if err := json.Unmarshal(data, &tmp); err != nil {
				return result, err
			}
			arrResult = append(arrResult, tmp)
		}
		result[record.Account] = arrResult
	}

	return result, nil
}

// GetLastStoredMarginTradeID return last stored margin trade id
func (bd *BinanceStorage) GetLastStoredMarginTradeID(symbol, account string) (uint64, error) {
	var (
		logger = bd.sugar.With("func", caller.GetCurrentFunctionName())
		result uint64
	)
	const selectStmt = `SELECT COALESCE(MAX(id), 0) FROM binance_margin_trades WHERE symbol=$1 AND account=$2`

	if err := bd.db.Get(&result, selectStmt, symbol, account); err != nil {
		logger.Errorw("failed to get last stored margin trade id", "error", err)
		return 0, err
	}

	return result, nil
}

// GetTradeByTimestamp ...
func (bd *BinanceStorage) GetTradeByTimestamp(symbol string, timestamp time.Time) (binance.TradeHistory, error) {
	var (
		logger = bd.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"symbol", symbol,
			"timestamp", timestamp,
		)
		result   binance.TradeHistory
		resultDB []byte
	)
	const query = "SELECT data FROM binance_trades where timestamp <= $1::TIMESTAMP AND data->>'symbol' = $2 LIMIT 1;"
	if err := bd.db.Get(&resultDB, query, timestamp, symbol); err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}
		logger.Errorw("failed to get trade", "error", err)
		return result, err
	}
	err := json.Unmarshal(resultDB, &result)
	return result, err
}

// UpdateConvertToETHPrice ...
func (bd *BinanceStorage) UpdateConvertToETHPrice(originalSymbol, symbol string, prices []float64, timestamps []uint64, originalTrades, trades []binance.TradeHistory) error {
	var (
		logger = bd.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"symbol", symbol,
		)
		tradesJSON, originalTradesJSON [][]byte
	)
	logger.Info("update eth trade")
	const query = `INSERT INTO binance_convert_to_eth_price (original_symbol, symbol, price, timestamp, original_trade, trade)
				   VALUES (
					   $1, 
					   $2, 
					   unnest($3::FLOAT[]), 
					   unnest($4::BIGINT[]), 
					   unnest($5::JSONB[]), 
					   unnest($6::JSONB[])
					   ) ON CONFLICT (symbol, price, timestamp) DO NOTHING;`
	for _, trade := range trades {
		tradeJSON, err := json.Marshal(trade)
		if err != nil {
			logger.Errorw("failed to marshal trade", "error", err)
			return err
		}
		tradesJSON = append(tradesJSON, tradeJSON)
	}
	for _, originalTrade := range originalTrades {
		originalTradeJSON, err := json.Marshal(originalTrade)
		if err != nil {
			logger.Errorw("failed to marshal trade", "error", err)
			return err
		}
		originalTradesJSON = append(originalTradesJSON, originalTradeJSON)
	}
	if _, err := bd.db.Exec(query, originalSymbol, symbol, pq.Array(prices), pq.Array(timestamps), pq.Array(originalTradesJSON), pq.Array(tradesJSON)); err != nil {
		logger.Errorw("failed to update eth trade", "error", err)
		return err
	}
	return nil
}

// GetConvertToETHPrice ...
func (bd *BinanceStorage) GetConvertToETHPrice(fromTime, toTime uint64) ([]binance.ConvertToETHPrice, error) {
	var (
		logger = bd.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"fromTime", fromTime,
			"toTime", toTime,
		)
		result []binance.ConvertToETHPrice
		err    error
	)
	const query = `SELECT symbol, price, timestamp FROM binance_convert_to_eth_price WHERE timestamp >= $1 AND timestamp <= $2;`
	if err := bd.db.Select(&result, query, fromTime, toTime); err != nil {
		logger.Errorw("failed to get convert eth price", "error", err)
		return result, err
	}
	return result, err
}
