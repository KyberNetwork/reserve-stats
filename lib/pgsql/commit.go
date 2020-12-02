package pgsql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
)

// CommitOrRollback is a deferred function that:
// - commit if given error is nil
// - rollback if given error is not nil
// Any error happens when commit/rollback is assigned to the given error instance.
// Caller function should be named function when calling this function with defer
func CommitOrRollback(tx *sqlx.Tx, sugar *zap.SugaredLogger, err *error) {
	var logger = sugar.With("func", caller.GetCurrentFunctionName())

	if *err == nil {
		if cErr := tx.Commit(); cErr != nil {
			*err = fmt.Errorf("failed to commit transaction: %v", cErr)
		}
		return
	}

	logger.Debugw("rolling back transaction")
	if rErr := tx.Rollback(); rErr != nil {
		logger.Errorw("failed to rollback transaction", "error", rErr)
		return
	}
	logger.Infow("transaction roll backed", "error", *err)
}
