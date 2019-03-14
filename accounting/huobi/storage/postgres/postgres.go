package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	huobiTradesTableName = "huobi_trades"
)

//Option define init behaviour for db storage.
type Option func(*HuobiStorage) error

//WithTradeTableName return Option to set trade table Name
func WithTradeTableName(name string) Option {
	return func(hs *HuobiStorage) error {
		if hs.tableNames == nil {
			hs.tableNames = make(map[string]string)
		}
		hs.tableNames[huobiTradesTableName] = name
		return nil
	}
}

//HuobiStorage defines the object to store Huobi data
type HuobiStorage struct {
	sugar      *zap.SugaredLogger
	db         *sqlx.DB
	tableNames map[string]string
}

// NewDB return the HuobiStorage instance. User must call Close() before exit.
// tableNames is a list of string for tablename[huobitrades]. It can be optional
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB, options ...Option) (*HuobiStorage, error) {
	const schemaFMT = `
	CREATE TABLE IF NOT EXISTS %[1]s
(
	id bigint NOT NULL,
	data JSONB,
	CONSTRAINT %[1]s_pk PRIMARY KEY(id)
) ;
CREATE INDEX IF NOT EXISTS %[1]s_time_idx ON %[1]s ((data ->> 'created-at'));
`
	var (
		logger     = sugar.With("func", "reserverates/storage/postgres/Newdb")
		tableNames = map[string]string{huobiTradesTableName: huobiTradesTableName}
	)
	hs := &HuobiStorage{
		sugar:      sugar,
		db:         db,
		tableNames: tableNames,
	}
	for _, opt := range options {
		if err := opt(hs); err != nil {
			return nil, err
		}
	}

	query := fmt.Sprintf(schemaFMT, hs.tableNames[huobiTradesTableName])
	logger.Debugw("initializing database schema", "query", query)
	if _, err := hs.db.Exec(query); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")
	return hs, nil
}

//TearDown removes all the tables
func (hdb *HuobiStorage) TearDown() error {
	const dropFMT = `
	DROP TABLE %[1]s;
	DROP INDEX IF EXISTS %[1]s_time_idx
	`
	query := fmt.Sprintf(dropFMT, hdb.tableNames[huobiTradesTableName])
	hdb.sugar.Debugw("tearingdown", "query", dropFMT, "table name", hdb.tableNames[huobiTradesTableName])
	_, err := hdb.db.Exec(query)
	return err
}

//Close close DB connection
func (hdb *HuobiStorage) Close() error {
	if hdb.db != nil {
		return hdb.db.Close()
	}
	return nil
}

//UpdateTradeHistory store the TradeHistory rate at that blockInfo
func (hdb *HuobiStorage) UpdateTradeHistory(trade huobi.TradeHistory) error {
	var (
		timestamp = timeutil.TimestampMsToTime(trade.CreateAt)
		logger    = hdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
			"trade_ID", trade.ID,
			"timestamp", timestamp,
		)
	)
	const updateStmt = `INSERT INTO %[1]s(id, data)
	VALUES ( 
		$1,
		$2
	)
	ON CONFLICT ON CONSTRAINT %[1]s_pk DO NOTHING;`
	query := fmt.Sprintf(updateStmt,
		hdb.tableNames[huobiTradesTableName],
	)
	data, err := json.Marshal(trade)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	logger.Debugw("updating tradeHistory...", "query", query)
	_, err = hdb.db.Exec(query,
		trade.ID,
		data,
	)
	if err != nil {
		return err
	}
	return err
}

//GetTradeHistory return tradehistory between from.. to.. in its json []byte form
func (hdb *HuobiStorage) GetTradeHistory(from, to time.Time) ([]huobi.TradeHistory, error) {
	var (
		dbResult [][]byte
		result   []huobi.TradeHistory
		logger   = hdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
			"from", from.String(),
			"to", to.String(),
		)
		tmp huobi.TradeHistory
	)
	const selectStmt = `SELECT data FROM %[1]s WHERE data->>'created-at'>=$1 AND data->>'created-at'<$2`
	query := fmt.Sprintf(selectStmt, hdb.tableNames[huobiTradesTableName])
	logger.Debugw("querying trade history...", "query", query)
	if err := hdb.db.Select(&dbResult, query, timeutil.TimeToTimestampMs(from), timeutil.TimeToTimestampMs(to)); err != nil {
		return result, err
	}
	for _, data := range dbResult {
		if err := json.Unmarshal(data, &tmp); err != nil {
			return result, err
		}
		result = append(result, tmp)
	}
	return result, nil
}
