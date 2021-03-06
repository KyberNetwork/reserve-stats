package storage

import "github.com/KyberNetwork/reserve-stats/burnedfees/common"

const (
	// PostgresDefaultDB default db name when choosing Postgres
	PostgresDefaultDB = "burned_fees"
	// PostgresDBEngine is value for flags dbEngine
	PostgresDBEngine = "postgres"
)

// Interface is the database interaction of burned-fees-crawler service.
type Interface interface {
	Store([]common.BurnAssignedFeesEvent) error
	LastBlock() (int64, error)
}
