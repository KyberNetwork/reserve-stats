package storage

import (
	"database/sql"
	"fmt"

	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	addressesTableName = "addresses"
	usersTableName     = "users"
)

//UserDB is storage of user data
type UserDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//DeleteAllTables delete all table from schema using for test only
func (udb *UserDB) DeleteAllTables() error {
	_, err := udb.db.Exec(fmt.Sprintf(`DROP TABLE "%s", "%s"`, addressesTableName, usersTableName))
	return err
}

//NewDB open a new database connection
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*UserDB, error) {
	const schemaFmt = `
CREATE TABLE IF NOT EXISTS "%s" (
  id    SERIAL PRIMARY KEY,
  email text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS "%s" (
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

	logger.Debug("initializing database schema")
	if _, err = tx.Exec(fmt.Sprintf(schemaFmt, usersTableName, addressesTableName)); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")

	if err = tx.Commit(); err != nil {
		return nil, err
	}

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
			"func", "users/storage.NewDB",
			"email", userData.Email,
		)
		userID int
	)

	tx, err := udb.db.Beginx()
	if err != nil {
		return err
	}

	err = tx.Get(&userID, fmt.Sprintf(`SELECT id FROM "%s" WHERE email = $1;`, usersTableName), userData.Email)
	if err == sql.ErrNoRows {
		logger.Debug("user does not exist, creating")
		row := tx.QueryRowx(`INSERT INTO users (email) VALUES ($1) RETURNING id;`, userData.Email)
		if err = row.Scan(&userID); err != nil {
			return err
		}
		logger = logger.With("user_id", userID)
		logger.Debug("user is created")
	} else if err != nil {
		return err
	} else {
		logger = logger.With("user_id", userID)
		logger.Debug("user already exists")
	}

	// client will submit all registered addresses every time
	_, err = tx.Exec(fmt.Sprintf(`
	DELETE FROM "%s" WHERE user_id = $1
`, addressesTableName), userID)
	if err != nil {
		return err
	}
	for _, info := range userData.UserInfo {
		logger.Debugw("updating user address",
			"address", info.Address,
			"timestamp", info.Timestamp,
		)
		_, err = tx.Exec(fmt.Sprintf(`
INSERT INTO "%s" (address, timestamp, user_id)
VALUES ($1, (TO_TIMESTAMP($2::double precision/1000)), $3);
`, addressesTableName),
			info.Address,
			info.Timestamp,
			userID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

//IsKYCed returns true when given address is found in database,
// means that user is already KYCed.
func (udb *UserDB) IsKYCed(address string) (bool, error) {
	var count int
	if err := udb.db.Get(&count, fmt.Sprintf(`SELECT COUNT(1) FROM "%s" WHERE address = $1`, addressesTableName), address); err != nil {
		return false, err
	}
	return count > 0, nil
}
