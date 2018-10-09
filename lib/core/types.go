package core

import "math/big"

// Token is a ERC20 token allowed to trade in core.
// Note: all fields below are valid, uncomment when needed.
type Token struct {
	//ID                   string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Decimals int64  `json:"decimals"`
	//Active               bool   `json:"active"`
	//Internal             bool   `json:"internal"`
	//LastActivationChange uint64 `json:"last_activation_change"`
}

// FormatAmount converts the given amount in Wei with following formatting rule.
// - Decimals: 3
//   FormatAmount(1100) = 1.1
// - Decimals: 2
//   FormatAmount(1100) = 11
// - Decimals: 5
//   FormatAmount(1100) = 0.11
func (t *Token) FormatAmount(amount *big.Int) float64 {
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
