package postgres

import (
	"database/sql"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/storage"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestRatesStorage(t *testing.T) {
	// assume that a test never takes more than this amount of time
	const truncateDuration = 10 * time.Second
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	_, db := testutil.MustNewDevelopmentDB()

	rs, err := NewDB(sugar, db, "test_rsv_table", "test_token_table", "test_quote_table", "test_rates_table", "test_usd_table")
	require.NoError(t, err)

	defer func(t *testing.T) {
		require.NoError(t, rs.TearDown())
		require.NoError(t, rs.Close())
	}(t)

	_, err = rs.GetLastResolvedBlockInfo(ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"))
	assert.Equal(t, sql.ErrNoRows, err)

	var tests = []struct {
		block           lastblockdaily.BlockInfo
		ethRates        map[string]map[string]float64
		usdRate         float64
		expectedUSDRate storage.AccountingReserveRates
		expectedETHRate map[string]storage.AccountingReserveRates
	}{
		{
			block: lastblockdaily.BlockInfo{
				Block:     uint64(3000),
				Timestamp: time.Now().Truncate(truncateDuration).UTC(),
			},
			ethRates: map[string]map[string]float64{
				"0x63825c174ab367968EC60f061753D3bbD36A0D8F": {
					"ETH-KNC": 0.4,
				},
				"0x818E6FECD516Ecc3849DAf6845e3EC868087B755": {
					"ETH-KNC": 0.2,
					"ETH-OMG": 0.3,
				},
			},
			usdRate: 0.1,
			expectedETHRate: map[string]storage.AccountingReserveRates{
				"0x63825c174ab367968EC60f061753D3bbD36A0D8F": {
					time.Now().Truncate(truncateDuration).UTC(): {
						"ETH": {
							"KNC": 0.4,
						},
					},
				},
				"0x818E6FECD516Ecc3849DAf6845e3EC868087B755": {
					time.Now().Truncate(truncateDuration).UTC(): {
						"ETH": {
							"KNC": 0.2,
							"OMG": 0.3,
						},
					},
				},
			},
			expectedUSDRate: storage.AccountingReserveRates{
				time.Now().Truncate(truncateDuration).UTC(): {
					"USD": {
						"ETH": 0.1,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		require.NoError(t, rs.UpdateRatesRecords(tc.block, tc.ethRates, tc.usdRate))

		ethRates, err := rs.GetRates(
			time.Now().AddDate(0, 0, -2),
			time.Now().AddDate(0, 0, 2),
		)
		require.NoError(t, err)
		assert.Equal(t, tc.expectedETHRate, ethRates)

		usdRate, err := rs.GetETHUSDRates(
			time.Now().AddDate(0, 0, -2),
			time.Now().AddDate(0, 0, 2),
		)
		require.NoError(t, err)
		assert.Equal(t, tc.expectedUSDRate, usdRate)

		lastBlock, err := rs.GetLastResolvedBlockInfo(ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"))
		assert.Equal(t, lastBlock, tc.block)
	}
}
