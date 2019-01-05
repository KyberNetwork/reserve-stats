package core

import (
	"fmt"
	"strings"
)

// Interface represents a client o interact with KyberNetwork core APIs.
type Interface interface {
	Tokens() ([]Token, error)
}

// LookupToken returns the token with given id from results of Tokens of given core client.
func LookupToken(client Interface, ID string) (Token, error) {
	tokens, err := client.Tokens()
	if err != nil {
		return Token{}, err
	}
	for _, token := range tokens {
		if strings.ToLower(token.ID) == strings.ToLower(ID) {
			return token, nil
		}
	}
	return Token{}, fmt.Errorf("cannot find token %s", ID)
}
