package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
)

// NewDB return the Ratestorage instance. User must call ratestorage.Close() before exit.
// tableNames is a list of 5 string for 5 tablename[reserve,token,base, rate,usdrate]. It can be optional
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB, customTableNames ...string) (*RatesStorage, error) {
	const schemaFMT = `--reserves table definition
	CREATE TABLE IF NOT EXISTS %[1]s
(
	id serial NOT NULL,
	address TEXT NOT NULL UNIQUE,
	CONSTRAINT %[1]s_pk PRIMARY KEY(id)
) ;

--tokens table definition
CREATE TABLE IF NOT EXISTS %[2]s
(
id serial NOT NULL,
	symbol TEXT NOT NULL UNIQUE,
	CONSTRAINT %[2]s_pk PRIMARY KEY(id)
);

--bases table definition
CREATE TABLE IF NOT EXISTS %[3]s
(
	id serial NOT NULL,
	symbol TEXT NOT NULL UNIQUE,
	CONSTRAINT %[3]s_pk PRIMARY KEY(id)
);

--rates table definition
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
	CONSTRAINT %[4]s_fk_reseve_id FOREIGN KEY(reserve_id) REFERENCES %[1]s(id),
	CONSTRAINT %[4]s_no_duplicate UNIQUE(token_id,base_id,block,reserve_id)
);
CREATE INDEX IF NOT EXISTS  %[4]s_time_idx ON %[4]s(time);

--usds table definition
CREATE TABLE IF NOT EXISTS %[5]s
(
	id serial NOT NULL,
	time TIMESTAMP NOT NULL UNIQUE,
	block integer NOT NULL UNIQUE,
	rate float8 NOT NULL,
	CONSTRAINT %[5]s_pk PRIMARY KEY(id),
	CONSTRAINT %[5]s_time_block UNIQUE(time,block)
);
CREATE INDEX IF NOT EXISTS %[5]s_time_idx ON %[5]s(time);

CREATE OR REPLACE VIEW rates_view AS 
	SELECT rt.time as time, tk.symbol as token, bs.symbol as base, rt.buy_reserve_rate as buy_rate, rs.address as reserve
		FROM %[4]s AS rt LEFT JOIN %[2]s AS tk ON rt.token_id = tk.id
		LEFT JOIN %[3]s AS bs ON rt.base_id=bs.id 
		LEFT JOIN %[1]s AS rs ON rt.reserve_id=rs.id;
`
	var (
		logger     = sugar.With("func", "reserverates/storage/postgres")
		tableNames = make(map[string]string)
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
		tableNames[reserveTableName] = customTableNames[0]
		tableNames[tokenTableName] = customTableNames[1]
		tableNames[baseTableName] = customTableNames[2]
		tableNames[rateTableName] = customTableNames[3]
		tableNames[usdTableName] = customTableNames[4]

	} else {
		tableNames[reserveTableName] = reserveTableName
		tableNames[tokenTableName] = tokenTableName
		tableNames[baseTableName] = baseTableName
		tableNames[rateTableName] = rateTableName
		tableNames[usdTableName] = usdTableName
	}
	query := fmt.Sprintf(schemaFMT, tableNames[reserveTableName], tableNames[tokenTableName], tableNames[baseTableName], tableNames[rateTableName], tableNames[usdTableName])
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
