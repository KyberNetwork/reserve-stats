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
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	walletErc20TxsTableName = "rsv_tx_erc20"
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
	data JSONB NOT NULL UNIQUE,
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
	const selectStmt = `SELECT data FROM %[1]s WHERE ((data->>'timestamp')>=$1::text AND (data->>'timestamp')<$2::text) AND
	($3 OR (data->>'from'=$4 OR data->>'to'=$4)) AND
	($5 OR data->>'contractAddress'=$6)`
	query := fmt.Sprintf(selectStmt, wdb.tableName)
	logger.Debugw("querying ERC20 transfers history...", "query", query)
	walletFilter := blockchain.IsZeroAddress(wallet)
	tokenFilter := blockchain.IsZeroAddress(token)
	if err := wdb.db.Select(&dbResult, query, timeutil.TimeToTimestampMs(from), timeutil.TimeToTimestampMs(to), walletFilter, wallet.Hex(), tokenFilter, token.Hex()); err != nil {
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
