package postgres

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
)

func TestSaveAccountingRates(t *testing.T) {
	var (
		blockInfo = lastblockdaily.BlockInfo{
			Block:     uint64(1000),
			Timestamp: time.Now(),
		}
		rateRecord = common.ReserveRateEntry{
			BuyReserveRate:  0.1,
			SellReserveRate: 0.2,
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

	arrs, err := NewDB(sugar, db, "test_rsv_table", "test_token_table", "test_bases_table", "test_rates_table")
	defer func() {
		// arrs.TearDown()
		arrs.Close()
	}()
	require.NoError(t, err)
	err = arrs.UpdateRatesRecords(blockInfo, rsvRates)
	require.NoError(t, err)
}
