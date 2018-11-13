package storage

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/jmoiron/sqlx"
)

const (
	addressesTableName = "addresses"
)

// UserPostgresStorage is storage saving kyced info
type UserPostgresStorage struct {
	db *sqlx.DB
}

//NewPostgresConnection return a new instance for UserPostgresStorage
func NewPostgresConnection(db *sqlx.DB) *UserPostgresStorage {
	return &UserPostgresStorage{
		db: db,
	}
}

// CountKYCEDAddresses is function return number of kyced user
func (ups UserPostgresStorage) CountKYCEDAddresses(ts uint64) (uint64, error) {
	var (
		result uint64
		err    error
	)
	fromTime := timeutil.TimestampMsToTime(ts)
	// one day time
	toTime := timeutil.TimestampMsToTime(ts + 86400000)
	if err = ups.db.Get(&result, fmt.Sprintf(`SELECT COUNT(1) FROM "%s" WHERE timestamp >= $1 AND timestamp < $2`, addressesTableName), fromTime.UTC(), toTime.UTC()); err != nil {
		return result, err
	}
	return result, err
}
