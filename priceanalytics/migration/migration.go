package migration

import (
	"encoding/json"
	"errors"

	"github.com/KyberNetwork/reserve-stats/priceanalytics/common"

	"github.com/boltdb/bolt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	boltPriceAnalyticBucket = "price_analytic"
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
CREATE TABLE IF NOT EXISTS "%s" (
	id SERIAL PRIMARY KEY,
	timestamp TIMESTAMP NOT NULL,
	block_expiration BOOL NOT NULL
);

CREATE TABLE IF NOT EXISTS "%s" (
	id SERIAL PRIMARY KEY,
	token text NOT NULL,
	ask_price float NOT NULL,
	bid_price float NOT NULL,
	mid_afp_price float NOT NULL,
	mid_afp_old_price float NOT NULL,
	min_spread float NOT NULL,
	trigger_update bool NOT NULL,
	price_analytic_id SERIAL NOT NULL REFERENCES %s (id)
);
`
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
			return nil
		})
		logger.Info("migration completed")
		return nil
	})
}
