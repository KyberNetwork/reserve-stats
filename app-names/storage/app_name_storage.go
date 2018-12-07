package storage

import (
	"errors"
	"fmt"
	"strings"

	"github.com/KyberNetwork/reserve-stats/app-names/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	appNameTable       string = "app_name"
	addressesTableName string = "addresses"
)

// AppNameDB storage app name with its correspond addresses
type AppNameDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//DeleteAllTables delete all table from schema using for test only
func (adb *AppNameDB) DeleteAllTables() error {
	_, err := adb.db.Exec(fmt.Sprintf(`DROP TABLE "%s", "%s"`, appNameTable, addressesTableName))
	return err
}

// NewAppNameDB return new app name constance
func NewAppNameDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*AppNameDB, error) {
	const schemaFmt = `
CREATE TABLE IF NOT EXISTS "%s" (
	id SERIAL PRIMARY KEY,
	name text NOT NULL UNIQUE,
	active boolean DEFAULT TRUE 
);

CREATE TABLE IF NOT EXISTS "%s" (
	id SERIAL PRIMARY KEY,
	address text NOT NULL UNIQUE,
	app_name_id SERIAL NOT NULL REFERENCES app_name (id),
	active boolean DEFAULT TRUE
);
`
	var logger = sugar.With("func", "appname/storage.NewAppNameDB")

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	logger.Debug("initializing app name database")
	if _, err = tx.Exec(fmt.Sprintf(schemaFmt, appNameTable, addressesTableName)); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	logger.Debug("initialized app name database")

	return &AppNameDB{
		sugar: sugar,
		db:    db,
	}, nil
}

// CreateOrUpdate create or update app name
func (adb *AppNameDB) CreateOrUpdate(app common.AppObject) (common.AppObject, error) {
	var (
		logger = adb.sugar.With(
			"func", "appname/storage.CreateOrUpdate",
			"app name", app.AppName,
		)
		appID int64
	)
	logger.Debug("Create or update an app")
	tx, err := adb.db.Beginx()
	if err != nil {
		return app, err
	}

	if app.ID == 0 {
		row := tx.QueryRowx(`INSERT INTO app_name (name) VALUES ($1) RETURNING id;`, app.AppName)
		if err = row.Scan(&appID); err != nil {
			return app, err
		}
		logger = logger.With("app id", appID)
		logger.Debug("app is created")
	} else {
		appID = app.ID
		if _, err = tx.Exec(`UPDATE app_name SET app_name=$1 WHERE id=$2`, app.AppName, app.ID); err != nil {
			return app, err
		}
	}

	// client will submit all registered addresses every time
	_, err = tx.Exec(fmt.Sprintf(`
	DELETE FROM "%s" WHERE app_name_id = $1
`, addressesTableName), appID)
	if err != nil {
		return app, err
	}
	for _, address := range app.Addresses {
		logger.Debugw("updating app address",
			"address", address,
		)
		_, err = tx.Exec(fmt.Sprintf(`
INSERT INTO "%s" (address, app_name_id)
VALUES ($1, $2);
`, addressesTableName),
			address,
			appID)
		if err != nil {
			return app, err
		}
	}

	app.ID = appID

	return app, tx.Commit()
}

//UpdateAppAddress add addresses to list address of appID
func (adb *AppNameDB) UpdateAppAddress(appID int64, app common.AppObject) error {
	var (
		logger = adb.sugar.With("func", "appname/storage.UpdateAppAddress")
	)
	logger.Debugw("udpate app address", "udpate app address", appID)
	tx, err := adb.db.Beginx()
	if err != nil {
		return err
	}
	var qr []AppQueryResult
	if err := tx.Select(&qr, fmt.Sprintf(`SELECT name FROM app_name WHERE id=%d`, appID)); err != nil {
		return err
	}
	if len(qr) == 0 {
		err := tx.Commit()
		if err != nil {
			return err
		}
		return errors.New("app does not exist")
	}

	for _, address := range app.Addresses {
		var qresult []AppQueryResult
		if err := tx.Select(&qresult, fmt.Sprintf(`SELECT address FROM addresses WHERE address='%s'`, address)); err != nil {
			return err
		}
		if len(qresult) != 0 {
			logger.Debugw("address already exist", "address already exist", address)
			continue
		}
		logger.Debug(address)
		if _, err := tx.Exec(`INSERT INTO addresses (address, app_name_id) VALUES ($1, $2)`, address, appID); err != nil {
			logger.Debug(err)
			return err
		}
	}
	return tx.Commit()
}

//AppQueryResult result from query
type AppQueryResult struct {
	ID      int64  `json:"id"`
	AppName string `json:"name" db:"name"`
	Address string `json:"address"`
}

func appExist(app AppQueryResult, apps []common.AppObject) (bool, int) {
	for index, a := range apps {
		if a.ID == app.ID {
			return true, index
		}
	}
	return false, 0
}

// GetAllApp return all app in storage
func (adb *AppNameDB) GetAllApp(name string) ([]common.AppObject, error) {
	var (
		logger = adb.sugar.With(
			"func", "appname/storage.GetAllApp",
		)
		result  []common.AppObject
		qresult []AppQueryResult
	)
	logger.Debug("Get all apps")

	query := `SELECT app_name.id, app_name.name, addresses.address from 
	app_name JOIN addresses ON app_name.id = addresses.app_name_id WHERE app_name.active IS TRUE `

	if strings.TrimSpace(name) != "" {
		query += fmt.Sprintf("AND app_name.name='%s'", name)
	}
	logger.Debug(query)
	if err := adb.db.Select(&qresult, query); err != nil {
		return result, err
	}

	for _, item := range qresult {
		exist, index := appExist(item, result)
		if exist {
			result[index].Addresses = append(result[index].Addresses, item.Address)
		} else {
			result = append(result, common.AppObject{
				ID:        item.ID,
				Addresses: []string{item.Address},
				AppName:   item.AppName,
			})
		}
	}

	return result, nil
}

// GetAppAddresses return app address rom an app id
func (adb *AppNameDB) GetAppAddresses(appID int64) (common.AppObject, error) {
	var (
		logger = adb.sugar.With(
			"func", "appname/storage.GetAppAddresses",
		)
		qresult []AppQueryResult
		result  common.AppObject
	)
	logger.Debug("get an app")
	if err := adb.db.Select(&qresult, fmt.Sprintf(`SELECT app_name.id, app_name.name, addresses.address 
	from app_name JOIN addresses ON app_name.id = addresses.app_name_id WHERE app_name.id = %d AND app_name.active IS TRUE`, appID)); err != nil {
		return result, err
	}
	if len(qresult) == 0 {
		return result, errors.New("app does not exist")
	}
	result.ID = appID
	result.AppName = qresult[0].AppName
	for _, item := range qresult {
		result.Addresses = append(result.Addresses, item.Address)
	}

	return result, nil
}

//DeleteApp set app active is false
func (adb *AppNameDB) DeleteApp(appID int64) error {
	var (
		logger = adb.sugar.With(
			"func", "appname/storage.DeleteApp",
		)
	)
	logger.Debug("start delete app")
	tx, err := adb.db.Beginx()
	if err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE app_name SET active=FALSE WHERE id=$1", appID); err != nil {
		return err
	}
	logger.Debug("finish delete app")
	return tx.Commit()
}

//DeleteAddress soft delete an address
func (adb *AppNameDB) DeleteAddress(address string) error {
	var logger = adb.sugar.With(
		"func", "appname/storage.DeleteAddress",
	)
	logger.Debug("start delete address")
	tx, err := adb.db.Beginx()
	if err != nil {
		return nil
	}
	if _, err := tx.Exec("UPDATE addresses SET active=FALSE WHERE address=$1", address); err != nil {
		return err
	}
	logger.Debug("finish delete address")
	return tx.Commit()
}
