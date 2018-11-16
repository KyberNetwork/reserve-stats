package migration

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/KyberNetwork/reserve-stats/priceanalytics/common"
	"github.com/boltdb/bolt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	boltPriceAnalyticBucket    = "price_analytic"
	priceAnalyticTableName     = "price_analytics"
	priceAnalyticDataTableName = "price_analytics_data"
)

//DBMigration price analytic storage from bolt db to postgres
type DBMigration struct {
	sugar *zap.SugaredLogger

	boltdb     *bolt.DB
	postgresDB *sqlx.DB
}

//NewDBMigration return a new DBMigration instance
func NewDBMigration(sugar *zap.SugaredLogger, dbPath string, postgres *sqlx.DB) (*DBMigration, error) {
	var (
		err    error
		boltDB *bolt.DB
	)
	boltDB, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	const schema = `
CREATE TABLE IF NOT EXISTS "price_analytics" (
	id SERIAL PRIMARY KEY,
	timestamp TIMESTAMP NOT NULL,
	block_expiration BOOL NOT NULL
);

CREATE TABLE IF NOT EXISTS "price_analytics_data" (
	id SERIAL PRIMARY KEY,
	token text NOT NULL,
	ask_price float NOT NULL,
	bid_price float NOT NULL,
	mid_afp_price float NOT NULL,
	mid_afp_old_price float NOT NULL,
	min_spread float NOT NULL,
	trigger_update bool NOT NULL,
	price_analytic_id SERIAL NOT NULL REFERENCES price_analytics (id)
);`
	tx, err := postgres.Beginx()
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(schema); err != nil {
		return nil, err
	}
	storage := &DBMigration{sugar: sugar, boltdb: boltDB, postgresDB: postgres}
	return storage, nil
}

// Migrate data from bolt db to postgres
func (dbm *DBMigration) Migrate() error {
	var logger = dbm.sugar.With("func", "priceanalytics/migration/DBMigration.Migrate")
	logger.Debug("start migrating")
	return dbm.boltdb.View(func(tx *bolt.Tx) error {
		priceAnalyticBK := tx.Bucket([]byte(boltPriceAnalyticBucket))
		if priceAnalyticBK == nil {
			return errors.New("price analytic bucket is empty")
		}
		var count = 0
		priceAnalyticBK.ForEach(func(k, v []byte) error {
			count++
			var priceAnalyticData common.OldPriceAnalyticFormat
			if err := json.Unmarshal(v, &priceAnalyticData); err != nil {
				return err
			}
			logger.Debug(priceAnalyticData)
			if err := dbm.insertOldPriceAnalytic(priceAnalyticData); err != nil {
				return err
			}
			return nil
		})
		logger.Info("migration completed")
		return nil
	})
}

func (dbm *DBMigration) insertOldPriceAnalytic(data common.OldPriceAnalyticFormat) error {
	var (
		logger          = dbm.sugar.With("func", "priceanalytics/migration/DBMigration.Migrate")
		priceAnalyticID int64
	)
	logger.Debug(data)
	tx, err := dbm.postgresDB.Beginx()
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	// insert into price_analytics
	row := tx.QueryRowx(`INSERT INTO price_analytics (timestamp, block_expiration) 
	VALUES ((TO_TIMESTAMP($1::double precision/1000)), $2) RETURNING id;`, data.Timestamp, data.Value.BlockExpiration)
	if err = row.Scan(&priceAnalyticID); err != nil {
		return err
	}
	logger = logger.With("price_analytic_id", priceAnalyticID)
	logger.Debug("price analytic data created")

	// insert into price_analytic_data
	for _, values := range data.Value.TriggeringTokensList {
		_, err = tx.Exec(fmt.Sprintf(`
INSERT INTO "%s" (token, ask_price, bid_price, mid_afp_price, mid_afp_old_price, min_spread, trigger_update, price_analytic_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, priceAnalyticDataTableName),
			values.Token,
			values.AskPrice,
			values.BidPrice,
			values.MidAfpPrice,
			values.MidAfpOldPrice,
			values.MinSpread,
			values.TriggerUpdate,
			priceAnalyticID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
