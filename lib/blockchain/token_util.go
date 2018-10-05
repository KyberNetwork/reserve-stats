package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// EthAddr is ethereum address
var EthAddr = common.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")

// KncAddr is KNC token address
var KncAddr = common.HexToAddress("0xdd974D5C2e2928deA5F71b9825b8b646686BD200")

type token struct {
	Name     string
	Address  common.Address
	Decimals int64
}

// TokenUtil allow to look up token info (name, decimals) by token address
type TokenUtil struct {
	info map[common.Address]*token
}

// NewTokenUtil return new instance of TokenUtil
func NewTokenUtil() (*TokenUtil, error) {
	var tokens []*token

	err := json.NewDecoder(bytes.NewReader([]byte(tokenData))).Decode(&tokens)
	if err != nil {
		return nil, err
	}

	info := make(map[common.Address]*token)
	for _, tk := range tokens {
		info[tk.Address] = tk
	}

	return &TokenUtil{info: info}, nil
}

// GetTokenAmount return token amount given the token address & token wei
func (ti *TokenUtil) GetTokenAmount(tokenAddr common.Address, tokenWei *big.Int) (float64, error) {
	tk, ok := ti.info[tokenAddr]
	if !ok {
		return 0, errors.New("token info not exists")
	}
	return BigToFloat(tokenWei, tk.Decimals), nil
}

// BigToFloat converts a big int to float according to its number of decimal digits
// Example:
// - BigToFloat(1100, 3) = 1.1
// - BigToFloat(1100, 2) = 11
// - BigToFloat(1100, 5) = 0.11
func BigToFloat(b *big.Int, decimal int64) float64 {
	if b == nil {
		return 0
	}
	f := new(big.Float).SetInt(b)
	power := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10), big.NewInt(decimal), nil,
	))
	res := new(big.Float).Quo(f, power)
	result, _ := res.Float64()
	return result
}
