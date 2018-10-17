package migratedb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/app"
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
	boltdb     *bolt.DB
	postgresdb *sqlx.DB
}

//DeleteAllTables delete all table from schema using for test only
func (dbm *DBMigration) DeleteAllTables() error {
	_, err := dbm.postgresdb.Exec(fmt.Sprintf(`DROP TABLE "%s", "%s"`, addressesTableName, usersTableName))
	return err
}

//DeleteTable delete a table from  schema
func (dbm *DBMigration) DeleteTable(tableName string) error {
	_, err := dbm.postgresdb.Exec(fmt.Sprintf(`DROP TABLE "%s"`, tableName))
	return err
}

//NewMigrateStorage connect to bolt db
func NewMigrateStorage(dbPath string, postgres *sqlx.DB) (*DBMigration, error) {
	var err error
	var boltDB *bolt.DB
	// open bolt db for migration
	boltDB, err = bolt.Open(dbPath, 0600, nil)
	if boltDB == nil {
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
	storage := &DBMigration{boltDB, postgres}
	return storage, err
}

//MigrateDB read data from bolt database file and input into postgres database
func (dbm *DBMigration) MigrateDB() error {
	var err error
	err = dbm.boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(kyced))
		if b == nil {
			return errors.New("bolt db is empty")
		}
		timeBucket := tx.Bucket([]byte(addressTime))
		// migrate user and address
		b.ForEach(func(k, v []byte) error {
			email := string(v)
			address := string(k)
			userID, err := dbm.InsertIntoUserTable(email)
			if err != nil {
				return err
			}
			timestampByte := timeBucket.Get(k)
			timestamp := app.BytesToUint64(timestampByte)
			if err := dbm.InsertAddress(address, timestamp, userID); err != nil {
				return err
			}
			return nil
		})
		return nil
	})
	return err
}

//InsertIntoUserTable insert email into
func (dbm *DBMigration) InsertIntoUserTable(email string) (int, error) {
	var userID int
	ptx, err := dbm.postgresdb.Beginx()
	if err != nil {
		return 0, err
	}
	err = ptx.Get(&userID, fmt.Sprintf(`SELECT id FROM "%s" WHERE email = $1;`, usersTableName), email)
	if err == sql.ErrNoRows {
		row := ptx.QueryRowx(`INSERT INTO users (email) VALUES ($1) RETURNING id;`, email)
		if err = row.Scan(&userID); err != nil {
			return 0, err
		}
		if err = ptx.Commit(); err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}
	return userID, nil
}

//InsertAddress insert address into address table
func (dbm *DBMigration) InsertAddress(address string, timestamp uint64, userID int) error {
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
