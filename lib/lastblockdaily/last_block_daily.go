package lastblockdaily

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"go.uber.org/zap"
)

const (
	// average block time, use for estimation.
	avgBlock = 15 * time.Second
	// blockStep is the number of blocks we try until find last block of the day. There
	// i assumption that a day never has less than blockStep number of blocks.
	blockStep = 200
)

var (
	firstBlock = BlockInfo{
		Block: 5000000,
		// Jan-30-2018 01:41:33 PM +UTC
		Timestamp: timeutil.TimestampMsToTime(1517319693000),
	}
)

// LastBlockResolver return daily last block between two timestamp
type LastBlockResolver struct {
	sugar        *zap.SugaredLogger
	client       *ethclient.Client
	resolver     *blockchain.BlockTimeResolver
	start        time.Time
	end          time.Time
	lastResolved BlockInfo
	avgBlockTime time.Duration
}

// NewLastBlockResolver return a new instance of lastBlock resolver
func NewLastBlockResolver(client *ethclient.Client, resolver *blockchain.BlockTimeResolver, start, end time.Time, sugar *zap.SugaredLogger) *LastBlockResolver {
	return &LastBlockResolver{
		client:       client,
		resolver:     resolver,
		start:        start,
		end:          end,
		lastResolved: BlockInfo{Timestamp: start},
		sugar:        sugar,
		avgBlockTime: avgBlock,
	}
}

// isSameDay returns true if both t1 an t2 are in a same day in UTC timezone.
func isSameDay(t1, t2 time.Time) bool {
	t1, t2 = t1.UTC(), t2.UTC()
	midnight1 := timeutil.Midnight(t1)
	midnight2 := timeutil.Midnight(t2)
	return midnight1.Equal(midnight2)
}

// isNextDay returns true if t2 is chronologically in the next day of t1 in UTC timezone.
func isNextDay(t1, t2 time.Time) bool {
	t1, t2 = t1.UTC(), t2.UTC()
	midnight1 := timeutil.Midnight(t1)
	midnight2 := timeutil.Midnight(t2)
	if midnight1.AddDate(0, 0, 1).Equal(midnight2) {
		return true
	}
	return false
}

func nextDay(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}

// nextDayBlock returns a Ethereum block number in the next day of given timestamp.
// The known block info if provided will be used for estimation.
func nextDayBlock(
	sugar *zap.SugaredLogger,
	resolver blockchain.BlockTimeResolverInterface,
	current time.Time,
	known BlockInfo) (BlockInfo, error) {
	var (
		logger = sugar.With(
			"func", "lib/lastblockdaily/nextDayBlock",
		)
		next uint64
	)

	if known.Block == 0 {
		known = firstBlock
	}

	logger = logger.With(
		"known_block", known.Block,
		"known_block_time", known.Timestamp.String(),
	)

	if current.Before(known.Timestamp) {
		return BlockInfo{}, fmt.Errorf("start time is too old: requested=%s, oldest supported:%s",
			current.String(),
			firstBlock.Timestamp.String())
	}

	// estimate approximate block gap between firstBlockEver and start time.
	estimatedBlocks := uint64(nextDay(current).Sub(known.Timestamp) / avgBlock)
	logger.Debugw("estimated number of blocks gap", "estimated_blocks", estimatedBlocks)

	next = known.Block + estimatedBlocks
	for {
		nextTimestamp, err := resolver.Resolve(next)
		if err != nil {
			return BlockInfo{}, err
		}

		// estimated block is older than start time or in the same day as start, trying newer block
		if nextTimestamp.Before(current) || isSameDay(current, nextTimestamp) {
			next += blockStep
			continue
		}

		if isNextDay(current, nextTimestamp) {
			return BlockInfo{
				Block:     next,
				Timestamp: nextTimestamp,
			}, nil
		}

		// estimated block is newer than next day of start, trying older block
		next -= blockStep
	}
}

// Next returns next last block number of the day from lastResolvedBlockInfo. The resolver works by find a block that
// belongs to next day of lastResolvedBlockInfo, and do binary search between this block and lastResolvedBlockInfo
// until found the last block of the day.
// return ethereum.NotFound
func (lbr *LastBlockResolver) Next() (BlockInfo, error) {
	if isSameDay(lbr.lastResolved.Timestamp, lbr.end) {
		return BlockInfo{}, ethereum.NotFound
	}
	var (
		logger = lbr.sugar.With(
			"func", "lib/lastblockdaily/LastBlockResolver.Next",
		)
		start BlockInfo
		end   BlockInfo
		err   error
	)

	if lbr.lastResolved.Block == 0 {
		// no last block of the day resolved, starts searching with:
		// - from: first ever block timestamp
		// - to: next day of configured start timestamp
		start = firstBlock

		if end, err = nextDayBlock(lbr.sugar, lbr.resolver, lbr.start, start); err != nil {
			return BlockInfo{}, err
		}
	} else {
		// last resolved block is last block of the day d1
		// search for last block of the day d1 + 1
		//   - from: last resolved + 1 block timestamp (in d1 + 1 day)
		//   - to: next day block of resolved + 1 (in d1 + 2 day)
		start = BlockInfo{Block: lbr.lastResolved.Block + 1}
		var startTime time.Time
		startTime, err = lbr.resolver.Resolve(start.Block)
		if err != nil {
			return BlockInfo{}, err
		}
		start.Timestamp = startTime

		if end, err = nextDayBlock(lbr.sugar, lbr.resolver, start.Timestamp, start); err != nil {
			return BlockInfo{}, err
		}
	}

	logger = logger.With(
		"start", start.Block,
		"start_time", start.Timestamp.String(),
		"stop", lbr.end.String(),
	)

	logger.Debug("getting next last block of the day")

	lastBlock, err := lbr.searchLastBlock(start, end)
	if err != nil {
		return BlockInfo{}, err
	}
	lbr.lastResolved = lastBlock
	return lastBlock, nil
}

// isLastBlock check if the block is the last block of the day.
func (lbr *LastBlockResolver) isLastBlock(block BlockInfo) (bool, error) {
	nextBlockTimestamp, err := lbr.resolver.Resolve(block.Block + 1)
	if err != nil {
		return false, err
	}

	if isNextDay(block.Timestamp, nextBlockTimestamp) {
		return true, nil
	}

	return false, nil
}

// searchLastBlock returns the last block of the day before end block.
func (lbr *LastBlockResolver) searchLastBlock(start, end BlockInfo) (BlockInfo, error) {
	var logger = lbr.sugar.With(
		"func", "lib/lastblockdaily.LastBlockResolver.searchLastBlock",
		"start", start.Block,
		"start_time", start.Timestamp.String(),
		"end", end.Block,
		"end_time", end.Timestamp.String(),
	)

	if start.Block == end.Block || end.Block-start.Block == 1 {
		logger.Infow("found last block of the day",
			"block", start.Block,
			"block_time", start.Timestamp,
		)
		return start, nil
	}

	midBlockNum := (start.Block + end.Block) / 2
	midBlockTimestamp, err := lbr.resolver.Resolve(midBlockNum)
	if err != nil {
		return BlockInfo{}, err
	}

	mid := BlockInfo{
		Block:     midBlockNum,
		Timestamp: midBlockTimestamp,
	}

	logger = logger.With("mid", mid.Block, "mid_time", mid.Timestamp.String())
	logger.Debugw("searching for last block...")

	if isSameDay(mid.Timestamp, end.Timestamp) {
		return lbr.searchLastBlock(start, mid)
	}
	return lbr.searchLastBlock(mid, end)
}

// Run push the result/ error into channels
func (lbr *LastBlockResolver) Run(resultChn chan BlockInfo, errChn chan error) {
	var (
		lastBlockInfo BlockInfo
		err           error
	)
	for {
		lastBlockInfo, err = lbr.Next()
		if err != nil {
			errChn <- err
			break
		}
		resultChn <- lastBlockInfo
	}
}
