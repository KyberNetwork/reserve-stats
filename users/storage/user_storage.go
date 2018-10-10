package storage

import (
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

//UserDB is storage of user data
type UserDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB

	stms []*sqlx.Stmt
}

//DeleteAllTables delete all table from schema using for test only
func (udb *UserDB) DeleteAllTables() error {
	//userTable := &common.User{}
	//if err := udb.db.DropTable(userTable, &orm.DropTableOptions{}); err != nil {
	//	return err
	//}
	return nil
}

//NewDB open a new database connection
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*UserDB, error) {
	const schema = `
CREATE TABLE IF NOT EXISTS users (
  id    SERIAL PRIMARY KEY,
  email text   NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS addresses (
  id      SERIAL PRIMARY KEY,
  address text   NOT NULL UNIQUE,
  user_id SERIAL NOT NULL REFERENCES users (id)
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

	tx.Prepare(`INSERT INTO users (email) VALUES (:email) RETURNING id`)

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
func (udb *UserDB) StoreUserInfo(userData common.UserData) error {
	var logger = udb.sugar.With(
		"func", "users/storage.NewDB",
		"email", userData.Email,
	)

	tx, err := udb.db.Beginx()
	if err != nil {
		logger.Errorw("failed to acquire transaction",
			"err", err)
		return err
	}

	rows, err := tx.Prepare(`INSERT INTO users (email) VALUES (:email) RETURNING id`, userData)
	if err != nil {
		return err
	}

	//for rows.Next() {
	//	var user_id int
	//	rows.Scan(&user_id)
	//	for _, userInfo := range userData.UserInfo {
	//		_, err = tx.NamedQuery(`INSERT INTO addresses (address, user_id) VALUES (:address, :user_id)`,
	//			userInfo.Address,
	//			user_id,
	//			// TODO: add missing timestamp field
	//		)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}

	return tx.Commit()

	//userModel := common.User{}
	////remove all old user
	//if _, err := udb.db.Model(&userModel).Where("email = ?", userData.Email).Delete(); err != nil {
	//	return err
	//}
	//
	//// insert updated address value
	//for _, address := range userData.UserInfo {
	//	userEntry := common.User{
	//		Email:     userData.Email,
	//		Address:   address.Address,
	//		Timestamp: time.Unix(address.Timestamp/1000, 0),
	//	}
	//	if err := udb.db.Insert(&userEntry); err != nil {
	//		return err
	//	}
	//}
}

//GetUserInfo return info of an user
func (udb *UserDB) GetUserInfo(address string) (common.User, error) {
	user := common.User{}
	return user, nil
	//
	//err := udb.db.Model(&user).Where("address = ?", address).Select()
	//return user, err
}
