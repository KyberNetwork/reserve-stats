package postgres

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

//HuobiStorage defines the object to store Huobi data
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
`
	var (
		logger = sugar.With("func", "reserverates/storage/postgres/Newdb")
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

//Close close DB connection
func (hdb *HuobiStorage) Close() error {
	if hdb.db != nil {
		return hdb.db.Close()
	}
	return nil
}

//UpdateTradeHistory store the TradeHistory rate at that blockInfo
func (hdb *HuobiStorage) UpdateTradeHistory(trades map[int64]huobi.TradeHistory) error {
	var (
		nTrades = len(trades)
		logger  = hdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
			"number of records", nTrades,
		)
		ids      []int64
		dataJSON [][]byte
	)

	const updateStmt = `INSERT INTO huobi_trades (id, data)
	VALUES ( 
		unnest($1::BIGINT[]),
		unnest($2::JSONB[])
	)
	ON CONFLICT ON CONSTRAINT huobi_trades_pk DO NOTHING;`
	logger.Debugw("updating tradeHistory...", "query", updateStmt)

	tx, err := hdb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)
	for _, trade := range trades {
		data, err := json.Marshal(trade)
		if err != nil {
			return err
		}
		ids = append(ids, trade.ID)
		dataJSON = append(dataJSON, data)
	}
	_, err = tx.Exec(updateStmt, pq.Array(ids), pq.Array(dataJSON))
	if err != nil {
		return err
	}

	return err
}

//GetTradeHistory return tradehistory between from.. to.. in its json []byte form
func (hdb *HuobiStorage) GetTradeHistory(from, to time.Time) ([]huobi.TradeHistory, error) {
	var (
		dbResult [][]byte
		result   []huobi.TradeHistory
		logger   = hdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
			"from", from.String(),
			"to", to.String(),
		)
		tmp huobi.TradeHistory
	)
	const selectStmt = `SELECT data FROM huobi_trades WHERE data->>'created-at'>=$1 AND data->>'created-at'<$2`
	logger.Debugw("querying trade history...", "query", selectStmt)
	if err := hdb.db.Select(&dbResult, selectStmt, timeutil.TimeToTimestampMs(from), timeutil.TimeToTimestampMs(to)); err != nil {
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

//GetLastStoredTimestamp return the last stored timestamp in database
func (hdb *HuobiStorage) GetLastStoredTimestamp() (time.Time, error) {
	var (
		dbResult uint64
		result   = time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
		logger   = hdb.sugar.With(
			"func", "accounting/huobi/storage/postgres/RateStorage.GetLastStoredTimestamp",
		)
	)
	const selectStmt = `SELECT COALESCE(MAX(data->>'created-at'), '0') FROM huobi_trades`
	logger.Debugw("querying trade history...", "query", selectStmt)
	if err := hdb.db.Get(&dbResult, selectStmt); err != nil {
		return result, err
	}
	if dbResult != 0 {
		result = timeutil.TimestampMsToTime(dbResult)
	}
	return result, nil
}
