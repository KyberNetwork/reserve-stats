package schema

// TradeLogsSchema is postgres schema for tradelog
const TradeLogsSchema = `
CREATE TABLE IF NOT EXISTS "users" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL,
	timestamp TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS "wallet" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL,
	name TEXT
);
CREATE TABLE IF NOT EXISTS "token" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL
);

DO $$ 
    BEGIN
        BEGIN
            ALTER TABLE "token" ADD COLUMN symbol TEXT DEFAULT '';
        EXCEPTION
            WHEN duplicate_column THEN RAISE NOTICE 'column symbol already exists in token.';
        END;
    END;
$$;

CREATE TABLE IF NOT EXISTS "reserve" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL
);

ALTER TABLE "reserve" 
	ADD COLUMN IF NOT EXISTS "reserve_id" TEXT DEFAULT '',
	ADD COLUMN IF NOT EXISTS "rebate_wallet" TEXT DEFAULT '', 
	ADD COLUMN IF NOT EXISTS "block_number" INTEGER;


DO $$ 
    BEGIN
        BEGIN
            ALTER TABLE "reserve" ADD COLUMN name TEXT DEFAULT '';
        EXCEPTION
            WHEN duplicate_column THEN RAISE NOTICE 'column name already exists in reserve.';
        END;
    END;
$$;

CREATE TABLE IF NOT EXISTS "` + TradeLogsTableName + `" (
	id SERIAL,
	timestamp TIMESTAMPTZ,
	block_number INTEGER,
	tx_hash TEXT,
	eth_amount FLOAT(32),
	original_eth_amount FLOAT(32),
	user_address_id BIGINT NOT NULL REFERENCES users,
	src_address_id BIGINT NOT NULL REFERENCES token,
	dst_address_id BIGINT NOT NULL REFERENCES token,
	src_reserve_address_id BIGINT NOT NULL REFERENCES reserve,
	dst_reserve_address_id BIGINT NOT NULL REFERENCES reserve,
	src_amount FLOAT(32),
	dst_amount FLOAT(32),
	wallet_address_id BIGINT NOT NULL REFERENCES wallet,
	src_burn_amount FLOAT(32),
	dst_burn_amount FLOAT(32),
	src_wallet_fee_amount FLOAT(32),
	dst_wallet_fee_amount FLOAT(32),
	integration_app TEXT,
	ip TEXT,
	country TEXT,
	eth_usd_rate FLOAT(32),
	eth_usd_provider TEXT,
	index INTEGER,
	kyced BOOLEAN,
	is_first_trade BOOLEAN,
	tx_sender	TEXT,
	receiver_address	TEXT,
	gas_used INTEGER,
	gas_price FLOAT(32),
	transaction_fee FLOAT(32),
	PRIMARY KEY (tx_hash,index)
);

CREATE UNIQUE INDEX IF NOT EXISTS "tradelogs_id_index" ON "` + TradeLogsTableName + `"(id);

ALTER TABLE "` + TradeLogsTableName + `"
	ADD COLUMN IF NOT EXISTS gas_used INTEGER,
	ADD COLUMN IF NOT EXISTS transaction_fee FLOAT(32),
	ADD COLUMN IF NOT EXISTS gas_price FLOAT(32);

CREATE TABLE IF NOT EXISTS "` + BigTradeLogsTableName + `" (
	id SERIAL PRIMARY KEY,
	tradelog_id INTEGER UNIQUE NOT NULL REFERENCES tradelogs (id),
	twitted BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS "trade_timestamp" ON "` + TradeLogsTableName + `"(timestamp);
CREATE INDEX IF NOT EXISTS "trade_user_address" ON "` + TradeLogsTableName + `"(user_address_id);
CREATE INDEX IF NOT EXISTS "trade_src_address" ON "` + TradeLogsTableName + `"(src_address_id);
CREATE INDEX IF NOT EXISTS "trade_dst_address" ON "` + TradeLogsTableName + `"(dst_address_id);
CREATE INDEX IF NOT EXISTS "trade_src_reserve_address" ON "` + TradeLogsTableName + `"(src_reserve_address_id);
CREATE INDEX IF NOT EXISTS "trade_dst_reserve_address" ON "` + TradeLogsTableName + `"(dst_reserve_address_id);
CREATE INDEX IF NOT EXISTS "trade_wallet_address" ON "` + TradeLogsTableName + `"(wallet_address_id);
CREATE INDEX IF NOT EXISTS "trade_tx_hash" ON "` + TradeLogsTableName + `"(tx_hash);


CREATE TABLE IF NOT EXISTS  "tradelog_v4" (
	id SERIAL,
	timestamp TIMESTAMPTZ,
	block_number INTEGER,
	tx_hash TEXT,
	eth_amount FLOAT(32),
	original_eth_amount FLOAT(32),
	user_address_id BIGINT NOT NULL REFERENCES users,
	src_address_id BIGINT NOT NULL REFERENCES token,
	dst_address_id BIGINT NOT NULL REFERENCES token,
	src_amount FLOAT(32),
	dst_amount FLOAT(32),
	platform_wallet_address_id BIGINT NOT NULL REFERENCES wallet,
	t2e_reserves_id BIGINT[] NOT NULL REFERENCES reserve,
	e2t_reserves_id BIGINT[] NOT NULL REFERENCES reserve, 
	integration_app TEXT,
	ip TEXT,
	country TEXT,
	eth_usd_rate FLOAT(32),
	eth_usd_provider TEXT,
	index INTEGER,
	kyced BOOLEAN,
	is_first_trade BOOLEAN,
	tx_sender	TEXT,
	receiver_address	TEXT,
	gas_used INTEGER,
	gas_price FLOAT(32),
	transaction_fee FLOAT(32),
	PRIMARY KEY (tx_hash,index)
);
`

// DefaultDateFormat ...
const DefaultDateFormat = "2006-01-02 15:04:05"
