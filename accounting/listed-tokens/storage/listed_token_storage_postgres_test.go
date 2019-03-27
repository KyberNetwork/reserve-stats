package storage

import (
	"math/big"
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const tokenTableTest = "listed_tokens_test"

func newListedTokenDB(sugar *zap.SugaredLogger) (*ListedTokenDB, error) {
	_, db := testutil.MustNewDevelopmentDB()
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

	var (
		blockNumber  = big.NewInt(7442895)
		reserve      = ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F")
		listedTokens = []common.ListedToken{
			{
				Address:   ethereum.HexToAddress("0xdd974D5C2e2928deA5F71b9825b8b646686BD200"),
				Symbol:    "KNC",
				Name:      "Kyber Network Crystal",
				Timestamp: timeutil.TimestampMsToTime(1553241328394).UTC(),
			},
			{
				Address:   ethereum.HexToAddress("0x1a7a8BD9106F2B8D977E08582DC7d24c723ab0DB"),
				Symbol:    "APPC",
				Name:      "AppCoins",
				Timestamp: timeutil.TimestampMsToTime(1509977454000).UTC(),
				Old: []common.OldListedToken{
					{
						Address:   ethereum.HexToAddress("0x27054b13b1B798B345b591a4d22e6562d47eA75a"),
						Timestamp: timeutil.TimestampMsToTime(1507599220000).UTC(),
					},
				},
			},
		}
		blockNumberNew  = big.NewInt(7442899)
		listedTokensNew = []common.ListedToken{
			{
				Address:   ethereum.HexToAddress("0xdd974D5C2e2928deA5F71b9825b8b646686BD200"),
				Symbol:    "KNC",
				Name:      "Kyber Network Crystal",
				Timestamp: timeutil.TimestampMsToTime(1553241328394).UTC(),
			},
			{
				Address:   ethereum.HexToAddress("0x406F1CddcFe308cf815Ce2914e15f96036230884"),
				Symbol:    "APPC",
				Name:      "AppCoins",
				Timestamp: timeutil.TimestampMsToTime(1509977458000).UTC(),
				Old: []common.OldListedToken{
					{
						Address:   ethereum.HexToAddress("0x1a7a8BD9106F2B8D977E08582DC7d24c723ab0DB"),
						Timestamp: timeutil.TimestampMsToTime(1509977454000).UTC(),
					},
					{
						Address:   ethereum.HexToAddress("0x27054b13b1B798B345b591a4d22e6562d47eA75a"),
						Timestamp: timeutil.TimestampMsToTime(1507599220000).UTC(),
					},
				},
			},
		}
	)

	storage, err := newListedTokenDB(logger)
	assert.NoError(t, err)

	defer teardown(t, storage)

	err = storage.CreateOrUpdate(listedTokens, blockNumber, reserve)
	require.NoError(t, err)

	storedListedTokens, version, storedBlockNumber, err := storage.GetTokens()
	require.NoError(t, err)
	assert.ElementsMatch(t, listedTokens, storedListedTokens)
	assert.Equal(t, version, uint64(1))
	assert.Equal(t, blockNumber.Uint64(), storedBlockNumber)

	err = storage.CreateOrUpdate(listedTokensNew, blockNumberNew, reserve)
	require.NoError(t, err)

	storedNewListedTokens, version, storedBlockNumber, err := storage.GetTokens()
	assert.NoError(t, err)
	assert.Equal(t, version, uint64(2))
	assert.Equal(t, blockNumberNew.Uint64(), storedBlockNumber)
	assert.ElementsMatch(t, listedTokensNew, storedNewListedTokens)
}
