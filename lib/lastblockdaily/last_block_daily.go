package lastblockdaily

import (
	"context"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"

	"go.uber.org/zap"
)

const (
	// this is the block number on 30th JAN 2018.
	firstBlockEver      = 5000000
	firstBlockTimestamp = 1517319693000
	// TODO  :  make this dynamic
	avgBlockSeconds = 15 * time.Second
	blockStep       = 200
)

// LastBlockResolver return daily last block between two timestamp
type LastBlockResolver struct {
	client                *ethclient.Client
	resolver              *blockchain.BlockTimeResolver
	start                 time.Time
	end                   time.Time
	lastResolvedBlockInfo BlockInfo
	sugar                 *zap.SugaredLogger
}

// NewLastBlockResolver return a new instance of lastBlock resolver
func NewLastBlockResolver(client *ethclient.Client, resolver *blockchain.BlockTimeResolver, start, end time.Time, sugar *zap.SugaredLogger) *LastBlockResolver {

	return &LastBlockResolver{
		client:   client,
		resolver: resolver,
		start:    start,
		end:      end,
		lastResolvedBlockInfo: BlockInfo{
			Block:     firstBlockEver,
			Timestamp: timeutil.TimestampMsToTime(uint64(firstBlockTimestamp)).UTC(),
		},
		sugar: sugar,
	}
}

// getNextDayblock return the block that is definitely belong to the next day
// of lastResolvedBlock
func (lbr *LastBlockResolver) getNextDayBlock() (BlockInfo, error) {
	var (
		blockTime time.Time
		err       error
		result    BlockInfo
	)
	estimate := lbr.lastResolvedBlockInfo.Block + uint64(24*time.Hour.Seconds()/avgBlockSeconds.Seconds()) + blockStep
	for {
		blockTime, err = lbr.resolver.Resolve(uint64(estimate))

		if err != nil {
			return result, err
		}
		lbr.sugar.Debugw("getting next day block", "block", estimate, "block_time:", blockTime.String())

		if timeutil.Midnight(blockTime) == timeutil.Midnight(lbr.lastResolvedBlockInfo.Timestamp).AddDate(0, 0, 1) {
			estimate += blockStep
			continue
		} else if timeutil.Midnight(lbr.lastResolvedBlockInfo.Timestamp).AddDate(0, 0, 2).Before(timeutil.Midnight(blockTime)) {
			estimate -= blockStep
			continue
		}
		break
	}
	return BlockInfo{
		Block:     estimate,
		Timestamp: blockTime,
	}, nil
}

func (lbr *LastBlockResolver) getLatestETHBlock() (BlockInfo, error) {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	header, err := lbr.client.HeaderByNumber(timeout, nil)
	if err != nil {
		return BlockInfo{}, err
	}
	return BlockInfo{
		Block:     header.Number.Uint64(),
		Timestamp: time.Unix(header.Time.Int64(), 0).UTC(),
	}, nil
}

func (lbr *LastBlockResolver) isLastBlock(blockNum uint64, expectedNextBlockTime time.Time) (bool, time.Time, error) {
	thisBlockTime, err := lbr.resolver.Resolve(uint64(blockNum))
	if err != nil {
		return false, thisBlockTime, err
	}

	if timeutil.Midnight(thisBlockTime).AddDate(0, 0, 1) != timeutil.Midnight(expectedNextBlockTime) {
		return false, thisBlockTime, nil
	}
	nextBlockTime, err := lbr.resolver.Resolve(uint64(blockNum + 1))
	if err != nil {
		return false, nextBlockTime, err
	}
	if timeutil.Midnight(nextBlockTime) == timeutil.Midnight(expectedNextBlockTime) {
		return true, thisBlockTime, nil
	}
	return false, thisBlockTime, nil
}

func (lbr *LastBlockResolver) binSearch(start, end BlockInfo, expectedNextBlockTime time.Time) (BlockInfo, error) {
	lbr.sugar.Debugw("searching for last block...", "start", start.Block, "end", end.Block, "start_time", start.Timestamp.String(), "end_time", end.Timestamp.String())

	midBlockNum := (start.Block + end.Block) / 2
	isLastBlock, midBlockTime, err := lbr.isLastBlock(midBlockNum, expectedNextBlockTime)
	if err != nil {
		return BlockInfo{}, nil
	}
	midBlockInfo := BlockInfo{
		Block:     midBlockNum,
		Timestamp: midBlockTime,
	}
	if isLastBlock {
		lbr.sugar.Debugw("found last block of the day", "time", midBlockTime.String(), "block number:", midBlockNum)
		return midBlockInfo, nil
	}
	if expectedNextBlockTime.Before(midBlockTime) {
		return lbr.binSearch(start, midBlockInfo, expectedNextBlockTime)
	}
	return lbr.binSearch(midBlockInfo, end, expectedNextBlockTime)
}

//FetchLastBlock continuously push lastBlock of days between start/end
//Will return ethereum.NotFound error if the block is already lastBlock of Ethereum
func (lbr *LastBlockResolver) FetchLastBlock(errCh chan error, queue chan BlockInfo) {
	lbr.start = timeutil.Midnight(lbr.start.UTC())
	lbr.end = timeutil.Midnight(lbr.end.UTC())
	var (
		toResolve      = lbr.start
		endBlockInfo   BlockInfo
		startBlockInfo BlockInfo
		err            error
	)
	for {
		startBlockInfo = lbr.lastResolvedBlockInfo
		lbr.sugar.Debugw("Last block resolve",
			"last resolved timestamp", startBlockInfo.Timestamp.String(),
			"last block +1 day timestamp", timeutil.Midnight(lbr.lastResolvedBlockInfo.Timestamp).AddDate(0, 0, 1).String(),
			"to resolve", toResolve.String(),
		)

		if timeutil.Midnight(lbr.lastResolvedBlockInfo.Timestamp).AddDate(0, 0, 1) == timeutil.Midnight(toResolve) {
			lbr.sugar.Debugw("calculating next block")
			// if the last Resolve is yesterday of toResolve, estimate the next day block
			endBlockInfo, err = lbr.getNextDayBlock()
			if err != nil {
				errCh <- err
				return
			}
			lbr.sugar.Debugw("next block to solve", "to resolve", toResolve.String())

		} else {
			//if it is unknown, get Latest ETH Block
			endBlockInfo, err = lbr.getLatestETHBlock()
			if err != nil {
				errCh <- err
				return
			}
			if endBlockInfo.Timestamp.Before(toResolve) {
				errCh <- ethereum.NotFound
				return
			}
		}
		nextBlockInfo, err := lbr.binSearch(startBlockInfo, endBlockInfo, toResolve.AddDate(0, 0, 1))
		if err != nil {
			errCh <- err
			return
		}
		queue <- nextBlockInfo
		lbr.lastResolvedBlockInfo = nextBlockInfo
		if timeutil.Midnight(lbr.lastResolvedBlockInfo.Timestamp) == timeutil.Midnight(lbr.end) {
			errCh <- ethereum.NotFound
			return
		}
		toResolve = toResolve.AddDate(0, 0, 1)
	}
}
