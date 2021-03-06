package workers

import (
	"fmt"
	"sync"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
)

type mockStorage struct {
	counter int
	m       sync.Mutex
}

func newMockStorage() *mockStorage {
	return &mockStorage{counter: 0}
}

func (s *mockStorage) Counter() int {
	s.m.Lock()
	defer s.m.Unlock()
	return s.counter
}

func (s *mockStorage) UpdateRatesRecords(uint64, map[string]map[string]common.ReserveRateEntry) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.counter++
	return nil
}
func (s *mockStorage) GetRatesByTimePoint(addrs []ethereum.Address, fromTime, toTime uint64) (map[string]map[string][]common.ReserveRates, error) {
	return nil, nil
}

func (s *mockStorage) LastBlock() (int64, error) {
	return 0, nil
}

type mockJob struct {
	order   int
	failure bool
}

func (j *mockJob) execute(sugar *zap.SugaredLogger) (map[string]map[string]common.ReserveRateEntry, error) {
	if j.failure {
		return nil, fmt.Errorf("failed to execute job %d", j.order)
	}
	return nil, nil
}

func (j *mockJob) info() (order int, block uint64) {
	return j.order, uint64(0)
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
			lastOrder, _ = j.info()
		}

		for pool.GetLastCompleteJobOrder() < lastOrder {
			time.Sleep(time.Millisecond)
		}

		doneCh <- struct{}{}
	}()
}

type assertFn func(t *testing.T, pool *Pool, err error)

func checkWorkerPoolError(t *testing.T, pool *Pool, doneCh <-chan struct{}, fn assertFn) {
	for {
		toBreak := false
		select {
		case <-doneCh:
			// all job success, shut down the pool
			pool.Shutdown()
		case err := <-pool.errCh:
			fn(t, pool, err)
			toBreak = true
		}

		if toBreak {
			break
		}
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

		ms, ok := pool.rateStorage.(*mockStorage)
		require.True(t, ok)
		assert.True(t, ms.Counter() < 2, "no job with order > 2 should trigger database saving")
	})
}
