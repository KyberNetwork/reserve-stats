package postgres

import (
	"database/sql"
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const updateTokenSymbolTemplate = `INSERT INTO %[1]s(
	address,
	symbol
) VALUES (
	unnest($1::TEXT[]), 
	unnest($2::TEXT[])
) ON CONFLICT ON CONSTRAINT %[1]s_address_key DO UPDATE SET symbol = EXCLUDED.symbol`

func (tldb *TradeLogDB) saveTokens(tx *sqlx.Tx, tokensArray []string) error {
	var logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
	query := fmt.Sprintf(insertionAddressTemplate, schema.TokenTableName)
	logger.Debugw("updating tokens...", "query", query)
	_, err := tx.Exec(query, pq.StringArray(tokensArray))
	return err
}

// GetTokenSymbol return symbol of provided token address
func (tldb *TradeLogDB) GetTokenSymbol(address string) (string, error) {
	var (
		logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		symbol string
	)
	query := fmt.Sprintf("SELECT symbol FROM %1s WHERE address = $1;", schema.TokenTableName)
	logger.Debugw("get token symbol", "token", address, "query", query)
	if err := tldb.db.Get(&symbol, query, ethereum.HexToAddress(address).Hex()); err != nil {
		if err != sql.ErrNoRows {
			return symbol, fmt.Errorf("failed to get token symbol: %s", err.Error())
		}
	}
	return symbol, nil
}

// UpdateTokens update token symbol, insert new record if the token have not yet added to the table
func (tldb *TradeLogDB) UpdateTokens(tokensArray []string, symbolArray []string) error {
	var logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
	query := fmt.Sprintf(updateTokenSymbolTemplate, schema.TokenTableName)
	logger.Debugw("updating token symbols ...", "query", query)
	_, err := tldb.db.Exec(query, pq.StringArray(tokensArray), pq.StringArray(symbolArray))
	return err
}
