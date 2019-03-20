package listedtokenstorage

import (
	"fmt"
	"testing"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	postgresHost     = "127.0.0.1"
	postgresPort     = 5432
	postgresUser     = "reserve_stats"
	postgresPassword = "reserve_stats"
	postgresDatabase = "reserve_stats"
	tokenTableTest   = "tokens_test"
)

func newListedTokenDB(sugar *zap.SugaredLogger) (*ListedTokenDB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDatabase,
	)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	storage, err := NewDB(sugar, db, tokenTableTest)
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func teardown(t *testing.T, storage *ListedTokenDB) {
	err := storage.DeleteTable()
	assert.NoError(t, err)
	err = storage.Close()
	assert.NoError(t, err)
}

func TestListedTokenStorage(t *testing.T) {
	logger := testutil.MustNewDevelopmentSugaredLogger()
	logger.Info("start testing")

	var listedTokens = map[string]common.ListedToken{
		"APPC-AppCoins": {
			Address:   "0x1a7a8BD9106F2B8D977E08582DC7d24c723ab0DB",
			Symbol:    "APPC",
			Name:      "AppCoins",
			Timestamp: 1509977454000,
			Old: []common.OldListedToken{
				{
					Address:   "0x27054b13b1B798B345b591a4d22e6562d47eA75a",
					Timestamp: 1507599220000,
				},
			},
		},
	}

	storage, err := newListedTokenDB(logger)
	assert.NoError(t, err)

	defer teardown(t, storage)

	err = storage.CreateOrUpdate(listedTokens)
	assert.NoError(t, err)

	storedListedTokens, err := storage.GetTokens()
	assert.NoError(t, err)
	assert.Equal(t, listedTokens, storedListedTokens)
}
