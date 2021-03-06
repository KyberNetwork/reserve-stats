package lastblockdaily

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily/common"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily/storage"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	// average block time, use for estimation.
	avgBlock = 15 * time.Second
	// blockStep is the number of blocks we try until find last block of the day. There
	// i assumption that a day never has less than blockStep number of blocks.
	blockStep = 200
)

var (
	firstBlock = common.BlockInfo{
		Block: 5000000,
		// Jan-30-2018 01:41:33 PM +UTC
		Timestamp: timeutil.TimestampMsToTime(1517319693000),
	}
)

// LastBlockResolver return daily last block between two timestamp
type LastBlockResolver struct {
	sugar        *zap.SugaredLogger
	client       *ethclient.Client
	resolver     blockchain.BlockTimeResolverInterface
	LastResolved common.BlockInfo
	avgBlockTime time.Duration
	db           storage.Interface
}

// NewLastBlockResolver return a new instance of lastBlock resolver
func NewLastBlockResolver(client *ethclient.Client, resolver *blockchain.BlockTimeResolver, sugar *zap.SugaredLogger, db storage.Interface) *LastBlockResolver {
	return &LastBlockResolver{
		client:       client,
		resolver:     resolver,
		sugar:        sugar,
		avgBlockTime: avgBlock,
		db:           db,
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
	return nextDay(midnight1).Equal(midnight2)
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
	known common.BlockInfo) (common.BlockInfo, error) {
	var (
		logger = sugar.With(
			"func", caller.GetCurrentFunctionName(),
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
		return common.BlockInfo{}, fmt.Errorf("start time is too old: requested=%s, oldest supported:%s",
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
			return common.BlockInfo{}, err
		}

		// estimated block is older than start time or in the same day as start, trying newer block
		if nextTimestamp.Before(current) || isSameDay(current, nextTimestamp) {
			next += blockStep
			continue
		}

		if isNextDay(current, nextTimestamp) {
			return common.BlockInfo{
				Block:     next,
				Timestamp: nextTimestamp,
			}, nil
		}

		// estimated block is newer than next day of start, trying older block
		next -= blockStep
	}
}

// Next returns next last block number of the day from lastResolvedcommon.BlockInfo. The resolver works by find a block that
// belongs to next day of lastResolvedcommon.BlockInfo, and do binary search between this block and lastResolvedcommon.BlockInfo
// until found the last block of the day.
// return ethereum.NotFound
func (lbr *LastBlockResolver) Next() (common.BlockInfo, error) {

	var (
		logger = lbr.sugar.With(
			"func", caller.GetCurrentFunctionName(),
		)
		start common.BlockInfo
		end   common.BlockInfo
		err   error
	)

	if lbr.LastResolved.Block == 0 {
		// no last block of the day resolved
		return common.BlockInfo{}, fmt.Errorf("no last resolved block to call Next()")
	}
	//look DB first
	nextDay := nextDay(lbr.LastResolved.Timestamp)
	blockInfo, err := lbr.db.GetBlockInfo(nextDay)
	switch err {
	case sql.ErrNoRows:
		logger.Debug("no such result in db, proceed forward")
	case nil:
		logger.Debug("got result from db, return...")
		lbr.LastResolved = blockInfo
		return blockInfo, nil
	default:
		return common.BlockInfo{}, err
	}

	// last resolved block is last block of the day d1
	// search for last block of the day d1 + 1
	//   - from: last resolved + 1 block timestamp (in d1 + 1 day)
	//   - to: next day block of resolved + 1 (in d1 + 2 day)
	start = common.BlockInfo{Block: lbr.LastResolved.Block + 1}
	var startTime time.Time
	startTime, err = lbr.resolver.Resolve(start.Block)
	if err != nil {
		return common.BlockInfo{}, err
	}
	start.Timestamp = startTime

	if end, err = nextDayBlock(lbr.sugar, lbr.resolver, start.Timestamp, start); err != nil {
		return common.BlockInfo{}, err
	}

	logger = logger.With(
		"start", start.Block,
		"start_time", start.Timestamp.String(),
		"stop", end.Timestamp.String(),
	)

	logger.Debug("getting next last block of the day")

	lastBlock, err := lbr.searchLastBlock(start, end)
	if err != nil {
		return common.BlockInfo{}, err
	}
	if err = lbr.db.UpdateBlockInfo(lastBlock); err != nil {
		return lastBlock, err
	}

	lbr.LastResolved = lastBlock
	return lastBlock, nil
}

// searchLastBlock returns the last block of the day before end block.
func (lbr *LastBlockResolver) searchLastBlock(start, end common.BlockInfo) (common.BlockInfo, error) {
	var logger = lbr.sugar.With(
		"func", caller.GetCurrentFunctionName(),
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
		return common.BlockInfo{}, err
	}

	mid := common.BlockInfo{
		Block:     midBlockNum,
		Timestamp: midBlockTimestamp,
	}

	if isSameDay(mid.Timestamp, end.Timestamp) {
		return lbr.searchLastBlock(start, mid)
	}
	return lbr.searchLastBlock(mid, end)
}

// Run push the result/ error into channels
func (lbr *LastBlockResolver) Run(from, to time.Time, resultChn chan common.BlockInfo, errChn chan error) {
	var (
		lastBlockInfo common.BlockInfo
		err           error
	)
	for {
		if lbr.LastResolved.Block == 0 {
			lastBlockInfo, err = lbr.Resolve(from)
		} else {
			lastBlockInfo, err = lbr.Next()
		}
		if err != nil {
			errChn <- err
			break
		}
		resultChn <- lastBlockInfo
		if to.Before(lastBlockInfo.Timestamp) {
			errChn <- ethereum.NotFound
		}
	}
}

//Resolve return last block of the day with timeInput
//Call this function when there is a need to resolve a date that is not Next()
//  or when there is no know block except for firstBlock (5000000)
func (lbr *LastBlockResolver) Resolve(date time.Time) (common.BlockInfo, error) {
	var (
		start  = firstBlock
		logger = lbr.sugar.With(
			"start", start.Block,
			"start_time", start.Timestamp.String(),
		)
	)
	logger.Debug("getting next last block of the day")

	blockInfo, err := lbr.db.GetBlockInfo(date)
	switch err {
	case sql.ErrNoRows:
		logger.Debug("no such result in db, proceed forward")
	case nil:
		logger.Debug("got result from db, return...")
		lbr.LastResolved = blockInfo
		return blockInfo, nil
	default:
		return common.BlockInfo{}, err
	}
	//if lastResolvedBlock is available (>0) and is before date, use lastResovled
	if lbr.LastResolved.Block > start.Block && lbr.LastResolved.Timestamp.Before(date) {
		start = lbr.LastResolved
	}

	end, err := nextDayBlock(lbr.sugar, lbr.resolver, date, start)
	if err != nil {
		return common.BlockInfo{}, err
	}

	lastBlock, err := lbr.searchLastBlock(start, end)
	if err != nil {
		return common.BlockInfo{}, err
	}
	if err = lbr.db.UpdateBlockInfo(lastBlock); err != nil {
		return lastBlock, err
	}

	lbr.LastResolved = lastBlock
	return lastBlock, err
}
