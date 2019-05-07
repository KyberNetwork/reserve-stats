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
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*HuobiStorage, error) {
	var (
		logger = sugar.With("func", "reserverates/storage/postgres/NewDB")
	)

	const schemaFMT = `
	CREATE TABLE IF NOT EXISTS huobi_withdrawals
(
	id bigint NOT NULL,
	data JSONB,
	CONSTRAINT huobi_withdrawals_pk PRIMARY KEY(id)
) ;

CREATE INDEX IF NOT EXISTS huobi_withdrawals_time_idx ON huobi_withdrawals ((data ->> 'created-at'));

`

	hs := &HuobiStorage{
		sugar: sugar,
		db:    db,
	}
	logger.Debugw("initializing database schema", "query", schemaFMT)

	if _, err := db.Exec(schemaFMT); err != nil {
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

//UpdateWithdrawHistory store the WithdrawHistory rate at that blockInfo
func (hdb *HuobiStorage) UpdateWithdrawHistory(withdraws []huobi.WithdrawHistory) error {
	var (
		logger = hdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
			"len(withdraws)", len(withdraws),
		)
		ids      []uint64
		dataJSON [][]byte
	)
	const updateStmt = `INSERT INTO huobi_withdrawals(id, data)
	VALUES ( 
		unnest($1::BIGINT[]),
		unnest($2::JSONB[])
	)
	ON CONFLICT ON CONSTRAINT huobi_withdrawals_pk DO NOTHING;`
	logger.Debugw("updating tradeHistory...", "query", updateStmt)

	tx, err := hdb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)
	//prepare data to insert into db
	for _, withdraw := range withdraws {
		data, err := json.Marshal(withdraw)
		if err != nil {
			return err
		}
		ids = append(ids, withdraw.ID)
		dataJSON = append(dataJSON, data)
	}
	_, err = tx.Exec(updateStmt, pq.Array(ids), pq.Array(dataJSON))
	if err != nil {
		return err
	}

	return err
}

//GetWithdrawHistory return tradehistory between from.. to.. in its json []byte form
func (hdb *HuobiStorage) GetWithdrawHistory(from, to time.Time) ([]huobi.WithdrawHistory, error) {
	var (
		dbResult [][]byte
		result   []huobi.WithdrawHistory
		logger   = hdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
			"from", from.String(),
			"to", to.String(),
		)
		tmp huobi.WithdrawHistory
	)
	const selectStmt = `SELECT data FROM huobi_withdrawals WHERE data->>'created-at'>=$1 AND data->>'created-at'<$2`
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

//GetLastIDStored return lastest id stored in database
func (hdb *HuobiStorage) GetLastIDStored() (uint64, error) {
	var (
		result uint64
		logger = hdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.GetLastIDStored",
		)
	)
	const selectStmt = `SELECT COALESCE(MAX(id),0) FROM huobi_withdrawals`
	logger.Debugw("querying trade history...", "query", selectStmt)

	if err := hdb.db.Get(&result, selectStmt); err != nil {
		return 0, err
	}

	return result, nil
}
