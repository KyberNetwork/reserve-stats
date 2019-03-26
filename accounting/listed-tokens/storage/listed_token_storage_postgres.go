package storage

import (
	"fmt"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
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
)
	`
	var logger = sugar.With("func", "accounting/storage.NewDB")

	logger.Debug("initializing database schema")
	if _, err := db.Exec(fmt.Sprintf(schemaFmt, tableName)); err != nil {
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
func (ltd *ListedTokenDB) CreateOrUpdate(tokens []common.ListedToken) (err error) {
	var (
		logger = ltd.sugar.With("func", "accounting/lisetdtokenstorage.CreateOrUpdate")
	)
	upsertQuery := fmt.Sprintf(`INSERT INTO "%[1]s" (address, name, symbol, timestamp, parent_id)
VALUES ($1,
        $2,
        $3,
        $4,
        CASE WHEN $5::text IS NOT NULL THEN (SELECT id FROM "%[1]s" WHERE address = $5) ELSE NULL END)
ON CONFLICT (address) DO UPDATE SET parent_id = EXCLUDED.parent_id`,
		ltd.tableName)

	logger.Debugw("upsert token", "value", upsertQuery)

	tx, err := ltd.db.Beginx()
	if err != nil {
		return
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)

	for _, token := range tokens {
		if _, err = tx.Exec(upsertQuery,
			token.Address.Hex(),
			token.Name,
			token.Symbol,
			token.Timestamp.UTC(),
			nil); err != nil {
			return
		}

		for _, oldToken := range token.Old {
			if _, err = tx.Exec(upsertQuery,
				oldToken.Address.Hex(),
				token.Name,
				token.Symbol,
				oldToken.Timestamp.UTC(),
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
func (ltd *ListedTokenDB) GetTokens() ([]common.ListedToken, error) {
	var (
		logger = ltd.sugar.With(
			"func",
			"accounting/listed-token-storage/listedtokenstorage.GetTokens",
		)
		result  []common.ListedToken
		records []listedTokenRecord
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

	if err := ltd.db.Select(&records, getQuery); err != nil {
		logger.Errorw("error query token", "error", err)
		return nil, err
	}

	logger.Debugw("result from listed token", "result", records)

	for _, record := range records {
		token, err := record.ListedToken()
		if err != nil {
			return nil, err
		}
		result = append(result, token)
	}

	return result, nil
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
	const dropQuery = `DROP TABLE %s;`
	query := fmt.Sprintf(dropQuery, ltd.tableName)

	ltd.sugar.Infow("Drop token table", "query", query)
	_, err := ltd.db.Exec(query)
	return err
}
