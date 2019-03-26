package tradestorage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const defaultTradeTableName = "binance_trades"

//BinanceStorage is storage for binance fetcher including trade history and withdraw history
type BinanceStorage struct {
	sugar     *zap.SugaredLogger
	db        *sqlx.DB
	tableName string
}

// Option is the option for BinanceStorage constructor.
type Option func(*BinanceStorage)

// WithTableName is the option to create BinanceStorage.
func WithTableName(tableName string) Option {
	return func(storage *BinanceStorage) {
		storage.tableName = tableName
	}
}

//NewDB return a new instance of binance storage
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB, options ...Option) (*BinanceStorage, error) {
	var (
		logger = sugar.With("func", "accounting/binance-storage/binancestorage.NewDB")
	)

	const schemaFmt = `CREATE TABLE IF NOT EXISTS "%[1]s"
	(
	  id   bigint NOT NULL,
	  data JSONB,
	  CONSTRAINT %[1]s_pk PRIMARY KEY(id)
	);
	CREATE INDEX IF NOT EXISTS %[1]s_time_idx ON %[1]s ((data ->> 'time'));
	`

	s := &BinanceStorage{
		sugar: sugar,
		db:    db,
	}

	for _, option := range options {
		option(s)
	}

	if len(s.tableName) == 0 {
		s.tableName = defaultTradeTableName
	}

	query := fmt.Sprintf(schemaFmt, s.tableName)
	logger.Debugw("create table query", "query", query)

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	logger.Info("binance trades table init successfully")

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
	query := fmt.Sprintf("DROP TABLE %s", bd.tableName)
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
	const updateQuery = `INSERT INTO %[1]s (id, data)
	VALUES(
		$1,
		$2
	) ON CONFLICT ON CONSTRAINT %[1]s_pk DO NOTHING;
	`

	tx, err := bd.db.Beginx()
	if err != nil {
		return err
	}

	defer pgsql.CommitOrRollback(tx, bd.sugar, &err)

	query := fmt.Sprintf(updateQuery, bd.tableName)
	logger.Debugw("query update trade history", "query", query)
	for _, trade := range trades {
		tradeJSON, err = json.Marshal(trade)
		if err != nil {
			return
		}
		if _, err = tx.Exec(query, trade.ID, tradeJSON); err != nil {
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
	const selectStmt = `SELECT data FROM %s WHERE data->>'time'>=$1 AND data->>'time'<=$2`
	query := fmt.Sprintf(selectStmt, bd.tableName)

	logger.Debugw("querying trade history...", "query", query)

	from := timeutil.TimeToTimestampMs(fromTime)
	to := timeutil.TimeToTimestampMs(toTime)
	if err := bd.db.Select(&dbResult, query, from, to); err != nil {
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
