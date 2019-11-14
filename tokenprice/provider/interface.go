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

// PriceProvider provide token price
type PriceProvider interface {
	Price(token, currency string, timestamp time.Time) (float64, error)
	Name() string
}

// NewPriceProvider return provider interface
func NewPriceProvider(provider string) (PriceProvider, error) {
	switch provider {
	case Coinbase:
		return coinbase.New(), nil
	default:
		return nil, fmt.Errorf("invalide provider provider=%s", provider)
	}
}
