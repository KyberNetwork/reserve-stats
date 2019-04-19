package tradestorage

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
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
		logger = sugar.With("func", "accounting/binance-storage/binancestorage.NewDB")
	)

	const schemaFmt = `CREATE TABLE IF NOT EXISTS "binance_trades"
	(
	  id   bigint NOT NULL,
	  data JSONB,
	  CONSTRAINT binance_trades_pk PRIMARY KEY(id)
	);
	CREATE INDEX IF NOT EXISTS binance_trades_time_idx ON binance_trades ((data ->> 'time'));
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

//DeleteTable remove trades table
func (bd *BinanceStorage) DeleteTable() error {
	query := "DROP TABLE binance_trades;"
	if _, err := bd.db.Exec(query); err != nil {
		return err
	}
	return nil
}

//UpdateTradeHistory save trade history into a postgres db
func (bd *BinanceStorage) UpdateTradeHistory(trades []binance.TradeHistory) (err error) {
	var (
		logger    = bd.sugar.With("func", "accounting/binance_storage.UpdateTradeHistory")
		tradeJSON []byte
	)
	const updateQuery = `INSERT INTO binance_trades (id, data)
	VALUES(
		$1,
		$2
	) ON CONFLICT ON CONSTRAINT binance_trades_pk DO NOTHING;
	`

	tx, err := bd.db.Beginx()
	if err != nil {
		return err
	}

	defer pgsql.CommitOrRollback(tx, bd.sugar, &err)

	logger.Debugw("query update trade history", "query", updateQuery)
	for _, trade := range trades {
		tradeJSON, err = json.Marshal(trade)
		if err != nil {
			return
		}
		if _, err = tx.Exec(updateQuery, trade.ID, tradeJSON); err != nil {
			return
		}
	}

	return err
}

//GetTradeHistory return trade history from binance storage
func (bd *BinanceStorage) GetTradeHistory(fromTime, toTime time.Time) ([]binance.TradeHistory, error) {
	var (
		logger   = bd.sugar.With("func", "account/binance_storage.GetTradeHistory")
		result   []binance.TradeHistory
		dbResult [][]byte
		tmp      binance.TradeHistory
	)
	const selectStmt = `SELECT data FROM binance_trades WHERE data->>'time'>=$1 AND data->>'time'<=$2`

	logger.Debugw("querying trade history...", "query", selectStmt)

	from := timeutil.TimeToTimestampMs(fromTime)
	to := timeutil.TimeToTimestampMs(toTime)
	if err := bd.db.Select(&dbResult, selectStmt, from, to); err != nil {
		return result, err
	}

	for _, data := range dbResult {
		if err := json.Unmarshal(data, &tmp); err != nil {
			return result, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

//GetLastStoredID return last stored id
func (bd *BinanceStorage) GetLastStoredID() (uint64, error) {
	var (
		logger = bd.sugar.With("func", "account/binance_storage.GetLastStoredID")
		result uint64
	)
	const selectStmt = `SELECT COALESCE(MAX(id), 0) FROM binance_trades`

	logger.Debugw("querying last stored id", "query", selectStmt)

	if err := bd.db.Get(&result, selectStmt); err != nil {
		return 0, err
	}

	return result, nil
}
