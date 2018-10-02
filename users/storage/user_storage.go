package storage

import (
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

//UserDB is storage of user data
type UserDB struct {
	db *pg.DB
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
func NewDB(address, user, password, database string) *UserDB {
	db := pg.Connect(&pg.Options{
		Addr:     address,
		User:     user,
		Password: password,
		Database: database,
	})
	// create schema
	if err := createSchema(db); err != nil {
		panic(err)
	}
	return &UserDB{
		db,
	}
}

//CloseDBConnection close db connection and return error if any
func (udb *UserDB) CloseDBConnection() error {
	return udb.db.Close()
}

//StoreUserInfo store user info to persist in database
func (udb *UserDB) StoreUserInfo(email string, addresses []common.UserAddress) error {
	userModel := common.User{}
	//remove all old user
	if _, err := udb.db.Model(&userModel).Where("email = ?", email).Delete(); err != nil {
		return err
	}

	// insert updated address value
	for _, address := range addresses {
		userEntry := common.User{
			Email:     email,
			Address:   address.Address,
			Timestamp: address.Timestamp,
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
