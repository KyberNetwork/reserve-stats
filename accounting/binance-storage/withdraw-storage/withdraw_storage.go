package withdrawstorage

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
	binanceWithdrawTable = "binance_withdraws"
)

//BinanceStorage is storage for binance fetcher including trade history and withdraw history
type BinanceStorage struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//NewDB return a new instance of binance storage
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB, tableName, timeIndexField, idType string) (*BinanceStorage, error) {
	var (
		logger = sugar.With("func", "accounting/binance-storage/binancestorage.NewDB")
	)

	const schemaFmt = `CREATE TABLE IF NOT EXISTS "%[1]s"
	(
	  id   %[2]s NOT NULL,
	  data JSONB,
	  CONSTRAINT %[1]s_pk PRIMARY KEY(id)
	);
	CREATE INDEX IF NOT EXISTS %[1]s_time_idx ON %[1]s ((data ->> '%[3]s'));
	`

	query := fmt.Sprintf(schemaFmt, tableName, idType, timeIndexField)
	logger.Debugw("create table query", "query", query)

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	logger.Info("binance table init successfully")

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
func (bd *BinanceStorage) DeleteTable(tableName string) error {
	query := fmt.Sprintf("DROP TABLE %s", tableName)
	if _, err := bd.db.Exec(query); err != nil {
		return err
	}
	return nil
}

//UpdateWithdrawHistory save withdraw history to db
func (bd *BinanceStorage) UpdateWithdrawHistory(withdrawHistories map[string]binance.WithdrawHistory) error {
	var (
		logger = bd.sugar.With("func", "accounting/binance_storage.UpdateWithdrawHistory")
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

	query := fmt.Sprintf(updateQuery, binanceWithdrawTable)
	logger.Debugw("query update withdraw history", "query", query)

	for _, withdraw := range withdrawHistories {
		withdrawJSON, err := json.Marshal(withdraw)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(query, withdraw.ID, withdrawJSON); err != nil {
			return err
		}
	}

	return err
}

//GetWithdrawHistory return list of withdraw fromTime to toTime
func (bd *BinanceStorage) GetWithdrawHistory(fromTime, toTime time.Time) ([]binance.WithdrawHistory, error) {
	var (
		logger   = bd.sugar.With("func", "account/binance_storage.GetTradeHistory")
		result   []binance.WithdrawHistory
		dbResult [][]byte
		tmp      binance.WithdrawHistory
	)
	const selectStmt = `SELECT data FROM %s WHERE data->>'applyTime'>=$1 AND data->>'applyTime'<=$2`
	query := fmt.Sprintf(selectStmt, binanceWithdrawTable)

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
