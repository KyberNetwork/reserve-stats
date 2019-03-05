package storage

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
)

const (
	tokenTable = "tokens"
)

//ListedTokenDB is storage for listed token
type ListedTokenDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//NewDB open a new database connection an create initiated table if it is not exist
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*ListedTokenDB, error) {
	const schemaFmt = `CREATE TABLE IF NOT EXIST "%s"
(
	id SERIAL PRIMARY KEY,
	address text NOT NULL UNIQUE,
	symbol text NOT NULL UNIQUE,
	timestamp TIMESTAMP NOT NULL,
	parent_id SERIAL REFERENCES "%s" (id)
)
	`
	var logger = sugar.With("func", "accounting/storage.NewDB")

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	defer pgsql.CommitOrRollback(tx, logger, &err)

	logger.Debug("initializing database schema")
	if _, err = tx.Exec(fmt.Sprintf(schemaFmt, tokenTable, tokenTable)); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")

	return &ListedTokenDB{
		sugar: sugar,
		db:    db,
	}, nil
}

//CreateOrUpdate add or edit an record in the tokens table
func (ltd *ListedTokenDB) CreateOrUpdate(tokens []common.ListedToken) error {
	var (
		logger = ltd.sugar.With("func", "accounting/storage.CreateOrUpdate")
	)
	logger.Info("create or upate token listed")
	tx, err := ltd.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)
	for _, token := range tokens {
		// check if token already exist in db

		// insert token into table
		query := fmt.Sprintf(`INSERT INTO "%s" (address,symbol, timestamp) VALUES($1, $2, $3)`, tokenTable)
		logger.Debugw("Insert token into table", "query", query)
		_, err = tx.Exec(query, token.Address, token.Symbol, token.Timestamp)
		if err != nil {
			return nil
		}
	}
	return nil
}
