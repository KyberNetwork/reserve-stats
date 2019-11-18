package provider

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/tokenprice/provider/coinbase"
	"github.com/KyberNetwork/reserve-stats/tokenprice/provider/coingecko"
)

const (
	// Coinbase coinbase provider
	Coinbase = "coinbase"
	// Coingecko coingecko provider
	Coingecko = "coingecko"
)

// PriceProvider provide token price
type PriceProvider interface {
	ETHPrice(timestamp time.Time) (float64, error)
	Name() string
}

// NewPriceProvider return provider interface
func NewPriceProvider(provider string) (PriceProvider, error) {
	switch provider {
	case Coinbase:
		return coinbase.New(), nil
	case Coingecko:
		return coingecko.New(), nil
	default:
		return nil, fmt.Errorf("invalide provider provider=%s", provider)
	}
}

// AllProvider return all provider interface
func AllProvider() []PriceProvider {
	return []PriceProvider{
		coinbase.New(),
		coingecko.New(),
	}
}
