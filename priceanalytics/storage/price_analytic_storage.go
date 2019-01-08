package storage

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/priceanalytics/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	priceAnalyticTableName     = "price_analytics"
	priceAnalyticDataTableName = "price_analytics_data"
)

//PriceAnalyticDB represent db for price analytic data
type PriceAnalyticDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//NewPriceStorage return new storage for price analytics
func NewPriceStorage(sugar *zap.SugaredLogger, db *sqlx.DB) (*PriceAnalyticDB, error) {
	const schemaFmt = `
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
	var logger = sugar.With("func", "priceanalytics/storage.NewPriceStorage")

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	defer pgsql.CommitOrRollback(tx, logger, &err)

	logger.Debug("initializing database schema")

	if _, err = tx.Exec(fmt.Sprintf(schemaFmt,
		priceAnalyticTableName, priceAnalyticDataTableName, priceAnalyticTableName)); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")

	return &PriceAnalyticDB{
		sugar: sugar,
		db:    db,
	}, nil
}

//Close close db connection and return error if any
func (pad *PriceAnalyticDB) Close() error {
	return pad.db.Close()
}

// UpdatePriceAnalytic store price analytic to db
func (pad *PriceAnalyticDB) UpdatePriceAnalytic(data common.PriceAnalytic) error {
	var (
		logger = pad.sugar.With(
			"func", "priceanalytics/storage.UpdatePriceAnalytic",
		)
		priceAnalyticID int
	)

	tx, err := pad.db.Beginx()
	if err != nil {
		logger.Debug(err.Error())
		return err
	}

	defer pgsql.CommitOrRollback(tx, logger, &err)

	// insert into price_analytics
	row := tx.QueryRowx(`INSERT INTO price_analytics (timestamp, block_expiration) 
	VALUES ((TO_TIMESTAMP($1::double precision/1000)), $2) RETURNING id;`, data.Timestamp, data.BlockExpiration)
	if err = row.Scan(&priceAnalyticID); err != nil {
		return err
	}
	logger = logger.With("price_analytic_id", priceAnalyticID)
	logger.Debug("price analytic data created")

	// insert into price_analytic_data
	for _, values := range data.TriggeringTokensList {
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

	return err
}

// GetPriceAnalytic get price analytic data to return to api
func (pad *PriceAnalyticDB) GetPriceAnalytic(fromTime, toTime time.Time) ([]common.PriceAnalytic, error) {
	logger := pad.sugar.With("func", "priceanalytics/storage.GetPriceAnalytic",
		"fromTime", fromTime,
		"toTime", toTime)

	logger.Debug("get price analytic data")
	var result []common.PriceAnalytic
	if err := pad.db.Select(&result, fmt.Sprintf(`SELECT id, cast (extract(epoch from timestamp)*1000 as bigint) as timestamp 
	FROM %s WHERE timestamp >= $1 AND timestamp <= $2`, priceAnalyticTableName), fromTime.UTC(), toTime.UTC()); err != nil {
		logger.Debug("error get price analytics")
		return result, err
	}

	logger.Debug(result)

	if len(result) > 0 {
		for index, data := range result {
			var priceAnalyticData []common.PriceAnalyticData
			if err := pad.db.Select(&priceAnalyticData, fmt.Sprintf(`SELECT * from "%s" WHERE price_analytic_id = $1`, priceAnalyticDataTableName), data.ID); err != nil {
				pad.sugar.Debug("error getting detail data")
				return result, err
			}
			if len(priceAnalyticData) > 0 {
				result[index].TriggeringTokensList = priceAnalyticData
			}
		}
	}

	return result, nil
}

//DeleteAllTables delete all table from schema using for test only
func (pad *PriceAnalyticDB) DeleteAllTables() error {
	_, err := pad.db.Exec(fmt.Sprintf(`DROP TABLE "%s", "%s"`, priceAnalyticTableName, priceAnalyticDataTableName))
	return err
}
