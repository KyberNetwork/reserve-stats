package blockchain

import "time"

// MockBlockTimeResolve is the mock implementation of blockTimeResolver for testing
type MockBlockTimeResolve struct {
	ts time.Time
}

// NewMockBlockTimeResolve creates new instance of MockBlockTimeResolve
func NewMockBlockTimeResolve(ts time.Time) *MockBlockTimeResolve {
	return &MockBlockTimeResolve{ts: ts}
}

// Resolve return current time as mock result for block time
func (btr *MockBlockTimeResolve) Resolve(_ uint64) (time.Time, error) {
	return btr.ts, nil
}
