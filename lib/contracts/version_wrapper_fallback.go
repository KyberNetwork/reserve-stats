package contracts

import (
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	emptyErrMsg = "abi: unmarshalling empty output"

	// KyberNetwork smart contracts return this address
	// in case it don't have one
	// e.g:
	// get sanity address from reserve, but the sanity was not set at that time
	voidAddr = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

// VersionedWrapperFallback is a wrapper around VersionedWrapper with fallback
// if getSanityRates throw an exception.
// For example, at block: 6000744, wrapper will returns an exception for BNB --> ETH,
// because sanity rate for BNB was not set at this block.
type VersionedWrapperFallback struct {
	sugar  *zap.SugaredLogger
	client bind.ContractBackend
	vw     *VersionedWrapper
}

// NewVersionedWrapperFallback creates a new instance of VersionedWrapperFallback.
func NewVersionedWrapperFallback(sugar *zap.SugaredLogger, client bind.ContractBackend) (*VersionedWrapperFallback, error) {
	vw, err := NewVersionedWrapper(client)
	if err != nil {
		return nil, err
	}
	return &VersionedWrapperFallback{sugar: sugar, client: client, vw: vw}, nil
}

func (vwf *VersionedWrapperFallback) getReserveRateFallback(block uint64, rsvAddr common.Address, srcs, dsts []common.Address) ([]*big.Int, []*big.Int, error) {
	var (
		logger = vwf.sugar.With("func", "lib/contracts/VersionedWrapperFallback.GetReserveRate",
			"block", block,
			"reserve_addr", rsvAddr.Hex(),
		)
		callOpts    = &bind.CallOpts{BlockNumber: big.NewInt(0).SetUint64(block)}
		rates       = make(map[int]*big.Int, len(srcs))
		sanityRates = make(map[int]*big.Int, len(srcs))
		m           sync.Mutex

		g           errgroup.Group
		resourcesCh = make(chan struct{}, 10) // resources limiter, thread need to acquire release resource

		ratesResults       []*big.Int
		sanityRatesResults []*big.Int
	)
	if len(srcs) != len(dsts) {
		return nil, nil, errors.New("sources, destinations length are not equal")
	}

	reserveContract, err := NewReserve(rsvAddr, vwf.client)
	if err != nil {
		logger.Errorw("failed to create reserve contract", "err", err)
		return nil, nil, err
	}

	for i := range srcs {
		var (
			i, src, dst = i, srcs[i], dsts[i]
			logger      = logger.With(
				"src", src.Hex(),
				"dst", dst.Hex(),
			)
		)

		g.Go(func() error {
			resourcesCh <- struct{}{}
			defer func() { <-resourcesCh }()
			logger.Debugw("calling reserve contract GetConversionRate")
			rate, gErr := reserveContract.GetConversionRate(
				callOpts,
				src, dst, big.NewInt(0),
				big.NewInt(0).SetUint64(block))
			if gErr != nil {
				if gErr.Error() != emptyErrMsg {
					return gErr
				}
				logger.Infow("got exception when calling reserve contract")
				gErr = nil
				rate = big.NewInt(0)
			}

			sanityRateAddr, gErr := reserveContract.SanityRatesContract(callOpts)
			if gErr != nil {
				return gErr
			}

			var sanityRate *big.Int
			if sanityRateAddr != voidAddr {
				sanityRateContract, gErr := NewSanityRates(sanityRateAddr, vwf.client)
				if gErr != nil {
					return gErr
				}

				logger = logger.With("sanity_contract", sanityRateAddr.Hex())
				logger.Debugw("calling sanity rates contract GetSanityRate")
				sanityRate, gErr = sanityRateContract.GetSanityRate(callOpts, src, dst)
				if gErr != nil {
					if gErr.Error() != emptyErrMsg {
						return gErr

					}
					logger.Infow("got exception when calling sanity rate contract")
					gErr = nil
					sanityRate = big.NewInt(0)
				}
			} else {
				logger.Infow("sanity_rate smart contract not available")
				sanityRate = big.NewInt(0)
			}

			m.Lock()
			rates[i] = rate
			sanityRates[i] = sanityRate
			logger.Debugw("got rates successfully",
				"rate", rate,
				"sanity_rate", sanityRate,
			)
			m.Unlock()
			return nil
		})
	}

	if err = g.Wait(); err != nil {
		return nil, nil, err
	}

	for i := range srcs {
		ratesResults = append(ratesResults, rates[i])
	}

	for i := range srcs {
		sanityRatesResults = append(sanityRatesResults, sanityRates[i])
	}

	return ratesResults, sanityRatesResults, nil
}

// GetReserveRate is the same as VersionedWrapper.GetReserveRate but fallback to calling each reserve contract
// directly in case of an exception happens.
func (vwf *VersionedWrapperFallback) GetReserveRate(block uint64, rsvAddr common.Address, srcs, dsts []common.Address) ([]*big.Int, []*big.Int, error) {
	var logger = vwf.sugar.With("func", "lib/contracts/VersionedWrapperFallback.GetReserveRate",
		"block", block,
		"reserve_addr", rsvAddr.Hex(),
	)
	rates, sanityRates, err := vwf.vw.GetReserveRate(block, rsvAddr, srcs, dsts)
	if err != nil {
		if err.Error() == emptyErrMsg {
			logger.Infow("wrapper contract execution failed, fallback to calling reserve directly", "err", err)
			return vwf.getReserveRateFallback(block, rsvAddr, srcs, dsts)
		}
		return nil, nil, err
	}

	return rates, sanityRates, nil
}
