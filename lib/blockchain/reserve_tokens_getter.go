package blockchain

import (
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"math/big"
)

// TokenInfo is information of a supported token of a reserve.
type TokenInfo struct {
	Address common.Address
	Symbol  string
}

// reserveTokenFetcher returns supported tokens information by
// calling GetListedTokens of pricing contract attached to reserve.
type ReserveTokenFetcher struct {
	sugar          *zap.SugaredLogger
	ethClient      bind.ContractBackend
	symbolResolver TokenSymbolResolver
}

// NewReserveTokenFetcher creates a new instance of ReserveTokenFetcher.
func NewReserveTokenFetcher(sugar *zap.SugaredLogger, ethClient bind.ContractBackend, symbolResolver TokenSymbolResolver) *ReserveTokenFetcher {
	return &ReserveTokenFetcher{sugar: sugar, ethClient: ethClient, symbolResolver: symbolResolver}
}

// Tokens returns list of supported token of given reserve at given block.
func (rtf *ReserveTokenFetcher) Tokens(reserve common.Address, block uint64) ([]TokenInfo, error) {
	var tokens []TokenInfo
	st, err := rtf.supportedTokens(reserve, block)
	if err != nil {
		return nil, err
	}
	for _, addr := range st {
		symbol, err := rtf.symbolResolver.Symbol(addr)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, TokenInfo{
			Address: addr,
			Symbol:  symbol,
		})
	}

	return tokens, nil
}

func (rtf *ReserveTokenFetcher) supportedTokens(rsvAddr common.Address, block uint64) ([]common.Address, error) {
	var callOpts = &bind.CallOpts{BlockNumber: big.NewInt(0).SetUint64(block)}
	reserveContract, err := contracts.NewReserve(rsvAddr, rtf.ethClient)
	if err != nil {
		return nil, err
	}

	conversionRatesAddr, err := reserveContract.ConversionRatesContract(callOpts)
	if err != nil {
		return nil, err
	}

	conversionRatesContract, err := contracts.NewConversionRates(conversionRatesAddr, rtf.ethClient)
	if err != nil {
		return nil, err
	}

	listedTokens, err := conversionRatesContract.GetListedTokens(callOpts)
	if err != nil {
		return nil, err
	}
	return listedTokens, nil
}
