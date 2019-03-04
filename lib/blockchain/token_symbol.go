package blockchain

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/KyberNetwork/reserve-stats/lib/deployment"

	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

var cachedSymbols = map[deployment.Deployment]map[common.Address]string{
	deployment.Production: {
		ETHAddr: "ETH",
		BQXAddr: "BQX",
		OSTAddr: "OST",
	},
	deployment.Staging: {
		ETHAddr: "ETH",
	},
}

// TokenSymbol is a helper to convert token amount from/to wei
type TokenSymbol struct {
	ethClient    bind.ContractBackend // eth client
	cachedSymbol sync.Map
}

// TokenSymbolOption is the option to configure TokenSymbol constructor.
type TokenSymbolOption func(*TokenSymbol)

// TokenSymbolWithSymbols configures TokenSymbol constructor to use a predefined cached symbol mapping.
func TokenSymbolWithSymbols(symbols map[common.Address]string) TokenSymbolOption {
	return func(ts *TokenSymbol) {
		for k, v := range symbols {
			ts.cachedSymbol.Store(k, v)
		}
	}
}

// NewTokenSymbol returns a new TokenSymbol instance.
func NewTokenSymbol(client bind.ContractBackend, options ...TokenSymbolOption) *TokenSymbol {
	ts := &TokenSymbol{
		ethClient: client,
	}

	for _, option := range options {
		option(ts)
	}

	return ts
}

// NewTokenSymbolFromContext return new instance of TokenSymbol
func NewTokenSymbolFromContext(c *cli.Context) (*TokenSymbol, error) {
	var options []TokenSymbolOption
	client, err := NewEthereumClientFromFlag(c)
	if err != nil {
		return nil, err
	}

	dpl := deployment.MustGetDeploymentFromContext(c)
	symbols, ok := cachedSymbols[dpl]
	if ok {
		options = append(options, TokenSymbolWithSymbols(symbols))
	}

	return NewTokenSymbol(client, options...), nil
}

// Symbol return symbol of a token
func (t *TokenSymbol) Symbol(address common.Address) (string, error) {
	var (
		symbol string
		err    error
	)

	if val, ok := t.cachedSymbol.Load(address); ok {
		if symbol, ok = val.(string); !ok {
			return "", errors.New("invalid value stored in cached symbol")
		}
		return symbol, nil
	}
	for _, fn := range getSymbolFns {
		if symbol, err = fn(address, t.ethClient); err != nil {
			if strings.Contains(err.Error(), "abi: cannot marshal") { // only ignore error when can not unpack symbol to string
				continue
			}
			return symbol, err
		}
		break
	}
	if err != nil {
		return symbol, err
	}
	symbol = strings.ToUpper(symbol)
	t.cachedSymbol.Store(address, symbol)
	return symbol, nil
}

var getSymbolFns = []func(common.Address, bind.ContractBackend) (string, error){
	getSymbol1,
	getSymbol2,
}

func getSymbol1(address common.Address, client bind.ContractBackend) (string, error) {
	tokenContract, err := contracts.NewERC20(address, client)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return tokenContract.Symbol(&bind.CallOpts{Context: ctx})
}

func getSymbol2(address common.Address, client bind.ContractBackend) (string, error) {
	tokenContract, err := contracts.NewERC20Type2(address, client)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	symbols, err := tokenContract.Symbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", err
	}
	return bytes32ToString(symbols), nil
}

func bytes32ToString(input [32]byte) string {
	var i = 0
	for _, b := range input {
		if b == 0 {
			break
		}
		i = i + 1
	}
	return string(input[:i])
}
