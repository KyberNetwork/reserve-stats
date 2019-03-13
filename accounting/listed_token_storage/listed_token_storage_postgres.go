package listedtokenstorage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
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
	const schemaFmt = `CREATE TABLE IF NOT EXISTS "%s"
(
	id SERIAL PRIMARY KEY,
	address text NOT NULL UNIQUE,
	name text NOT NULL,
	symbol text NOT NULL,
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
func (ltd *ListedTokenDB) CreateOrUpdate(tokens map[string]common.ListedToken) error {
	var (
		logger = ltd.sugar.With("func", "accounting/lisetdtokenstorage.CreateOrUpdate")
	)
	upsertQuery := fmt.Sprintf(`INSERT INTO "%s" (address, name, symbol, timestamp)
	VALUES (
		$1, 
		$2, 
		$3,
		to_timestamp($4::double precision / 100)
	)
	ON CONFLICT (address) DO NOTHING`,
		tokenTable)

	upsertOldTokenQuery := fmt.Sprintf(`INSERT INTO "%[1]s" (address, name, symbol, timestamp, parent_id)
	VALUES (
		$1, 
		$2, 
		$3,
		to_timestamp($4::double precision / 1000),
		SELECT id FROM "%[1]s" WHERE symbol = $2
	)
	ON CONFLICT (address) DO NOTHING`,
		tokenTable)

	tx, err := ltd.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)

	for _, token := range tokens {
		if _, err = tx.Exec(upsertQuery,
			token.Address,
			token.Name,
			token.Symbol,
			token.Timestamp); err != nil {
			return err
		}

		if len(token.Old) != 0 {
			for _, oldToken := range token.Old {
				if _, err = tx.Exec(upsertOldTokenQuery,
					oldToken.Address,
					token.Name,
					token.Symbol,
					oldToken.Timestamp); err != nil {
					return err
				}
			}
		}
	}

	return err
}
