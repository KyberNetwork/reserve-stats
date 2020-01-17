package postgres

import (
	"log"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveBigTrades(t *testing.T) {
	// save tradelogs
	const (
		dbName = "test_save_trade_log"
	)
	testStorage, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testStorage.tearDown(dbName))
	}()

	tradeLogs, err := utils.GetSampleTradeLogs("../testdata/trade_logs.json")
	require.NoError(t, err)
	require.NoError(t, testStorage.SaveTradeLogs(tradeLogs))
	log.Printf("len trade logs: %d", len(tradeLogs))

	// save big trades
	require.NoError(t, testStorage.SaveBigTrades(float32(100), 6100010))

	// get big trades
	var fromTime time.Time
	toTime := time.Now()
	bigTrades, err := testStorage.GetNotTwittedTrades(fromTime, toTime)
	require.NoError(t, err)
	// expect len(bigTrades) > 0
	assert.Greater(t, len(bigTrades), 0)

	bigTradeIDs := []uint64{}
	for _, trade := range bigTrades {
		bigTradeIDs = append(bigTradeIDs, trade.TradelogID)
	}

	require.NoError(t, testStorage.UpdateBigTradesTwitted(bigTradeIDs))
}
