package crawler

import (
	"fmt"
	"sync"

	"github.com/KyberNetwork/reserve-stats/common"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// ResreveRatesCrawler contains two wrapper contracts for V1 and V2 contract,
// a set of addresses to crawl rates from and setting object to query for reserve's token settings
type ResreveRatesCrawler struct {
	wrapperContract *contracts.VersionedWrapper
	Addresses       []ethereum.Address
	setting         Setting
	logger          *zap.SugaredLogger
}

// NewReserveRatesCrawler returns an instant of ReserveRatesCrawler.
func NewReserveRatesCrawler(addrs []string, client *ethclient.Client, sett Setting, lger *zap.SugaredLogger) (*ResreveRatesCrawler, error) {
	wrpContract, err := contracts.NewVersionedWrapper(client)
	if err != nil {
		return nil, err
	}
	var ethAddrs []ethereum.Address
	for _, addr := range addrs {
		ethAddrs = append(ethAddrs, ethereum.HexToAddress(addr))
	}
	return &ResreveRatesCrawler{
		wrapperContract: wrpContract,
		Addresses:       ethAddrs,
		setting:         sett,
		logger:          lger,
	}, nil
}

func (rrc *ResreveRatesCrawler) getSupportedTokens(rsvAddr ethereum.Address) ([]common.Token, error) {
	return rrc.setting.GetInternalTokens()
}

func (rrc *ResreveRatesCrawler) getEachReserveRate(block uint64, rsvAddr ethereum.Address, data *sync.Map, wg *sync.WaitGroup) error {
	defer wg.Done()
	tokens, err := rrc.getSupportedTokens(rsvAddr)
	if err != nil {
		return fmt.Errorf("cannot get supported tokens for reserve %s. Error: %s", rsvAddr.Hex(), err)
	}
	ETH := common.ETHToken
	srcAddresses := []ethereum.Address{}
	destAddresses := []ethereum.Address{}
	for _, token := range tokens {
		srcAddresses = append(srcAddresses, ethereum.HexToAddress(token.Address), ethereum.HexToAddress(ETH.Address))
		destAddresses = append(destAddresses, ethereum.HexToAddress(ETH.Address), ethereum.HexToAddress(token.Address))
	}

	reserveRate, sanityRate, callError := rrc.wrapperContract.GetReserveRate(block, rsvAddr, srcAddresses, destAddresses)
	if callError != nil {
		return fmt.Errorf("cannot get rates for reserve %s. Error: %s", rsvAddr.Hex(), callError)
	}
	rates := common.ReserveRates{}
	rsvTokenRateEntry := common.ReserveTokenRateEntry{}
	rates.Timestamp = common.GetTimepoint()
	rates.BlockNumber = block - 1
	rates.ToBlockNumber = block
	rates.ReturnTime = common.GetTimepoint()
	for index, token := range tokens {
		rateEntry := common.ReserveRateEntry{}
		rateEntry.BuyReserveRate = common.BigToFloat(reserveRate[index*2+1], 18)
		rateEntry.BuySanityRate = common.BigToFloat(sanityRate[index*2+1], 18)
		rateEntry.SellReserveRate = common.BigToFloat(reserveRate[index*2], 18)
		rateEntry.SellSanityRate = common.BigToFloat(sanityRate[index*2], 18)
		rsvTokenRateEntry[fmt.Sprintf("ETH-%s", token.ID)] = rateEntry
	}
	rates.Data = rsvTokenRateEntry
	data.Store(rsvAddr, rates)
	return nil
}

// GetReserveRates returns the map[ReserveAddress]ReserveRates at the given block number.
// It will only return rates from the set of addresses within its definition.
func (rrc *ResreveRatesCrawler) GetReserveRates(block uint64) (map[string]common.ReserveRates, error) {
	result := make(map[string]common.ReserveRates)
	data := sync.Map{}
	wg := sync.WaitGroup{}
	errs := make(chan error, len(rrc.Addresses))
	for _, rsvAddr := range rrc.Addresses {
		wg.Add(1)
		go func(addr ethereum.Address) {
			err := rrc.getEachReserveRate(block, addr, &data, &wg)
			errs <- err
		}(rsvAddr)
	}
	wg.Wait()
	defer close(errs)
	for i := 0; i < len(rrc.Addresses); i++ {
		err := <-errs
		if err != nil {
			return result, err
		}
	}
	var cErr error
	data.Range(func(key, value interface{}) bool {
		reserveAddr, ok := key.(ethereum.Address)
		//if there is conversion error, continue to next key,val
		if !ok {
			cErr = fmt.Errorf("key (%v) cannot be asserted to ethereum.Address", key)
			return false
		}
		rates, ok := value.(common.ReserveRates)
		if !ok {
			cErr = fmt.Errorf("value (%v) cannot be asserted to reserveRates", value)
			return true
		}
		result[reserveAddr.Hex()] = rates
		return true
	})
	return result, cErr
}
