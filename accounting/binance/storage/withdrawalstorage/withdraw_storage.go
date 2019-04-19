package withdrawalstorage

import (
	"encoding/json"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
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

	const schemaFmt = `CREATE TABLE IF NOT EXISTS "binance_withdrawals"
	(
	  id   text NOT NULL,
	  data JSONB,
	  CONSTRAINT binance_withdrawals_pk PRIMARY KEY(id)
	);
	CREATE INDEX IF NOT EXISTS binance_withdrawals_time_idx ON binance_withdrawals ((data ->> 'applyTime'));
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
	query := "DROP TABLE binance_withdrawals"
	if _, err := bd.db.Exec(query); err != nil {
		return err
	}
	return nil
}

//UpdateWithdrawHistory save withdraw history to db
func (bd *BinanceStorage) UpdateWithdrawHistory(withdrawHistories []binance.WithdrawHistory) (err error) {
	var (
		logger       = bd.sugar.With("func", "accounting/binance_storage.UpdateWithdrawHistory")
		withdrawJSON []byte
	)
	const updateQuery = `INSERT INTO binance_withdrawals (id, data)
	VALUES(
		$1,
		$2
	) ON CONFLICT ON CONSTRAINT binance_withdrawals_pk DO NOTHING;
	`

	tx, err := bd.db.Beginx()
	if err != nil {
		return
	}

	defer pgsql.CommitOrRollback(tx, bd.sugar, &err)

	logger.Debugw("query update withdraw history", "query", updateQuery)

	for _, withdraw := range withdrawHistories {
		withdrawJSON, err = json.Marshal(withdraw)
		if err != nil {
			return
		}
		if _, err = tx.Exec(updateQuery, withdraw.ID, withdrawJSON); err != nil {
			return
		}
	}

	return
}

//GetWithdrawHistory return list of withdraw fromTime to toTime
func (bd *BinanceStorage) GetWithdrawHistory(fromTime, toTime time.Time) ([]binance.WithdrawHistory, error) {
	var (
		logger   = bd.sugar.With("func", "account/binance_storage.GetTradeHistory")
		result   []binance.WithdrawHistory
		dbResult [][]byte
		tmp      binance.WithdrawHistory
	)
	const selectStmt = `SELECT data FROM binance_withdrawals WHERE data->>'applyTime'>=$1 AND data->>'applyTime'<=$2`

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

//GetLastStoredTimestamp return last timestamp stored in database
func (bd *BinanceStorage) GetLastStoredTimestamp() (time.Time, error) {
	var (
		logger   = bd.sugar.With("func", "account/binance_storage.GetLastStoredTimestamp")
		result   = time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
		dbResult uint64
	)
	const selectStmt = `SELECT COALESCE(MAX(data->>'applyTime'), '0') FROM binance_withdrawals`

	logger.Debugw("querying last stored timestamp", "query", selectStmt)

	if err := bd.db.Get(&dbResult, selectStmt); err != nil {
		return result, err
	}

	if dbResult != 0 {
		result = timeutil.TimestampMsToTime(dbResult)
	}

	return result, nil
}
