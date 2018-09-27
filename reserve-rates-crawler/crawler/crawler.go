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
	var result []common.Token
	return result, nil
}

func (rrc *ResreveRatesCrawler) GetEachReserveRate(block uint64, rsvAddr ethereum.Address, data *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	tokens, err := rrc.GetSupportedTokens(rsvAddr)
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
	for _, rsvAddr := range rrc.Addresses {
		wg.Add(1)
		go rrc.GetEachReserveRate(block, rsvAddr, &data, &wg)
	}
	wg.Wait()
	return result
}
