package binancestorage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	binanceTradeTable = "binance_trades"
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

	const schemaFmt = `CREATE TABLE IF NOT EXISTS "%[1]s"
	(
	  id   bigint NOT NULL,
	  data JSONB,
	  CONSTRAINT %[1]s_pk PRIMARY KEY(id)
	);
	CREATE INDEX IF NOT EXISTS %[1]s_time_idx ON %[1]s ((data ->> 'time'));
	`

	query := fmt.Sprintf(schemaFmt, binanceTradeTable)
	logger.Debugw("create table query", "query", query)

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	logger.Info("binance trades table init successfully")

	return &BinanceStorage{
		sugar: sugar,
		db:    db,
	}, nil
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
	query := fmt.Sprintf("DROP TABLE %s", binanceTradeTable)
	if _, err := bd.db.Exec(query); err != nil {
		return err
	}
	return nil
}

//UpdateTradeHistory save trade history into a postgres db
func (bd *BinanceStorage) UpdateTradeHistory(trades []binance.TradeHistory) error {
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

	for _, trade := range trades {
		query := fmt.Sprintf(updateQuery, binanceTradeTable)
		tradeJSON, err := json.Marshal(trade)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(query, trade.ID, tradeJSON); err != nil {
			return err
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
	query := fmt.Sprintf(selectStmt, binanceTradeTable)

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
