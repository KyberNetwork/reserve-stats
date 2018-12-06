package storage

import (
	"fmt"

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

// NewAppNameDB return new app name constance
func NewAppNameDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*AppNameDB, error) {
	const schemaFmt = `
CREATE TABLE IF NOT EXISTS "%s" (
	id SERIAL PRIMARY KEY,
	name text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS "%s" (
	id SERIAL PRIMARY KEY,
	address text NOT NULL UNIQUE,
	app_name_id SERIAL NOT NULL REFERENCES app_name (id)
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

//AppQueryResult result from query
type AppQueryResult struct {
	ID      int64  `json:"id"`
	AppName string `json:"name" db:"name"`
	Address string `json:"address"`
}

// GetAllApp return all app in storage
func (adb *AppNameDB) GetAllApp() ([]AppQueryResult, error) {
	var (
		logger = adb.sugar.With(
			"func", "appname/storage.CreateOrUpdate",
		)
		result []AppQueryResult
	)
	logger.Debug("Get all apps")

	if err := adb.db.Select(&result, `SELECT app_name.id, app_name.name, addresses.address from 
	app_name JOIN addresses ON app_name.id = addresses.app_name_id`); err != nil {
		return result, err
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
	from app_name JOIN addresses ON app_name.id = addresses.app_name_id WHERE app_name.id = %d`, appID)); err != nil {
		return result, err
	}
	result.ID = appID
	result.AppName = qresult[0].AppName
	for _, item := range qresult {
		result.Addresses = append(result.Addresses, item.Address)
	}

	return result, nil
}
