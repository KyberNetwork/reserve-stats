package postgres

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
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

//TearDown removes all the tables
func (hdb *HuobiStorage) TearDown() error {
	const dropFMT = `
	DROP TABLE huobi_withdrawals;
	`
	hdb.sugar.Debugw("tearingdown", "query", dropFMT)
	_, err := hdb.db.Exec(dropFMT)
	return err
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
	)
	const updateStmt = `INSERT INTO huobi_withdrawals(id, data)
	VALUES ( 
		$1,
		$2
	)
	ON CONFLICT ON CONSTRAINT huobi_withdrawals_pk DO NOTHING;`
	logger.Debugw("updating tradeHistory...", "query", updateStmt)

	tx, err := hdb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)
	for _, withdraw := range withdraws {
		data, err := json.Marshal(withdraw)
		if err != nil {
			return err
		}
		_, err = tx.Exec(updateStmt, withdraw.ID, data)
		if err != nil {
			return err
		}
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
	const selectStmt = `SELECT COALESCE(MAX(id),0) FROM %[1]s`
	logger.Debugw("querying trade history...", "query", selectStmt)

	if err := hdb.db.Get(&result, selectStmt); err != nil {
		return 0, err
	}

	return result, nil
}
