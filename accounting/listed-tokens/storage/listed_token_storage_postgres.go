package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
)

//ListedTokenDB is storage for listed token
type ListedTokenDB struct {
	sugar     *zap.SugaredLogger
	db        *sqlx.DB
	tableName string
}

//NewDB open a new database connection an create initiated table if it is not exist
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB, tableName string) (*ListedTokenDB, error) {
	const schemaFmt = `CREATE TABLE IF NOT EXISTS "%[1]s"
(
	id SERIAL PRIMARY KEY,
	address text NOT NULL UNIQUE,
	name text NOT NULL,
	symbol text NOT NULL,
	timestamp TIMESTAMP NOT NULL,
	parent_id INT REFERENCES "%[1]s" (id)
)
	`
	var logger = sugar.With("func", "accounting/storage.NewDB")

	logger.Debug("initializing database schema")
	if _, err := db.Exec(fmt.Sprintf(schemaFmt, tableName)); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")

	return &ListedTokenDB{
		sugar:     sugar,
		db:        db,
		tableName: tableName,
	}, nil
}

//CreateOrUpdate add or edit an record in the tokens table
func (ltd *ListedTokenDB) CreateOrUpdate(tokens map[string]common.ListedToken) (err error) {
	var (
		logger = ltd.sugar.With("func", "accounting/lisetdtokenstorage.CreateOrUpdate")
	)
	upsertQuery := fmt.Sprintf(`INSERT INTO "%[1]s" (address, name, symbol, timestamp, parent_id)
	VALUES (
		$1, 
		$2, 
		$3,
		$4,
		(SELECT id FROM "%[1]s" WHERE address = $5)
	)
	ON CONFLICT (address) DO UPDATE SET parent_id = EXCLUDED.parent_id`,
		ltd.tableName)

	logger.Debugw("upsert token", "value", upsertQuery)

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
			token.Timestamp.UTC(),
			token.Address); err != nil {
			return
		}

		for _, oldToken := range token.Old {
			if _, err = tx.Exec(upsertQuery,
				oldToken.Address,
				token.Name,
				token.Symbol,
				oldToken.Timestamp.UTC(),
				token.Address); err != nil {
				return
			}
		}
	}

	return
}

// GetTokens return all tokens listed
func (ltd *ListedTokenDB) GetTokens() (map[string]common.ListedToken, error) {
	var (
		logger = ltd.sugar.With(
			"func",
			"accounting/listed-token-storage/listedtokenstorage.GetTokens",
		)
		result       []common.ListedToken
		listedTokens = make(map[string]common.ListedToken)
	)

	getQuery := fmt.Sprintf(`SELECT address, name, symbol, timestamp FROM %[1]s ORDER BY timestamp DESC`, ltd.tableName)
	logger.Debugw("get tokens query", "query", getQuery)

	if err := ltd.db.Select(&result, getQuery); err != nil {
		logger.Errorw("error query token", "error", err)
	}

	logger.Debugw("result from listed token", "result", result)

	// Assemble old tokens
	for _, token := range result {
		key := fmt.Sprintf("%s-%s", token.Symbol, token.Name)
		if listedToken, exist := listedTokens[key]; exist {
			if token.Timestamp.Before(listedToken.Timestamp) {
				listedToken.Old = append(listedToken.Old, common.OldListedToken{
					Address:   token.Address,
					Timestamp: token.Timestamp.UTC(),
				})
				listedTokens[key] = listedToken
			} else {
				listedToken.Old = append(listedToken.Old, common.OldListedToken{
					Address:   listedToken.Address,
					Timestamp: listedToken.Timestamp.UTC(),
				})
				listedToken.Address = token.Address
				listedToken.Timestamp = token.Timestamp.UTC()
				listedTokens[key] = listedToken
			}
		} else {
			listedTokens[key] = common.ListedToken{
				Address:   token.Address,
				Name:      token.Name,
				Symbol:    token.Symbol,
				Timestamp: token.Timestamp.UTC(),
			}
		}
	}

	return listedTokens, nil
}

//Close db connection
func (ltd *ListedTokenDB) Close() error {
	if ltd.db != nil {
		return ltd.db.Close()
	}
	return nil
}

//DeleteTable remove tables use for test
func (ltd *ListedTokenDB) DeleteTable() error {
	const dropQuery = `DROP TABLE %s;`
	query := fmt.Sprintf(dropQuery, ltd.tableName)

	ltd.sugar.Infow("Drop token table", "query", query)
	_, err := ltd.db.Exec(query)
	return err
}
