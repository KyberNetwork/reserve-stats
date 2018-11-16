package crawler

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	rsvRateCommon "github.com/KyberNetwork/reserve-stats/reserverates/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// ResreveRatesCrawler contains two wrapper contracts for V1 and V2 contract,
// a set of addresses to crawl rates from and setting object to query for reserve's token settings
type ResreveRatesCrawler struct {
	sugar *zap.SugaredLogger

	wrapperContract     reserveRateGetter
	addresses           []ethereum.Address
	stg                 supportedTokensGetter
	internalReserveAddr ethereum.Address
	blkTimeRsv          blockchain.BlockTimeResolverInterface
}

// NewReserveRatesCrawler returns an instant of ReserveRatesCrawler.
func NewReserveRatesCrawler(sugar *zap.SugaredLogger, addrs []string, client *ethclient.Client, coreClient core.Interface, internalReserveAddr ethereum.Address, bl blockchain.BlockTimeResolverInterface) (*ResreveRatesCrawler, error) {
	wrpContract, err := contracts.NewVersionedWrapperFallback(sugar, client)
	if err != nil {
		return nil, err
	}

	var ethAddrs []ethereum.Address
	for _, addr := range addrs {
		ethAddrs = append(ethAddrs, ethereum.HexToAddress(addr))
	}
	return &ResreveRatesCrawler{
		sugar:               sugar,
		wrapperContract:     wrpContract,
		addresses:           ethAddrs,
		stg:                 newCoreSupportedTokens(sugar, client, coreClient),
		internalReserveAddr: internalReserveAddr,
		blkTimeRsv:          bl,
	}, nil
}

func (rrc *ResreveRatesCrawler) getEachReserveRate(block uint64, rsvAddr ethereum.Address) (*rsvRateCommon.ReserveRates, error) {
	var (
		err   error
		rates = &rsvRateCommon.ReserveRates{
			BlockNumber: block,
			Data:        make(map[string]rsvRateCommon.ReserveRateEntry),
		}
		srcAddresses  []ethereum.Address
		destAddresses []ethereum.Address
	)

	logger := rrc.sugar.With(
		"func", "reserverates/reserve-rates-crawler/ResreveRatesCrawler.getEachReserveRate",
		"block", block,
		"reserve_address", rsvAddr.Hex(),
	)
	logger.Debug("fetching reserve rates")

	if rates.Timestamp, err = rrc.blkTimeRsv.Resolve(block); err != nil {
		return nil, err
	}

	tokens, err := rrc.stg.supportedTokens(rsvAddr, block)
	if err != nil {
		if err.Error() == bind.ErrNoCode.Error() {
			logger.Infow("reserve contract does not exist")
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get supported tokens for reserve %s. Error: %s", rsvAddr.Hex(), err)
	}

	for _, token := range tokens {
		srcAddresses = append(srcAddresses, ethereum.HexToAddress(token.Address), ethereum.HexToAddress(core.ETHToken.Address))
		destAddresses = append(destAddresses, ethereum.HexToAddress(core.ETHToken.Address), ethereum.HexToAddress(token.Address))
	}

	reserveRates, sanityRates, err := rrc.wrapperContract.GetReserveRate(block, rsvAddr, srcAddresses, destAddresses)
	if err != nil {
		return nil, fmt.Errorf("cannot get rates for reserve %s. Error: %s", rsvAddr.Hex(), err)
	}

	for index, token := range tokens {
		rates.Data[fmt.Sprintf("ETH-%s", token.ID)] = rsvRateCommon.NewReserveRateEntry(reserveRates, sanityRates, index)
	}

	logger.Debug("reserve rates fetched successfully")
	return rates, err
}

// GetReserveRates returns the map[ReserveAddress]ReserveRates at the given block number.
// It will only return rates from the set of addresses within its definition.
func (rrc *ResreveRatesCrawler) GetReserveRates(block uint64) (map[string]rsvRateCommon.ReserveRates, error) {
	var (
		err    error
		g      errgroup.Group
		data   = sync.Map{}
		result = make(map[string]rsvRateCommon.ReserveRates)
	)

	logger := rrc.sugar.With(
		"func", "reserverates/reserve-rates-crawler/ResreveRatesCrawler.GetReserveRates",
		"block", block,
		"reserves", len(rrc.addresses),
	)
	logger.Debug("fetching rates for all reserves")

	for _, rsvAddr := range rrc.addresses {
		// copy to local variables to avoid race condition
		block, rsvAddr := block, rsvAddr
		g.Go(func() error {
			rates, err := rrc.getEachReserveRate(block, rsvAddr)
			if err != nil {
				return err
			}

			if rates == nil {
				logger.Info("rates is not available, skipping")
				return nil
			}

			data.Store(rsvAddr, *rates)
			return nil
		})
	}

	if err = g.Wait(); err != nil {
		return nil, err
	}

	data.Range(func(key, value interface{}) bool {
		reserveAddr, ok := key.(ethereum.Address)
		//if there is conversion error, continue to next key,val
		if !ok {
			err = fmt.Errorf("key (%v) cannot be asserted to ethereum.Address", key)
			return false
		}
		rates, ok := value.(rsvRateCommon.ReserveRates)
		if !ok {
			err = fmt.Errorf("value (%v) cannot be asserted to reserveRates", value)
			return false
		}
		result[reserveAddr.Hex()] = rates
		return true
	})
	if err != nil {
		return nil, err
	}

	return result, err
}
