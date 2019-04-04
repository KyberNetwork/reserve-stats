package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	walletErc20TxsTableName = "wallet_erc20_txs"
)

//Option define init behaviour for db storage.
type Option func(*WalletErc20Storage) error

//WithTableName return Option to set trade table Name
func WithTableName(name string) Option {
	return func(ws *WalletErc20Storage) error {
		ws.tableName = name
		return nil
	}
}

//WalletErc20Storage defines the object to store WalletErc20 data
type WalletErc20Storage struct {
	sugar     *zap.SugaredLogger
	db        *sqlx.DB
	tableName string
}

// NewDB return the WalletErc20Storage instance. User must call Close() before exit.
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB, options ...Option) (*WalletErc20Storage, error) {
	var (
		logger = sugar.With("func", "accounting/wallet-erc20/storage/postgres.NewDB")
		//set Default table name
	)

	const schemaFMT = `
	CREATE TABLE IF NOT EXISTS %[1]s
(
	id BIGSERIAL NOT NULL,
	data JSONB NOT NULL UNIQUE,
	CONSTRAINT %[1]s_pk PRIMARY KEY(id)
) ;

CREATE INDEX IF NOT EXISTS %[1]s_time_idx ON %[1]s 
((data ->> 'timestamp'),(data ->> 'contractAddress'),(data ->> 'from'),(data ->> 'from'));
`

	ws := &WalletErc20Storage{
		sugar:     sugar,
		db:        db,
		tableName: walletErc20TxsTableName,
	}
	for _, opt := range options {
		if err := opt(ws); err != nil {
			return nil, err
		}
	}

	query := fmt.Sprintf(schemaFMT, ws.tableName)
	logger.Debugw("initializing database schema", "query", query)

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")
	return ws, nil
}

//TearDown removes all the tables
func (wdb *WalletErc20Storage) TearDown() error {
	const dropFMT = `
	DROP TABLE %s;
	`
	query := fmt.Sprintf(dropFMT, wdb.tableName)
	wdb.sugar.Debugw("tearingdown", "query", dropFMT, "table name", wdb.tableName)
	_, err := wdb.db.Exec(query)
	return err
}

//Close close DB connection
func (wdb *WalletErc20Storage) Close() error {
	if wdb.db != nil {
		return wdb.db.Close()
	}
	return nil
}

//UpdateERC20Transfers store the ERC20Transfer rate at that blockInfo
func (wdb *WalletErc20Storage) UpdateERC20Transfers(erc20Txs []common.ERC20Transfer) error {
	var (
		logger = wdb.sugar.With(
			"func", "accounting/wallet-erc20/storage/postgres..UpdateRatesRecords",
			"len(erc20Txs)", len(erc20Txs),
		)
	)
	const updateStmt = `INSERT INTO %[1]s(data)
	VALUES ( 
		$1
	)
	ON CONFLICT ON CONSTRAINT %[1]s_data_key DO NOTHING
	`
	query := fmt.Sprintf(updateStmt,
		wdb.tableName,
	)
	logger.Debugw("updating ERC20transfers...", "query", query)

	tx, err := wdb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)
	for _, erc20Tx := range erc20Txs {
		var data []byte
		data, err = json.Marshal(erc20Tx)
		if err != nil {
			return err
		}
		logger.Debugf("%s contract %s", data, erc20Tx.ContractAddress.Hex())
		_, err = tx.Exec(query, data)
		if err != nil {
			return err
		}
	}

	return err
}

//GetERC20Transfers return erc20 transfer between from.. to.. in its json []byte form
func (wdb *WalletErc20Storage) GetERC20Transfers(wallet, token ethereum.Address, from, to time.Time) ([]common.ERC20Transfer, error) {
	var (
		dbResult [][]byte
		result   []common.ERC20Transfer
		logger   = wdb.sugar.With(
			"func", "accounting/wallet-erc20/storage/postgres..UpdateRatesRecords",
			"from", from.UTC(),
			"to", to.UTC(),
			"wallet", wallet.Hex(),
			"token", token.Hex(),
		)
		tmp common.ERC20Transfer
	)
	const selectStmt = `SELECT data FROM %[1]s WHERE ((data->>'timestamp')>=$1::text AND (data->>'timestamp')<$2::text)`
	query := fmt.Sprintf(selectStmt, wdb.tableName)
	if !blockchain.IsZeroAddress(wallet) {
		query += fmt.Sprintf(` AND (data->>'from'='%[1]s' OR data->>'to'='%[1]s')`, wallet.Hex())
	}

	if !blockchain.IsZeroAddress(token) {
		query += fmt.Sprintf(` AND data->>'contractAddress'='%[1]s'`, token.Hex())
	}

	logger.Debugw("querying ERC20 transfers history...", "query", query)
	if err := wdb.db.Select(&dbResult, query, timeutil.TimeToTimestampMs(from), timeutil.TimeToTimestampMs(to)); err != nil {
		return result, err
	}
	logger.Debugw("result", "len", len(dbResult))
	for _, data := range dbResult {
		if err := json.Unmarshal(data, &tmp); err != nil {
			return result, err
		}
		result = append(result, tmp)
	}
	return result, nil
}

//GetLastStoredBlock return last stored block of a wallet
//return sql.NoRows error if there is no row
func (wdb *WalletErc20Storage) GetLastStoredBlock(wallet ethereum.Address) (int, error) {
	var (
		result int
		logger = wdb.sugar.With(
			"func", "accounting/wallet-erc20/storage/postgres.GetLastStoredBlock",
			"wallet", wallet.Hex(),
		)
	)
	const selectStmt = `SELECT data->>'blockNumber' FROM %[1]s WHERE data->>'timestamp'=
	(SELECT MAX(data->>'timestamp') FROM %[1]s WHERE (data->>'from'=$1 OR data->>'to'=$1)) LIMIT 1 
	`
	query := fmt.Sprintf(selectStmt, wdb.tableName)
	logger.Debugw("querying trade history...", "query", query)
	if err := wdb.db.Get(&result, query, wallet.Hex()); err != nil {
		return result, err
	}
	return result, nil
}
