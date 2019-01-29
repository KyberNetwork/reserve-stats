package migration

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/boltutil"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
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
  email text NOT NULL UNIQUE,
  last_updated TIMESTAMP NOT NULL
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

//InsertOrUpdate insert or update userData from a record
func (dbm *DBMigration) InsertOrUpdate(address, email, usersTableName, addressTableName string, timestamp uint64) error {
	var (
		logger = dbm.sugar.With(
			"func", "users/migrate.InserOrUpdate",
			"email", email,
		)
		stmt string
	)

	tx, err := dbm.postgresdb.Beginx()
	if err != nil {
		return err
	}

	defer pgsql.CommitOrRollback(tx, logger, &err)
	timepoint := timeutil.TimestampMsToTime(timestamp)
	stmt = fmt.Sprintf(`WITH u AS (
  INSERT INTO "%s" (email, last_updated)
    VALUES ($1, $2)
    ON CONFLICT ON CONSTRAINT users_email_key
      DO UPDATE SET last_updated = $2 RETURNING id
),
     a AS (
       SELECT $3::text            AS address,
              $2 AS timestamp
     )
INSERT
INTO "%s"(address, timestamp, user_id)
SELECT a.address, a.timestamp, u.id
FROM u NATURAL JOIN a
ON CONFLICT ON CONSTRAINT addresses_address_key DO UPDATE SET timestamp = EXCLUDED.timestamp, user_id = EXCLUDED.user_id 
`,
		usersTableName,
		addressTableName)
	logger.Debugw("upsert email and Ethereum addresses",
		"stmt", stmt)
	_, err = tx.Exec(stmt,
		email,
		timepoint,
		address,
	)
	return err
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

			timestampByte := addressesBucket.Get(k)
			timestamp := boltutil.BytesToUint64(timestampByte)
			logger.Debugw("inserting to users table")
			return dbm.InsertOrUpdate(address, email, usersTableName, addressTableName, timestamp)
		}); err != nil {
			return err
		}

		logger.Infow("migration completed", "count", count)
		return nil
	})
}
