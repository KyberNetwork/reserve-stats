package blockchain

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// MockContractTimestampResolver is the mock implementation of ContractTimestampResolver
// to use in tests.
type MockContractTimestampResolver struct {
	ts time.Time
}

// NewMockContractTimestampResolver returns a mock ContractTimestampResolver that
// always returns the given timestamp.
func NewMockContractTimestampResolver(ts time.Time) *MockContractTimestampResolver {
	return &MockContractTimestampResolver{ts: ts}
}

// Resolve always returns a fixed timestamp.
func (r *MockContractTimestampResolver) Resolve(_ common.Address) (time.Time, error) {
	return r.ts, nil
}
