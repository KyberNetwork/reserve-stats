package fetcher

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	lbdCommon "github.com/KyberNetwork/reserve-stats/lib/lastblockdaily/common"
)

func (fc *Fetcher) isFailed() bool {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	return fc.failed
}

func (fc *Fetcher) markAsFailed() {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	fc.failed = true
}

func (fc *Fetcher) getLastCompletedJobOrder() uint64 {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	return fc.lastCompletedJobOrder
}

func (fc *Fetcher) serialDataStore(blockInfo lbdCommon.BlockInfo, rates map[string]map[string]float64, ethUSDRate float64, jobOrder uint64) error {
	var logger = fc.sugar.With("func", caller.GetCurrentFunctionName(),
		"block", blockInfo.Block, "job_order", jobOrder)
	for {
		if fc.isFailed() {
			return fmt.Errorf("fetcher has failed before job order %d", jobOrder)
		}
		lastCompleted := fc.getLastCompletedJobOrder()
		if lastCompleted+1 == jobOrder {
			if err := fc.storage.UpdateRatesRecords(blockInfo, rates, ethUSDRate); err != nil {
				fc.markAsFailed()
				return err
			}
			fc.mutex.Lock()
			fc.lastCompletedJobOrder = jobOrder
			fc.mutex.Unlock()
			return nil
		}
		logger.Debugw("waiting for previous job to completed", "last_completed", lastCompleted)
		time.Sleep(time.Second)
	}
}
