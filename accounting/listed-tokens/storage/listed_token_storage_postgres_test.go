package storage

import (
	"math/big"
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestListedTokenStorage(t *testing.T) {
	logger := testutil.MustNewDevelopmentSugaredLogger()
	logger.Info("start testing")

	var (
		blockNumber       = big.NewInt(7442895)
		reserve           = ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F")
		secondReserve     = ethereum.HexToAddress("0x21433Dec9Cb634A23c6A4BbcCe08c83f5aC2EC18")
		notExistedReserve = ethereum.HexToAddress("0x21433Dec9Cb634A23c6A4BbcCe08c83f5aC2EC18")
		zeroReserve       = ethereum.HexToAddress("0x0000000000000000000000000000000000000000")
		listedTokens      = []common.ListedToken{
			{
				Address:   ethereum.HexToAddress("0xdd974D5C2e2928deA5F71b9825b8b646686BD200"),
				Symbol:    "KNC",
				Name:      "Kyber Network Crystal",
				Timestamp: timeutil.TimestampMsToTime(1553241328394).UTC(),
				Decimals:  18,
			},
			{
				Address:   ethereum.HexToAddress("0x1a7a8BD9106F2B8D977E08582DC7d24c723ab0DB"),
				Symbol:    "APPC",
				Name:      "AppCoins",
				Timestamp: timeutil.TimestampMsToTime(1509977454000).UTC(),
				Decimals:  18,
				Old: []common.OldListedToken{
					{
						Address:   ethereum.HexToAddress("0x27054b13b1B798B345b591a4d22e6562d47eA75a"),
						Timestamp: timeutil.TimestampMsToTime(1507599220000).UTC(),
						Decimals:  10,
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
				Decimals:  18,
			},
			{
				Address:   ethereum.HexToAddress("0x406F1CddcFe308cf815Ce2914e15f96036230884"),
				Symbol:    "APPC",
				Name:      "AppCoins",
				Timestamp: timeutil.TimestampMsToTime(1509977458000).UTC(),
				Decimals:  18,
				Old: []common.OldListedToken{
					{
						Address:   ethereum.HexToAddress("0x1a7a8BD9106F2B8D977E08582DC7d24c723ab0DB"),
						Timestamp: timeutil.TimestampMsToTime(1509977454000).UTC(),
						Decimals:  10,
					},
					{
						Address:   ethereum.HexToAddress("0x27054b13b1B798B345b591a4d22e6562d47eA75a"),
						Timestamp: timeutil.TimestampMsToTime(1507599220000).UTC(),
						Decimals:  10,
					},
				},
			},
		}
	)

	db, teardown := testutil.MustNewRandomDevelopmentDB()
	storage, err := NewDB(logger, db)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, teardown())
	}()

	err = storage.CreateOrUpdate(listedTokens, blockNumber, reserve)
	require.NoError(t, err)

	storedListedTokens, version, storedBlockNumber, err := storage.GetTokens(reserve)
	require.NoError(t, err)
	assert.ElementsMatch(t, listedTokens, storedListedTokens)
	assert.Equal(t, uint64(1), version)
	assert.Equal(t, blockNumber.Uint64(), storedBlockNumber)

	err = storage.CreateOrUpdate(listedTokensNew, blockNumberNew, reserve)
	require.NoError(t, err)
	storedNewListedTokens, version, storedBlockNumber, err := storage.GetTokens(reserve)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), version)
	assert.Equal(t, blockNumberNew.Uint64(), storedBlockNumber)
	assert.ElementsMatch(t, listedTokensNew, storedNewListedTokens)

	// assert provided reserve is zero
	zeroReserveTokens, version, storedBlockNumber, err := storage.GetTokens(zeroReserve)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), version)
	assert.Equal(t, blockNumberNew.Uint64(), storedBlockNumber)
	assert.ElementsMatch(t, listedTokensNew, zeroReserveTokens)

	noTokens, _, _, err := storage.GetTokens(notExistedReserve)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(noTokens))

	//
	err = storage.CreateOrUpdate(listedTokensNew, blockNumber, secondReserve)
	require.NoError(t, err)

	testDuplicateSavedTokens, version, storedBlockNumber, err := storage.GetTokens(zeroReserve)
	require.NoError(t, err)
	assert.ElementsMatch(t, listedTokensNew, testDuplicateSavedTokens)
	assert.Equal(t, uint64(3), version)
	assert.Equal(t, blockNumber.Uint64(), storedBlockNumber)
}
