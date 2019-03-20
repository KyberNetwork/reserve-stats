package crawler

import (
	"testing"

	"github.com/KyberNetwork/reserve-stats/burnedfees/common"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage struct{}

func newMockStorage() *mockStorage {
	return &mockStorage{}
}

func (*mockStorage) Store([]common.BurnAssignedFeesEvent) error {
	return nil
}

func (*mockStorage) LastBlock() (int64, error) {
	return 0, nil
}

func TestBurnedFeesCrawlerExecute(t *testing.T) {
	testutil.SkipExternal(t)
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	ethClient := testutil.MustNewDevelopmentwEthereumClient()

	burners := []ethereum.Address{ethereum.HexToAddress("0xed4f53268bfdff39b36e8786247ba3a02cf34b04")}
	crawler := NewBurnedFeesCrawler(sugar, ethClient, newMockStorage(), burners)
	events, err := crawler.crawl(7019633, 7025313)
	require.NoError(t, err)
	require.Len(t, events, 3)
	assert.Equal(t, events[0].BlockNumber, uint64(7019633))
	assert.Equal(t, events[1].BlockNumber, uint64(7019643))
	assert.Equal(t, events[2].BlockNumber, uint64(7019666))
}
