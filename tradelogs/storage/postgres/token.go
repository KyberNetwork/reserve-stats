package postgres

import (
	"database/sql"
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const updateTokenSymbolTemplate = `INSERT INTO token(
	address,
	symbol
) VALUES (
	unnest($1::TEXT[]), 
	unnest($2::TEXT[])
) ON CONFLICT ON CONSTRAINT token_address_key DO UPDATE SET symbol = EXCLUDED.symbol`

func (tldb *TradeLogDB) saveTokens(tx *sqlx.Tx, tokensArray []string, decimals []int64) error {
	var logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
	logger.Debugw("updating tokens...", "query", insertionAddressTemplate)
	_, err := tx.Exec(insertionAddressTemplate, pq.StringArray(tokensArray), pq.Array(decimals))
	return err
}

// GetTokenSymbol return symbol of provided token address
func (tldb *TradeLogDB) GetTokenSymbol(address string) (string, error) {
	var (
		logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		symbol string
	)
	query := "SELECT symbol FROM token WHERE address = $1;"
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
	logger.Debugw("updating token symbols ...", "query", updateTokenSymbolTemplate)
	_, err := tldb.db.Exec(updateTokenSymbolTemplate, pq.StringArray(tokensArray), pq.StringArray(symbolArray))
	return err
}
