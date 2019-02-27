package storage

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
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
