package crawler

import (
	"fmt"
	"go.uber.org/zap"
	"math/big"

	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type supportedTokensGetter interface {
	supportedTokens(common.Address, uint64) ([]core.Token, error)
}

// coreSupportedTokens uses the configuration tokens from Kyber core.
type coreSupportedTokens struct {
	sugar      *zap.SugaredLogger
	ethClient  bind.ContractBackend
	coreClient core.Interface
}

func newCoreSupportedTokens(sugar *zap.SugaredLogger, ethClient bind.ContractBackend, coreClient core.Interface) *coreSupportedTokens {
	return &coreSupportedTokens{
		sugar:      sugar,
		ethClient:  ethClient,
		coreClient: coreClient,
	}
}

func (cst *coreSupportedTokens) supportedTokens(rsvAddr common.Address, block uint64) ([]core.Token, error) {
	var (
		logger = cst.sugar.With(
			"func", "reserverates/crawler/coreSupportedTokens.supportedTokens",
			"rsv_addr", rsvAddr.Hex(),
			"block_number", block,
		)
		callOpts = &bind.CallOpts{BlockNumber: big.NewInt(0).SetUint64(block)}
		results  []core.Token
	)
	reserveContract, err := contracts.NewReserve(rsvAddr, cst.ethClient)
	if err != nil {
		return nil, err
	}

	conversionRatesAddr, err := reserveContract.ConversionRatesContract(callOpts)
	if err != nil {
		return nil, err
	}

	conversionRatesContract, err := contracts.NewConversionRates(conversionRatesAddr, cst.ethClient)
	if err != nil {
		return nil, err
	}

	listedTokens, err := conversionRatesContract.GetListedTokens(callOpts)
	if err != nil {
		return nil, err
	}

	configuredTokens, err := cst.coreClient.Tokens()
	if err != nil {
		return nil, err
	}

	configured := make(map[string]core.Token, len(configuredTokens))
	for i := range configuredTokens {
		configured[common.HexToAddress(configuredTokens[i].Address).Hex()] = configuredTokens[i]
	}

	for i := range listedTokens {
		token, ok := configured[listedTokens[i].Hex()]
		if !ok {
			return nil, fmt.Errorf("no configured token found: %s", listedTokens[i].Hex())
		}
		results = append(results, token)
	}

	logger.Debugw("listed tokens", "listed_tokens", len(results), "configured_tokens", len(configuredTokens))
	return results, nil
}
