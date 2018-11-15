package crawler

import (
	"fmt"
	"math/big"

	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (rrc *ResreveRatesCrawler) getSupportedTokens(rsvAddr common.Address, block uint64) ([]core.Token, error) {
	var (
		logger = rrc.sugar.With(
			"func", "reserverates/crawler/ResreveRatesCrawler.getSupportedTokens",
			"rsv_addr", rsvAddr.Hex(),
			"block_number", block,
		)
		callOpts = &bind.CallOpts{BlockNumber: big.NewInt(0).SetUint64(block)}
		results  []core.Token
	)
	reserveContract, err := contracts.NewReserve(rsvAddr, rrc.client)
	if err != nil {
		return nil, err
	}

	conversionRatesAddr, err := reserveContract.ConversionRatesContract(callOpts)
	if err != nil {
		return nil, err
	}

	conversionRatesContract, err := contracts.NewConversionRates(conversionRatesAddr, rrc.client)
	if err != nil {
		return nil, err
	}

	listedTokens, err := conversionRatesContract.GetListedTokens(callOpts)
	if err != nil {
		return nil, err
	}

	configuredTokens, err := rrc.coreClient.Tokens()
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
