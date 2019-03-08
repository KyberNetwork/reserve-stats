package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	huobiWithdrawalTableName = "huobi_withdrawal"
)

//HuobiStorage defines the object to store Huobi data
type HuobiStorage struct {
	sugar      *zap.SugaredLogger
	db         *sqlx.DB
	tableNames map[string]string
}

// NewDB return the HuobiStorage instance. User must call Close() before exit.
// tableNames is a list of string for tablename[huobitrades]. It can be optional
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB, customTableNames ...string) (*HuobiStorage, error) {
	const schemaFMT = `
	CREATE TABLE IF NOT EXISTS %[1]s
(
	id bigint NOT NULL,
	time TIMESTAMP NOT NULL,
	currency TEXT,
	data JSONB,
	CONSTRAINT %[1]s_pk PRIMARY KEY(id)
) ;
`
	var (
		logger     = sugar.With("func", "reserverates/storage/postgres/NewDB")
		tableNames = make(map[string]string)
	)
	if len(customTableNames) > 0 {
		if len(customTableNames) != 1 {
			return nil, fmt.Errorf("expect 1 tables name [trades], got %v", customTableNames)
		}
		tableNames[huobiWithdrawalTableName] = customTableNames[0]

	} else {
		tableNames[huobiWithdrawalTableName] = huobiWithdrawalTableName
	}

	query := fmt.Sprintf(schemaFMT, tableNames[huobiWithdrawalTableName])
	logger.Debugw("initializing database schema", "query", query)

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")
	return &HuobiStorage{
		sugar:      sugar,
		db:         db,
		tableNames: tableNames,
	}, nil
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
func (hdb *HuobiStorage) UpdateWithdrawHistory(withdraw huobi.WithdrawHistory) error {
	var (
		timestamp = timeutil.TimestampMsToTime(withdraw.CreatedAt)
		logger    = hdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
			"withdrawal_ID", withdraw.ID,
			"timestamp", timestamp,
		)
	)
	const updateStmt = `INSERT INTO %[1]s(id,time, currency, data)
	VALUES ( 
		$1,
		$2, 
		$3,
		$4
	)
	ON CONFLICT ON CONSTRAINT %[1]s_pk DO NOTHING;`
	query := fmt.Sprintf(updateStmt,
		hdb.tableNames[huobiWithdrawalTableName],
	)
	data, err := json.Marshal(withdraw)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	logger.Debugw("updating tradeHistory...", "query", query)
	_, err = hdb.db.Exec(query,
		withdraw.ID,
		timestamp,
		withdraw.Currency,
		data,
	)
	if err != nil {
		return err
	}
	return err
}
