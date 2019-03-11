package blockchain

import (
	"context"
	"errors"
	"fmt"
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

var cachedName = map[common.Address]string{
	common.HexToAddress("0x86Fa049857E0209aa7D9e616F7eb3b3B78ECfdb0"): "EOS Token", // special case for EOS cos it does not return name in name function
}

// TokenInfoGetter is a helper to get token info
type TokenInfoGetter struct {
	ethClient    bind.ContractBackend // eth client
	cachedSymbol sync.Map
	cachedName   sync.Map
}

// TokenSymbolOption is the option to configure TokenSymbol constructor.
type TokenSymbolOption func(*TokenInfoGetter)

// TokenSymbolWithSymbols configures TokenSymbol constructor to use a predefined cached symbol mapping.
func TokenSymbolWithSymbols(symbols map[common.Address]string) TokenSymbolOption {
	return func(ts *TokenInfoGetter) {
		for k, v := range symbols {
			ts.cachedSymbol.Store(k, v)
		}
	}
}

// NewTokenSymbol returns a new TokenSymbol instance.
func NewTokenSymbol(client bind.ContractBackend, options ...TokenSymbolOption) *TokenInfoGetter {
	ts := &TokenInfoGetter{
		ethClient: client,
	}

	for _, option := range options {
		option(ts)
	}

	for k, v := range cachedName {
		ts.cachedName.Store(k, v)
	}

	return ts
}

// NewTokenInfoGetterFromContext return new instance of TokenInfoGetter
func NewTokenInfoGetterFromContext(c *cli.Context) (*TokenInfoGetter, error) {
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
func (t *TokenInfoGetter) Symbol(address common.Address) (string, error) {
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

//Name return name of token
func (t *TokenInfoGetter) Name(address common.Address) (string, error) {
	var (
		name string
		err  error
	)
	if val, ok := t.cachedName.Load(address); ok {
		if name, ok = val.(string); !ok {
			return "", errors.New("invalid token name stored in cached symbol")
		}
		return name, nil
	}

	name, err = getName1(address, t.ethClient)
	if err != nil {
		err1 := err
		name, err = getName2(address, t.ethClient)
		if err != nil {
			// combine 2 errors if we cannot get name
			err = fmt.Errorf("%v + %v", err1, err)
		} else if name == "" {
			// if we cannot get name then return first error
			err = err1
		}
	}

	if err == nil {
		t.cachedName.Store(address, name)
	}

	return name, err
}

func getName1(address common.Address, ethClient bind.ContractBackend) (string, error) {
	tokenContract, err := contracts.NewERC20(address, ethClient)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	callOpts := &bind.CallOpts{Context: ctx}
	return tokenContract.Name(callOpts)
}

func getName2(address common.Address, ethClient bind.ContractBackend) (string, error) {
	tokenContractType2, err := contracts.NewERC20Type2(address, ethClient)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	callOpts := &bind.CallOpts{Context: ctx}
	name2, err := tokenContractType2.Name(callOpts)
	return bytes32ToString(name2), nil
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
