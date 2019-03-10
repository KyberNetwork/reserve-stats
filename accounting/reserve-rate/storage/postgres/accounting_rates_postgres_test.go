package postgres

import (
	"encoding/json"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
)

func TestSaveAndGetAccountingRates(t *testing.T) {
	var (
		blockInfo = lastblockdaily.BlockInfo{
			Block:     uint64(3000),
			Timestamp: time.Now(),
		}
		rateRecord = common.ReserveRateEntry{
			BuyReserveRate:  0.4,
			SellReserveRate: 0.5,
		}
		aRsvRate = map[string]common.ReserveRateEntry{
			"ETH-KNC": rateRecord,
		}
		rsvRates = map[string]map[string]common.ReserveRateEntry{
			"0x63825c174ab367968EC60f061753D3bbD36A0D8F": aRsvRate,
		}
	)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	sugar := logger.Sugar()

	db, err := sqlx.Connect("postgres", "host=127.0.0.1 port=5432 user=reserve_stats password=reserve_stats dbname=reserve_stats sslmode=disable")
	require.NoError(t, err)

	arrs, err := NewDB(sugar, db, "test_rsv_table", "test_token_table", "test_bases_table", "test_rates_table", "test_usd_table")
	defer func() {
		arrs.TearDown()
		arrs.Close()
	}()
	require.NoError(t, err)
	err = arrs.UpdateRatesRecords(blockInfo, rsvRates)
	require.NoError(t, err)
	err = arrs.UpdateETHUSDPrice(blockInfo, 0.1)
	require.NoError(t, err)

	result, err := arrs.GetRates(
		[]ethereum.Address{ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F")},
		time.Now().AddDate(0, 0, -2),
		time.Now().AddDate(0, 0, 2))
	require.NoError(t, err)
	sugar.Debugw("query completed", "result", result)
	jsonData, err := json.Marshal(result)
	require.NoError(t, err)
	sugar.Debugf("Result json is %s", jsonData)
	lastBlock, err := arrs.GetLastResolvedBlockInfo()
	require.NoError(t, err)
	assert.Equal(t, lastBlock.Block, blockInfo.Block)
}
