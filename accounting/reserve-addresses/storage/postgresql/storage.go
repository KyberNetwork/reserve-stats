package postgresql

import (
	"database/sql"
	"fmt"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/storage"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
)

const addressesTableName = "addresses"

// Storage implements accounting reserve addresses storage.Interface with PostgreSQL as storage engine.
type Storage struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

// NewStorage creates a new instance of Storage.
func NewStorage(sugar *zap.SugaredLogger, db *sqlx.DB) (*Storage, error) {
	var logger = sugar.With("func", "accounting/reserve-addresses/storage/postgresql.NewStorage")
	const schemaFmt = `CREATE TABLE IF NOT EXISTS "%s"
(
  id           SERIAL PRIMARY KEY,
  address      TEXT      NOT NULL UNIQUE,
  type         TEXT      NOT NULL,
  description  TEXT,
  timestamp    TIMESTAMP,
  last_updated TIMESTAMP NOT NULL
)`

	logger.Debugw("initializing database schema")
	if _, err := db.Exec(fmt.Sprintf(schemaFmt, addressesTableName)); err != nil {
		return nil, err
	}

	return &Storage{sugar: sugar, db: db}, nil
}

// Create creates a new address and store to database.
func (s *Storage) Create(address ethereum.Address, addressType common.AddressType, description string, ts time.Time) (uint64, error) {
	var (
		logger = s.sugar.With(
			"func", "accounting/reserve-addresses/storage/postgresql.Storage.Create",
			"address", address.String(),
			"type", addressType.String(),
			"description", description,
		)
		insertFmt = `INSERT INTO "%s" (address, type, description, timestamp, last_updated)
VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
		id uint64
	)
	logger.Debugw("creating new address")

	params := []interface{}{
		address.String(),
		addressType.String(),
		description,
	}

	if ts.IsZero() {
		params = append(params, nil)
	} else {
		params = append(params, ts)
	}

	err := s.db.Get(&id, fmt.Sprintf(insertFmt, addressesTableName), params...)
	if err != nil {
		// check if return error is a known pq error
		pErr, ok := err.(*pq.Error)
		if !ok {
			return 0, err
		}

		// https://www.postgresql.org/docs/9.3/errcodes-appendix.html
		// 23505: unique_violation
		if pErr.Code == "23505" {
			return 0, storage.ErrExists
		}

		return 0, pErr
	}
	return id, err
}

// Get returns the stored reserve address with matching id.
func (s *Storage) Get(id uint64) (*common.ReserveAddress, error) {
	var (
		logger = s.sugar.With("func", "accounting/reserve-addresses/storage/postgresql/Storage.Get",
			"id", id,
		)
		addr      = &ReserveAddress{}
		queryStmt = `SELECT id, address, type, description, timestamp
FROM addresses
WHERE id = $1`
	)

	if err := s.db.Get(addr, queryStmt, id); err == sql.ErrNoRows {
		logger.Infow("no record found in database")
		return nil, storage.ErrNotExists
	} else if err != nil {
		return nil, err
	}

	ra, err := addr.Common()
	if err != nil {
		return nil, err
	}
	return ra, nil
}
