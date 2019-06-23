package schema

// TradeLogsSchema is postgres schema for tradelog
const TradeLogsSchema = `
CREATE TABLE IF NOT EXISTS "users" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL,
	timestamp TIMESTAMP
);
CREATE TABLE IF NOT EXISTS "wallet" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS "token" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS "reserve" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS "` + TradeLogsTableName + `" (
	id SERIAL,
	timestamp TIMESTAMP,
	block_number INTEGER,
	tx_hash TEXT,
	eth_amount FLOAT(32),
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
	PRIMARY KEY (tx_hash,index)
);

CREATE INDEX IF NOT EXISTS "trade_timestamp" ON "` + TradeLogsTableName + `"(timestamp);
CREATE INDEX IF NOT EXISTS "trade_user_address" ON "` + TradeLogsTableName + `"(user_address_id);
CREATE INDEX IF NOT EXISTS "trade_src_address" ON "` + TradeLogsTableName + `"(src_address_id);
CREATE INDEX IF NOT EXISTS "trade_dst_address" ON "` + TradeLogsTableName + `"(dst_address_id);
CREATE INDEX IF NOT EXISTS "trade_src_reserve_address" ON "` + TradeLogsTableName + `"(src_reserve_address_id);
CREATE INDEX IF NOT EXISTS "trade_dst_reserve_address" ON "` + TradeLogsTableName + `"(dst_reserve_address_id);
`

const DefaultDateFormat = "2006-01-02 15:04:05"
