package withdrawalstorage

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/lib/pq"
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

	const schemaFmt = `CREATE TABLE IF NOT EXISTS "binance_withdrawals"
	(
	  id   text NOT NULL,
	  data JSONB,
	  CONSTRAINT binance_withdrawals_pk PRIMARY KEY(id)
	);
	CREATE INDEX IF NOT EXISTS binance_withdrawals_time_idx ON binance_withdrawals ((data ->> 'applyTime'));
	CREATE INDEX IF NOT EXISTS binance_withdrawals_txid_idx ON binance_withdrawals ((data ->> 'txId'));

	ALTER TABLE binance_withdrawals ADD COLUMN IF NOT EXISTS account TEXT;
	ALTER TABLE binance_withdrawals ADD COLUMN IF NOT EXISTS timestamp TIMESTAMP;
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

//UpdateWithdrawHistory save withdraw history to db
func (bd *BinanceStorage) UpdateWithdrawHistory(withdrawHistories []binance.WithdrawHistory, account string) (err error) {
	var (
		logger       = bd.sugar.With("func", caller.GetCurrentFunctionName())
		withdrawJSON []byte
		ids          []string
		dataJSON     [][]byte
		timestamps   []time.Time
	)
	const updateQuery = `INSERT INTO binance_withdrawals (id, data, timestamp, account)
	VALUES(
		unnest($1::TEXT[]),
		unnest($2::JSONB[]),
		unnest($3::TIMESTAMP[]),
		$4
	) ON CONFLICT ON CONSTRAINT binance_withdrawals_pk DO UPDATE SET data = EXCLUDED.data;
	`

	tx, err := bd.db.Beginx()
	if err != nil {
		return
	}

	defer pgsql.CommitOrRollback(tx, bd.sugar, &err)

	logger.Debugw("query update withdraw history", "query", updateQuery)

	// prepare data to insert into db
	for _, withdraw := range withdrawHistories {
		withdrawJSON, err = json.Marshal(withdraw)
		if err != nil {
			return
		}
		ids = append(ids, withdraw.ID)
		timestamp := timeutil.TimestampMsToTime(withdraw.ApplyTime)
		timestamps = append(timestamps, timestamp)
		dataJSON = append(dataJSON, withdrawJSON)
	}

	if _, err = tx.Exec(updateQuery, pq.Array(ids), pq.Array(dataJSON), pq.Array(timestamps), account); err != nil {
		return
	}

	return
}

//WithdrawRecord represent a record of binace withdraw
type WithdrawRecord struct {
	Account string        `db:"account"`
	Data    pq.ByteaArray `db:"data"`
}

//GetWithdrawHistory return list of withdraw fromTime to toTime
func (bd *BinanceStorage) GetWithdrawHistory(fromTime, toTime time.Time) (map[string][]binance.WithdrawHistory, error) {
	var (
		logger   = bd.sugar.With("func", caller.GetCurrentFunctionName())
		result   = make(map[string][]binance.WithdrawHistory)
		dbResult []WithdrawRecord
		tmp      binance.WithdrawHistory
	)
	const selectStmt = `SELECT account, ARRAY_AGG(data) as data FROM binance_withdrawals WHERE data->>'applyTime'>=$1 AND data->>'applyTime'<=$2 GROUP BY account;`

	logger.Debugw("querying trade history...", "query", selectStmt)

	from := timeutil.TimeToTimestampMs(fromTime)
	to := timeutil.TimeToTimestampMs(toTime)
	if err := bd.db.Select(&dbResult, selectStmt, from, to); err != nil {
		return result, err
	}

	for _, record := range dbResult {
		arrResult := []binance.WithdrawHistory{}
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

//GetLastStoredTimestamp return last timestamp stored in database
func (bd *BinanceStorage) GetLastStoredTimestamp(account string) (time.Time, error) {
	var (
		logger   = bd.sugar.With("func", caller.GetCurrentFunctionName())
		result   = time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
		dbResult uint64
		statuses = []string{strconv.Itoa(int(common.AwaitingApproval)), strconv.Itoa(int(common.Processing))}
	)
	const (
		selectStmt = `SELECT COALESCE(MAX(data->>'applyTime'), '0') FROM binance_withdrawals WHERE account = $1`
		//handle not completed withdraw
		latestNotCompleted = `SELECT COALESCE(MIN(data->>'applyTime'), '0') FROM binance_withdrawals WHERE data->>'status' = any($1) AND account = $2`
	)
	logger.Debugw("querying last stored timestamp", "query", selectStmt)

	if err := bd.db.Get(&dbResult, latestNotCompleted, pq.Array(statuses), account); err != nil {
		return result, err
	}
	logger.Debugw("min processing record time", "time", dbResult)

	if dbResult == 0 {
		if err := bd.db.Get(&dbResult, selectStmt, account); err != nil {
			return result, err
		}
	}

	if dbResult != 0 {
		result = timeutil.TimestampMsToTime(dbResult)
	}

	return result, nil
}

//UpdateWithdrawHistoryWithFee update fee into withdraw history table
func (bd *BinanceStorage) UpdateWithdrawHistoryWithFee(withdrawHistories []binance.WithdrawHistory, account string) (err error) {
	var (
		logger = bd.sugar.With("func", caller.GetCurrentFunctionName())
	)

	const (
		updateStmt = `UPDATE binance_withdrawals SET data = $1 WHERE data->>'txId' = $2`
	)
	logger.Debugw("update withdraw history", "query", updateStmt)
	tx, err := bd.db.Beginx()
	if err != nil {
		return
	}
	defer pgsql.CommitOrRollback(tx, bd.sugar, &err)

	for _, withdraw := range withdrawHistories {
		withdrawJSON, err := json.Marshal(withdraw)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(updateStmt, withdrawJSON, withdraw.TxID); err != nil {
			return err
		}
	}

	return nil
}
