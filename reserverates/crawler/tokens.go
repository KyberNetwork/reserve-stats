package crawler

import (
	"math/big"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type supportedTokensGetter interface {
	supportedTokens(common.Address, uint64) ([]common.Address, error)
	symbol(common.Address) (string, error)
}

// coreSupportedTokens uses the configuration tokens from Kyber core.
type coreSupportedTokens struct {
	sugar       *zap.SugaredLogger
	ethClient   bind.ContractBackend
	tokenSymbol blockchain.TokenSymbolInterface
}

func newCoreSupportedTokens(sugar *zap.SugaredLogger, ethClient bind.ContractBackend) *coreSupportedTokens {
	return &coreSupportedTokens{
		sugar:       sugar,
		ethClient:   ethClient,
		tokenSymbol: blockchain.NewTokenSymbol(ethClient),
	}
}

func (cst *coreSupportedTokens) supportedTokens(rsvAddr common.Address, block uint64) ([]common.Address, error) {
	var callOpts = &bind.CallOpts{BlockNumber: big.NewInt(0).SetUint64(block)}
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
	return listedTokens, nil
}

func (cst *coreSupportedTokens) symbol(address common.Address) (string, error) {
	return cst.tokenSymbol.Symbol(address)
}
