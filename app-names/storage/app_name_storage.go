package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/KyberNetwork/reserve-stats/app-names/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	appNameTable       string = "app_name"
	addressesTableName string = "addresses"
)

var (

	// ErrAppNotExist exported error for checking
	ErrAppNotExist = errors.New("app does not exist")

	// ErrAddrExisted exported error for checking
	ErrAddrExisted = errors.New("address already exists")
)

// AppNameDB storage app name with its correspond addresses
type AppNameDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//DeleteAllTables delete all table from schema using for test only
func (adb *AppNameDB) DeleteAllTables() (err error) {
	_, err = adb.db.Exec(fmt.Sprintf(`DROP TABLE "%s", "%s"`, appNameTable, addressesTableName))
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
	app_name_id SERIAL NOT NULL REFERENCES app_name (id)
);
`
	var logger = sugar.With("func", "appname/storage.NewAppNameDB")

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != err {
				logger.Debugw("rollback error", "err", rollBackErr)
			}
			return
		}
		err = tx.Commit()
	}()
	logger.Debug("initializing app name database")
	if _, err = tx.Exec(fmt.Sprintf(schemaFmt, appNameTable, addressesTableName)); err != nil {
		return nil, err
	}
	logger.Debug("initialized app name database")

	return &AppNameDB{
		sugar: sugar,
		db:    db,
	}, err
}

// CreateOrUpdate create or update app name
func (adb *AppNameDB) CreateOrUpdate(app common.Application) (id int64, update bool, err error) {
	var (
		logger = adb.sugar.With(
			"func", "appname/storage.CreateOrUpdate",
			"app name", app.Name,
		)
		appID int64
	)
	logger.Debug("Create or update an app")
	tx, err := adb.db.Beginx()
	if err != nil {
		return id, update, err
	}

	defer func() {
		if err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != nil && rollBackErr != err {
				logger.Debugw("rollback error", "err", rollBackErr)
			}
			return
		}
		err = tx.Commit()
	}()

	if app.ID == 0 {
		// check if app_name.name exist
		row := tx.QueryRowx(`SELECT id FROM app_name WHERE name=$1`, app.Name)
		err = row.Scan(&appID)
		if err == sql.ErrNoRows {
			row = tx.QueryRowx(`INSERT INTO app_name (name) VALUES ($1) RETURNING id;`, app.Name)
			if err = row.Scan(&appID); err != nil {
				return id, update, err
			}
			logger = logger.With("app id", appID)
			logger.Debug("app is created")
		} else {
			update = true
			if _, err = tx.Exec(`UPDATE app_name SET name=$1, active=$2 WHERE id=$3`, app.Name, "TRUE", appID); err != nil {
				return id, update, err
			}
		}
	} else {
		appID = app.ID
		update = true
		if _, err = tx.Exec(`UPDATE app_name SET name=$1, active=$2 WHERE id=$3`, app.Name, "TRUE", app.ID); err != nil {
			return id, update, err
		}
	}

	// client will submit all registered addresses every time
	_, err = tx.Exec(fmt.Sprintf(`
	DELETE FROM "%s" WHERE app_name_id = $1
`, addressesTableName), appID)
	if err != nil {
		return id, update, err
	}
	for _, address := range app.Addresses {
		addressStr := address.Hex()
		logger.Debugw("updating app address",
			"address", addressStr,
		)
		var count int
		if err := tx.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE address='%s'", addressesTableName, addressStr)); err != nil {
			return id, update, err
		}
		logger.Debug(count)
		if count > 0 {
			return id, update, ErrAddrExisted
		}

		_, err = tx.Exec(fmt.Sprintf(`
INSERT INTO "%s" (address, app_name_id)
VALUES ($1, $2);
`, addressesTableName),
			addressStr,
			appID)
		if err != nil {
			return id, update, err
		}
	}

	return appID, update, err
}

//UpdateAppAddress add addresses to list address of appID
func (adb *AppNameDB) UpdateAppAddress(appID int64, app common.Application) (err error) {
	var (
		logger = adb.sugar.With("func", "appname/storage.UpdateAppAddress")
	)
	logger.Debugw("update app address", "id", appID)
	tx, err := adb.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		logger.Debug("commit transaction before return")
		if err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != err {
				logger.Debugw("rollback error", "err", rollBackErr)
			}
			return
		}
		err = tx.Commit()
	}()

	var qr []AppQueryResult
	if err := tx.Select(&qr, fmt.Sprintf(`SELECT name FROM app_name WHERE id=%d`, appID)); err != nil {
		return err
	}
	if len(qr) == 0 {
		return ErrAppNotExist
	}
	if _, err := tx.Exec("UPDATE app_name SET name=$1 WHERE id=$2", app.Name, appID); err != nil {
		return err
	}

	for _, address := range app.Addresses {
		var qresult []AppQueryResult
		addressStr := address.Hex()
		if err := tx.Select(&qresult, fmt.Sprintf(`SELECT address FROM addresses WHERE address='%s'`, addressStr)); err != nil {
			return err
		}
		if len(qresult) != 0 {
			logger.Debugw("address already exist", "address already exist", addressStr)
			continue
		}
		if _, err := tx.Exec(`INSERT INTO addresses (address, app_name_id) VALUES ($1, $2)`, addressStr, appID); err != nil {
			return err
		}
	}

	return err
}

//AppQueryResult result from query
type AppQueryResult struct {
	ID      int64  `json:"id"`
	Name    string `json:"name" db:"name"`
	Address string `json:"address"`
}

func appExist(app AppQueryResult, apps []common.Application) (bool, int) {
	for index, a := range apps {
		if a.ID == app.ID {
			return true, index
		}
	}
	return false, 0
}

// GetAllApp return all app in storage
func (adb *AppNameDB) GetAllApp(name, active string) ([]common.Application, error) {
	var (
		logger = adb.sugar.With(
			"func", "appname/storage.GetAllApp",
		)
		result  []common.Application
		qresult []AppQueryResult
	)
	logger.Debug("Get all apps")

	query := `SELECT app_name.id, app_name.name, addresses.address from 
	app_name JOIN addresses ON app_name.id = addresses.app_name_id WHERE `

	if strings.TrimSpace(active) != "" && strings.TrimSpace(active) == "false" {
		query += fmt.Sprintf("app_name.active IS FALSE ")
	} else {
		query += fmt.Sprintf("app_name.active IS TRUE ")
	}

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
			result[index].Addresses = append(result[index].Addresses, ethereum.HexToAddress(item.Address))
		} else {
			result = append(result, common.Application{
				ID:        item.ID,
				Addresses: []ethereum.Address{ethereum.HexToAddress(item.Address)},
				Name:      item.Name,
			})
		}
	}

	return result, nil
}

// GetApp return app address rom an app id
func (adb *AppNameDB) GetApp(appID int64) (common.Application, error) {
	var (
		logger = adb.sugar.With(
			"func", "appname/storage.GetAppAddresses",
		)
		qresult []AppQueryResult
		result  common.Application
	)
	logger.Debug("get an app")
	if err := adb.db.Select(&qresult, fmt.Sprintf(`SELECT app_name.id, app_name.name, addresses.address 
	from app_name JOIN addresses ON app_name.id = addresses.app_name_id WHERE app_name.id = %d AND app_name.active IS TRUE`, appID)); err != nil {
		return result, err
	}
	if len(qresult) == 0 {
		return result, ErrAppNotExist
	}
	result.ID = appID
	result.Name = qresult[0].Name
	for _, item := range qresult {
		result.Addresses = append(result.Addresses, ethereum.HexToAddress(item.Address))
	}

	return result, nil
}

//DeleteApp set app active is false
func (adb *AppNameDB) DeleteApp(appID int64) (err error) {
	var (
		logger = adb.sugar.With(
			"func", "appname/storage.DeleteApp",
		)
		count int
	)
	logger.Debugw("start delete app", "id", appID)
	tx, err := adb.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if err != rollbackErr {
				logger.Debugw("rollback error", "err", err)
			}
			return
		}
		err = tx.Commit()
	}()

	if err := tx.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM app_name WHERE id=%d", appID)); err != nil {
		return err
	}

	if count == 0 {
		return ErrAppNotExist
	}

	if _, err := tx.Exec("UPDATE app_name SET active=FALSE WHERE id=$1", appID); err != nil {
		return err
	}
	logger.Debug("finish delete app")
	return err
}
