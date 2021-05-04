package blockchain

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
)

var cachedSymbols = map[deployment.Deployment]map[common.Address]string{
	deployment.Production: {
		ETHAddr:  "ETH",
		USDTAddr: "USDT",
	},
	deployment.Staging: {
		ETHAddr:  "ETH",
		USDTAddr: "USDT",
	},
}

var cachedName = map[common.Address]string{
	common.HexToAddress("0x86Fa049857E0209aa7D9e616F7eb3b3B78ECfdb0"): "EOS Token", // special case for EOS cos it does not return name in name function
}

var cachedDecimals = map[common.Address]uint8{
	common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"): 18,
}

// TokenInfoGetter is a helper to get token info
type TokenInfoGetter struct {
	ethClient      bind.ContractBackend // eth client
	tokenStorage   tokenStorageInterface
	cachedSymbol   sync.Map
	cachedName     sync.Map
	cachedDecimals sync.Map
}

// TokenSymbolOption is the option to configure TokenSymbol constructor.
type TokenSymbolOption func(*TokenInfoGetter)

type getNameOrSymbolFunc func(common.Address, bind.ContractBackend) (string, error)

// TokenSymbolWithSymbols configures TokenSymbol constructor to use a predefined cached symbol mapping.
func TokenSymbolWithSymbols(symbols map[common.Address]string) TokenSymbolOption {
	return func(ts *TokenInfoGetter) {
		for k, v := range symbols {
			ts.cachedSymbol.Store(k, v)
		}
	}
}

// NewTokenSymbol returns a new TokenSymbol instance.
func NewTokenSymbol(client bind.ContractBackend, tokenStorage tokenStorageInterface, options ...TokenSymbolOption) *TokenInfoGetter {
	ts := &TokenInfoGetter{
		ethClient:    client,
		tokenStorage: tokenStorage,
	}

	for _, option := range options {
		option(ts)
	}

	for k, v := range cachedName {
		ts.cachedName.Store(k, v)
	}

	for k, v := range cachedDecimals {
		ts.cachedDecimals.Store(k, v)
	}

	return ts
}

// NewTokenInfoGetterFromContext return new instance of TokenInfoGetter
func NewTokenInfoGetterFromContext(c *cli.Context, tokenStorage tokenStorageInterface) (*TokenInfoGetter, error) {
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

	return NewTokenSymbol(client, tokenStorage, options...), nil
}

func (t *TokenInfoGetter) getNameOrSymbol(fns []getNameOrSymbolFunc, address common.Address) (string, error) {
	var (
		err    error
		result string
	)
	for _, fn := range fns {
		if result, err = fn(address, t.ethClient); err != nil {
			if strings.Contains(err.Error(), "abi: cannot marshal") { // only ignore error when can not unpack symbol to string
				continue
			}
			return result, err
		}
		break
	}
	return result, err
}

// Symbol return symbol of a token
func (t *TokenInfoGetter) Symbol(address common.Address) (string, error) {
	var (
		symbol string
		err    error
	)

	// get symbol from cached
	if val, ok := t.cachedSymbol.Load(address); ok {
		if symbol, ok = val.(string); !ok {
			return "", errors.New("invalid value stored in cached symbol")
		}
		return symbol, nil
	}

	if t.tokenStorage != nil {
		// get symbol from storage
		symbol, err = t.tokenStorage.GetTokenSymbol(address.Hex())
		if err == nil && symbol != "" {
			// save to cached
			t.cachedSymbol.Store(address, symbol)
			return symbol, nil
		}
		// just logging
		// TODO: using zap log with sentry
		log.Printf("cannot get token symbol from storage: %s", err)
	}

	// get symbol from blockchain
	symbol, err = t.getNameOrSymbol(getSymbolFns, address)
	if err != nil {
		return symbol, err
	}
	symbol = strings.ToUpper(symbol)
	// save to cached
	t.cachedSymbol.Store(address, symbol)

	if t.tokenStorage != nil {
		// save to persistent storage
		err = t.tokenStorage.UpdateTokens([]string{address.Hex()}, []string{symbol})
		if err != nil {
			// TODO: using zap log with sentry
			log.Printf("cannot update token symbol to persistent storage")
		}
	}

	return symbol, nil
}

//Decimals return token decimals
func (t *TokenInfoGetter) Decimals(address common.Address) (uint8, error) {
	var (
		decimals uint8
	)
	if val, ok := t.cachedDecimals.Load(address); ok {
		if decimals, ok = val.(uint8); !ok {
			return 0, errors.New("invalid value stored in cached decimals")
		}
		return decimals, nil
	}
	tokenContract, err := contracts.NewERC20(address, t.ethClient)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	decimals, err = tokenContract.Decimals(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, err
	}
	t.cachedDecimals.Store(address, decimals)
	return decimals, nil
}

var getNameFns = []getNameOrSymbolFunc{
	getName1,
	getName2,
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

	name, err = t.getNameOrSymbol(getNameFns, address)
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
	if err != nil {
		return "", err
	}
	return bytes32ToString(name2), nil
}

var getSymbolFns = []getNameOrSymbolFunc{
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
		i++
	}
	return string(input[:i])
}
