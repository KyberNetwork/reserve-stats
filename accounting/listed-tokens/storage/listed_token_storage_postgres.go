package storage

import (
	"fmt"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

//ListedTokenDB is storage for listed token
type ListedTokenDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//NewDB open a new database connection an create initiated table if it is not exist
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*ListedTokenDB, error) {
	var (
		logger = sugar.With("func", "accounting/storage.NewDB")
		ltd    = &ListedTokenDB{
			sugar: sugar,
			db:    db,
		}
	)

	logger.Debugw("initializing database schema", "query", schemaFmt)
	if _, err := db.Exec(schemaFmt); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")
	return ltd, nil
}

//CreateOrUpdate add or edit an record in the tokens table
func (ltd *ListedTokenDB) CreateOrUpdate(tokens []common.ListedToken, blockNumber *big.Int, reserve ethereum.Address) (err error) {
	var (
		logger  = ltd.sugar.With("func", "accounting/listed_tokens/storage/ltd.CreateOrUpdate")
		changed = false
	)
	saveTokenQuery := fmt.Sprintf(`SELECT save_token($1, $2, $3, $4, $5, $6)`)

	updateVersionQuery := `UPDATE "listed_tokens_version"
SET version      = CASE WHEN $1 THEN version + 1 ELSE version END,
    block_number = $2;`
	logger.Debugw("update version", "query", updateVersionQuery)

	tx, err := ltd.db.Beginx()
	if err != nil {
		return
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)

	for _, token := range tokens {
		var dbChanged = false
		if err = tx.Get(&dbChanged, saveTokenQuery,
			token.Address.Hex(),
			token.Name,
			token.Symbol,
			token.Timestamp.UTC(),
			nil,
			reserve.Hex()); err != nil {
			return
		}

		if dbChanged {
			changed = true
		}

		for _, oldToken := range token.Old {
			if err = tx.Get(&dbChanged, saveTokenQuery,
				oldToken.Address.Hex(),
				token.Name,
				token.Symbol,
				oldToken.Timestamp.UTC(),
				token.Address.Hex(),
				reserve.Hex()); err != nil {
				return
			}

			if dbChanged {
				changed = true
			}
		}
	}

	if changed {
		if _, err = tx.Exec(updateVersionQuery, changed, blockNumber.Uint64()); err != nil {
			return
		}
	}

	return
}

type listedTokenRecord struct {
	Address       string         `db:"address"`
	Symbol        string         `db:"symbol"`
	Name          string         `db:"name"`
	Timestamp     time.Time      `db:"timestamp"`
	OldAddresses  pq.StringArray `db:"old_addresses"`
	OldTimestamps pq.Int64Array  `db:"old_timestamps"`
}

type listedTokenVersion struct {
	Version     uint64 `db:"version"`
	BlockNumber uint64 `db:"block_number"`
}

// ListedToken converts listedTokenRecord instance to a common.ListedToken.
func (r *listedTokenRecord) ListedToken() (common.ListedToken, error) {
	token := common.ListedToken{
		Address:   ethereum.HexToAddress(r.Address),
		Symbol:    r.Symbol,
		Name:      r.Name,
		Timestamp: r.Timestamp.UTC(),
	}

	if len(r.OldAddresses) != len(r.OldTimestamps) {
		return common.ListedToken{}, fmt.Errorf(
			"malformed old data record: old_addresses=%d, old_timestamps=%d",
			len(r.OldAddresses), len(r.OldTimestamps))
	}

	for i := range r.OldAddresses {
		oldToken := common.OldListedToken{
			Address:   ethereum.HexToAddress(r.OldAddresses[i]),
			Timestamp: timeutil.TimestampMsToTime(uint64(r.OldTimestamps[i])).UTC(),
		}
		if token.Old == nil {
			token.Old = []common.OldListedToken{oldToken}
		} else {
			token.Old = append(token.Old, oldToken)
		}
	}
	return token, nil
}

// GetTokens return all tokens listed
func (ltd *ListedTokenDB) GetTokens(reserve ethereum.Address) (result []common.ListedToken, version, blockNumber uint64, err error) {
	var (
		logger = ltd.sugar.With(
			"func",
			"accounting/listed-token-storage/listedtokenstorage.GetTokens",
			"reserve", reserve,
		)
		records       []listedTokenRecord
		versionRecord listedTokenVersion
	)

	getQuery := `SELECT address,
       name,
       symbol,
       timestamp,
       old_addresses,
	   old_timestamps
FROM "tokens_view" WHERE ( $1 OR reserve_address = $2 );`
	logger.Debugw("get tokens query", "query", getQuery)

	getVersionQuery := `SELECT version, block_number FROM "listed_tokens_version" LIMIT 1`
	logger.Debugw("get token version", "query", getVersionQuery)

	tx, err := ltd.db.Beginx()
	if err != nil {
		return nil, 0, 0, err
	}

	defer pgsql.CommitOrRollback(tx, logger, &err)

	if err := tx.Select(&records, getQuery, blockchain.IsZeroAddress(reserve), reserve.Hex()); err != nil {
		logger.Errorw("error query token", "error", err)
		return nil, 0, 0, err
	}

	for _, record := range records {
		token, err := record.ListedToken()
		if err != nil {
			return nil, 0, 0, err
		}
		result = append(result, token)
	}

	if err := tx.Get(&versionRecord, getVersionQuery); err != nil {
		logger.Error("error query token version", "error", err)
		return nil, 0, 0, err
	}

	version = versionRecord.Version
	blockNumber = versionRecord.BlockNumber

	return result, version, blockNumber, nil
}

//Close db connection
func (ltd *ListedTokenDB) Close() error {
	if ltd.db != nil {
		return ltd.db.Close()
	}
	return nil
}

//DeleteTable remove tables use for test
func (ltd *ListedTokenDB) DeleteTable() error {
	const dropQuery = `DROP VIEW "tokens_view";
drop function save_token(text, text, text, timestamp, text, text);
DROP TABLE "listed_tokens_reserves_tokens", "listed_tokens", "listed_tokens_version", "listed_tokens_reserves";`
	ltd.sugar.Infow("Drop token table", "query", dropQuery)
	_, err := ltd.db.Exec(dropQuery)
	return err
}
