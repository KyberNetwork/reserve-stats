package lastblockdaily

import (
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

func TestGetNextDayBlock(t *testing.T) {
	var (
		lastResolve   uint64 = 6255278
		expectedBlock uint64 = 6261438
		start                = timeutil.TimestampMsToTime(uint64(1535806920000))
		end                  = timeutil.TimestampMsToTime(uint64(1536386520000))
	)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	sugar := logger.Sugar()

	ethClient, err := ethclient.Dial("https://mainnet.infura.io/")
	require.NoError(t, err)

	blkTimeRsv, err := blockchain.NewBlockTimeResolver(sugar, ethClient)
	require.NoError(t, err)

	lbResolver := NewLastBlockResolver(ethClient, blkTimeRsv, start, end, sugar)

	lastResolvedTimeStamp, err := lbResolver.resolver.Resolve(uint64(lastResolve))
	require.NoError(t, err)
	lbResolver.lastResolvedBlockInfo = BlockInfo{
		Block:     lastResolve,
		Timestamp: lastResolvedTimeStamp,
	}
	require.NoError(t, err)

	nextBlockInfo, err := lbResolver.getNextDayBlock()
	require.NoError(t, err)

	sugar.Debugw("result", "next block", nextBlockInfo.Block, "next time", nextBlockInfo.Timestamp.String())
	assert.Equal(t, expectedBlock, nextBlockInfo.Block)
}

func TestLastBlockDaily(t *testing.T) {
	testutil.SkipExternal(t)
	var (
		start        = timeutil.TimestampMsToTime(uint64(1535806920000))
		end          = timeutil.TimestampMsToTime(uint64(1535954520000))
		errCh        = make(chan error)
		blCh         = make(chan BlockInfo)
		resultsBlock = []uint64{}
		expectBlocks = []uint64{6255278, 6261305, 6267192}
		tobreak      = false
	)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	sugar := logger.Sugar()

	ethClient, err := ethclient.Dial("https://mainnet.infura.io/")
	require.NoError(t, err)

	blkTimeRsv, err := blockchain.NewBlockTimeResolver(sugar, ethClient)
	require.NoError(t, err)

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
			resultsBlock = append(resultsBlock, block.Block)
		}
		if tobreak {
			break
		}
	}
	assert.Equal(t, resultsBlock, expectBlocks)
}
