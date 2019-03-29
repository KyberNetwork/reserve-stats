package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	huobiWithdrawalTableName = "huobi_withdrawal"
)

//Option define init behaviour for db storage.
type Option func(*HuobiStorage) error

//WithWithdrawalTableName return Option to set trade table Name
func WithWithdrawalTableName(name string) Option {
	return func(hs *HuobiStorage) error {
		if hs.tableNames == nil {
			hs.tableNames = make(map[string]string)
		}
		hs.tableNames[huobiWithdrawalTableName] = name
		return nil
	}
}

//HuobiStorage defines the object to store Huobi data
type HuobiStorage struct {
	sugar      *zap.SugaredLogger
	db         *sqlx.DB
	tableNames map[string]string
}

// NewDB return the HuobiStorage instance. User must call Close() before exit.
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB, options ...Option) (*HuobiStorage, error) {
	var (
		logger = sugar.With("func", "reserverates/storage/postgres/NewDB")
		//set Default table name
		tableNames = map[string]string{
			huobiWithdrawalTableName: huobiWithdrawalTableName,
		}
	)

	const schemaFMT = `
	CREATE TABLE IF NOT EXISTS %[1]s
(
	id bigint NOT NULL,
	data JSONB,
	CONSTRAINT %[1]s_pk PRIMARY KEY(id)
) ;

CREATE INDEX IF NOT EXISTS %[1]s_time_idx ON %[1]s ((data ->> 'created-at'));

`

	hs := &HuobiStorage{
		sugar:      sugar,
		db:         db,
		tableNames: tableNames,
	}
	for _, opt := range options {
		if err := opt(hs); err != nil {
			return nil, err
		}
	}

	query := fmt.Sprintf(schemaFMT, hs.tableNames[huobiWithdrawalTableName])
	logger.Debugw("initializing database schema", "query", query)

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")
	return hs, nil
}

//TearDown removes all the tables
func (hdb *HuobiStorage) TearDown() error {
	const dropFMT = `
	DROP TABLE %s;
	`
	query := fmt.Sprintf(dropFMT, hdb.tableNames[huobiWithdrawalTableName])
	hdb.sugar.Debugw("tearingdown", "query", dropFMT, "table name", hdb.tableNames[huobiWithdrawalTableName])
	_, err := hdb.db.Exec(query)
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
	const updateStmt = `INSERT INTO %[1]s(id, data)
	VALUES ( 
		$1,
		$2
	)
	ON CONFLICT ON CONSTRAINT %[1]s_pk DO NOTHING;`
	query := fmt.Sprintf(updateStmt,
		hdb.tableNames[huobiWithdrawalTableName],
	)
	logger.Debugw("updating tradeHistory...", "query", query)

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
		_, err = tx.Exec(query, withdraw.ID, data)
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
	const selectStmt = `SELECT data FROM %[1]s WHERE data->>'created-at'>=$1 AND data->>'created-at'<$2`
	query := fmt.Sprintf(selectStmt, hdb.tableNames[huobiWithdrawalTableName])
	logger.Debugw("querying trade history...", "query", query)
	if err := hdb.db.Select(&dbResult, query, timeutil.TimeToTimestampMs(from), timeutil.TimeToTimestampMs(to)); err != nil {
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
	const selectStmt = `SELECT MAX(id) FROM %[1]s`
	query := fmt.Sprintf(selectStmt, hdb.tableNames[huobiWithdrawalTableName])
	logger.Debugw("querying trade history...", "query", query)

	if err := hdb.db.Get(&result, query); err != nil {
		return 0, err
	}

	return result, nil
}
