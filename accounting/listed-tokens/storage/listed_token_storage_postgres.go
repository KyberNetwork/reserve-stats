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
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	versionTable        = "listed_token_version"
	reservesTable       = "reserves"
	reservesTokensTable = "reserves_tokens"
)

//ListedTokenDB is storage for listed token
type ListedTokenDB struct {
	sugar     *zap.SugaredLogger
	db        *sqlx.DB
	tableName string
}

//NewDB open a new database connection an create initiated table if it is not exist
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB, tableName string) (*ListedTokenDB, error) {
	const schemaFmt = `CREATE TABLE IF NOT EXISTS "%[1]s"
(
	id SERIAL PRIMARY KEY,
	address text NOT NULL UNIQUE,
	name text NOT NULL,
	symbol text NOT NULL,
	timestamp TIMESTAMP NOT NULL,
	parent_id INT REFERENCES "%[1]s" (id)
);
CREATE TABLE IF NOT EXISTS "%[2]s"
(
	id SERIAL PRIMARY KEY,
	version INT NOT NULL,
	block_number bigint NOT NULL,
	reserve text NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS "%[3]s"
(
	id serial NOT NULL,
	address TEXT NOT NULL UNIQUE,
	CONSTRAINT %[1]s_pk PRIMARY KEY(id)
);
CREATE TABLE IF NOT EXISTS "%[4]s"
(
	id SERIAL PRIMARY KEY,
	token_id INT REFERENCE "%[1]s" (id),
	reserve_id INT REFERENCE "%[4]s" (id)
)
	`
	var logger = sugar.With("func", "accounting/storage.NewDB")

	logger.Debug("initializing database schema")
	if _, err := db.Exec(fmt.Sprintf(schemaFmt, tableName, versionTable, reservesTable, reservesTokensTable)); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")

	return &ListedTokenDB{
		sugar:     sugar,
		db:        db,
		tableName: tableName,
	}, nil
}

//CreateOrUpdate add or edit an record in the tokens table
func (ltd *ListedTokenDB) CreateOrUpdate(tokens []common.ListedToken, blockNumber *big.Int, reserve ethereum.Address) (err error) {
	var (
		logger             = ltd.sugar.With("func", "accounting/lisetdtokenstorage.CreateOrUpdate")
		tokenID, reserveID int64
	)
	upsertQuery := fmt.Sprintf(`INSERT INTO "%[1]s" (address, name, symbol, timestamp, parent_id)
VALUES ($1,
        $2,
        $3,
		$4,
        CASE WHEN $5::text IS NOT NULL THEN (SELECT id FROM "%[1]s" WHERE address = $5) ELSE NULL END)
ON CONFLICT (address) DO UPDATE SET parent_id = EXCLUDED.parent_id RETURNING id`,
		ltd.tableName)
	logger.Debugw("upsert token", "value", upsertQuery)

	updateVersionQuery := fmt.Sprintf(`INSERT INTO "%[1]s" (block_number, reserve)
	VALUES($1, $2)
	ON CONFLICT (reserve) DO UPDATE version = %[1]s.version+1, block_number = EXCLUDED.block_number`, versionTable)
	logger.Debugw("update version", "query", updateVersionQuery)

	updateReserveQuery := fmt.Sprintf(`INSERT INTO "%[1]s" (address) VALUE($1) ON CONFLICT (address) DO NOTHING RETURNING id`, reservesTable)
	logger.Debugw("update reserve", "query", updateReserveQuery)

	updateReserveTokenQuery := fmt.Sprintf(`INSERT INTO "%[1]s" (token_id, reserve_id) values($1, $2)`, reservesTokensTable)
	logger.Debugw("update reserve token", "query", updateReserveTokenQuery)

	tx, err := ltd.db.Beginx()
	if err != nil {
		return
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)

	if err = tx.Get(&reserveID, updateReserveQuery, reserve); err != nil {
		return
	}

	if _, err = tx.Exec(updateVersionQuery, blockNumber.Uint64(), reserve); err != nil {
		return
	}

	for _, token := range tokens {
		if err = tx.Get(&tokenID, upsertQuery,
			token.Address.Hex(),
			token.Name,
			token.Symbol,
			token.Timestamp.UTC(),
			reserve.Hex(),
			nil); err != nil {
			return
		}

		for _, oldToken := range token.Old {
			if _, err = tx.Exec(upsertQuery,
				oldToken.Address.Hex(),
				token.Name,
				token.Symbol,
				oldToken.Timestamp.UTC(),
				reserve.Hex(),
				token.Address.Hex()); err != nil {
				return
			}
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
func (ltd *ListedTokenDB) GetTokens() (result []common.ListedToken, version, blockNumber uint64, err error) {
	var (
		logger = ltd.sugar.With(
			"func",
			"accounting/listed-token-storage/listedtokenstorage.GetTokens",
		)
		records        []listedTokenRecord
		versionRecords []listedTokenVersion
	)

	getQuery := fmt.Sprintf(`SELECT joined.address,
       joined.name,
       joined.symbol,
       joined.timestamp,
       array_agg(joined.old_address) FILTER ( WHERE joined.old_address IS NOT NULL)::text[] AS old_addresses,
       array_agg(extract(EPOCH FROM joined.old_timestamp) * 1000)
                 FILTER ( WHERE joined.old_timestamp IS NOT NULL)::BIGINT[]                 AS old_timestamps
FROM (SELECT toks.address,
             toks.name,
             toks.symbol,
             toks.timestamp,
             olds.address   AS old_address,
             olds.timestamp AS old_timestamp
      FROM "%[1]s" AS toks
             LEFT JOIN "%[1]s" AS olds
                       ON toks.id = olds.parent_id
      WHERE toks.parent_id IS NULL
      ORDER BY timestamp DESC) AS joined
GROUP BY joined.address, joined.name, joined.symbol, joined.timestamp`, ltd.tableName)
	logger.Debugw("get tokens query", "query", getQuery)

	getVersionQuery := fmt.Sprintf(`SELECT version, block_number FROM "%[1]s" LIMIT 1`, versionTable)
	logger.Debugw("get token version", "query", getVersionQuery)

	tx, err := ltd.db.Beginx()
	if err != nil {
		return nil, 0, 0, err
	}

	defer pgsql.CommitOrRollback(tx, logger, &err)

	if err := tx.Select(&records, getQuery); err != nil {
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

	if err := tx.Select(&versionRecords, getVersionQuery); err != nil {
		logger.Error("error query token verson", "error", err)
		return nil, 0, 0, err
	}

	version = versionRecords[0].Version
	blockNumber = versionRecords[0].BlockNumber

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
	const dropQuery = `DROP TABLE %[1]s, %[2]s;`
	query := fmt.Sprintf(dropQuery, ltd.tableName, versionTable)

	ltd.sugar.Infow("Drop token table", "query", query)
	_, err := ltd.db.Exec(query)
	return err
}
