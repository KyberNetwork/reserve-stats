CREATE TABLE IF NOT EXISTS tradelogs(
	tx TEXT PRIMARY KEY,
	timestamp BIGINT NOT NULL,
	input_token_symbol TEXT NOT NULL,
	input_token_address TEXT NOT NULL,
	input_token_amount FLOAT NOT NULL,
	output_token_symbol TEXT NOT NULL,
	output_token_address TEXT NOT NULL,
	output_token_amount FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS convert_trades(
	original_symbol TEXT NOT NULL,
	symbol TEXT NOT NULL,
	price FLOAT NOT NULL,
	timestamp BIGINT NOT NULL,
	in_token TEXT NOT NULL,
	in_token_amount FLOAT NOT NULL,
	out_token TEXT NOT NULL,
	out_token_amount FLOAT NOT NULL,
	original_trade JSONB,
	convert_trade JSONB,
	UNIQUE(symbol, timestamp)
)
