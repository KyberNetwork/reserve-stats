package postgresql

import (
	"database/sql"
	"fmt"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/storage"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
)

const (
	addressesTableName      = "addresses"
	addressVersionTableName = "addresses_version"
)

// Storage implements accounting reserve addresses storage.Interface with PostgreSQL as storage engine.
type Storage struct {
	sugar  *zap.SugaredLogger
	resolv blockchain.ContractTimestampResolver
	db     *sqlx.DB
}

// NewStorage creates a new instance of Storage.
func NewStorage(sugar *zap.SugaredLogger, db *sqlx.DB, resolv blockchain.ContractTimestampResolver) (*Storage, error) {
	var logger = sugar.With("func", "accounting/reserve-addresses/storage/postgresql.NewStorage")
	const schemaFmt = `CREATE TABLE IF NOT EXISTS "%[1]s"
(
  id           SERIAL PRIMARY KEY,
  address      TEXT      NOT NULL UNIQUE,
  type         TEXT      NOT NULL,
  description  TEXT,
  timestamp    TIMESTAMP,
  last_updated TIMESTAMP NOT NULL
);
--create version table
CREATE TABLE IF NOT EXISTS "%[2]s"
(
	id SERIAL PRIMARY KEY,
	version integer NOT NULL,
	timestamp TIMESTAMP
);
--create trigger function
CREATE OR REPLACE FUNCTION inc_version() RETURNS TRIGGER AS
$$
DECLARE
    inc BOOLEAN = false;
BEGIN
    IF tg_op = 'INSERT' THEN
        inc = TRUE;
	END IF;
	IF tg_op = 'UPDATE' AND OLD IS DISTINCT FROM NEW THEN
		inc = TRUE;
	END IF;
	IF inc THEN
		INSERT INTO "%[2]s" (id, version, timestamp)
		VALUES (1, 1, now()) ON CONFLICT (id) DO UPDATE SET version = %[2]s.version+1, timestamp = EXCLUDED.timestamp;
	END IF;
	RETURN NULL;
END;
$$ LANGUAGE PLPGSQL;
DROP TRIGGER IF EXISTS version_trigger ON %[1]s;
CREATE TRIGGER version_trigger
    AFTER INSERT OR UPDATE
    ON %[1]s 
    FOR EACH ROW
EXECUTE PROCEDURE inc_version();`

	logger.Debugw("initializing database schema")
	if _, err := db.Exec(fmt.Sprintf(schemaFmt, addressesTableName, addressVersionTableName)); err != nil {
		return nil, err
	}

	return &Storage{sugar: sugar, db: db, resolv: resolv}, nil
}

// Create creates a new address and store to database.
func (s *Storage) Create(address ethereum.Address, addressType common.AddressType, description string) (uint64, error) {
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

	ts, err := s.resolv.Resolve(address)
	switch err {
	case blockchain.ErrNotAvailable:
		logger.Debugw("address contract creation time is not available", "err", err)
		params = append(params, nil)
	case nil:
		params = append(params, ts.UTC())
	default:
		return 0, err
	}

	if err = s.db.Get(&id, fmt.Sprintf(insertFmt, addressesTableName), params...); err != nil {
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

// GetAll returns all stored reserve addresses in database.
// It returns no error if there is nothing in database.
func (s *Storage) GetAll() ([]*common.ReserveAddress, int64, error) {
	var (
		logger    = s.sugar.With("func", "accounting/reserve-addresses/storage/postgresql/Storage.GetAll")
		stored    []*ReserveAddress
		results   []*common.ReserveAddress
		queryStmt = `SELECT id, address, type, description, timestamp
FROM addresses`
		queryVersionStmt = fmt.Sprintf(`SELECT version FROM %[1]s WHERE id = 1`, addressVersionTableName)
		version          int64
	)

	logger.Debug("querying all stored reserve addresses")
	if err := s.db.Select(&stored, queryStmt); err != nil {
		return nil, 0, err
	}

	for _, r := range stored {
		result, err := r.Common()
		if err != nil {
			return nil, 0, err
		}
		results = append(results, result)
	}
	if len(results) > 0 {
		if err := s.db.Get(&version, queryVersionStmt); err != nil {
			return nil, 0, err
		}
	}
	return results, version, nil
}

// Update updates the reserve address with given information. If given data is zero, it won't be updated to database.
func (s *Storage) Update(id uint64, address ethereum.Address, addressType *common.AddressType, description string) error {
	var (
		logger = s.sugar.With("func", "accounting/reserve-addresses/storage/postgresql/Storage.Update",
			"id", id,
			"address", address.String(),
			"description", description,
		)
		queryStmt = `UPDATE addresses
SET address     = COALESCE($1, address),
    type        = COALESCE($2, type),
    description = COALESCE($3, description),
    timestamp   = $4
WHERE id = $5 RETURNING id;`
		params []interface{}
	)

	// fill address value
	if !blockchain.IsZeroAddress(address) {
		params = append(params, address.String())
	} else {
		params = append(params, nil)
	}

	// fill type value
	if addressType != nil {
		logger = logger.With("type", addressType.String())
		params = append(params, addressType.String())
	} else {
		params = append(params, nil)
	}

	// fill description value
	if len(description) != 0 {
		params = append(params, description)
	} else {
		params = append(params, nil)
	}

	// fill timestamp value
	ts, err := s.resolv.Resolve(address)
	switch err {
	case blockchain.ErrNotAvailable:
		logger.Debugw("address contract creation time is not available", "err", err)
		params = append(params, nil)
	case nil:
		params = append(params, ts.UTC())
	default:
		return err
	}

	// fill query condition param
	params = append(params, id)

	logger.Debug("updating reserve address record in database")
	var updatedID uint64
	if err = s.db.Get(&updatedID, queryStmt, params...); err == sql.ErrNoRows {
		return storage.ErrNotExists
	} else if err != nil {
		return err
	}

	return nil
}
