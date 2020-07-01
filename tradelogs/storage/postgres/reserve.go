package postgres

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// AddUpdateReserve update reserve address with reserve id
// func (tldb *TradeLogDB) AddUpdateReserve(reserve common.Reserve) error {
// 	var (
// 		logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
// 	)
// 	logger.Info("update reserve address")
// 	query := "INSERT INTO reserve (address, rebate_wallet, reserve_id, block_number) VALUES ($1, $2, $3, $4);" // TODO: handle ON CONFLICT
// 	if _, err := tldb.db.Exec(query,
// 		reserve.Reserve.Hex(),
// 		reserve.RebateWallet.Hex(),
// 		common.BytesToHash(reserve.ReserveID[:]).Hex(),
// 		reserve.BlockNumber,
// 	); err != nil {
// 		return err
// 	}
// 	return nil
// }

// AddUpdateRebateWallet update rebase wallet
// func (tldb *TradeLogDB) AddUpdateRebateWallet(rebateWallet *tradelogs.UpdateRebateWallet) error {
// 	var (
// 		logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
// 	)
// 	logger.Info("update wallet address")
// 	query := `
// 		WITH r as (SELECT address as reserve_address FROM reserve WHERE reserve_id = $1)
// 		INSERT INTO reserve (address, rebate_wallet, reserve_id, block_number) VALUES (r.reserve_address, $2, $3, $4);
// 	`
// 	logger.Info(query)
// 	return nil
// }

// use for version before v4
func (tldb *TradeLogDB) saveReserveAddress(tx *sqlx.Tx, reserveAddressArray []string) error {
	var (
		logger = tldb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"reserves", reserveAddressArray,
		)
	)
	// query := fmt.Sprintf(insertionAddressTemplate, schema.ReserveTableName)
	query := fmt.Sprintf(`INSERT INTO %[1]s(address) 
	VALUES (UNNEST($1::TEXT[])) ON CONFLICT ON CONSTRAINT reserve_pk DO NOTHING;`, schema.ReserveTableName) // TODO: define conflict here
	logger.Debugw("updating rsv...", "query", query)

	_, err := tx.Exec(query, pq.StringArray(reserveAddressArray))
	if err != nil {
		logger.Errorw("failed to update reserve", "error", err)
	}
	return err
}
