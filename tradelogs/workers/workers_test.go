package workers

import (
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

type mockStorage struct {
	counter int
	m       sync.Mutex
}

func (s *mockStorage) GetTokenSymbol(address string) (string, error) {
	return "", nil
}

func (s *mockStorage) UpdateTokens(address, symbol []string) error {
	return nil
}

func (s *mockStorage) Counter() int {
	s.m.Lock()
	defer s.m.Unlock()
	return s.counter
}

func newMockStorage() *mockStorage {
	return &mockStorage{counter: 0}
}

func (s *mockStorage) LastBlock() (int64, error) {
	return 0, nil
}

func (s *mockStorage) GetIntegrationVolume(fromTime, toTime time.Time) (map[uint64]*common.IntegrationVolume, error) {
	return nil, nil
}

func (s *mockStorage) SaveTradeLogs(log *common.CrawlResult) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.counter++
	return nil
}

func (s *mockStorage) LoadTradeLogs(from, to time.Time) ([]common.Tradelog, error) {
	return nil, nil
}

func (s *mockStorage) GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddrs []ethereum.Address) (map[ethereum.Address]map[string]float64, error) {
	return nil, nil
}

func (s *mockStorage) GetAssetVolume(token ethereum.Address, fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error) {
	return nil, nil
}

func (s *mockStorage) GetReserveVolume(rsvAddr ethereum.Address, token ethereum.Address, fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error) {
	return nil, nil
}

func (s *mockStorage) GetTradeSummary(fromTime, toTime time.Time, timezone int8) (map[uint64]*common.TradeSummary, error) {
	return nil, nil
}

func (s *mockStorage) GetUserVolume(userAddr ethereum.Address, fromTime, toTime time.Time, freq string) (map[uint64]common.UserVolume, error) {
	return nil, nil
}

func (s *mockStorage) GetWalletStats(fromTime, toTime time.Time, walletAddr string, timezone int8) (map[uint64]common.WalletStats, error) {
	return nil, nil
}

func (s *mockStorage) GetUserList(fromTime, toTime time.Time) ([]common.UserInfo, error) {
	return nil, nil
}

func (s *mockStorage) GetCountryStats(country string, fromTime, toTime time.Time, timezone int8) (map[uint64]*common.CountryStats, error) {
	return nil, nil
}

func (s *mockStorage) GetTokenHeatmap(token ethereum.Address, fromTime, toTime time.Time, timezone int8) (map[string]common.Heatmap, error) {
	return nil, nil
}

func (s *mockStorage) GetMonthlyVolume(rsvAddr ethereum.Address, from, to time.Time) (map[uint64]*common.VolumeStats, error) {
	return nil, nil
}

func (s *mockStorage) LoadTradeLogsByTxHash(tx ethereum.Hash) ([]common.Tradelog, error) {
	return nil, nil
}

func (s *mockStorage) GetStats(from, to time.Time) (common.StatsResponse, error) {
	return common.StatsResponse{}, nil
}

func (s *mockStorage) GetTopTokens(from, to time.Time, limit uint64) (common.TopTokens, error) {
	return common.TopTokens{}, nil
}

func (s *mockStorage) GetTopIntegrations(from, to time.Time, limit uint64) (common.TopIntegrations, error) {
	return common.TopIntegrations{}, nil
}

func (s *mockStorage) GetTopReserves(from, to time.Time, limit uint64) (common.TopReserves, error) {
	return common.TopReserves{}, nil
}

type mockJob struct {
	order   int
	failure bool
}

func (j *mockJob) execute(sugar *zap.SugaredLogger) (*common.CrawlResult, error) {
	if j.failure {
		return nil, fmt.Errorf("failed to execute job %d", j.order)
	}
	return &common.CrawlResult{
		Trades: []common.Tradelog{{
			Timestamp: time.Now(),
		}},
	}, nil
}

func (j *mockJob) info() (order int, from, to *big.Int) {
	return j.order, big.NewInt(0), big.NewInt(0)
}

func newTestWorkerPool(maxWorkers int) *Pool {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	return NewPool(sugar, maxWorkers, newMockStorage())
}

func sendJobsToWorkerPool(pool *Pool, jobs []job, doneCh chan<- struct{}) {
	go func() {
		var lastOrder int
		for _, j := range jobs {
			pool.Run(j)
			lastOrder, _, _ = j.info()
		}

		for pool.GetLastCompleteJobOrder() < lastOrder {
			time.Sleep(time.Millisecond)
		}

		doneCh <- struct{}{}
	}()
}

type assertFn func(t *testing.T, pool *Pool, err error)

func checkWorkerPoolError(t *testing.T, pool *Pool, doneCh <-chan struct{}, fn assertFn) {
	select {
	case <-doneCh:
		pool.Shutdown()
	case err := <-pool.errCh:
		fn(t, pool, err)
	}
}

func TestWorkersPoolHandleAllJobSuccessful(t *testing.T) {
	maxWorkers := 2
	pool := newTestWorkerPool(maxWorkers)

	lastCompleteJobOrder := pool.GetLastCompleteJobOrder()
	numberOfJobs := 3

	doneCh := make(chan struct{})
	var jobs []job
	for i := lastCompleteJobOrder; i < numberOfJobs; i++ {
		jobs = append(jobs, &mockJob{order: i + 1})
	}
	sendJobsToWorkerPool(pool, jobs, doneCh)

	checkWorkerPoolError(t, pool, doneCh, func(t *testing.T, pool *Pool, err error) {
		assert.Equal(t, err, nil)
		assert.Equal(t, pool.GetLastCompleteJobOrder(), lastCompleteJobOrder+numberOfJobs)
	})
}

func TestWorkerPoolEncounterErr(t *testing.T) {
	maxWorkers := 2
	pool := newTestWorkerPool(maxWorkers)

	doneCh := make(chan struct{})
	jobs := []job{
		&mockJob{order: 1},
		&mockJob{order: 2, failure: true},
		&mockJob{order: 3},
		&mockJob{order: 4},
	}

	sendJobsToWorkerPool(pool, jobs, doneCh)

	checkWorkerPoolError(t, pool, doneCh, func(t *testing.T, pool *Pool, err error) {
		assert.Equal(t, err.Error(), "failed to execute job 2")
		assert.True(t, pool.GetLastCompleteJobOrder() < 2, "job with order > 2 should not aborted")

		ms, ok := pool.storage.(*mockStorage)
		require.True(t, ok)
		assert.True(t, ms.Counter() < 2, "no job with order > 2 should trigger database saving")
	})
}
