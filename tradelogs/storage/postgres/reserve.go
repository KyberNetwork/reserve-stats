package postgres

import (
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// use for version before v4
func (tldb *TradeLogDB) saveReserveAddress(tx *sqlx.Tx, reserveAddressArray []string) error {
	var (
		logger = tldb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"reserves", reserveAddressArray,
		)
	)
	query := `INSERT INTO reserve(address) 
	VALUES (UNNEST($1::TEXT[])) 
	ON CONFLICT ON CONSTRAINT reserve_pk DO NOTHING;`
	logger.Debugw("updating rsv...", "query", query)

	_, err := tx.Exec(query, pq.StringArray(reserveAddressArray))
	if err != nil {
		logger.Errorw("failed to update reserve", "error", err)
	}
	return err
}
