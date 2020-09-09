package tradestorage

import (
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
	logger.Debugw("query update trade history", "query", updateQuery)

	// prepare data for insert into db
	for _, trade := range trades {
		tradeJSON, err = json.Marshal(trade)
		if err != nil {
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
	const selectStmt = `SELECT account, ARRAY_AGG(data) as data FROM binance_trades WHERE extract(epoch from timestamp)>=$1 AND extract(epoch from timestamp)<=$2 GROUP BY account`

	logger.Debugw("querying trade history...", "query", selectStmt)

	from := timeutil.TimeToTimestampMs(fromTime) / 1000
	to := timeutil.TimeToTimestampMs(toTime) / 1000
	if err := bd.db.Select(&dbResult, selectStmt, from, to); err != nil {
		return result, err
	}
	logger.Debugw("db result", "length", len(dbResult))
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
