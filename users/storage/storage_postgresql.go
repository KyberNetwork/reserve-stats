package storage

import (
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/users/common"
)

//UserDB is storage of user data
type UserDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//NewDB open a new database connection
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*UserDB, error) {
	const schema = `CREATE TABLE IF NOT EXISTS "users"
(
  id           SERIAL PRIMARY KEY,
  email        text      NOT NULL UNIQUE,
  last_updated TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS "addresses"
(
  id        SERIAL PRIMARY KEY,
  address   text      NOT NULL UNIQUE,
  timestamp TIMESTAMP NOT NULL,
  user_id   SERIAL    NOT NULL REFERENCES users (id)
);
`
	var logger = sugar.With("func", "users/storage.NewDB")

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	defer pgsql.CommitOrRollback(tx, logger, &err)

	logger.Debug("initializing database schema")
	if _, err = tx.Exec(schema); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")

	return &UserDB{
		sugar: sugar,
		db:    db,
	}, nil
}

//Close close db connection and return error if any
func (udb *UserDB) Close() error {
	return udb.db.Close()
}

//CreateOrUpdate store user info to persist in database
func (udb *UserDB) CreateOrUpdate(userData common.UserData) error {
	var (
		logger = udb.sugar.With(
			"func", "users/storage.CreateOrUpdate",
			"email", userData.Email,
		)
		stmt       string
		addresses  []string
		timestamps []int64
	)

	for _, ui := range userData.UserInfo {
		address := ethereum.HexToAddress(ui.Address)
		addresses = append(addresses, address.Hex())
	}
	for _, ui := range userData.UserInfo {
		timestamps = append(timestamps, ui.Timestamp)
	}

	tx, err := udb.db.Beginx()
	if err != nil {
		return err
	}

	defer pgsql.CommitOrRollback(tx, logger, &err)

	stmt = `WITH u AS (
    INSERT INTO "users" (email, last_updated)
        VALUES ($1, NOW())
        ON CONFLICT ON CONSTRAINT users_email_key
            DO UPDATE SET last_updated = NOW() RETURNING id
),
     a AS (
         SELECT unnest($2::text[])             AS address,
                unnest($3::double precision[]) AS timestamp
     )
INSERT
INTO "addresses"(address, timestamp, user_id)
SELECT a.address, to_timestamp(a.timestamp / 1000), u.id
FROM u
         NATURAL JOIN a
ON CONFLICT ON CONSTRAINT addresses_address_key DO UPDATE SET timestamp = EXCLUDED.timestamp,
                                                              user_id   = EXCLUDED.user_id
`
	logger.Debugw("upsert email and Ethereum addresses",
		"stmt", stmt)
	_, err = tx.Exec(stmt,
		userData.Email,
		pq.StringArray(addresses),
		pq.Int64Array(timestamps),
	)
	if err != nil {
		return err
	}

	logger.Debugw("delete removed Ethereum addresses",
		"stmt", stmt)
	stmt = `DELETE
FROM "addresses"
WHERE user_id IN (SELECT id AS user_id FROM "users" WHERE email = $1)
  AND address NOT IN (SELECT unnest($2::text[]) as address)
`
	_, err = tx.Exec(stmt,
		userData.Email,
		pq.StringArray(addresses))
	if err != nil {
		return err
	}
	return err
}

//GetAllAddresses return all user address info from addresses table
func (udb *UserDB) GetAllAddresses() ([]string, error) {
	var result []string
	if err := udb.db.Select(&result, `SELECT address FROM "addresses"`); err != nil {
		return result, err
	}
	return result, nil
}

//IsKYCedAtTime returns true when given address is found in database before a given timestamp,
// means that user is already KYCed before that time
func (udb *UserDB) IsKYCedAtTime(userAddr string, ts time.Time) (bool, error) {
	var (
		logger = udb.sugar.With(
			"func", "users/storage/UserDB.IsKYCedAtTime",
			"user_addr", userAddr,
			"timestamp", ts.String(),
		)
		result uint64
	)
	stmt := `SELECT COUNT(1)
FROM "addresses"
WHERE address = $1
  AND timestamp <= $2
`
	logger = logger.With("query", stmt)
	if err := udb.db.Get(&result, stmt, userAddr, ts.UTC()); err != nil {
		return false, err
	}
	logger.Debugw("got result from database", "result", result)
	return result != 0, nil
}
