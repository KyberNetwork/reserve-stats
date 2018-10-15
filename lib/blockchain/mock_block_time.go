package blockchain

import "time"

// MockBlockTimeResolve is the mock implementation of blockTimeResolver for testing
type MockBlockTimeResolve struct {
}

// Resolve return current time as mock result for blocktime
func (btr *MockBlockTimeResolve) Resolve(_ uint64) (time.Time, error) {
	return time.Now(), nil
}
