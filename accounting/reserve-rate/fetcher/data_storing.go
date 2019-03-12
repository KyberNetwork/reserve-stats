package fetcher

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	rsvRateCommon "github.com/KyberNetwork/reserve-stats/reserverates/common"
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

func (fc *Fetcher) serialDataStore(blockInfo lastblockdaily.BlockInfo, rates map[string]map[string]rsvRateCommon.ReserveRateEntry, ethUSDRate float64, jobOrder uint64) error {
	var logger = fc.sugar.With("func", "accounting/reserve-rate/fetcher/fetcher.serialDataStore", "block", blockInfo.Block, "job_order", jobOrder)
	for {
		if fc.isFailed() {
			return fmt.Errorf("fetcher has failed before job order %d", jobOrder)
		}
		lastCompleted := fc.getLastCompletedJobOrder()
		if lastCompleted+1 == jobOrder {
			if err := fc.storage.UpdateRatesRecords(blockInfo, rates); err != nil {
				fc.markAsFailed()
				return err
			}
			if err := fc.storage.UpdateETHUSDPrice(blockInfo, ethUSDRate); err != nil {
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
