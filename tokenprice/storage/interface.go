package storage

import (
	"time"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/tokenprice/storage/postgres"
)

const (
	// DefaultDB default db
	DefaultDB = "reserve_stats"
)

// Storage storage interface
type Storage interface {
	SaveTokenRate(token, currency, source string, timestamp time.Time, rate float64) error
	GetTokenRate(token, currency string, timestamp time.Time) (float64, error)
}

// NewStorageFromContext return storage interface from context
func NewStorageFromContext(sugar *zap.SugaredLogger, c *cli.Context) (Storage, error) {
	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return nil, err
	}
	return postgres.NewTokenRateDB(sugar, db)
}
