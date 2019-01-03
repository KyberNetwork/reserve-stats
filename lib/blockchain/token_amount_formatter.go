package blockchain

import (
	"context"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
)

var timeout = 30 * time.Second

// TokenAmountFormatter is a helper to convert token amount from/to wei
type TokenAmountFormatter struct {
	mu             *sync.RWMutex
	ethClient      *ethclient.Client // eth client
	cachedDecimals map[common.Address]int64
	tokenAddress   tokenAddresses
}

// NewTokenAmountFormatter returns a new TokenAmountFormatter instance.
func NewTokenAmountFormatter(client *ethclient.Client, tAddr tokenAddresses) (*TokenAmountFormatter, error) {
	var cachedDecimals = make(map[common.Address]int64)
	cachedDecimals[tAddr.ETHAddr] = 18

	return &TokenAmountFormatter{
		mu:             &sync.RWMutex{},
		ethClient:      client,
		cachedDecimals: cachedDecimals,
		tokenAddress:   tAddr,
	}, nil
}

// NewToKenAmountFormatterFromContext return new instance of TokenAmountFormatter
func NewToKenAmountFormatterFromContext(c *cli.Context) (*TokenAmountFormatter, error) {
	client, err := app.NewEthereumClientFromFlag(c)
	if err != nil {
		return nil, err
	}
	tokenAddrs, err := getTokenAddressFromContext(c)
	if err != nil {
		return nil, err
	}
	return NewTokenAmountFormatter(client, tokenAddrs)
}

// FromWei formats the given amount in wei to human friendly
// number with token decimals from contract.
func (f *TokenAmountFormatter) FromWei(address common.Address, amount *big.Int) (float64, error) {
	if amount == nil {
		return 0, nil
	}
	decimals, err := f.getDecimals(address)
	if err != nil {
		return 0, err
	}
	floatAmount := new(big.Float).SetInt(amount)
	power := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10), big.NewInt(decimals), nil,
	))
	res := new(big.Float).Quo(floatAmount, power)
	result, _ := res.Float64()
	return result, nil
}

// ToWei return the given human friendly number to wei unit.
func (f *TokenAmountFormatter) ToWei(address common.Address, amount float64) (*big.Int, error) {
	decimals, err := f.getDecimals(address)
	if err != nil {
		return nil, err
	}
	// 6 is our smallest precision,
	if decimals < 6 {
		return big.NewInt(int64(amount * math.Pow10(int(decimals)))), nil
	}
	result := big.NewInt(int64(amount * math.Pow10(6)))
	return result.Mul(result, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimals-6), nil)), nil
}

func (f *TokenAmountFormatter) getDecimals(address common.Address) (int64, error) {
	f.mu.RLock()
	if decimals, ok := f.cachedDecimals[address]; ok {
		f.mu.RUnlock()
		return decimals, nil
	}
	f.mu.RUnlock()

	tokenContract, err := contracts.NewERC20(address, f.ethClient)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	decimals, err := tokenContract.Decimals(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, err
	}
	f.mu.Lock()
	f.cachedDecimals[address] = int64(decimals)
	f.mu.Unlock()
	return int64(decimals), err
}

// IsBurnable indicate if the burn fee event was emitted when
// the given token was trade on KyberNetwork
func (f *TokenAmountFormatter) IsBurnable(token common.Address) bool {
	var notBurnTokens = map[common.Address]struct{}{
		f.tokenAddress.ETHAddr:  {},
		f.tokenAddress.WETHAddr: {},
		f.tokenAddress.KCCAddr:  {},
	}

	_, notBurn := notBurnTokens[token]
	return !notBurn
}

// KNCAddr return KNC address on current deployment mode
func (f *TokenAmountFormatter) KNCAddr() common.Address {
	return f.tokenAddress.KNCAddr
}

// ETHAddr return ETH address on current deployment mode
func (f *TokenAmountFormatter) ETHAddr() common.Address {
	return f.tokenAddress.ETHAddr
}
