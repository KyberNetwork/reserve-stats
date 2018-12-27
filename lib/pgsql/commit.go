package pgsql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// CommitOrRollback is a deferred function that:
// - commit if given error is nil
// - rollback if given error is not nil
// Any error happens when commit/rollback is assigned to the given error instance.
func CommitOrRollback(tx *sqlx.Tx, sugar *zap.SugaredLogger, err *error) {
	var logger = sugar.With("func", "lib/pgsql/CommitOrRollback")

	if *err == nil {
		logger.Debugw("committing transaction")
		if cErr := tx.Commit(); cErr != nil {
			*err = fmt.Errorf("failed to commit transaction: %v", cErr)
		}
		return
	}

	logger.Debugw("rolling back transaction")
	if rErr := tx.Rollback(); rErr != nil {
		logger.Errorw("failed to rollback transaction", "error", rErr)
	}
	logger.Infow("transaction roll backed", "error", *err)
}
