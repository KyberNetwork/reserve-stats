package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	storage "github.com/KyberNetwork/reserve-stats/accounting/reserve-rate-storage"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

const (
	reserveTableName = "reserves"
	tokenTableName   = "tokens"
	baseTableName    = "bases"
	rateTableName    = "rates"
	usdTableName     = "usds"
)

//RatesStorage defines the object to store rates
type RatesStorage struct {
	sugar      *zap.SugaredLogger
	db         *sqlx.DB
	tableNames []string
}

// NewDB return the Ratestorage instance. User must call ratestorage.Close() before exit.
// tableNames is a list of 5 string for 5 tablename[reserve,token,base, rate,usdrate]. It can be optional
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

CREATE TABLE IF NOT EXISTS %[5]s
(
	id serial NOT NULL,
	time TIMESTAMP NOT NULL,
	block integer NOT NULL,
	rate float8 NOT NULL,
	CONSTRAINT %[5]s_pk PRIMARY KEY(id)
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
		if len(customTableNames) != 5 {
			return nil, fmt.Errorf("expect 5 tables name [reserve,token,base,rates], got %v", customTableNames)
		}
		tableNames = customTableNames
	} else {
		tableNames = []string{reserveTableName, tokenTableName, baseTableName, rateTableName, usdTableName}
	}
	query := fmt.Sprintf(schemaFMT, tableNames[0], tableNames[1], tableNames[2], tableNames[3], tableNames[4])
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
	if rdb.db != nil {
		return rdb.db.Close()
	}
	return nil
}

//getTokenBase from pair return token and base from base-token pair
func getTokenBaseFromPair(pair string) (string, string, error) {
	ids := strings.Split(strings.TrimSpace(pair), "-")
	if len(ids) != 2 {
		return "", "", fmt.Errorf("Pair %s is malformed. Expected token_base, got", pair)
	}
	return ids[1], ids[0], nil
}

//UpdateRatesRecords update mutiple rate records from a block with mutiple reserve address into the DB
func (rdb *RatesStorage) UpdateRatesRecords(blockInfo lastblockdaily.BlockInfo, rateRecords map[string]map[string]common.ReserveRateEntry) error {
	var logger = rdb.sugar.With(
		"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
		"block_number", blockInfo.Block,
		"timestamp", blockInfo.Timestamp.String(),
	)
	const updateStmt = `
	WITH
rs AS(
	INSERT INTO %[1]s(address)
	VALUES ($1)
 	ON CONFLICT ON CONSTRAINT %[1]s_address_key DO UPDATE SET address=EXCLUDED.address RETURNING %[1]s.id
),
tk AS (
	INSERT INTO %[2]s(symbol)
	VALUES($2)
	ON CONFLICT ON CONSTRAINT %[2]s_symbol_key DO UPDATE SET symbol=EXCLUDED.symbol RETURNING %[2]s.id
),
bs AS( 
 	INSERT INTO %[3]s(symbol) 
 	VALUES($3) 
 	ON CONFLICT ON CONSTRAINT %[3]s_symbol_key DO UPDATE SET symbol=EXCLUDED.symbol RETURNING %[3]s.id
)

INSERT INTO %[4]s(time, token_id, base_id, block, buy_reserve_rate, sell_reserve_rate, reserve_id)
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
		rdb.tableNames[0],
		rdb.tableNames[1],
		rdb.tableNames[2],
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
				rsvAddr,
				token,
				base,
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

//GetRates return the reserve rates in a period of times
func (rdb *RatesStorage) GetRates(reservers []ethereum.Address, from time.Time, to time.Time) (map[string]storage.AccountingReserveRates, error) {
	var (
		result    = make(map[string]storage.AccountingReserveRates)
		rowData   = ratesFromDB{}
		rsvsAddrs = []string{}
		logger    = rdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.GetRates",
			"from", from.String(),
			"to", to.String(),
			"reservers", reservers,
		)
	)
	const (
		selectStmt = `SELECT rt.time as time, tk.symbol as token, bs.symbol as base, rt.buy_reserve_rate as buy_rate, rs.address as reserve
		FROM %[1]s AS rt LEFT JOIN %[2]s AS tk ON rt.token_id = tk.id
		LEFT JOIN %[3]s AS bs ON rt.base_id=bs.id 
		LEFT JOIN %[4]s AS rs ON rt.reserve_id=rs.id
		WHERE  time>=$1 AND time<$2 AND rs.address IN (SELECT unnest($3::text[]));`
		shortForm = "2006-01-02"
	)
	for _, rsv := range reservers {
		rsvsAddrs = append(rsvsAddrs, rsv.Hex())
	}
	query := fmt.Sprintf(selectStmt,
		rdb.tableNames[3],
		rdb.tableNames[1],
		rdb.tableNames[2],
		rdb.tableNames[0],
	)
	logger.Debugw("Querrying rate...", "query", query)

	rows, err := rdb.db.Queryx(query, from, to, pq.StringArray(rsvsAddrs))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.StructScan(&rowData); err != nil {
			return nil, err
		}
		if _, ok := result[rowData.Reserve]; !ok {
			result[rowData.Reserve] = map[time.Time]map[string]map[string]float64{
				rowData.Time: map[string]map[string]float64{
					rowData.Base: map[string]float64{
						rowData.Token: rowData.BuyRate,
					},
				},
			}
		}
		if _, ok := result[rowData.Reserve][rowData.Time]; !ok {
			result[rowData.Reserve][rowData.Time] = map[string]map[string]float64{
				rowData.Base: map[string]float64{
					rowData.Token: rowData.BuyRate,
				},
			}
		}
		if _, ok := result[rowData.Reserve][rowData.Time][rowData.Base]; !ok {
			result[rowData.Reserve][rowData.Time][rowData.Base] = map[string]float64{
				rowData.Token: rowData.BuyRate,
			}
		}
		result[rowData.Reserve][rowData.Time][rowData.Base][rowData.Token] = rowData.BuyRate
	}
	return result, nil
}

//UpdateETHUSDPrice store the ETHUSD rate at that blockInfo
func (rdb *RatesStorage) UpdateETHUSDPrice(blockInfo lastblockdaily.BlockInfo, ethusdRate float64) error {
	var logger = rdb.sugar.With(
		"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
		"block_number", blockInfo.Block,
		"timestamp", blockInfo.Timestamp.String(),
	)
	const updateStmt = `INSERT INTO %[1]s(time, block, rate)
	VALUES ( 
		$1,
		$2, 
		$3);`
	query := fmt.Sprintf(updateStmt,
		rdb.tableNames[4],
	)

	logger.Debugw("updating eth-usdrates...", "query", query)
	tx, err := rdb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, rdb.sugar, &err)
	_, err = tx.Exec(query,
		blockInfo.Timestamp,
		blockInfo.Block,
		ethusdRate,
	)
	if err != nil {
		return err
	}

	return err
}
