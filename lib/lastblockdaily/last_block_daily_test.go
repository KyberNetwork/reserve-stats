package lastblockdaily

import (
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

func TestGetNextDayBlock(t *testing.T) {
	var (
		lastResolve   int64 = 6255278
		expectedBlock int64 = 6261438
		start               = timeutil.TimestampMsToTime(uint64(1535806920000))
		end                 = timeutil.TimestampMsToTime(uint64(1536386520000))
	)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("can't create logger")
	}
	sugar := logger.Sugar()

	ethClient, err := ethclient.Dial("https://mainnet.infura.io/")
	if err != nil {
		sugar.Fatalw("can't create ethClient", "error", err)
	}
	blkTimeRsv, err := blockchain.NewBlockTimeResolver(sugar, ethClient)
	if err != nil {
		sugar.Fatalw("can't create blocktime resolver", "error", err)
	}
	lbResolver := NewLastBlockResolver(ethClient, blkTimeRsv, start, end, sugar)
	lbResolver.lastResolvedBlock = lastResolve
	lbResolver.lastResolvedTimeStamp, err = lbResolver.resolver.Resolve(uint64(lastResolve))
	if err != nil {
		t.Fatalf("can't resolve")
	}
	nextblock, nextTime, err := lbResolver.getNextDayBlock()
	if err != nil {
		t.Fatalf("can't find next d ay block")
	}
	sugar.Debugw("result", "next block", nextblock, "next time", nextTime.String())
	if nextblock != expectedBlock {
		sugar.Errorw("Wrong result", "expected", expectedBlock, "result", nextblock)
		t.Fail()
	}
}

func TestLastBlockDaily(t *testing.T) {
	testutil.SkipExternal(t)
	var (
		start        = timeutil.TimestampMsToTime(uint64(1535806920000))
		end          = timeutil.TimestampMsToTime(uint64(1535954520000))
		errCh        = make(chan error)
		blCh         = make(chan int64, 10)
		resultsBlock = []int64{}
		expectBlocks = []int64{6255278, 6261305, 6267192}
		tobreak      = false
	)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("can't create logger")
	}
	sugar := logger.Sugar()

	ethClient, err := ethclient.Dial("https://mainnet.infura.io/")
	if err != nil {
		sugar.Fatalw("can't create ethClient", "error", err)
	}
	blkTimeRsv, err := blockchain.NewBlockTimeResolver(sugar, ethClient)
	if err != nil {
		sugar.Fatalw("can't create blocktime resolver", "error", err)
	}
	lbResolver := NewLastBlockResolver(ethClient, blkTimeRsv, start, end, sugar)
	go lbResolver.FetchLastBlock(errCh, blCh)
	for {
		select {
		case err := <-errCh:
			if err == ethereum.NotFound {
				tobreak = true
			} else {
				sugar.Fatalw("error in fetching")
			}
		case block := <-blCh:
			resultsBlock = append(resultsBlock, block)
		}
		if tobreak {
			break
		}
	}
	if len(resultsBlock) != len(expectBlocks) {
		sugar.Errorw("wrong result", "expected", expectBlocks, "result", resultsBlock)
		t.Fail()
	}
	for i, rb := range resultsBlock {
		if expectBlocks[i] != rb {
			sugar.Errorw("wrong result", "expected", expectBlocks, "result", resultsBlock)
			t.Fail()
		}
	}
}
