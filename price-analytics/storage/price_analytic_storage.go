package storage

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/price-analytics/common"
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

//StorePriceAnalytic store information from send from analytic to stat
func (pad *PriceAnalyticDB) StorePriceAnalytic() error {
	return nil
}

//NewPriceStorage return new storage for price analytics
func NewPriceStorage(sugar *zap.SugaredLogger, db *sqlx.DB) (*PriceAnalyticDB, error) {
	const schemaFmt = `
CREATE TABLE IF NOT EXISTS "%s" (
	id SERIAL PRIMARY KEY,
	timestamp TIMESTAMP NOT NULL,
	block_expriration BOOL NOT NULL
);

CREATE TABLE IF NOT EXISTS "%s" (
	id SERIAL PRIMARY KEY,
	token text NOT NULL,
	ask_price float NOT NULL,
	bid_price float NOT NULL,
	min_afp_price float NOT NULL,
	min_afp_old_price float NOT NULL,
	min_spread float NOT NULL,
	trigger_update bool NOT NULL,
	price_analatic_id NOT NULL REFERENCES "%s" (id)
);
`
	var logger = sugar.With("func", "price-analytics/storage.NewPriceStorage")

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	logger.Debug("initilizing database schema")

	if _, err = tx.Exec(fmt.Sprintf(schemaFmt, priceAnalyticTableName, priceAnalyticDataTableName, priceAnalyticTableName)); err != nil {
		return nil, err
	}
	logger.Debug("database schema initilized successfully")

	return &PriceAnalyticDB{
		sugar: sugar,
		db:    db,
	}, nil
}

// UpdatePriceAnalytic store price analytic to db
func (pad *PriceAnalyticDB) UpdatePriceAnalytic(data common.PriceAnalytic) error {
	var (
		logger = pad.sugar.With(
			"func", "price-analytics/storage.UpdatePriceAnalytic",
		)
		priceAnalyticID int
	)

	tx, err := pad.db.Beginx()
	if err != nil {
		return err
	}
	// insert into price_analytics
	row := tx.QueryRowx(`INSERt INTO price_analytics (timestamp, block_expiration) 
	VALUES ((TO_TIMESTAMP($1::double precision/1000)), $2) RETURNING id;`, data.Timestamp, data.BlockExpiration)
	if err = row.Scan(&priceAnalyticID); err != nil {
		return err
	}
	logger = logger.With("price_analytic_id", priceAnalyticID)
	logger.Debug("price analytic data created")

	// insert into price_analytic_data
	for _, values := range data.TriggeringTokensList {
		_, err = tx.Exec(fmt.Sprintf(`
INSERT INTO "%s" (token, ask_price, bid_price, min_afp_price, min_afp_old_price, min_spread, trigger_update, price_analytic_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, priceAnalyticDataTableName),
			values.Token,
			values.AskPrice,
			values.BidPrice,
			values.MinAfpPrice,
			values.MinAfpOldPrice,
			values.MinSpread,
			values.TriggerUpdate,
			priceAnalyticID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetPriceAnalytic get price analytic data to return to api
func (pad *PriceAnalyticDB) GetPriceAnalytic(fromTime, toTime time.Time) ([]common.PriceAnalytic, error) {
	var (
		result []common.PriceAnalytic
		err    error
	)

	if err := pad.db.Get(&result, fmt.Sprintf(`SELECT * FROM "%s" WHERE timestamp >= $1 AND timestamp <= $2`, priceAnalyticTableName), fromTime, toTime); err != nil {
		return result, err
	}

	if len(result) > 0 {
		for index, data := range result {
			priceAnalyticData := []common.PriceAnalyticData{}
			if err := pad.db.Get(&priceAnalyticData, fmt.Sprintf(`SELECT * from "%s" WHERE price_analytic_id = $1`, priceAnalyticDataTableName), data.ID); err != nil {
				return result, err
			}
			if len(priceAnalyticData) > 0 {
				result[index].TriggeringTokensList = priceAnalyticData
			}
		}
	}

	return result, err
}
