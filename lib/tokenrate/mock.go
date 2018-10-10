package tokenrate

import "time"

// NewMock returns a new Mock instance.
func NewMock() *Mock {
	return &Mock{}
}

// Mock is a mock tokenrate.Provider to use in tests.
type Mock struct{}

func (m *Mock) Rate(_, _ string, _ time.Time) (float64, error) {
	return 100, nil
}

func (m *Mock) USDRate(_ time.Time) (float64, error) {
	return m.Rate("", "", time.Time{})
}
