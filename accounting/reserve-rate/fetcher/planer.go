package fetcher

import (
	"database/sql"
	"sort"
	"time"

	lbdCommon "github.com/KyberNetwork/reserve-stats/lib/lastblockdaily/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

//ReserveBlockInfo contain the reserveAddress and the blockInfo where it's last resolved
type ReserveBlockInfo struct {
	*lbdCommon.BlockInfo
	address ethereum.Address
}

//ReserveBlockInfos abstract the list of ReserveBlockInfo
type ReserveBlockInfos []ReserveBlockInfo

//Len reutnr ReserveBlockInfos len
func (RBIs ReserveBlockInfos) Len() int {
	return len(RBIs)
}

//Swap swaps two ReserveBlockInfo member by index i and j
func (RBIs ReserveBlockInfos) Swap(i, j int) {
	RBIs[i], RBIs[j] = RBIs[j], RBIs[i]
}

//Less return if member i is less than member j
func (RBIs ReserveBlockInfos) Less(i, j int) bool {
	return RBIs[i].Timestamp.Before(RBIs[j].Timestamp)
}

func (fc *Fetcher) getLastFetchedBlockPerReserve() (ReserveBlockInfos, error) {
	var (
		result ReserveBlockInfos
		logger = fc.sugar.With("func", "accounting/reserve-rate/fetcher/planer.go/getLastFetchedBlockPerReserve")
	)
	addresses, err := fc.addressClient.GetAllReserveAddress()
	if err != nil {
		return result, err
	}

	for _, rsv := range addresses {
		fromBlockInfo, err := fc.storage.GetLastResolvedBlockInfo(rsv.Address)
		switch err {
		case sql.ErrNoRows:
			fromBlockInfo = lbdCommon.BlockInfo{Timestamp: rsv.Timestamp}
		case nil:
			logger.Debugw("Got last block info from DB", "reserve", rsv.Address.Hex(), "block", fromBlockInfo.Block)
		default:
			return result, err
		}
		result = append(result, ReserveBlockInfo{
			BlockInfo: &fromBlockInfo,
			address:   rsv.Address,
		})
	}
	sort.Sort(result)
	return result, nil
}

//Run start the fetcher in Daemon mode
func (fc *Fetcher) Run() error {
	var (
		toTime   time.Time
		fromTime time.Time
		logger   = fc.sugar.With("func", "accounting/reserve-rate/Fetcher.Run")
	)
	for {
		rbis, err := fc.getLastFetchedBlockPerReserve()
		if err != nil {
			return err
		}
		index := 0
		rsvAddrs := []ethereum.Address{}
		//rbis is a sorted array of (rsvAddress, block)
		//For example, rbis is [
		//						{(KyberReseve), 10},
		// 						{(MyReseve), 10},
		//						{(NewReseve), 20}
		//					   ]
		//Fetcher will fetch (KyberResreve,MyReserve) from 10 to 20
		//Then it will fetch (KyberResreve,MyReserve,NewReserve) from 20 to now
		for index < len(rbis) {
			rsvAddrs = append(rsvAddrs, rbis[index].address)
			fromTime = rbis[index].Timestamp
			for index+1 < len(rbis) && rbis[index+1].Timestamp == rbis[index].Timestamp {
				index++
				rsvAddrs = append(rsvAddrs, rbis[index].address)
			}

			if index+1 < len(rbis) {
				toTime = rbis[index+1].Timestamp
			} else {
				toTime = time.Now()
			}
			fc.crawler.Addresses = rsvAddrs
			logger.Infow("calling fetch", "fromTime", fromTime.String(), "toTime", toTime.String(), "addresses", rsvAddrs)
			if err := fc.fetch(fromTime, toTime); err != nil {
				return err
			}
			index++
		}

		time.Sleep(fc.sleepTime)
	}
}
