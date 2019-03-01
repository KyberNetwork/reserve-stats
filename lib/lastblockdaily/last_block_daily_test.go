package lastblockdaily

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

//func TestGetNextDayBlock(t *testing.T) {
//	var (
//		lastResolve   uint64 = 6255278
//		expectedBlock uint64 = 6261438
//		start                = timeutil.TimestampMsToTime(uint64(1535806920000))
//		end                  = timeutil.TimestampMsToTime(uint64(1536386520000))
//	)
//	logger, err := zap.NewDevelopment()
//	require.NoError(t, err)
//	sugar := logger.Sugar()
//
//	ethClient, err := ethclient.Dial("https://mainnet.infura.io/")
//	require.NoError(t, err)
//
//	blkTimeRsv, err := blockchain.NewBlockTimeResolver(sugar, ethClient)
//	require.NoError(t, err)
//
//	lbResolver := NewLastBlockResolver(ethClient, blkTimeRsv, start, end, sugar)
//
//	lastResolvedTimeStamp, err := lbResolver.resolver.Resolve(uint64(lastResolve))
//	require.NoError(t, err)
//	lbResolver.lastResolvedBlockInfo = BlockInfo{
//		Block:     lastResolve,
//		Timestamp: lastResolvedTimeStamp,
//	}
//	require.NoError(t, err)
//
//	nextBlockInfo, err := lbResolver.getNextDayBlock()
//	require.NoError(t, err)
//
//	sugar.Debugw("result", "next block", nextBlockInfo.Block, "next time", nextBlockInfo.Timestamp.String())
//	assert.Equal(t, expectedBlock, nextBlockInfo.Block)
//}

//func TestLastBlockDaily(t *testing.T) {
//	testutil.SkipExternal(t)
//	var (
//		start        = timeutil.TimestampMsToTime(uint64(1535806920000))
//		end          = timeutil.TimestampMsToTime(uint64(1535954520000))
//		errCh        = make(chan error)
//		blCh         = make(chan BlockInfo)
//		resultsBlock []uint64
//		expectBlocks = []uint64{6255278, 6261305, 6267192}
//		toBreak      = false
//	)
//	logger, err := zap.NewDevelopment()
//	require.NoError(t, err)
//
//	sugar := logger.Sugar()
//
//	ethClient, err := ethclient.Dial("https://mainnet.infura.io/")
//	require.NoError(t, err)
//
//	blkTimeRsv, err := blockchain.NewBlockTimeResolver(sugar, ethClient)
//	require.NoError(t, err)
//
//	lbResolver := NewLastBlockResolver(ethClient, blkTimeRsv, start, end, sugar)
//	go lbResolver.FetchLastBlock(errCh, blCh)
//	for {
//		select {
//		case err := <-errCh:
//			if err == ethereum.NotFound {
//				toBreak = true
//			} else {
//				sugar.Fatalw("error in fetching")
//			}
//		case block := <-blCh:
//			resultsBlock = append(resultsBlock, block.Block)
//		}
//		if toBreak {
//			break
//		}
//	}
//	assert.Equal(t, resultsBlock, expectBlocks)
//}

func TestIsNextDay(t *testing.T) {
	var tests = []struct {
		// RFC3339 string format
		ts1      string
		ts2      string
		expected bool
	}{
		{
			ts1:      "2006-01-02T15:04:05-07:00",
			ts2:      "2006-01-02T15:04:05-07:00",
			expected: false,
		},
		{
			ts1:      "2006-01-02T15:04:05-07:00",
			ts2:      "2006-01-03T15:04:05-07:00",
			expected: true,
		},
		{
			ts1:      "2006-01-02T15:04:05-07:00",
			ts2:      "2006-01-04T15:04:05-07:00",
			expected: false,
		},
		{
			ts1:      "2006-01-02T15:04:05-07:00",
			ts2:      "2009-01-03T15:04:05-07:00",
			expected: false,
		},
	}

	for _, tc := range tests {
		t1, err := time.Parse(time.RFC3339, tc.ts1)
		require.NoError(t, err)

		t2, err := time.Parse(time.RFC3339, tc.ts2)
		require.NoError(t, err)

		assert.Equal(t, tc.expected, isNextDay(t1, t2))
	}
}

func TestNext(t *testing.T) {
	testutil.SkipExternal(t)

	var (
		// Saturday, September 1, 2018 1:02:00 PM
		start = timeutil.TimestampMsToTime(uint64(1535806920000))
		// Saturday, September 8, 2018 6:02:00 AM
		end          = timeutil.TimestampMsToTime(uint64(1536386520000))
		expectBlocks = []uint64{
			6255278, // Sep-01-2018 11:59:50 PM +UTC
			6261305, // Sep-02-2018 11:59:49 PM +UTC
			6267192, // Sep-03-2018 11:59:58 PM +UTC
			6273160, // Sep-04-2018 11:58:30 PM +UTC
			6279116, // Sep-05-2018 11:59:38 PM +UTC
			6285164, // Sep-06-2018 11:59:41 PM +UTC
			6291077, // Sep-07-2018 11:59:47 PM +UTC
		}
	)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	sugar := logger.Sugar()

	ethClient, err := ethclient.Dial("https://mainnet.infura.io/")
	require.NoError(t, err)

	blkTimeRsv, err := blockchain.NewBlockTimeResolver(sugar, ethClient)
	require.NoError(t, err)

	lbResolver := NewLastBlockResolver(ethClient, blkTimeRsv, start, end, sugar)

	var results []uint64
	for {
		var lastBlock uint64
		lastBlock, err = lbResolver.Next()
		if err != nil {
			require.Equal(t, err, ethereum.NotFound)
			break
		} else {
			results = append(results, lastBlock)
		}
	}
	assert.Equal(t, expectBlocks, results)
}
