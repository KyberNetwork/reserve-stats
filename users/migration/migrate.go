package migration

import (
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/boltutil"
	"github.com/boltdb/bolt"
	"github.com/jmoiron/sqlx"
)

const (
	addressTime        = "address_time"
	kyced              = "kyced"
	addressesTableName = "addresses"
	usersTableName     = "users"
)

//DBMigration user storage in bolt
type DBMigration struct {
	sugar *zap.SugaredLogger

	boltdb     *bolt.DB
	postgresdb *sqlx.DB
}

//NewDBMigration connect to bolt db
func NewDBMigration(sugar *zap.SugaredLogger, dbPath string, postgres *sqlx.DB) (*DBMigration, error) {
	var (
		err    error
		boltDB *bolt.DB
	)
	// open bolt db for migration
	boltDB, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	// initiate postgres for migration
	const schema = `
CREATE TABLE IF NOT EXISTS "users" (
  id    SERIAL PRIMARY KEY,
  email text NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS "addresses" (
  id        SERIAL PRIMARY KEY,
  address   text      NOT NULL UNIQUE,
  timestamp TIMESTAMP NOT NULL DEFAULT now(),
  user_id   SERIAL    NOT NULL REFERENCES users (id)
);
`

	tx, err := postgres.Beginx()
	if err != nil {
		return nil, err
	}

	if _, err = tx.Exec(schema); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	storage := &DBMigration{sugar: sugar, boltdb: boltDB, postgresdb: postgres}
	return storage, err
}

//Migrate read data from bolt database file and input into postgres database
func (dbm *DBMigration) Migrate() error {
	var logger = dbm.sugar.With("func", "users/migration/DBMigration.Migrate")
	return dbm.boltdb.View(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte(kyced))
		if usersBucket == nil {
			return errors.New("users bucket is empty")
		}

		addressesBucket := tx.Bucket([]byte(addressTime))
		if addressesBucket == nil {
			return errors.New("addresses bucket is empty")
		}

		logger.Debugw("starting migration process from legacy BoltDB to PostgreSQL")
		var count = 0
		usersBucket.ForEach(func(k, v []byte) error {
			count++
			var (
				email   = string(v)
				address = string(k)
				logger  = logger.With("email", email, "address", address, "count", count)
			)

			logger.Debugw("inserting to users table")
			userID, eErr := dbm.insertIntoUserTable(email)
			if eErr != nil {
				return eErr
			}

			timestampByte := addressesBucket.Get(k)
			timestamp := boltutil.BytesToUint64(timestampByte)

			logger.Infow("inserting to addresses table")
			return dbm.insertAddress(address, timestamp, userID)
		})

		logger.Infow("migration completed", "count", count)
		return nil
	})
}

//insertIntoUserTable insert email into
func (dbm *DBMigration) insertIntoUserTable(email string) (int, error) {
	var (
		logger = dbm.sugar.With(
			"func", "users/migration/DBMigration.insertIntoUserTable",
			"email", email,
		)
		userID int
	)
	ptx, err := dbm.postgresdb.Beginx()
	if err != nil {
		return 0, err
	}

	err = ptx.Get(&userID, fmt.Sprintf(`SELECT id FROM "%s" WHERE email = $1;`, usersTableName), email)
	if err == sql.ErrNoRows {
		logger.Debugw("user does not exist, creating")
		row := ptx.QueryRowx(`INSERT INTO users (email) VALUES ($1) RETURNING id;`, email)
		if err = row.Scan(&userID); err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	logger.Debugw("user already exists, skipping")

	if err = ptx.Commit(); err != nil {
		return 0, err
	}
	return userID, nil
}

//insertAddress insert address into address table
func (dbm *DBMigration) insertAddress(address string, timestamp uint64, userID int) error {
	ptx, err := dbm.postgresdb.Beginx()
	_, err = ptx.Exec(fmt.Sprintf(`
INSERT INTO "%s" (address, timestamp, user_id)
VALUES ($1, (TO_TIMESTAMP($2::double precision/1000)), $3);
`, addressesTableName),
		address,
		timestamp,
		userID)
	if err != nil {
		return err
	}
	if err = ptx.Commit(); err != nil {
		return err
	}
	return nil
}
