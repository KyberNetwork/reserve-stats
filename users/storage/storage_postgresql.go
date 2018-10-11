package storage

import (
	"database/sql"
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

//UserDB is storage of user data
type UserDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//DeleteAllTables delete all table from schema using for test only
func (udb *UserDB) DeleteAllTables() error {
	_, err := udb.db.Exec(`DROP TABLE "users", "addresses"`)
	return err
}

//NewDB open a new database connection
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*UserDB, error) {
	const schema = `
CREATE TABLE IF NOT EXISTS "users" (
  id    SERIAL PRIMARY KEY,
  email text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS "addresses" (
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
	if _, err = tx.Exec(schema); err != nil {
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

//StoreUserInfo store user info to persist in database
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

	err = tx.Get(&userID, `SELECT id FROM "users" WHERE email = $1;`, userData.Email)
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

	for _, info := range userData.UserInfo {
		logger.Debugw("updating user address",
			"address", info.Address,
			"timestamp", info.Timestamp,
		)
		_, err = tx.Exec(`
INSERT INTO "addresses" (address, timestamp, user_id)
VALUES ($1, (TO_TIMESTAMP($2::double precision/1000)), $3);
`,
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
	if err := udb.db.Get(&count, `SELECT COUNT(1) FROM "addresses" WHERE address = $1`, address); err != nil {
		return false, err
	}
	return count > 0, nil
}
