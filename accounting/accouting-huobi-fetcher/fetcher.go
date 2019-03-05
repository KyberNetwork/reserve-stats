package fetcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"go.uber.org/zap"
)

//Fetcher is the struct to fetch/store Data from huobi
type Fetcher struct {
	sugar  *zap.SugaredLogger
	client huobi.Interface
}

//NewFetcher return a fetcher object
func NewFetcher(sugar *zap.SugaredLogger, client huobi.Interface) *Fetcher {
	return &Fetcher{
		sugar:  sugar,
		client: client,
	}
}

func (fc *Fetcher) getTradeHistoryWithSymbol(symbol string, from, to time.Time, wg *sync.WaitGroup, tradeHistory *sync.Map) error {
	var (
		startTime = from
		endTime   = to
		result    []huobi.TradeHistory
	)
	defer wg.Done()
	for {
		tradeHistoriesResponse, err := fc.client.GetTradeHistory(symbol, startTime, endTime)
		if err != nil {
			return err
		}
		// while result != empty, get trades latest time to toTime
		if len(tradeHistoriesResponse.Data) == 0 {
			break
		}
		result = append(result, tradeHistoriesResponse.Data...)
		lastTrade := tradeHistoriesResponse.Data[len(tradeHistoriesResponse.Data)-1]
		startTime = timeutil.TimestampMsToTime(lastTrade.CreateAt + 1)
	}
	if len(result) != 0 {
		tradeHistory.Store(symbol, result)
	}
	return nil
}

//GetTradeHistory return all trade history between from-to and
func (fc *Fetcher) GetTradeHistory(from, to time.Time) (map[string]huobi.TradeHistoryList, error) {
	var (
		result      = make(map[string]huobi.TradeHistoryList)
		fetchResult = sync.Map{}
		assertError error
	)

	symbols, err := fc.client.GetSymbolsPair()
	if err != nil {
		return result, err
	}
	wg := &sync.WaitGroup{}
	for _, sym := range symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			fc.getTradeHistoryWithSymbol(symbol, from, to, wg, &fetchResult)
		}(sym.SymBol)
	}
	wg.Wait()

	fetchResult.Range(func(key, value interface{}) bool {
		symbol, ok := key.(string)
		if !ok {
			assertError = fmt.Errorf("cannot assert key (value: %v) of map result to symbol (string)", key)
			return false
		}
		historyList, ok := value.(huobi.TradeHistoryList)
		if !ok {
			assertError = fmt.Errorf("cannot assert key (value: %v) of map result to symbol (string)", key)
			return false
		}
		result[symbol] = historyList
		return true
	})
	return result, assertError
}
