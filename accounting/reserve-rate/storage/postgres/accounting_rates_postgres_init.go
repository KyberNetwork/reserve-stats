package postgres

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
)

// NewDB return the Ratestorage instance. User must call ratestorage.Close() before exit.
// tableNames is a list of 5 string for 5 tablename[reserve,token,quote, rate,usdrate]. It can be optional
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*RatesStorage, error) {
	const schemaFMT = `--reserves table definition
	CREATE TABLE IF NOT EXISTS reserves
(
	id serial NOT NULL,
	address TEXT NOT NULL UNIQUE,
	CONSTRAINT reserves_pk PRIMARY KEY(id)
) ;

--tokens table definition
CREATE TABLE IF NOT EXISTS tokens
(
id serial NOT NULL,
	symbol TEXT NOT NULL UNIQUE,
	CONSTRAINT tokens_pk PRIMARY KEY(id)
);

--quotes table definition
CREATE TABLE IF NOT EXISTS quotes
(
	id serial NOT NULL,
	symbol TEXT NOT NULL UNIQUE,
	CONSTRAINT quotes_pk PRIMARY KEY(id)
);

--rates table definition
CREATE TABLE IF NOT EXISTS token_rates
(
	id serial NOT NULL,
	time TIMESTAMP NOT NULL,
	token_id serial NOT NULL,
    quote_id serial NOT NULL,
	block integer NOT NULL,
	rate float8 NOT NULL,
	reserve_id integer NOT NULL,
	CONSTRAINT token_rates_pk PRIMARY KEY(id),
	CONSTRAINT token_rates_fk_token_id FOREIGN KEY(token_id) REFERENCES tokens(id),
    CONSTRAINT token_rates_fk_quote_id FOREIGN KEY(quote_id) REFERENCES quotes(id),
	CONSTRAINT token_rates_fk_reseve_id FOREIGN KEY(reserve_id) REFERENCES reserves(id),
	CONSTRAINT token_rates_no_duplicate UNIQUE(token_id,quote_id,block,reserve_id)
);
CREATE INDEX IF NOT EXISTS  token_rates_time_idx ON token_rates(time);

--usds table definition
CREATE TABLE IF NOT EXISTS usd_rates
(
	id serial NOT NULL,
	time TIMESTAMP NOT NULL UNIQUE,
	block integer NOT NULL UNIQUE,
	rate float8 NOT NULL,
	CONSTRAINT usd_rates_pk PRIMARY KEY(id),
	CONSTRAINT usd_rates_time_block UNIQUE(time,block)
);
CREATE INDEX IF NOT EXISTS usd_rates_time_idx ON usd_rates(time);

CREATE OR REPLACE VIEW rates_view AS 
	SELECT rt.time as time, tk.symbol as token, bs.symbol as quote, rt.rate as rate, rs.address as reserve
		FROM token_rates AS rt LEFT JOIN tokens AS tk ON rt.token_id = tk.id
		LEFT JOIN quotes AS bs ON rt.quote_id=bs.id 
		LEFT JOIN reserves AS rs ON rt.reserve_id=rs.id;
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
	logger.Debugw("initializing database schema", "query", schemaFMT)

	if _, err = tx.Exec(schemaFMT); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")
	return &RatesStorage{
		sugar:      sugar,
		db:         db,
		tableNames: tableNames,
	}, nil
}
