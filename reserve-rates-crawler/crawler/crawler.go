package crawler

import (
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/KyberNetwork/reserve-stats/common"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ResreveRatesCrawler contains two wrapper contracts for V1 and V2 contract,
// a set of addresses to crawl rates from and setting object to query for reserve's token settings
type ResreveRatesCrawler struct {
	WrapperContractV1 *contracts.Wrapper
	WrapperContractV2 *contracts.Wrapper
	Addresses         []ethereum.Address
	setting           Setting
}

// NewReserveRatesCrawler returns an instant of ReserveRatesCrawler.
func NewReserveRatesCrawler(addrs []string, client *ethclient.Client, sett Setting) (*ResreveRatesCrawler, error) {
	wrapperContractV1, err := contracts.NewWrapper(common.WrapperAddrV1, client)
	wrapperContractV2, err := contracts.NewWrapper(common.WrapperAddrV2, client)

	if err != nil {
		return nil, err
	}
	var ethAddrs []ethereum.Address
	for _, addr := range addrs {
		ethAddrs = append(ethAddrs, ethereum.HexToAddress(addr))
	}
	return &ResreveRatesCrawler{
		WrapperContractV1: wrapperContractV1,
		WrapperContractV2: wrapperContractV2,
		Addresses:         ethAddrs,
		setting:           sett,
	}, nil
}

func (rrc *ResreveRatesCrawler) getSupportedTokens(rsvAddr ethereum.Address) ([]common.Token, error) {
	return rrc.setting.GetInternalTokens()
}

func (rrc *ResreveRatesCrawler) getEachReserveRate(block uint64, rsvAddr ethereum.Address, data *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	tokens, err := rrc.getSupportedTokens(rsvAddr)
	if err != nil {
		log.Printf("can not get supported tokens for reserve %s", rsvAddr.Hex())
		return
	}
	ETH := common.ETHToken
	srcAddresses := []ethereum.Address{}
	destAddresses := []ethereum.Address{}
	for _, token := range tokens {
		srcAddresses = append(srcAddresses, ethereum.HexToAddress(token.Address), ethereum.HexToAddress(ETH.Address))
		destAddresses = append(destAddresses, ethereum.HexToAddress(ETH.Address), ethereum.HexToAddress(token.Address))
	}
	var (
		reserveRate []*big.Int
		sanityRate  []*big.Int
		callError   error
	)
	if block < common.StartingBlockV2 {
		reserveRate, sanityRate, callError = rrc.WrapperContractV1.GetReserveRate(&bind.CallOpts{BlockNumber: big.NewInt(int64(block))}, rsvAddr, srcAddresses, destAddresses)
	} else {
		reserveRate, sanityRate, callError = rrc.WrapperContractV2.GetReserveRate(&bind.CallOpts{BlockNumber: big.NewInt(int64(block))}, rsvAddr, srcAddresses, destAddresses)
	}
	if callError != nil {
		log.Printf("can not get reserve rate for reserve %s, error : %s", rsvAddr.Hex(), callError)
		return
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
}

// GetReserveRates returns the map[ReserveAddress]ReserveRates at the given block number.
// It will only return rates from the set of addresses within its definition.
func (rrc *ResreveRatesCrawler) GetReserveRates(block uint64) map[string]common.ReserveRates {
	result := make(map[string]common.ReserveRates)
	data := sync.Map{}
	wg := sync.WaitGroup{}
	for _, rsvAddr := range rrc.Addresses {
		wg.Add(1)
		go rrc.getEachReserveRate(block, rsvAddr, &data, &wg)
	}
	wg.Wait()
	data.Range(func(key, value interface{}) bool {
		reserveAddr, ok := key.(ethereum.Address)
		//if there is conversion error, continue to next key,val
		if !ok {
			log.Printf("key (%v) cannot be asserted to ethereum.Address", key)
			return true
		}
		rates, ok := value.(common.ReserveRates)
		if !ok {
			log.Printf("value (%v) cannot be asserted to reserveRates", value)
			return true
		}
		result[reserveAddr.Hex()] = rates
		return true
	})
	return result
}
