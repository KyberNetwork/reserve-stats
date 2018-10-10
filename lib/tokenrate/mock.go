package tokenrate

import "time"

// NewMock returns a new Mock instance.
func NewMock() *Mock {
	return &Mock{}
}

// Mock is a mock tokenrate.Provider to use in tests.
type Mock struct{}

// Rate is a mock method to satisfy the tokenrate.Provider interface.
func (m *Mock) Rate(_, _ string, _ time.Time) (float64, error) {
	return 100, nil
}

// USDRate is a mock method to satisfy the tokenrate.USDRateProvider interface.
func (m *Mock) USDRate(_ time.Time) (float64, error) {
	return m.Rate("", "", time.Time{})
}
