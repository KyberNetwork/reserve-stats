package storage

import (
	"database/sql"
	"errors"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/app-names/common"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
)

var (
	// ErrNotExists exported error for checking
	ErrNotExists = errors.New("app does not exist")
)

// AppNameDB storage app name with its correspond addresses
type AppNameDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

// NewAppNameDB return new app name constance
func NewAppNameDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*AppNameDB, error) {
	var logger = sugar.With("func", caller.GetCurrentFunctionName())

	logger.Debug("initializing app name database")
	if _, err := db.Exec(schemaFmt); err != nil {
		return nil, err
	}
	logger.Debug("initialized app name database")

	return &AppNameDB{
		sugar: sugar,
		db:    db,
	}, nil
}

type createOrUpdateResult struct {
	ID      int64 `db:"id"`
	Updated bool  `db:"updated"`
}

// CreateOrUpdate create or update app name
func (adb *AppNameDB) CreateOrUpdate(app common.Application) (int64, bool, error) {
	var (
		logger = adb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"app name", app.Name,
		)
		addresses []string
		result    = createOrUpdateResult{}
	)
	logger.Debug("Create or update an app")

	for _, address := range app.Addresses {
		addresses = append(addresses, address.String())
	}

	err := adb.db.Get(&result, `SELECT _id AS id, updated FROM create_or_update_app($1, $2, $3)`,
		app.ID,
		app.Name,
		pq.StringArray(addresses))
	if err != nil {
		return 0, false, err
	}
	return result.ID, result.Updated, nil
}

//Update add addresses to list address of appID
func (adb *AppNameDB) Update(app common.Application) error {
	var (
		logger    = adb.sugar.With("func", caller.GetCurrentFunctionName())
		addresses []string
	)
	logger.Debugw("update app address", "id", app.ID)

	for _, address := range app.Addresses {
		addresses = append(addresses, address.String())
	}

	_, err := adb.db.Exec(`SELECT update_app($1, $2, $3)`,
		app.ID,
		app.Name,
		pq.StringArray(addresses))

	if err != nil {
		pErr, ok := err.(*pq.Error)
		if !ok {
			return err
		}
		if pErr != nil && pErr.Code == "P0002" {
			return ErrNotExists
		}
		return pErr
	}

	return nil
}

type getAppResult struct {
	ID        int64          `db:"id"`
	Name      string         `db:"name"`
	Addresses pq.StringArray `db:"addresses"`
}

// GetAll return all app in storage
func (adb *AppNameDB) GetAll(filters ...Filter) ([]common.Application, error) {
	var (
		logger = adb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
		)
		query = `SELECT joined.id,
       joined.name,
       joined.addresses
FROM (SELECT apps.id, apps.name, ARRAY_AGG(address) FILTER ( WHERE address IS NOT NULL) AS addresses
      FROM app_names AS apps
               LEFT JOIN addresses AS addrs on apps.id = addrs.app_name_id
      WHERE ($1::TEXT IS NULL OR apps.name = $1)
        AND ($3::BOOLEAN IS NULL OR apps.active = $3)
      GROUP BY apps.id, apps.name) AS joined
WHERE ($2::TEXT IS NULL OR $2 ILIKE ANY (joined.addresses));
`
		filterConf = &FilterConf{}
		result     []getAppResult
		apps       []common.Application
	)
	logger.With("query", query)

	for _, filter := range filters {
		filter(filterConf)
	}

	if filterConf.Name != nil {
		logger.With("name", *filterConf.Name)
	}

	if filterConf.Address != nil {
		logger.With("address", *filterConf.Address)
	}

	if filterConf.Active != nil {
		logger.With("active", *filterConf.Active)
	}

	logger.Debug("get all applications")
	if err := adb.db.Select(&result, query,
		filterConf.Name,
		filterConf.Address,
		filterConf.Active); err != nil {
		return nil, err
	}

	for _, r := range result {
		app := common.Application{
			ID:   r.ID,
			Name: r.Name,
		}
		for _, addr := range []string(r.Addresses) {
			app.Addresses = append(app.Addresses, ethereum.HexToAddress(addr))
		}
		apps = append(apps, app)
	}
	return apps, nil
}

// Get return app address rom an app id
func (adb *AppNameDB) Get(appID int64) (common.Application, error) {
	var (
		logger = adb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
		)
		result getAppResult
		app    common.Application
	)
	logger.Debug("get an application")
	err := adb.db.Get(&result, `SELECT apps.id, apps.name, ARRAY_AGG(address) FILTER ( WHERE address IS NOT NULL) AS addresses
FROM app_names AS apps
         LEFT JOIN addresses AS addrs on apps.id = addrs.app_name_id
WHERE apps.id = $1
  AND apps.active = TRUE
GROUP BY apps.id, apps.name;`, appID)

	switch {
	case err == sql.ErrNoRows:
		return common.Application{}, ErrNotExists
	case err != nil:
		return common.Application{}, err
	}

	app = common.Application{
		ID:   result.ID,
		Name: result.Name,
	}

	for _, addr := range result.Addresses {
		app.Addresses = append(app.Addresses, ethereum.HexToAddress(addr))
	}

	return app, nil
}

//Delete set app active is false
func (adb *AppNameDB) Delete(appID int64) (err error) {
	var (
		logger = adb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
		)
		storedID int64
	)
	logger.Debugw("start delete app", "id", appID)
	tx, err := adb.db.Beginx()
	if err != nil {
		return
	}

	defer pgsql.CommitOrRollback(tx, adb.sugar, &err)

	err = tx.Get(&storedID, `SELECT id FROM app_names WHERE id = $1`, appID)
	switch {
	case err == sql.ErrNoRows:
		err = ErrNotExists
		return
	case err != nil:
		return
	}

	if _, err = tx.Exec("UPDATE app_names SET active=FALSE WHERE id=$1", appID); err != nil {
		return
	}
	logger.Debug("finish delete app")
	return
}
