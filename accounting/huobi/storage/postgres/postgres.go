package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

// HuobiStorage defines the object to store Huobi data
type HuobiStorage struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

// NewDB return the HuobiStorage instance. User must call Close() before exit.
// tableNames is a list of string for tablename[huobitrades]. It can be optional
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*HuobiStorage, error) {
	const schemaFMT = `
	CREATE TABLE IF NOT EXISTS huobi_trades
(
	id bigint NOT NULL,
	data JSONB,
	CONSTRAINT huobi_trades_pk PRIMARY KEY(id)
) ;
CREATE INDEX IF NOT EXISTS huobi_trades_time_idx ON huobi_trades ((data ->> 'created-at'));
ALTER TABLE huobi_trades ADD COLUMN IF NOT EXISTS created TIMESTAMPTZ;
ALTER TABLE huobi_trades ADD COLUMN IF NOT EXISTS account TEXT;
`
	var (
		logger = sugar.With("func", caller.GetCurrentFunctionName())
	)
	hs := &HuobiStorage{
		sugar: sugar,
		db:    db,
	}
	logger.Debugw("initializing database schema", "query", schemaFMT)
	if _, err := hs.db.Exec(schemaFMT); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")
	return hs, nil
}

// Close close DB connection
func (hdb *HuobiStorage) Close() error {
	if hdb.db != nil {
		return hdb.db.Close()
	}
	return nil
}

// UpdateTradeHistory store the TradeHistory rate at that blockInfo
func (hdb *HuobiStorage) UpdateTradeHistory(trades map[int64]huobi.TradeHistory, account string) (err error) {
	var (
		nTrades = len(trades)
		logger  = hdb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"number of records", nTrades,
		)
		ids        []int64
		dataJSON   [][]byte
		timestamps []time.Time
	)

	const updateStmt = `INSERT INTO huobi_trades (id, data, created, account)
	VALUES ( 
		unnest($1::BIGINT[]),
		unnest($2::JSONB[]),
		unnest($3::TIMESTAMPTZ[]),
		$4
	)
	ON CONFLICT ON CONSTRAINT huobi_trades_pk DO UPDATE SET created = EXCLUDED.created;`
	logger.Debugw("updating tradeHistory...", "query", updateStmt)

	tx, err := hdb.db.Beginx()
	if err != nil {
		return
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)
	for _, trade := range trades {
		var data []byte
		data, err = json.Marshal(trade)
		if err != nil {
			return
		}
		createdAt := timeutil.TimestampMsToTime(trade.CreatedAt)
		ids = append(ids, trade.ID)
		dataJSON = append(dataJSON, data)
		timestamps = append(timestamps, createdAt)
	}
	_, err = tx.Exec(updateStmt, pq.Array(ids), pq.Array(dataJSON), pq.Array(timestamps), account)
	if err != nil {
		return
	}
	return err
}

// TradeHistoryDB ...
type TradeHistoryDB struct {
	Account string        `db:"account"`
	Data    pq.ByteaArray `db:"data"`
}

// GetTradeHistory return tradehistory between from.. to.. in its json []byte form
func (hdb *HuobiStorage) GetTradeHistory(from, to time.Time) (map[string][]huobi.TradeHistory, error) {
	var (
		dbResult []TradeHistoryDB
		result   = make(map[string][]huobi.TradeHistory)
		logger   = hdb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"from", from.String(),
			"to", to.String(),
		)
		tmp huobi.TradeHistory
	)
	const selectStmt = `SELECT account, ARRAY_AGG(data) as data FROM huobi_trades WHERE data->>'created-at'>=$1 AND data->>'created-at'<$2 GROUP BY account;`
	logger.Debugw("querying trade history...", "query", selectStmt)
	if err := hdb.db.Select(&dbResult, selectStmt, timeutil.TimeToTimestampMs(from), timeutil.TimeToTimestampMs(to)); err != nil {
		return result, err
	}
	for _, record := range dbResult {
		arrResult := []huobi.TradeHistory{}
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

// GetLastStoredTimestamp return the last stored timestamp in database
func (hdb *HuobiStorage) GetLastStoredTimestamp(account string) (time.Time, error) {
	var (
		dbResult uint64
		result   = time.Now().Add(time.Hour * 24 * 90 * -1) // 90 days ago
		logger   = hdb.sugar.With("func", caller.GetCurrentFunctionName())
	)
	const selectStmt = `SELECT COALESCE(MAX(data->>'created-at'), '0') FROM huobi_trades WHERE account = $1`
	logger.Debugw("querying trade history...", "query", selectStmt)
	if err := hdb.db.Get(&dbResult, selectStmt, account); err != nil {
		if err != sql.ErrNoRows {
			return result, err
		}
	}
	if dbResult != 0 {
		result = timeutil.TimestampMsToTime(dbResult)
	}
	return result, nil
}
