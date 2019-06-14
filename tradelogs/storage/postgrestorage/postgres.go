package postgrestorage

import (
	"database/sql"
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"time"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"
)

const (
	tradeLogsTableName = "tradelogs"
	reserveTableName   = "reserve"
	tokenTableName     = "token"
	walletTableName    = "wallet"
)

// TradeLogDB is storage of tradelog data
type TradeLogDB struct {
	sugar                *zap.SugaredLogger
	db                   *sqlx.DB
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface
}

//NewTradeLogDB create a new instance of TradeLogDB
func NewTradeLogDB(sugar *zap.SugaredLogger, db *sqlx.DB, tokenAmountFormatter blockchain.TokenAmountFormatterInterface) (*TradeLogDB, error) {
	var logger = sugar.With("func", "tradelogs/storage.NewTradeLogDB")
	var err error
	logger.Debug("initializing database schema")
	if _, err = db.Exec(schema.TradeLogsSchema); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")

	return &TradeLogDB{
		sugar:                sugar,
		db:                   db,
		tokenAmountFormatter: tokenAmountFormatter,
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
	stmt := fmt.Sprintf(`SELECT "MAX(block_number)" FROM "%s"`, tradeLogsTableName)
	logger = logger.With("query", stmt)
	logger.Debug("Start query")
	if err := tldb.db.Get(&result, stmt); err != nil {
		if err == sql.ErrNoRows {
			logger.Info("No log saved")
			return 0, nil
		}
		logger.Errorw("Get error ", "error", err)
		return 0, err
	}
	return result, nil
}

func (tldb *TradeLogDB) saveReserveAddress(tx *sqlx.Tx, reserveAddressArray []string) error {
	var logger = tldb.sugar.With(
		"func", "tradelogs/storage/postgrestorage/TradeLogDB.saveReserveAddress",
	)
	query := fmt.Sprintf(insertionAddressTemplate, reserveTableName)
	logger.Debugw("updating rsv...", "query", query)
	_, err := tx.Exec(query, pq.StringArray(reserveAddressArray))
	return err
}

func (tldb *TradeLogDB) saveTokens(tx *sqlx.Tx, tokensArray []string) error {
	var logger = tldb.sugar.With(
		"func", "tradelogs/storage/postgrestorage/TradeLogDB.saveTokens",
	)
	query := fmt.Sprintf(insertionAddressTemplate, tokenTableName)
	logger.Debugw("updating rsv...", "query", query)
	_, err := tx.Exec(query, pq.StringArray(tokensArray))
	return err
}

func (tldb *TradeLogDB) saveWallets(tx *sqlx.Tx, walletAddressArray []string) error {
	var logger = tldb.sugar.With(
		"func", "tradelogs/storage/postgrestorage/TradeLogDB.saveWallets",
	)
	query := fmt.Sprintf(insertionAddressTemplate, walletTableName)
	logger.Debugw("updating rsv...", "query", query)
	_, err := tx.Exec(query, pq.StringArray(walletAddressArray))
	return err
}

// SaveTradeLogs persist trade logs to DB
func (tldb *TradeLogDB) SaveTradeLogs(logs []common.TradeLog) error {
	var (
		logger = tldb.sugar.With(
			"func", "tradelogs/storage/postgrestorage/TradeLogDB.SaveTradeLogs",
		)
		reserveAddress      = make(map[string]struct{})
		reserveAddressArray []string
		tokens              = make(map[string]struct{})
		tokensArray         []string
		walletAddress       = make(map[string]struct{})
		walletAddressArray  []string
		records             []*record
	)
	for _, log := range logs {
		r, err := tldb.recordFromTradeLog(log)
		if err != nil {
			return err
		}
		records = append(records, r)
	}

	for _, r := range records {
		reserve := r.SrcReserveAddress
		if _, ok := reserveAddress[reserve]; !ok {
			reserveAddress[reserve] = struct{}{}
			reserveAddressArray = append(reserveAddressArray, reserve)
		}
		reserve = r.DstReserveAddress
		if _, ok := reserveAddress[reserve]; !ok {
			reserveAddress[reserve] = struct{}{}
			reserveAddressArray = append(reserveAddressArray, reserve)
		}
		token := r.SrcAddress
		if _, ok := tokens[reserve]; !ok {
			tokens[token] = struct{}{}
			tokensArray = append(tokensArray, token)
		}
		token = r.DestAddress
		if _, ok := tokens[reserve]; !ok {
			tokens[token] = struct{}{}
			tokensArray = append(tokensArray, token)
		}
		wallet := r.WalletAddress
		if _, ok := walletAddress[wallet]; !ok {
			walletAddress[wallet] = struct{}{}
			walletAddressArray = append(walletAddressArray, wallet)
		}
	}

	tx, err := tldb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)

	err = tldb.saveReserveAddress(tx, reserveAddressArray)
	if err != nil {
		return err
	}

	err = tldb.saveTokens(tx, tokensArray)
	if err != nil {
		return err
	}

	err = tldb.saveWallets(tx, walletAddressArray)
	if err != nil {
		return err
	}

	for _, r := range records {
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

	return err
}

// TODO: implement this
func (tldb *TradeLogDB) LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error) {
	return nil, nil
}

const insertionAddressTemplate = `INSERT INTO %[1]s(
	address
) VALUES(
	unnest($1::TEXT[])
)
ON CONFLICT ON CONSTRAINT %[1]s_address_key DO NOTHING`

const insertionUserTemplate string = `
INSERT INTO users(
	address,
	timestamp
) VALUES (
	:user_address,
	:timestamp
)
ON CONFLICT (address) 
DO NOTHING;`
const insertionTradelogsTemplate string = `
INSERT INTO tradelogs(
	timestamp,
 	block_number,
 	tx_hash,
 	eth_amount,
 	user_address_id,
 	src_address_id,
 	dest_address_id,
 	src_reserveaddress_id,
 	dst_reserveaddress_id,
 	src_amount,
 	dest_amount,
 	wallet_address_id,
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
 	(SELECT id FROM users WHERE address=:user_address),
 	(SELECT id FROM token WHERE address=:src_address),
 	(SELECT id FROM token WHERE address=:dest_address),
 	(SELECT id FROM reserve WHERE address=:src_reserveaddress),
 	(SELECT id FROM reserve WHERE address=:dst_reserveaddress),
 	:src_amount,
 	:dest_amount,
 	(SELECT id FROM wallet WHERE address=:wallet_address),
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
