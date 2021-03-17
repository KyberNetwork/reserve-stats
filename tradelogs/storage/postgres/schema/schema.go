package schema

// TradeLogsSchema is postgres schema for tradelog
const TradeLogsSchema = `
CREATE TABLE IF NOT EXISTS "users" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL,
	timestamp TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS "token" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL,
	symbol TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS "reserve" (
	id SERIAL PRIMARY KEY,
	address TEXT NOT NULL,
	reserve_id TEXT DEFAULT '',
	reserve_type INTEGER DEFAULT 0,
	block_number INTEGER DEFAULT 0,
	name TEXT DEFAULT '',
	CONSTRAINT reserve_pk UNIQUE (address, reserve_id, block_number)
);


CREATE TABLE IF NOT EXISTS "` + TradeLogsTableName + `" (
	id SERIAL PRIMARY KEY,
	timestamp TIMESTAMPTZ,
	block_number INTEGER,
	tx_hash TEXT,
	usdt_amount FLOAT(32),
	original_usdt_amount FLOAT(32),
	user_address_id BIGINT NOT NULL REFERENCES users,
	src_address_id BIGINT NOT NULL REFERENCES token,
	dst_address_id BIGINT NOT NULL REFERENCES token,
	src_amount FLOAT(32),
	dst_amount FLOAT(32),
	index INTEGER,
	is_first_trade BOOLEAN,
	tx_sender	TEXT,
	receiver_address	TEXT,
	gas_used INTEGER,
	gas_price FLOAT(32),
	transaction_fee FLOAT(32),
	version integer,
	CONSTRAINT tradelog_constraint UNIQUE (tx_hash, index)
);

CREATE UNIQUE INDEX IF NOT EXISTS "tradelogs_id_index" ON "` + TradeLogsTableName + `"(id);

CREATE INDEX IF NOT EXISTS "trade_timestamp" ON "` + TradeLogsTableName + `"(timestamp);
CREATE INDEX IF NOT EXISTS "trade_user_address" ON "` + TradeLogsTableName + `"(user_address_id);
CREATE INDEX IF NOT EXISTS "trade_src_address" ON "` + TradeLogsTableName + `"(src_address_id);
CREATE INDEX IF NOT EXISTS "trade_dst_address" ON "` + TradeLogsTableName + `"(dst_address_id);
CREATE INDEX IF NOT EXISTS "trade_tx_hash" ON "` + TradeLogsTableName + `"(tx_hash);

-- create_or_update_tradelogs creates or update tradelogs
CREATE OR REPLACE FUNCTION create_or_update_tradelogs(INOUT _id tradelogs.id%TYPE,
												_timestamp tradelogs.timestamp%TYPE,
												_block_number tradelogs.block_number%TYPE,
												_tx_hash tradelogs.tx_hash%TYPE,
												_usdt_amount tradelogs.usdt_amount%TYPE,
												_original_usdt_amount tradelogs.original_usdt_amount%TYPE,
												_user_address TEXT,
												_src_address TEXT,
												_dst_address TEXT,
												_src_amount tradelogs.src_amount%TYPE,
												_dst_amount tradelogs.dst_amount%TYPE,
												_index tradelogs.index%TYPE,
												_is_first_trade tradelogs.is_first_trade%TYPE,
												_tx_sender tradelogs.tx_sender%TYPE,
												_receiver_address tradelogs.receiver_address%TYPE,
												_gas_used tradelogs.gas_used%TYPE,
												_gas_price tradelogs.gas_price%TYPE,
												_transaction_fee tradelogs.transaction_fee%TYPE,
												_version tradelogs.version%TYPE
												) AS
$$
BEGIN
    IF _id = 0 THEN
		INSERT INTO tradelogs (timestamp, block_number, tx_hash, usdt_amount, 
			original_usdt_amount, user_address_id, src_address_id, dst_address_id, src_amount, dst_amount,
			index, is_first_trade, tx_sender,
			receiver_address, gas_used, gas_price, transaction_fee, version) 
		VALUES (_timestamp,
			_block_number,
			_tx_hash,
			_usdt_amount,
			_original_usdt_amount,
			(SELECT id FROM users WHERE address=_user_address),
			(SELECT id FROM token WHERE address=_src_address),
			(SELECT id FROM token WHERE address=_dst_address),
			_src_amount,
			_dst_amount,
			_index, 
			_is_first_trade,
			_tx_sender,
			_receiver_address,
			_gas_used,
			_gas_price,
			_transaction_fee,
			_version
		) ON CONFLICT (tx_hash, index) DO UPDATE SET 
			timestamp = _timestamp
		 RETURNING id INTO _id;
    END IF;

    RETURN;
END;
$$ LANGUAGE PLPGSQL;
`

// DefaultDateFormat ...
const DefaultDateFormat = "2006-01-02 15:04:05"
