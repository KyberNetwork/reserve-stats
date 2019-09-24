package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/priceanalytics/common"
)

//PriceAnalyticDB represent db for price analytic data
type PriceAnalyticDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//NewPriceStorage return new storage for price analytics
func NewPriceStorage(sugar *zap.SugaredLogger, db *sqlx.DB) (pa *PriceAnalyticDB, err error) {
	const schema = `CREATE TABLE IF NOT EXISTS "price_analytics"
(
    id               SERIAL PRIMARY KEY,
    timestamp        TIMESTAMP NOT NULL,
    block_expiration BOOL      NOT NULL
);

CREATE TABLE IF NOT EXISTS "price_analytics_data"
(
    id                SERIAL PRIMARY KEY,
    token             text   NOT NULL,
    ask_price         float  NOT NULL,
    bid_price         float  NOT NULL,
    mid_afp_price     float  NOT NULL,
    mid_afp_old_price float  NOT NULL,
    min_spread        float  NOT NULL,
    trigger_update    bool   NOT NULL,
    price_analytic_id SERIAL NOT NULL REFERENCES price_analytics (id)
);
`
	var logger = sugar.With("func", caller.GetCurrentFunctionName())

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	defer pgsql.CommitOrRollback(tx, logger, &err)

	logger.Debug("initializing database schema")

	if _, err = tx.Exec(schema); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")

	pa, err = &PriceAnalyticDB{
		sugar: sugar,
		db:    db,
	}, nil
	return
}

//Close close db connection and return error if any
func (pad *PriceAnalyticDB) Close() error {
	return pad.db.Close()
}

// UpdatePriceAnalytic store price analytic to db
func (pad *PriceAnalyticDB) UpdatePriceAnalytic(data common.PriceAnalytic) (err error) {
	var (
		logger          = pad.sugar.With("func", caller.GetCurrentFunctionName())
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
		_, err = tx.Exec(`INSERT INTO "price_analytics_data" (token,
                                    ask_price,
                                    bid_price,
                                    mid_afp_price,
                                    mid_afp_old_price,
                                    min_spread,
                                    trigger_update,
                                    price_analytic_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`,
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
	logger := pad.sugar.With("func", caller.GetCurrentFunctionName(),
		"fromTime", fromTime,
		"toTime", toTime)

	logger.Debug("get price analytic data")
	var result []common.PriceAnalytic
	if err := pad.db.Select(&result, `SELECT id, cast(extract(epoch from timestamp) * 1000 as bigint) as timestamp
FROM price_analytics
WHERE timestamp >= $1
  AND timestamp <= $2
`, fromTime.UTC(), toTime.UTC()); err != nil {
		logger.Debug("error get price analytics")
		return result, err
	}

	logger.Debug(result)

	if len(result) > 0 {
		for index, data := range result {
			var priceAnalyticData []common.PriceAnalyticData
			if err := pad.db.Select(&priceAnalyticData, `SELECT *
from "price_analytics_data"
WHERE price_analytic_id = $1`, data.ID); err != nil {
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
