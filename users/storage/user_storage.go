package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"go.uber.org/zap"
)

//UserDB is storage of user data
type UserDB struct {
	sugar *zap.SugaredLogger
	db    *pg.DB
}

//DeleteAllTables delete all table from schema using for test only
func (udb *UserDB) DeleteAllTables() error {
	userTable := &common.User{}
	if err := udb.db.DropTable(userTable, &orm.DropTableOptions{}); err != nil {
		return err
	}
	return nil
}

func createSchema(db *pg.DB) error {
	createTableOptions := &orm.CreateTableOptions{
		IfNotExists: true,
	}

	// create table user info
	userTable := &common.User{}
	if err := db.CreateTable(userTable, createTableOptions); err != nil {
		return err
	}

	return nil
}

//NewDB open a new database connection
func NewDB(sugar *zap.SugaredLogger, db *pg.DB) *UserDB {
	// create schema
	if err := createSchema(db); err != nil {
		panic(err)
	}
	return &UserDB{
		sugar,
		db,
	}
}

//CloseDBConnection close db connection and return error if any
func (udb *UserDB) CloseDBConnection() error {
	return udb.db.Close()
}

//StoreUserInfo store user info to persist in database
func (udb *UserDB) StoreUserInfo(userData common.UserData) error {
	userModel := common.User{}
	//remove all old user
	if _, err := udb.db.Model(&userModel).Where("email = ?", userData.Email).Delete(); err != nil {
		return err
	}

	// insert updated address value
	for _, address := range userData.UserInfo {
		userEntry := common.User{
			Email:     userData.Email,
			Address:   address.Address,
			Timestamp: time.Unix(address.Timestamp/1000, 0),
		}
		if err := udb.db.Insert(&userEntry); err != nil {
			return err
		}
	}
	return nil
}

//GetUserInfo return info of an user
func (udb *UserDB) GetUserInfo(address string) (common.User, error) {
	user := common.User{}

	err := udb.db.Model(&user).Where("address = ?", address).Select()
	return user, err
}
