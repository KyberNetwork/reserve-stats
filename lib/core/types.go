package core

import (
	"math"
	"math/big"
)

// ETHToken is the configuration of Ethereum, will never be changed.
var ETHToken = Token{
	ID:       "ETH",
	Name:     "Ethereum",
	Address:  "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
	Decimals: 18,
}

// WETHToken is the configuration of Ethereum, will hardly be changed.
var WETHToken = Token{
	ID:       "WETH",
	Name:     "Wrapped-ETH",
	Address:  "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
	Decimals: 18,
}

// Token is a ERC20 token allowed to trade in core.
// Note: all fields below are valid, uncomment when needed.
type Token struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Decimals int64  `json:"decimals"`
	//Active               bool   `json:"active"`
	//Internal             bool   `json:"internal"`
	//LastActivationChange uint64 `json:"last_activation_change"`
}

// FromWei converts the given amount in Wei with following formatting rule.
// - Decimals: 3
//   FromWei(1100) = 1.1
// - Decimals: 2
//   FromWei(1100) = 11
// - Decimals: 5
//   FromWei(1100) = 0.11
func (t *Token) FromWei(amount *big.Int) float64 {
	if amount == nil {
		return 0
	}
	f := new(big.Float).SetInt(amount)
	power := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10), big.NewInt(t.Decimals), nil,
	))
	res := new(big.Float).Quo(f, power)
	result, _ := res.Float64()
	return result
}

// ToWei return the given human friendly number to wei unit.
func (t *Token) ToWei(amount float64) *big.Int {
	decimals := t.Decimals
	// 6 is our smallest precision,
	if decimals < 6 {
		return big.NewInt(int64(amount * math.Pow10(int(decimals))))
	}

	result := big.NewInt(int64(amount * math.Pow10(6)))
	return result.Mul(result, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimals-6), nil))
}
