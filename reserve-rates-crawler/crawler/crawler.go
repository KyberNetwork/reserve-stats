package crawler

import (
	"log"
	"sync"

	"github.com/KyberNetwork/reserve-stats/common"
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/blockchain"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type ResreveRatesCrawler struct {
	Blockchain *blockchain.Blockchain
	Addresses  []ethereum.Address
	setting    Setting
}

func NewReserveRatesCrawler(addrs []string, endpoint string, sett Setting) (*ResreveRatesCrawler, error) {
	blockchain, err := blockchain.NewBlockchain(endpoint)
	if err != nil {
		return nil, err
	}
	var ethAddrs []ethereum.Address
	for _, addr := range addrs {
		ethAddrs = append(ethAddrs, ethereum.HexToAddress(addr))
	}
	return &ResreveRatesCrawler{
		Blockchain: blockchain,
		Addresses:  ethAddrs,
		setting:    sett,
	}, nil
}

func (rrc *ResreveRatesCrawler) GetSupportedTokens(rsvAddr ethereum.Address) ([]common.Token, error) {
	return rrc.setting.GetInternalTokens()
}

func (rrc *ResreveRatesCrawler) GetEachReserveRate(block uint64, rsvAddr ethereum.Address, data *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	tokens, err := rrc.GetSupportedTokens(rsvAddr)
	log.Printf("token is %v", tokens)
	if err != nil {
		log.Printf("can not get supported tokens for reserve %s", rsvAddr.Hex())
		return
	}
	rates, err := rrc.Blockchain.GetReserveRates(block, rsvAddr, tokens)
	data.Store(rsvAddr, rates)
}

func (rrc *ResreveRatesCrawler) GetReserveRates(block uint64) map[string]common.ReserveRates {
	result := make(map[string]common.ReserveRates)
	data := sync.Map{}
	wg := sync.WaitGroup{}
	log.Printf("rcc addrs is %v", rrc.Addresses)
	for _, rsvAddr := range rrc.Addresses {
		wg.Add(1)
		go rrc.GetEachReserveRate(block, rsvAddr, &data, &wg)
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
