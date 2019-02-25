package influxstorage

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

//KycChecker is the interface to abstract functionality for checking kyc status
type KycChecker interface {
	IsKYCedAtTime(common.Address, time.Time) (bool, error)
}

// NewUserKYCChecker creates a new instance of UserKYCChecker.
func NewUserKYCChecker(sugar *zap.SugaredLogger, db *sqlx.DB) *UserKYCChecker {
	return &UserKYCChecker{sugar: sugar, db: db}
}

// UserKYCChecker is an implementation of kycChecker interface that read the KYC status from users database.
type UserKYCChecker struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

// IsKYCedAtTime returns true if user is already KYCed at the given timestamp.
func (c *UserKYCChecker) IsKYCedAtTime(userAddr common.Address, ts time.Time) (bool, error) {
	const addressesTableName = "addresses"

	var (
		logger = c.sugar.With(
			"func", "tradelogs/UserKYCChecker.IsKYCed",
			"user_addr", userAddr,
			"timestamp", ts.String(),
		)
		result uint64
	)

	stmt := fmt.Sprintf(`SELECT COUNT(1) FROM "%s" WHERE address = $1 AND timestamp <= $2`, addressesTableName)
	logger = logger.With("query", stmt)
	if err := c.db.Get(&result, stmt, userAddr.Hex(), ts.UTC()); err != nil {
		return false, err
	}
	logger.Debugw("got result from database", "result", result)
	return result != 0, nil
}

type mocKYCChecker struct{}

func (*mocKYCChecker) IsKYCedAtTime(_ common.Address, _ time.Time) (bool, error) {
	return true, nil
}

func newMockKYCChecker() *mocKYCChecker {
	return &mocKYCChecker{}
}
