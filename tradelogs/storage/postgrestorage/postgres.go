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

// TradeLogDB is storage of tradelog data
type TradeLogDB struct {
	sugar                *zap.SugaredLogger
	db                   *sqlx.DB
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface
	dbName               string
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
		result sql.NullInt64
	)
	stmt := fmt.Sprintf(`SELECT MAX("block_number") FROM "%v"`, schema.TradeLogsTableName)
	logger = logger.With("query", stmt)
	logger.Debug("Start query")
	//if err := tldb.db.QueryRow(&result, stmt); err != nil {
	//	logger.Errorw("Get error ", "error", err)
	//	return 0, err
	//}
	row := tldb.db.QueryRow(stmt)
	if err := row.Scan(&result); err != nil {
		logger.Errorw("Get error ", "error", err)
		return 0, err
	}

	if !result.Valid {
		logger.Info("No log saved")
		return 0, nil
	}

	return result.Int64, nil
}

func (tldb *TradeLogDB) saveReserveAddress(tx *sqlx.Tx, reserveAddressArray []string) error {
	var logger = tldb.sugar.With(
		"func", "tradelogs/storage/postgrestorage/TradeLogDB.saveReserveAddress",
	)
	query := fmt.Sprintf(insertionAddressTemplate, schema.ReserveTableName)
	logger.Debugw("updating rsv...", "query", query)
	_, err := tx.Exec(query, pq.StringArray(reserveAddressArray))
	return err
}

func (tldb *TradeLogDB) saveTokens(tx *sqlx.Tx, tokensArray []string) error {
	var logger = tldb.sugar.With(
		"func", "tradelogs/storage/postgrestorage/TradeLogDB.saveTokens",
	)
	query := fmt.Sprintf(insertionAddressTemplate, schema.TokenTableName)
	logger.Debugw("updating rsv...", "query", query)
	_, err := tx.Exec(query, pq.StringArray(tokensArray))
	return err
}

func (tldb *TradeLogDB) saveWallets(tx *sqlx.Tx, walletAddressArray []string) error {
	var logger = tldb.sugar.With(
		"func", "tradelogs/storage/postgrestorage/TradeLogDB.saveWallets",
	)
	query := fmt.Sprintf(insertionAddressTemplate, schema.WalletTableName)
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
			logger.Debugw("Error while add trade logs", "error", err)
			return err
		}
	}

	return err
}

// TODO: implement this
func (tldb *TradeLogDB) LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error) {
	var queryResult []struct {
	}
	tldb.db.Select(&queryResult, `SELECT timestamp, block_number, eth_amount, eth_amount, user_addr`)

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
INSERT INTO "` + schema.TradeLogsTableName + `"(
	timestamp,
 	block_number,
 	tx_hash,
 	eth_amount,
 	user_address_id,
 	src_address_id,
 	dst_address_id,
 	src_reserve_address_id,
 	dst_reserve_address_id,
 	src_amount,
 	dst_amount,
 	wallet_address_id,
 	src_burn_amount,
 	dst_burn_amount,
 	src_wallet_fee_amount,
 	dst_wallet_fee_amount,
 	integration_app,
 	ip,
 	country,
 	eth_usd_rate,
 	eth_usd_provider,
	index
) VALUES (
 	:timestamp,
 	:block_number,
 	:tx_hash,
 	:eth_amount,
 	(SELECT id FROM users WHERE address=:user_address),
 	(SELECT id FROM token WHERE address=:src_address),
 	(SELECT id FROM token WHERE address=:dst_address),
 	(SELECT id FROM reserve WHERE address=:src_reserve_address),
 	(SELECT id FROM reserve WHERE address=:dst_reserve_address),
 	:src_amount,
 	:dst_amount,
 	(SELECT id FROM wallet WHERE address=:wallet_address),
 	:src_burn_amount,
 	:dst_burn_amount,
 	:src_wallet_fee_amount,
 	:dst_wallet_fee_amount,
 	:integration_app,
 	:ip,
 	:country,
 	:eth_usd_rate,
 	:eth_usd_provider,
 	:index
)
ON CONFLICT
DO NOTHING;`

const selectTradeLogsQuery = `
SELECT timestamp, block_number, eth_amount, d.address AS user_addr, e.address AS src_addr, f.address AS dst_addr,
src_amount, dst_amount, ip, country, uid, integration_app, src_burn_amount, dst_burn_amount,
log_index, tx_hash, b.address AS src_rsv_addr, c.address AS dst_rsv_addr, src_wallet_fee_amount, dst_wallet_fee_amount,
g.address AS wallet_addr, tx_sender, receiver_addr 
FROM "` + schema.TradeLogsTableName + `" AS a
INNER JOIN reserve AS b where a.src_reserve_address_id = b.id
INNER JOIN reserve AS c where a.dst_reserve_address_id = c.id
INNER JOIN user AS d where a.user_address_id = d.id
INNER JOIN token AS e where a.src_address_id = e.id
INNER JOIN token AS f where f.dst_address_id = f.id
INNER JOIN wallet AS g where a.wallet_address_id = g.id
WHERE timestamp >= $1 and timestamp <= $2
`
