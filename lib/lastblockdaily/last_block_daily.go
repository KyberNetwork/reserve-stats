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

// LastBlockResolving define an interface to resolve last block of a day
type LastBlockResolving interface {
	LastBlock(ts time.Time) int64
}

// this is the block number on 30th JAN 2018.
const (
	firstblockever      = 5000000
	firstBlockTimestamp = 1517319693000
	avgBlockSeconds     = 15 * time.Second
	blockStep           = 200
)

// LastBlockResolver return daily last block between two timestamp
type LastBlockResolver struct {
	client                *ethclient.Client
	resolver              *blockchain.BlockTimeResolver
	start                 time.Time
	end                   time.Time
	lastResolvedBlock     int64
	lastResolvedTimeStamp time.Time
	sugar                 *zap.SugaredLogger
}

// NewLastBlockResolver return a new instance of lastBlock resolver
func NewLastBlockResolver(client *ethclient.Client, resolver *blockchain.BlockTimeResolver, start, end time.Time, sugar *zap.SugaredLogger) *LastBlockResolver {
	return &LastBlockResolver{
		client:                client,
		resolver:              resolver,
		start:                 start,
		end:                   end,
		lastResolvedBlock:     firstblockever,
		lastResolvedTimeStamp: timeutil.TimestampMsToTime(uint64(firstBlockTimestamp)).UTC(),
		sugar:                 sugar,
	}
}

func (lbr *LastBlockResolver) getNextDayBlock() (int64, time.Time, error) {
	var (
		blockTime time.Time
		err       error
	)
	estimate := lbr.lastResolvedBlock + int64(24*time.Hour.Seconds()/avgBlockSeconds.Seconds()) + blockStep
	for {
		blockTime, err = lbr.resolver.Resolve(uint64(estimate))

		if err != nil {
			return 0, blockTime, err
		}
		lbr.sugar.Debugw("getting next day block", "block", estimate, "block time:", blockTime.String())

		if timeutil.Midnight(blockTime) == timeutil.Midnight(lbr.lastResolvedTimeStamp).AddDate(0, 0, 1) {
			estimate += blockStep
			continue
		} else if timeutil.Midnight(lbr.lastResolvedTimeStamp).AddDate(0, 0, 2).Before(timeutil.Midnight(blockTime)) {
			estimate -= blockStep
			continue
		}
		break
	}
	return estimate, blockTime, nil
}

func (lbr *LastBlockResolver) getLastETHBlock() (int64, time.Time, error) {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	header, err := lbr.client.HeaderByNumber(timeout, nil)
	if err != nil {
		return 0, time.Time{}, err
	}
	return header.Number.Int64(), time.Unix(header.Time.Int64(), 0).UTC(), nil
}

func (lbr *LastBlockResolver) isLastBlock(blockNum int64, expectedNextBlockTime time.Time) (bool, time.Time, error) {
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

func (lbr *LastBlockResolver) binSearch(startBlockNum, endBlockNum int64, startBlockTime, endBlockTime, expectedNextBlockTime time.Time) (int64, time.Time, error) {
	lbr.sugar.Debugw("seaching for last block...", "start", startBlockNum, "end", endBlockNum, "starttime", startBlockTime.String(), "endtime", endBlockTime.String())

	midBlockNum := (startBlockNum + endBlockNum) / 2
	isLastBlock, midBlockTime, err := lbr.isLastBlock(midBlockNum, expectedNextBlockTime)
	if err != nil {
		return 0, time.Time{}, nil
	}
	if isLastBlock {
		lbr.sugar.Debugw("found last block of the day", "time", midBlockTime.String(), "block number:", midBlockNum)
		return midBlockNum, midBlockTime, nil
	}
	if expectedNextBlockTime.Before(midBlockTime) {
		return lbr.binSearch(startBlockNum, midBlockNum, startBlockTime, midBlockTime, expectedNextBlockTime)
	}
	return lbr.binSearch(midBlockNum, endBlockNum, midBlockTime, endBlockTime, expectedNextBlockTime)
}

//FetchLastBlock continously push lastBlock of days between start/end
//Will return ethereum.NotFound error if the block is already lastBlock of Ethereum
func (lbr *LastBlockResolver) FetchLastBlock(errCh chan error, queue chan int64) {
	lbr.start = timeutil.Midnight(lbr.start.UTC())
	lbr.end = timeutil.Midnight(lbr.end.UTC())
	var (
		toResolve      time.Time = lbr.start
		endBlockNum    int64
		startBlockNum  int64
		startBlockTime time.Time
		endBlockTime   time.Time
		err            error
	)
	for {
		startBlockNum = lbr.lastResolvedBlock
		startBlockTime = lbr.lastResolvedTimeStamp
		lbr.sugar.Debugw("Last block resolve", "last resolved timestamp", startBlockTime.String(), " last block +1 day timestamp", timeutil.Midnight(lbr.lastResolvedTimeStamp).AddDate(0, 0, 1).String(), "to resolve", toResolve.String())

		if timeutil.Midnight(lbr.lastResolvedTimeStamp).AddDate(0, 0, 1) == timeutil.Midnight(toResolve) {
			lbr.sugar.Debugw("calculating next block")
			// if the last Resolve is yesterday of toResolve, estimate the next day block
			endBlockNum, endBlockTime, err = lbr.getNextDayBlock()
			if err != nil {
				errCh <- err
				return
			}
			lbr.sugar.Debugw("next block to solve", "to resolve", toResolve.String())

		} else {
			//if it is unknown
			endBlockNum, endBlockTime, err = lbr.getLastETHBlock()
			if err != nil {
				errCh <- err
				return
			}
			if endBlockTime.Before(toResolve) {
				errCh <- ethereum.NotFound
				return
			}
		}
		nextLastBlock, nextTimestamp, err := lbr.binSearch(startBlockNum, endBlockNum, startBlockTime, endBlockTime, toResolve.AddDate(0, 0, 1))
		if err != nil {
			errCh <- err
			return
		}
		queue <- nextLastBlock
		lbr.lastResolvedBlock = nextLastBlock
		lbr.lastResolvedTimeStamp = nextTimestamp
		if timeutil.Midnight(lbr.lastResolvedTimeStamp) == timeutil.Midnight(lbr.end) {
			lbr.sugar.Debugw("return to the fuckers", "end block time", endBlockTime.String(), "to resolve", toResolve.String())
			errCh <- ethereum.NotFound
			return
		}
		toResolve = toResolve.AddDate(0, 0, 1)
	}
}
