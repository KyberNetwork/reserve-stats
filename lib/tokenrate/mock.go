package tokenrate

import "time"

const mockRate float64 = 100

// NewMock returns a new Mock instance.
func NewMock() *Mock {
	return &Mock{}
}

// Mock is a mock tokenrate.Provider to use in tests.
type Mock struct{}

// Rate is a mock method to satisfy the tokenrate.Provider interface.
func (m *Mock) Rate(_, _ string, _ time.Time) (float64, error) {
	return mockRate, nil
}

// USDRate is a mock method to satisfy the tokenrate.USDRateProvider interface.
func (m *Mock) USDRate(_ time.Time) (float64, error) {
	return m.Rate("", "", time.Time{})
}

// Name is a mock method to satisfy the tokenrate.USDRatePrivder interface.
func (m *Mock) Name() string {
	return "tokenRateMock"
}
