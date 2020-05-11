package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tokenratefetcher/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// Storage is object represent postgres storage
type Storage struct {
	db    *sqlx.DB
	sugar *zap.SugaredLogger
}

const schema = `CREATE TABLE IF NOT EXISTS token_rates (
		id SERIAL,
		symbol TEXT NOT NULL,
		provider TEXT NOT NULL,
		timestamp TIMESTAMP NOT NULL,
		rate FLOAT NOT NULL,
		PRIMARY KEY (symbol, provider, timestamp)
	);`

// NewPostgresStorage return new instance for postgres storage of token rate
func NewPostgresStorage(sugar *zap.SugaredLogger, db *sqlx.DB) (*Storage, error) {
	if _, err := db.Exec(schema); err != nil {
		sugar.Errorw("failed to create table token_rates", "error", err)
		return nil, err
	}
	return &Storage{
		sugar: sugar,
		db:    db,
	}, nil
}

// LastTimePoint return last timepoint saved in database for a token and currency
func (s *Storage) LastTimePoint(providerName, tokenID, currencyID string) (time.Time, error) {
	var (
		result time.Time
	)
	query := `SELECT timestamp FROM token_rates WHERE symbol = $1 ORDER BY timestamp DESC LIMIT 1;`
	symbol := fmt.Sprintf("%s_%s", common.GetTokenSymbolFromProviderNameTokenID(providerName, tokenID), currencyID)
	if err := s.db.Get(&result, query, symbol); err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}
		return result, err
	}
	return result, nil
}

// SaveRates save rates to database
func (s *Storage) SaveRates(rates []common.TokenRate) error {
	var (
		logger = s.sugar.With(
			"func", caller.GetCurrentFunctionName(),
		)
		symbols, providers []string
		timestamps         []time.Time
		savedRates         []float64
	)
	query := `INSERT INTO token_rates (symbol, provider, timestamp, rate)
	VALUES (
		UNNEST($1::TEXT[]),
		UNNEST($2::TEXT[]),
		UNNEST($3::TIMESTAMP[]),
		UNNEST($4::FLOAT[])
	) ON CONFLICT DO NOTHING;`
	logger.Infow("saving rates", "query", query)
	for _, rate := range rates {
		symbols = append(symbols, fmt.Sprintf("%s_%s", common.GetTokenSymbolFromProviderNameTokenID(rate.Provider, rate.TokenID), rate.Currency))
		providers = append(providers, rate.Provider)
		timestamps = append(timestamps, rate.Timestamp)
		savedRates = append(savedRates, rate.Rate)
	}
	_, err := s.db.Exec(query, pq.Array(symbols), pq.Array(providers), pq.Array(timestamps), pq.Array(savedRates))
	return err
}
