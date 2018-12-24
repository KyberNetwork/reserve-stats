package migration

import (
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/boltutil"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/boltdb/bolt"
	"github.com/jmoiron/sqlx"
)

const (
	addressTime = "address_time"
	kyced       = "kyced"
	//DefaultAddressesTableName is the default address table name for migration
	DefaultAddressesTableName = "addresses"
	//DefaultUsersTableName is the default user table name for migration
	DefaultUsersTableName = "users"
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
	var schema = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS "%[1]s" (
  id    SERIAL PRIMARY KEY,
  email text NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS "%[2]s" (
  id        SERIAL PRIMARY KEY,
  address   text      NOT NULL UNIQUE,
  timestamp TIMESTAMP NOT NULL DEFAULT now(),
  user_id   SERIAL    NOT NULL REFERENCES "%[1]s" (id)
);
`, DefaultUsersTableName, DefaultAddressesTableName)

	tx, err := postgres.Beginx()
	if err != nil {
		return nil, err
	}

	defer pgsql.CommitOrRollback(tx, sugar, &err)

	if _, err = tx.Exec(schema); err != nil {
		return nil, err
	}

	storage := &DBMigration{sugar: sugar, boltdb: boltDB, postgresdb: postgres}
	return storage, err
}

//Migrate read data from bolt database file and input into postgres database
func (dbm *DBMigration) Migrate(usersTableName, addressTableName string) error {
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
		if err := usersBucket.ForEach(func(k, v []byte) error {
			count++
			var (
				email   = string(v)
				address = string(k)
				logger  = logger.With("email", email, "address", address, "count", count)
			)

			logger.Debugw("inserting to users table")
			userID, eErr := dbm.insertIntoUserTable(usersTableName, email)
			if eErr != nil {
				return eErr
			}

			timestampByte := addressesBucket.Get(k)
			timestamp := boltutil.BytesToUint64(timestampByte)

			logger.Infow("inserting to addresses table")
			return dbm.insertAddress(addressTableName, address, timestamp, userID)
		}); err != nil {
			return err
		}

		logger.Infow("migration completed", "count", count)
		return nil
	})
}

//insertIntoUserTable insert email into
func (dbm *DBMigration) insertIntoUserTable(tableName string, email string) (int, error) {
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

	defer pgsql.CommitOrRollback(ptx, logger, &err)

	err = ptx.Get(&userID, fmt.Sprintf(`SELECT id FROM "%s" WHERE email = $1;`, tableName), email)
	if err == sql.ErrNoRows {
		logger.Debugw("user does not exist, creating")
		row := ptx.QueryRowx(`INSERT INTO users (email) VALUES ($1) RETURNING id;`, email)
		if err = row.Scan(&userID); err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	} else if err == nil {
		logger.Debugw("user already exists, skipping")
	}
	logger.Debugw("user already exists, skipping")

	return userID, nil
}

//insertAddress insert address into address table
func (dbm *DBMigration) insertAddress(tableName, address string, timestamp uint64, userID int) error {
	ptx, err := dbm.postgresdb.Beginx()

	defer pgsql.CommitOrRollback(ptx, dbm.sugar, &err)

	var (
		logger = dbm.sugar.With(
			"func", "users/migration/DBMigration.insertAddress",
			"address", address,
		)
		userIDFromDB int
	)

	err = ptx.Get(&userIDFromDB, fmt.Sprintf(`SELECT user_id FROM "%s" WHERE address = $1;`, tableName), address)
	if err == nil {
		if userID != userIDFromDB {
			logger.Debugw("address already belong to a different user", "userID input", userID, "userID from DB", userIDFromDB)
			if err = ptx.Commit(); err != nil {
				return err
			}
			return errors.New("address is already existed with a different userID")
		}
		logger.Debugw("address already belong to the same user, skipping")
		return ptx.Commit()
	}
	_, err = ptx.Exec(fmt.Sprintf(`
INSERT INTO "%s" (address, timestamp, user_id)
VALUES ($1, (TO_TIMESTAMP($2::double precision/1000)), $3);
`, tableName),
		address,
		timestamp,
		userID)
	return err
}
