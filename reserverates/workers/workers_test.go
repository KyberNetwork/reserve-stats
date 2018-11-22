package workers

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/reserverates/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockStorage struct {
}

func (s *mockStorage) UpdateRatesRecords(uint64, map[string]map[string]common.ReserveRateEntry) error {
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
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	return NewPool(sugar, maxWorkers, &mockStorage{})
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

	lastCompleteJobOrder := pool.GetLastCompleteJobOrder()

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
		// expect the last completed job is job 1
		assert.Equal(t, pool.GetLastCompleteJobOrder(), lastCompleteJobOrder+1)
	})
}
