package postgres

import (
	"database/sql"
	"encoding/json"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

// Storage is the an implementation of storage.ReserveTransactionStorage interface using PostgresQL.
type Storage struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

// NewStorage creates new instance of Storage.
func NewStorage(sugar *zap.SugaredLogger, db *sqlx.DB) (*Storage, error) {
	var (
		logger = sugar.With("func", "accounting/reserve-transaction-fetcher/storage/postgres/NewStorage")
	)
	const schemaFmt = `
	-- create table tx normal
	CREATE TABLE IF NOT EXISTS "rsv_tx_normal"
(
	  id SERIAL PRIMARY KEY,
    tx_hash text  UNIQUE NOT NULL,
    data    JSONB NOT NULL
);
CREATE INDEX IF NOT EXISTS "rsv_tx_normal_time_idx" ON "rsv_tx_normal" ((data ->> 'timestamp'));

-- create table tx internal
CREATE TABLE IF NOT EXISTS "rsv_tx_internal"
(
		id SERIAL PRIMARY KEY,
    data JSONB NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS "rsv_tx_internal_time_idx" ON "rsv_tx_internal" ((data ->> 'timestamp'));

-- create table tx erc20
CREATE TABLE IF NOT EXISTS "rsv_tx_erc20"
(
	  id SERIAL PRIMARY KEY,
		data JSONB NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS "rsv_tx_erc20_time_idx" ON "rsv_tx_erc20" ((data ->> 'timestamp'),
(data ->> 'contractAddress'),(data ->> 'from'),(data ->> 'to'));

-- create table reserves
CREATE TABLE IF NOT EXISTS "rsv_tx_reserve"
(
	address text NOT NULL PRIMARY KEY,
	address_type text NOT NULL
);

-- create table last inserted
CREATE TABLE IF NOT EXISTS "rsv_tx_last_inserted"
(
	address_key text UNIQUE REFERENCES "rsv_tx_reserve" (address),
	last_inserted BIGINT NOT NULL
);

-- create table link from tx to reserve
CREATE TABLE IF NOT EXISTS "rsv_tx_normal_tx_reserve"
(
	tx_id int REFERENCES "rsv_tx_normal" (id),
	address_key text REFERENCES "rsv_tx_reserve" (address),
	PRIMARY KEY (tx_id, address_key)
);

-- create table link from tx to reserve
CREATE TABLE IF NOT EXISTS rsv_tx_internal_tx_reserve
(
	tx_id int REFERENCES "rsv_tx_internal" (id),
	address_key text REFERENCES "rsv_tx_reserve" (address),
	PRIMARY KEY (tx_id, address_key)
);

-- create table link from tx to reserve
CREATE TABLE IF NOT EXISTS rsv_tx_erc20_tx_reserve
(
	tx_id int REFERENCES "rsv_tx_erc20" (id),
	address_key text REFERENCES "rsv_tx_reserve" (address),
	PRIMARY KEY (tx_id, address_key)
);
`

	s := &Storage{sugar: sugar, db: db}

	logger.Infow("initializing database schema", "query", schemaFmt)
	if _, err := db.Exec(schemaFmt); err != nil {
		return nil, err
	}
	return s, nil
}

//StoreReserve save fetching reserve address into database
func (s *Storage) StoreReserve(reserve ethereum.Address, reserveType string) error {
	var (
		logger = s.sugar.With("func", "accounting/reserve-transaction-fetcher/storage/postgres/Storage.StoreReserve")
	)
	const storeReserve = `INSERT INTO "rsv_tx_reserve" (address, address_type)
	VALUES ($1, $2) 
	ON CONFLICT (address) DO UPDATE SET address_type = EXCLUDED.address_type;`

	logger.Debugw("query to store reserve into database", "query", storeReserve)

	if _, err := s.db.Exec(storeReserve, reserve.Hex(), reserveType); err != nil {
		return err
	}
	return nil
}

//StoreNormalTx store normal tx
func (s *Storage) StoreNormalTx(txs []common.NormalTx, reserve ethereum.Address) (err error) {
	var (
		logger = s.sugar.With("func", "accounting/reserve-transaction-fetcher/storage/postgres/Storage.StoreNormalTx")
		id     int64
	)
	const (
		updateStmt = `INSERT INTO "rsv_tx_normal"(tx_hash, data)
VALUES ($1, $2)
ON CONFLICT (tx_hash) DO UPDATE SET data = EXCLUDED.data RETURNING id;
`
		insertStmt = `INSERT INTO "rsv_tx_normal_tx_reserve" (tx_id, address_key)
	VALUES ($1, $2)
		ON CONFLICT DO NOTHING;`
	)

	logger.Debugw("storing normal transactions to database", "query", updateStmt)

	tx, err := s.db.Beginx()
	if err != nil {
		return
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)
	for _, t := range txs {
		var data []byte
		data, err = json.Marshal(t)
		if err != nil {
			return
		}
		if err = tx.Get(&id, updateStmt, t.BlockHash, data); err != nil && err != sql.ErrNoRows {
			return
		}
		if err == sql.ErrNoRows {
			err = nil
		}

		// insert into txs_reserves
		if id != 0 {
			if _, err = tx.Exec(insertStmt, id, reserve.Hex()); err != nil {
				return
			}
		}

	}

	return nil
}

//GetNormalTx get normal tx between a certain period of time
func (s *Storage) GetNormalTx(from time.Time, to time.Time) ([]common.NormalTx, error) {
	var (
		logger = s.sugar.With(
			"func", "accounting/reserve-transaction-fetcher/storage/postgres/Storage.GetNormalTx",
			"from", from.String(),
			"to", to.String(),
		)
		dbResult [][]byte
		results  []common.NormalTx
		t        common.NormalTx
	)
	const selectStmt = `SELECT data
FROM "rsv_tx_normal"
WHERE data ->> 'timestamp' >= $1
  AND data ->> 'timestamp' < $2`
	logger.Debugw("querying normal transactions from database", "query", selectStmt)
	if err := s.db.Select(
		&dbResult,
		selectStmt,
		timeutil.TimeToTimestampMs(from),
		timeutil.TimeToTimestampMs(to)); err != nil {
		return nil, err
	}
	for _, data := range dbResult {
		if err := json.Unmarshal(data, &t); err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	return results, nil
}

//StoreInternalTx stores internal tx
func (s *Storage) StoreInternalTx(txs []common.InternalTx, reserve ethereum.Address) (err error) {
	var (
		logger = s.sugar.With(
			"func", "accounting/reserve-transaction-fetcher/storage/postgres/Storage.StoreInternalTx",
		)
		id int64
	)

	const (
		updateStmt = `INSERT INTO "rsv_tx_internal"(data)
VALUES ($1)
ON CONFLICT DO NOTHING RETURNING id;
`
		insertStmt = `INSERT INTO "rsv_tx_internal_tx_reserve" (tx_id, address_key)
	VALUES ($1, $2)
		ON CONFLICT DO NOTHING;`
	)

	logger.Debugw("storing internal transactions to database", "query", updateStmt)

	tx, err := s.db.Beginx()
	if err != nil {
		return
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)
	for _, t := range txs {
		var data []byte
		data, err = json.Marshal(t)
		if err != nil {
			return
		}

		if err = tx.Get(&id, updateStmt, data); err != nil && err != sql.ErrNoRows {
			return
		}
		if err == sql.ErrNoRows {
			err = nil
		}
		// insert into txs_reserves
		if id != 0 {
			if _, err = tx.Exec(insertStmt, id, reserve.Hex()); err != nil {
				return
			}
		}
	}
	return
}

//GetInternalTx get internal txs between a period of time
func (s *Storage) GetInternalTx(from time.Time, to time.Time) ([]common.InternalTx, error) {
	var (
		logger = s.sugar.With(
			"func", "accounting/reserve-transaction-fetcher/storage/postgres/Storage.GetInternalTx",
			"from", from.String(),
			"to", to.String(),
		)
		dbResult [][]byte
		results  []common.InternalTx
		t        common.InternalTx
	)
	const selectStmt = `SELECT data
FROM "rsv_tx_internal"
WHERE data ->> 'timestamp' >= $1
  AND data ->> 'timestamp' < $2`
	logger.Debugw("querying internal transactions from database", "query", selectStmt)
	if err := s.db.Select(
		&dbResult,
		selectStmt,
		timeutil.TimeToTimestampMs(from),
		timeutil.TimeToTimestampMs(to)); err != nil {
		return nil, err
	}
	for _, data := range dbResult {
		if err := json.Unmarshal(data, &t); err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	return results, nil
}

//StoreERC20Transfer save ERC20 transfer
func (s *Storage) StoreERC20Transfer(txs []common.ERC20Transfer, reserve ethereum.Address) (err error) {
	var (
		logger = s.sugar.With(
			"func", "accounting/reserve-transaction-fetcher/storage/postgres/Storage.StoreERC20Transfer",
		)
		id int64
	)

	const (
		updateStmt = `INSERT INTO "rsv_tx_erc20"(data)
VALUES ($1)
ON CONFLICT DO NOTHING RETURNING id;
`
		insertStmt = `INSERT INTO "rsv_tx_erc20_tx_reserve" (tx_id, address_key)
	VALUES ($1, $2)
		ON CONFLICT DO NOTHING;`
	)

	logger.Debugw("storing ERC20 transfers to database", "query", updateStmt)

	tx, err := s.db.Beginx()
	if err != nil {
		return
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)
	for _, t := range txs {
		var data []byte
		data, err = json.Marshal(t)
		if err != nil {
			return
		}

		// insert into rsv_erc20_txs
		if err = tx.Get(&id, updateStmt, data); err != nil && err != sql.ErrNoRows {
			return
		}

		if err == sql.ErrNoRows {
			err = nil
		}

		// insert into txs_reserves
		if id != 0 {
			if _, err = tx.Exec(insertStmt, id, reserve.Hex()); err != nil {
				return
			}
		}
	}
	return
}

//GetERC20Transfer get ERC20 transfer between a period of time
func (s *Storage) GetERC20Transfer(from time.Time, to time.Time) ([]common.ERC20Transfer, error) {
	var (
		logger = s.sugar.With(
			"func", "accounting/reserve-transaction-fetcher/storage/postgres/Storage.GetERC20Transfer",
			"from", from.String(),
			"to", to.String(),
		)
		dbResult [][]byte
		results  []common.ERC20Transfer
		t        common.ERC20Transfer
	)
	const selectStmt = `SELECT data
FROM "rsv_tx_erc20"
JOIN "rsv_tx_erc20_tx_reserve" AS a ON a.tx_id = rsv_tx_erc20.id
JOIN "rsv_tx_reserve" AS reserve ON a.address_key = reserve.address 
WHERE data ->> 'timestamp' >= $1
	AND data ->> 'timestamp' < $2
	AND reserve.address_type <> $3`
	logger.Debugw("querying ERC20 transfers from database", "query", selectStmt)
	if err := s.db.Select(
		&dbResult,
		selectStmt,
		timeutil.TimeToTimestampMs(from),
		timeutil.TimeToTimestampMs(to),
		common.CompanyWallet.String()); err != nil {
		return nil, err
	}
	for _, data := range dbResult {
		if err := json.Unmarshal(data, &t); err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	return results, nil
}

//StoreLastInserted save last insert address and block number where it's last inserted
func (s *Storage) StoreLastInserted(addr ethereum.Address, blockNumber *big.Int) error {
	var (
		logger = s.sugar.With(
			"func", "accounting/reserve-transaction-fetcher/storage/postgres/Storage.StoreLastInserted",
			"address", addr.Hex(),
			"block_number", blockNumber.String(),
		)
	)
	const queryFmt = `INSERT INTO "rsv_tx_last_inserted"(address_key, last_inserted)
VALUES ($1, $2)
ON CONFLICT (address_key) DO UPDATE SET last_inserted = EXCLUDED.last_inserted;
`

	logger.Debugw("updating last inserted to database")

	_, err := s.db.Exec(queryFmt, addr.String(), blockNumber.Uint64())
	return err
}

//GetLastInserted return last inserted block of an address
func (s *Storage) GetLastInserted(addr ethereum.Address) (*big.Int, error) {
	var (
		logger = s.sugar.With(
			"func", "accounting/reserve-transaction-fetcher/storage/postgres/Storage.GetLastInserted",
			"address", addr.Hex(),
		)
		lastInserted uint64
	)
	const queryFmt = `SELECT last_inserted
FROM "rsv_tx_last_inserted"
WHERE address_key ILIKE $1`
	logger.Debugw("fetching last inserted to database")
	err := s.db.Get(&lastInserted, queryFmt, addr.String())
	switch err {
	case sql.ErrNoRows:
		logger.Infow("no last inserted record exists")
		return nil, nil
	case nil:
		return big.NewInt(0).SetUint64(lastInserted), nil
	default:
		return nil, err
	}
}

//GetWalletERC20Transfers return erc20 transfer between from.. to.. in its json []byte form
func (s *Storage) GetWalletERC20Transfers(wallet, token ethereum.Address, from, to time.Time) ([]common.ERC20Transfer, error) {
	var (
		dbResult [][]byte
		result   []common.ERC20Transfer
		logger   = s.sugar.With(
			"func", "accounting/wallet-erc20/storage/postgres..UpdateRatesRecords",
			"from", from.UTC(),
			"to", to.UTC(),
			"wallet", wallet.Hex(),
			"token", token.Hex(),
		)
		tmp common.ERC20Transfer
	)
	const selectStmt = `SELECT data FROM rsv_tx_erc20 
	JOIN "rsv_tx_erc20_tx_reserve" as a ON a.tx_id = rsv_tx_erc20.id
	JOIN "rsv_tx_reserve" as reserve ON a.address_key = reserve.address WHERE ((data->>'timestamp')>=$1::text AND (data->>'timestamp')<$2::text) AND
	($3 OR (data->>'from'=$4 OR data->>'to'=$4)) AND
	($5 OR data->>'contractAddress'=$6)
	AND reserve.address_type = $7`
	logger.Debugw("querying ERC20 transfers history...", "query", selectStmt)
	walletFilter := blockchain.IsZeroAddress(wallet)
	tokenFilter := blockchain.IsZeroAddress(token)
	if err := s.db.Select(&dbResult, selectStmt, timeutil.TimeToTimestampMs(from), timeutil.TimeToTimestampMs(to),
		walletFilter, wallet.Hex(), tokenFilter, token.Hex(), common.CompanyWallet.String()); err != nil {
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
