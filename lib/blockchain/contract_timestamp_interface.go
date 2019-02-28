package blockchain

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// ErrNotAvailable is returning when Resolve found out that given address
// does not have any transaction or not a contract.
var ErrNotAvailable = errors.New("address is not a contract, or does not exist")

// ContractTimestampResolverInterface is the common interface of contract
// creation timestamp resolver.
type ContractTimestampResolver interface {
	// Resolve returns creation time of given address.
	// If contract is not an account contract, or does not have any transaction yet,
	// return ErrNotAvailable error.
	Resolve(address common.Address) (time.Time, error)
}
