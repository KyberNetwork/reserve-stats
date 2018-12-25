package blockchain

import (
	"context"
	"strings"
	"sync"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

// TokenSymbol is a helper to convert token amount from/to wei
type TokenSymbol struct {
	mu           *sync.RWMutex
	ethClient    bind.ContractBackend // eth client
	cachedSymbol map[common.Address]string
}

// NewTokenSymbol returns a new TokenSymbol instance.
func NewTokenSymbol(client bind.ContractBackend) *TokenSymbol {
	var cachedSymbol = make(map[common.Address]string)
	cachedSymbol[ETHAddr] = "ETH"
	cachedSymbol[BQXAddr] = "BQX"
	cachedSymbol[OSTAddr] = "OST"

	return &TokenSymbol{
		mu:           &sync.RWMutex{},
		ethClient:    client,
		cachedSymbol: cachedSymbol,
	}
}

// NewTokenSymbolFromContext return new instance of TokenSymbol
func NewTokenSymbolFromContext(c *cli.Context) (*TokenSymbol, error) {
	client, err := app.NewEthereumClientFromFlag(c)
	if err != nil {
		return nil, err
	}
	return NewTokenSymbol(client), nil
}

// Symbol return symbol of a token
func (t *TokenSymbol) Symbol(address common.Address) (string, error) {
	var (
		symbol string
		err    error
	)
	t.mu.RLock()
	if symbol, ok := t.cachedSymbol[address]; ok {
		t.mu.RUnlock()
		return symbol, nil
	}
	t.mu.RUnlock()
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
	t.mu.Lock()
	t.cachedSymbol[address] = symbol
	t.mu.Unlock()
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
