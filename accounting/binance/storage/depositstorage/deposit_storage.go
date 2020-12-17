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

//Storage is storage for binance fetcher including trade history and deposit history
type Storage struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//NewDB return a new instance of binance storage
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*Storage, error) {
	var (
		logger = sugar.With("func", caller.GetCurrentFunctionName())
	)

	const schemaFmt = `CREATE TABLE IF NOT EXISTS "binance_deposits"
	(
	  id SERIAL,
	  asset text NOT NULL,
	  tx_id text NOT NULL,
	  data JSONB,
	  account TEXT NOT NULL,
	  timestamp timestamptz,
	  CONSTRAINT binance_deposits_pk PRIMARY KEY(tx_id)
	);
	CREATE INDEX IF NOT EXISTS binance_deposits_time_idx ON binance_deposits ((data ->> 'insertTime'));
	`

	s := &Storage{
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
func (bd *Storage) Close() error {
	if bd.db != nil {
		return bd.db.Close()
	}
	return nil
}

//UpdateDepositHistory save deposit history to db
func (bd *Storage) UpdateDepositHistory(depositHistories []binance.DepositHistory, account string) (err error) {
	var (
		logger        = bd.sugar.With("func", caller.GetCurrentFunctionName())
		depositJSON   []byte
		assets, txIDs []string
		dataJSON      [][]byte
		timestamps    []time.Time
	)
	const updateQuery = `INSERT INTO binance_deposits (asset, tx_id, data, timestamp, account)
	VALUES(
		unnest($1::TEXT[]),
		unnest($2::TEXT[]),
		unnest($3::JSONB[]),
		unnest($4::TIMESTAMP[]),
		$5
	) ON CONFLICT ON CONSTRAINT binance_deposits_pk DO UPDATE SET data = EXCLUDED.data;
	`

	tx, err := bd.db.Beginx()
	if err != nil {
		return
	}

	defer pgsql.CommitOrRollback(tx, bd.sugar, &err)

	logger.Debugw("query update deposit history", "query", updateQuery)

	// prepare data to insert into db
	for _, deposit := range depositHistories {
		depositJSON, err = json.Marshal(deposit)
		if err != nil {
			return
		}
		assets = append(assets, deposit.Asset)
		txIDs = append(txIDs, deposit.TxID)
		timestamp := timeutil.TimestampMsToTime(deposit.InsertTime)
		timestamps = append(timestamps, timestamp)
		dataJSON = append(dataJSON, depositJSON)
	}

	if _, err = tx.Exec(updateQuery, pq.Array(assets), pq.Array(txIDs), pq.Array(dataJSON), pq.Array(timestamps), account); err != nil {
		return
	}

	return
}

//DepositRecord represent a record of binace deposit
type DepositRecord struct {
	Account string        `db:"account"`
	Data    pq.ByteaArray `db:"data"`
}

//GetDepositHistory return list of deposit fromTime to toTime
func (bd *Storage) GetDepositHistory(fromTime, toTime time.Time) (map[string][]binance.DepositHistory, error) {
	var (
		logger   = bd.sugar.With("func", caller.GetCurrentFunctionName())
		result   = make(map[string][]binance.DepositHistory)
		dbResult []DepositRecord
		tmp      binance.DepositHistory
	)
	const selectStmt = `SELECT account, ARRAY_AGG(data) as data FROM binance_deposits WHERE timestamp>=$1 AND timstamp<=$2 GROUP BY account;`

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

//GetLastStoredTimestamp return last timestamp stored in database
func (bd *Storage) GetLastStoredTimestamp(account, asset string) (time.Time, error) {
	var (
		logger   = bd.sugar.With("func", caller.GetCurrentFunctionName())
		result   = time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
		dbResult uint64
		statuses = []string{strconv.Itoa(int(common.AwaitingApproval)), strconv.Itoa(int(common.Processing))}
	)
	const (
		selectStmt = `SELECT COALESCE(MAX(data->>'insertTime'), '0') FROM binance_deposits WHERE account = $1 AND asset = $2;`
		//handle not completed deposit
		latestNotCompleted = `SELECT COALESCE(MIN(data->>'insertTime'), '0') FROM binance_deposits WHERE data->>'status' = any($1) AND account = $2 AND asset = $3;`
	)
	logger.Debugw("querying last stored timestamp", "query", selectStmt, "account", account, "asset", asset)

	if err := bd.db.Get(&dbResult, latestNotCompleted, pq.Array(statuses), account, asset); err != nil {
		return result, err
	}
	logger.Debugw("min processing record time", "time", dbResult)

	if dbResult == 0 {
		if err := bd.db.Get(&dbResult, selectStmt, account, asset); err != nil {
			return result, err
		}
	}

	if dbResult != 0 {
		result = timeutil.TimestampMsToTime(dbResult)
	}

	return result, nil
}
