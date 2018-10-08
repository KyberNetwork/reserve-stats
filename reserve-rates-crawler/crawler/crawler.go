package crawler

import (
	"fmt"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/common"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

var (
	//InternalReserveAddr is the Kyber's own reserve address
	InternalReserveAddr = ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F")
	ethToken            = core.Token{
		ID:       "ETH",
		Name:     "Ethereum",
		Address:  "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		Decimals: 18,
	}
)

// ResreveRatesCrawler contains two wrapper contracts for V1 and V2 contract,
// a set of addresses to crawl rates from and setting object to query for reserve's token settings
type ResreveRatesCrawler struct {
	wrapperContract *contracts.VersionedWrapper
	Addresses       []ethereum.Address
	tokenSetting    TokenSetting
	logger          *zap.SugaredLogger
}

// NewReserveRatesCrawler returns an instant of ReserveRatesCrawler.
func NewReserveRatesCrawler(addrs []string, client *ethclient.Client, sett TokenSetting, lger *zap.SugaredLogger) (*ResreveRatesCrawler, error) {
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
		tokenSetting:    sett,
		logger:          lger,
	}, nil
}

func (rrc *ResreveRatesCrawler) callTokens(rsvAddr ethereum.Address) ([]core.Token, error) {
	if rsvAddr.Hex() == InternalReserveAddr.Hex() {
		return rrc.tokenSetting.GetInternalTokens()
	}
	return rrc.tokenSetting.GetActiveTokens()
}

func (rrc *ResreveRatesCrawler) getSupportedTokens(rsvAddr ethereum.Address) ([]core.Token, error) {
	tokens := []core.Token{}
	tokensFromCore, err := rrc.callTokens(rsvAddr)
	if err != nil {
		return tokens, err
	}
	for _, token := range tokensFromCore {
		if token.ID != "ETH" {
			tokens = append(tokens, token)
		}
	}
	return tokens, nil
}

func (rrc *ResreveRatesCrawler) getEachReserveRate(block int64, rsvAddr ethereum.Address, data *sync.Map, wg *sync.WaitGroup) error {
	defer wg.Done()
	tokens, err := rrc.getSupportedTokens(rsvAddr)
	if err != nil {
		return fmt.Errorf("cannot get supported tokens for reserve %s. Error: %s", rsvAddr.Hex(), err)
	}
	var (
		srcAddresses      []ethereum.Address    = []ethereum.Address{}
		destAddresses     []ethereum.Address    = []ethereum.Address{}
		rates             ReserveRates          = ReserveRates{}
		rsvTokenRateEntry ReserveTokenRateEntry = ReserveTokenRateEntry{}
	)
	for _, token := range tokens {
		srcAddresses = append(srcAddresses, ethereum.HexToAddress(token.Address), ethereum.HexToAddress(ethToken.Address))
		destAddresses = append(destAddresses, ethereum.HexToAddress(ethToken.Address), ethereum.HexToAddress(token.Address))
	}
	rates.Timestamp = time.Now()
	reserveRate, sanityRate, callError := rrc.wrapperContract.GetReserveRate(block, rsvAddr, srcAddresses, destAddresses)
	if callError != nil {
		return fmt.Errorf("cannot get rates for reserve %s. Error: %s", rsvAddr.Hex(), callError)
	}
	rates.BlockNumber = block
	rates.ReturnTime = time.Now()
	for index, token := range tokens {
		// the logic to get ReserveRate from conversion contract can be viewed here
		// https://developer.kyber.network/docs/ReservesGuide/#step-3-setting-token-conversion-rates-prices
		rateEntry := ReserveRateEntry{}
		rateEntry.BuyReserveRate = common.BigToFloat(reserveRate[index*2+1], ethToken.Decimals)
		rateEntry.BuySanityRate = common.BigToFloat(sanityRate[index*2+1], ethToken.Decimals)
		rateEntry.SellReserveRate = common.BigToFloat(reserveRate[index*2], ethToken.Decimals)
		rateEntry.SellSanityRate = common.BigToFloat(sanityRate[index*2], ethToken.Decimals)
		rsvTokenRateEntry[fmt.Sprintf("ETH-%s", token.ID)] = rateEntry
	}
	rates.Data = rsvTokenRateEntry
	data.Store(rsvAddr, rates)
	return nil
}

// GetReserveRates returns the map[ReserveAddress]ReserveRates at the given block number.
// It will only return rates from the set of addresses within its definition.
func (rrc *ResreveRatesCrawler) GetReserveRates(block int64) (map[string]ReserveRates, error) {
	var (
		result map[string]ReserveRates = make(map[string]ReserveRates)
		data   sync.Map                = sync.Map{}
		wg     sync.WaitGroup          = sync.WaitGroup{}
		errs   chan error              = make(chan error, len(rrc.Addresses))
	)

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
	var err error
	data.Range(func(key, value interface{}) bool {
		reserveAddr, ok := key.(ethereum.Address)
		//if there is conversion error, continue to next key,val
		if !ok {
			err = fmt.Errorf("key (%v) cannot be asserted to ethereum.Address", key)
			return false
		}
		rates, ok := value.(ReserveRates)
		if !ok {
			err = fmt.Errorf("value (%v) cannot be asserted to reserveRates", value)
			return true
		}
		result[reserveAddr.Hex()] = rates
		return true
	})
	return result, err
}
