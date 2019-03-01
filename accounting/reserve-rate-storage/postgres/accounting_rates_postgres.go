package postgres

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
)

const (
	reserveTableName = "reserves"
	tokenTableName   = "tokens"
	baseTableName    = "bases"
	rateTableName    = "rates"
)

type RatesStorage struct {
	sugar      *zap.SugaredLogger
	db         *sqlx.DB
	tableNames []string
}

// NewDB return the Ratestorage instance. User must call ratestorage.Close() before exit.
// tableNames is a list of 4 string for 4 tablename[reserve,token,base, rate]. It can be optional
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB, customTableNames ...string) (*RatesStorage, error) {
	const schemaFMT = `CREATE TABLE IF NOT EXISTS %[1]s
(
	id serial NOT NULL,
	address TEXT NOT NULL UNIQUE,
	CONSTRAINT %[1]s_pk PRIMARY KEY(id)
) ;

CREATE TABLE IF NOT EXISTS %[2]s
(
	id serial NOT NULL,
	symbol TEXT NOT NULL UNIQUE,
	CONSTRAINT %[2]s_pk PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS %[3]s
(
	id serial NOT NULL,
	symbol TEXT NOT NULL UNIQUE,
	CONSTRAINT %[3]s_pk PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS %[4]s
(
	id serial NOT NULL,
	time TIMESTAMP NOT NULL,
	token_id serial NOT NULL,
    base_id serial NOT NULL,
	block integer NOT NULL,
	buy_reserve_rate float8 NOT NULL,
	sell_reserve_rate float8 NOT NULL,
	reserve_id integer NOT NULL,
	CONSTRAINT %[4]s_pk PRIMARY KEY(id),
	CONSTRAINT %[4]s_fk_token_id FOREIGN KEY(token_id) REFERENCES %[2]s(id),
    CONSTRAINT %[4]s_fk_base_id FOREIGN KEY(base_id) REFERENCES %[3]s(id),
	CONSTRAINT %[4]s_fk_reseve_id FOREIGN KEY(reserve_id) REFERENCES %[1]s(id)
);
`
	var (
		logger     = sugar.With("func", "reserverates/storage/postgres")
		tableNames []string
	)

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)
	if len(customTableNames) > 0 {
		if len(customTableNames) != 4 {
			return nil, fmt.Errorf("expect 4 tables name [reserve,token,base,rates], got %v", customTableNames)
		}
		tableNames = customTableNames
	} else {
		tableNames = []string{reserveTableName, tokenTableName, baseTableName, rateTableName}
	}
	query := fmt.Sprintf(schemaFMT, tableNames[0], tableNames[1], tableNames[2], tableNames[3])
	logger.Debugw("initializing database schema", "query", query)

	if _, err = tx.Exec(query); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")
	return &RatesStorage{
		sugar:      sugar,
		db:         db,
		tableNames: tableNames,
	}, nil
}

//TearDown removes all the tables
func (rdb *RatesStorage) TearDown() error {
	const dropFMT = `
	DROP TABLE %[1]s,%[2]s,%[3]s,%[4]s;
	`
	tx, err := rdb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, rdb.sugar, &err)
	query := fmt.Sprintf(dropFMT,
		rdb.tableNames[0],
		rdb.tableNames[1],
		rdb.tableNames[2],
		rdb.tableNames[3],
	)
	_, err = tx.Exec(query)
	return err
}

//Close close DB connection
func (rdb *RatesStorage) Close() error {
	return rdb.db.Close()
}

//getTokenBase from pair return token and base from base-token pair
func getTokenBaseFromPair(pair string) (string, string, error) {
	ids := strings.Split(strings.TrimSpace(pair), "-")
	if len(ids) != 2 {
		return "", "", fmt.Errorf("Pair %s is malformed. Expected token_base, got", pair)
	}
	return ids[1], ids[0], nil
}

func (rdb *RatesStorage) UpdateRatesRecords(blockInfo lastblockdaily.BlockInfo, rateRecords map[string]map[string]common.ReserveRateEntry) error {
	var logger = rdb.sugar.With(
		"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
		"block_number", blockInfo.Block,
		"timestamp", blockInfo.Timestamp.String(),
	)
	const updateStmt = `
	WITH
tk AS (
	INSERT INTO %[1]s(symbol)
	VALUES($1)
	ON CONFLICT ON CONSTRAINT %[1]s_symbol_key DO UPDATE SET symbol=EXCLUDED.symbol RETURNING %[1]s.id
),
bs AS( 
 	INSERT INTO %[2]s(symbol) 
 	VALUES($2) 
 	ON CONFLICT ON CONSTRAINT %[2]s_symbol_key DO UPDATE SET symbol=EXCLUDED.symbol RETURNING %[2]s.id
),
rs AS(
	INSERT INTO %[3]s(address)
	VALUES ($3)
 	ON CONFLICT ON CONSTRAINT %[3]s_address_key DO UPDATE SET address=EXCLUDED.address RETURNING %[3]s.id
)

INSERT INTO %[4]s(time, "token_id", "base_id", block, buy_reserve_rate, sell_reserve_rate, "reserve_id")
VALUES (
	$4, 
	(SELECT id FROM tk),
	(SELECT id FROM bs),
	$5, 
	$6, 
	$7, 
	(SELECT id FROM rs)
	);`
	query := fmt.Sprintf(updateStmt,
		rdb.tableNames[1],
		rdb.tableNames[2],
		rdb.tableNames[0],
		rdb.tableNames[3],
	)

	logger.Debugw("updating rates...", "query", query)
	tx, err := rdb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, rdb.sugar, &err)
	for rsvAddr, rateRecord := range rateRecords {
		for pair, rate := range rateRecord {
			token, base, err := getTokenBaseFromPair(pair)
			if err != nil {
				return err
			}
			_, err = tx.Exec(query,
				token,
				base,
				rsvAddr,
				blockInfo.Timestamp,
				blockInfo.Block,
				rate.BuyReserveRate,
				rate.SellReserveRate,
			)
			if err != nil {
				return err
			}
		}
	}
	return err
}
