package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// Interface represents a client o interact with KyberNetwork core APIs.
type Interface interface {
	Tokens() ([]Token, error)
	FromWei(common.Address, *big.Int) (float64, error)
	ToWei(common.Address, float64) (*big.Int, error)
}

// LookupToken returns the token with given id from results of Tokens of given core client.
func LookupToken(client Interface, ID string) (Token, error) {
	var (
		err    error
		result Token
	)
	tokens, err := client.Tokens()
	if err != nil {
		return result, err
	}
	for _, token := range tokens {
		if token.ID == ID {
			return result, nil
		}
	}
	return result, fmt.Errorf("cannot find token %s", ID)
}
