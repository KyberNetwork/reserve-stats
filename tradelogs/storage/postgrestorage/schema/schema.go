package schema

// TradeLogsSchema is postgres schema for tradelog
const TradeLogsSchema = `
CREATE TABLE IF NOT EXISTS "users" (
	id SERIAL PRIMARY KEY,
	user_address TEXT UNIQUE NOT NULL,
	timestamp TIMESTAMP
);
CREATE TABLE IF NOT EXISTS "tradelogs" (
	id SERIAL PRIMARY KEY,
	timestamp TIMESTAMP,
	block_number INTEGER,
	tx_hash TEXT,
	eth_amount FLOAT(32),
	user_address_id BIGINT NOT NULL REFERENCES users,
	src_address TEXT,
	dest_address TEXT,
	src_reserveaddress TEXT,
	dst_reserveaddress TEXT,
	src_amount FLOAT(32),
	dest_amount FLOAT(32),
	fiat_amount FLOAT(32),
	wallet_address TEXT,
	src_burn_amount FLOAT(32),
	dst_burn_amount FLOAT(32),
	src_wallet_fee_amount FLOAT(32),
	dst_wallet_fee_amount FLOAT(32),
	integration_app TEXT,
	ip TEXT,
	country TEXT,
	ethusd_rate FLOAT(32),
	ethusd_provider TEXT,
	index INTEGER
);
`
