package postgrestorage

import (
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"strings"

	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"
)

const (
	tradeLogsTableName = "tradelogs"
)

// TradeLogDB is storage of tradelog data
type TradeLogDB struct {
	sugar                *zap.SugaredLogger
	db                   *sqlx.DB
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface

	// traded stored traded addresses to use in a single SaveTradeLogs
	traded map[string]struct{}
}

//NewTradeLogDB create a new instance of TradeLogDB
func NewTradeLogDB(sugar *zap.SugaredLogger, db *sqlx.DB, tokenAmountFormatter blockchain.TokenAmountFormatterInterface) (*TradeLogDB, error) {
	var logger = sugar.With("func", "tradelogs/storage.NewTradeLogDB")

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	defer pgsql.CommitOrRollback(tx, logger, &err)

	logger.Debug("initializing database schema")
	if _, err = tx.Exec(schema.TradeLogsSchema); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")

	return &TradeLogDB{
		sugar:                sugar,
		db:                   db,
		tokenAmountFormatter: tokenAmountFormatter,
		traded:               make(map[string]struct{}),
	}, err
}

// LastBlock returns last stored trade log block number from database.
func (tldb *TradeLogDB) LastBlock() (int64, error) {
	var (
		logger = tldb.sugar.With(
			"func", "tradelog/storage/postgrestorage/TradeLogDB.SaveTradeLogs",
		)
		result int64
	)
	stmt := fmt.Sprintf(`SELECT "block_number" FROM "%s" ORDER BY timestamp DESC limit 1`, tradeLogsTableName)
	logger = logger.With("query", stmt)
	logger.Debug("Start query")
	if err := tldb.db.Get(&result, stmt); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			logger.Info("No log saved")
			return 0, nil
		}
		logger.Errorw("Get error ", "error", err)
		return 0, err
	}
	return result, nil
}

// SaveTradeLogs persist trade logs to DB
func (tldb *TradeLogDB) SaveTradeLogs(logs []common.TradeLog) error {
	logger := tldb.sugar.With(
		"func", "tradelogs/storage/postgrestorage/TradeLogDB.SaveTradeLogs",
	)
	tx, err := tldb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)
	for _, log := range logs {
		r, err := tldb.recordFromTradeLog(log)
		if err != nil {
			return err
		}
		logger.Debugw("Record", "record", r)
		_, err = tx.NamedExec(insertionUserTemplate, r)
		if err != nil {
			logger.Debugw("Error while add users", "error", err)
			return err
		}
		_, err = tx.NamedExec(insertionTradelogsTemplate, r)
		if err != nil {
			logger.Debugw("Error while add tradelogs", "error", err)
			return err
		}
	}

	// reset traded map to avoid ever growing size
	tldb.traded = make(map[string]struct{})
	return err
}

const insertionUserTemplate string = `
INSERT INTO users(
	user_address,
	timestamp
) VALUES (
	:user_address,
	:timestamp
)
ON CONFLICT (user_address) 
DO NOTHING;`
const insertionTradelogsTemplate string = `
INSERT INTO tradelogs(
	timestamp,
 	block_number,
 	tx_hash,
 	eth_amount,
 	user_address_id,
 	src_address,
 	dest_address,
 	src_reserveaddress,
 	dst_reserveaddress,
 	src_amount,
 	dest_amount,
 	wallet_address,
 	src_burn_amount,
 	dst_burn_amount,
 	src_wallet_fee_amount,
 	dst_wallet_fee_amount,
 	integration_app,
 	ip,
 	country,
 	ethusd_rate,
 	ethusd_provider,
	index
) VALUES (
 	:timestamp,
 	:block_number,
 	:tx_hash,
 	:eth_amount,
 	(SELECT id FROM users WHERE user_address=:user_address),
 	:src_address,
 	:dest_address,
 	:src_reserveaddress,
 	:dst_reserveaddress,
 	:src_amount,
 	:dest_amount,
 	:wallet_address,
 	:src_burn_amount,
 	:dst_burn_amount,
 	:src_wallet_fee_amount,
 	:dst_wallet_fee_amount,
 	:integration_app,
 	:ip,
 	:country,
 	:ethusd_rate,
 	:ethusd_provider,
 	:index
);`
