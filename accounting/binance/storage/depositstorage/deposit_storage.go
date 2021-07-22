package depositstorage

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

// BinanceStorage is storage for binance fetcher including trade history and withdraw history
type BinanceStorage struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

// NewDB return a new instance of binance storage
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*BinanceStorage, error) {
	var (
		logger = sugar.With("func", caller.GetCurrentFunctionName())
	)

	const schemaFmt = `CREATE TABLE IF NOT EXISTS "binance_deposits"
	(
	  id   SERIAL,
	  account TEXT,
	  data JSONB,
	  timestamp TIMESTAMP,
	  CONSTRAINT binance_deposits_pk PRIMARY KEY(id)
	);
	CREATE INDEX IF NOT EXISTS binance_deposits_time_idx ON binance_deposits ((data ->> 'insertTime'));
	CREATE INDEX IF NOT EXISTS binance_deposits_txid_idx ON binance_deposits ((data ->> 'txId'));

	ALTER TABLE binance_deposits ADD COLUMN IF NOT EXISTS timestamp TIMESTAMP;
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

// Close database connection
func (bd *BinanceStorage) Close() error {
	if bd.db != nil {
		return bd.db.Close()
	}
	return nil
}

// UpdateDepositHistory save withdraw history to db
func (bd *BinanceStorage) UpdateDepositHistory(depositHistories []binance.DepositHistory, account string) (err error) {
	var (
		logger      = bd.sugar.With("func", caller.GetCurrentFunctionName())
		depositJSON []byte
		dataJSON    [][]byte
		timestamps  []time.Time
	)
	const updateQuery = `INSERT INTO binance_deposits (data, timestamp, account)
	VALUES(
		unnest($1::JSONB[]),
		unnest($2::TIMESTAMP[]),
		$3
	) ON CONFLICT ON CONSTRAINT binance_deposits_pk DO UPDATE SET data = EXCLUDED.data;
	`

	tx, err := bd.db.Beginx()
	if err != nil {
		logger.Errorw("failed to init tx", "error", err)
		return err
	}

	defer pgsql.CommitOrRollback(tx, bd.sugar, &err)

	logger.Debugw("query update deposit history", "query", updateQuery)

	// prepare data to insert into db
	for _, deposit := range depositHistories {
		depositJSON, err = json.Marshal(deposit)
		if err != nil {
			logger.Errorw("failed to marshal json", "error", err)
			return err
		}
		timestamp := timeutil.TimestampMsToTime(deposit.InsertTime)
		timestamps = append(timestamps, timestamp)
		dataJSON = append(dataJSON, depositJSON)
	}

	if _, err = tx.Exec(updateQuery, pq.Array(dataJSON), pq.Array(timestamps), account); err != nil {
		logger.Errorw("failed to update deposit history", "error", err)
		return err
	}

	return nil
}

// WithdrawRecord represent a record of binace withdraw
type DepositRecord struct {
	Account string        `db:"account"`
	Data    pq.ByteaArray `db:"data"`
}

// GetDepositHistory return list of withdraw fromTime to toTime
func (bd *BinanceStorage) GetDepositHistory(fromTime, toTime time.Time) (map[string][]binance.DepositHistory, error) {
	var (
		logger   = bd.sugar.With("func", caller.GetCurrentFunctionName())
		result   = make(map[string][]binance.DepositHistory)
		dbResult []DepositRecord
		tmp      binance.DepositHistory
	)
	const selectStmt = `SELECT account, ARRAY_AGG(data) as data FROM binance_deposits WHERE data->>'insertTime'>=$1 AND data->>'insertTime'<=$2 GROUP BY account;`

	logger.Debugw("querying trade history...", "query", selectStmt)

	from := timeutil.TimeToTimestampMs(fromTime)
	to := timeutil.TimeToTimestampMs(toTime)
	if err := bd.db.Select(&dbResult, selectStmt, from, to); err != nil {
		return result, err
	}

	for _, record := range dbResult {
		arrResult := []binance.DepositHistory{}
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

// GetLastStoredTimestamp return last timestamp stored in database
func (bd *BinanceStorage) GetLastStoredTimestamp(account string) (time.Time, error) {
	var (
		logger   = bd.sugar.With("func", caller.GetCurrentFunctionName())
		result   = time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
		dbResult uint64
		statuses = []string{strconv.Itoa(int(common.AwaitingApproval)), strconv.Itoa(int(common.Processing))}
	)
	const (
		selectStmt = `SELECT COALESCE(MAX(data->>'insertTime'), '0') FROM binance_deposits WHERE account = $1`
		// handle not completed withdraw
		latestNotCompleted = `SELECT COALESCE(MIN(data->>'insertTime'), '0') FROM binance_deposits WHERE data->>'status' = any($1) AND account = $2`
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
