package provider

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/tokenprice/provider/coinbase"
)

const (
	// Coinbase coinbase provider
	Coinbase = "coinbase"
)

// Provider provide token rate
type Provider interface {
	Rate(token, currency string, timestamp time.Time) (float64, error)
	Name() string
}

// NewProvider return provider interface
func NewProvider(provider string) (Provider, error) {
	switch provider {
	case Coinbase:
		return coinbase.New(), nil
	default:
		return nil, fmt.Errorf("invalide provider provider=%s", provider)
	}
}
