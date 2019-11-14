package postgres

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
)

var (
	// ErrExists return when data already exist
	ErrExists = errors.New("already exist")
	// ErrNotFound return when data not found
	ErrNotFound = errors.New("not found")
)

// TokenPriceDB is storage of token price
type TokenPriceDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

// NewTokenPriceDB return instance of TokenPriceDB
func NewTokenPriceDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*TokenPriceDB, error) {
	if err := initTable(db); err != nil {
		return nil, err
	}
	return &TokenPriceDB{
		sugar: sugar,
		db:    db,
	}, nil
}

func initTable(db *sqlx.DB) error {
	const (
		tokenPricesSchema = `
			CREATE TABLE IF NOT EXISTS "tokenprices" (
				date DATE,
				source TEXT NOT NULL,
				token TEXT NOT NULL,
				currency TEXT NOT NULL,
				value FLOAT(32) NOT NULL,
				PRIMARY KEY (date, source, token, currency)
			);
		`
	)
	if _, err := db.Exec(tokenPricesSchema); err != nil {
		return err
	}
	return nil
}

// SaveTokenPrice save token price data
func (trdb *TokenPriceDB) SaveTokenPrice(token, currency, source string, timestamp time.Time, price float64) error {
	var (
		logger = trdb.sugar.With("func", caller.GetCurrentFunctionName(),
			"date", timestamp,
			"token", token,
			"source", source,
			"currency", currency,
		)
		query = `INSERT INTO "tokenprices"(date, source, token, currency, value) 
		VALUES (DATE($1), $2, $3, $4, $5)`
	)
	logger.Infow("save new token price", "query", query)
	log.Println(query)
	_, err := trdb.db.Exec(query, timestamp, source, token, currency, price)
	if err != nil {
		// check if return error is a known pq error
		pErr, ok := err.(*pq.Error)
		if !ok {
			return errors.Wrap(err, "failed to save token price to db")
		}

		logger.Errorw("got error from database",
			"code", pErr.Code, "message", pErr.Message)

		// https://www.postgresql.org/docs/9.3/errcodes-appendix.html
		// 23505: unique_violation
		if pErr.Code == "23505" {
			return ErrExists
		}

		return errors.Wrap(err, "failed to store token price to database")
	}
	return nil
}

type tokenPriceDB struct {
	Price sql.NullFloat64 `db:"value"`
}

// GetTokenPrice save token price data
func (trdb *TokenPriceDB) GetTokenPrice(token, currency string, timestamp time.Time) (float64, error) {
	var (
		logger = trdb.sugar.With("func", caller.GetCurrentFunctionName(),
			"date", timestamp,
			"token", token,
			"currency", currency,
		)
		query = `SELECT value FROM "tokenprices" 
			WHERE token=$1 AND currency=$2 AND date=DATE($3)`

		dbResult tokenPriceDB
	)
	logger.Infow("get token price", "query", query)
	if err := trdb.db.Get(&dbResult, query, token, currency, timestamp); err == sql.ErrNoRows {
		return 0, ErrNotFound
	} else if err != nil {
		logger.Errorw("got error from database", "error", err)
		return 0, errors.New("failed to query token price in database")
	}
	return dbResult.Price.Float64, nil
}
