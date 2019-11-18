package provider

import (
	"fmt"
	"time"

	"github.com/urfave/cli"

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
	Wait()
}

// NewPriceProvider return provider interface
func NewPriceProvider(c *cli.Context, provider string) (PriceProvider, error) {
	switch provider {
	case Coinbase:
		return coinbase.NewCoinBaseFromContext(c), nil
	case Coingecko:
		return coingecko.NewCoinGeckoFromContext(c), nil
	default:
		return nil, fmt.Errorf("invalide provider provider=%s", provider)
	}
}

// AllProvider return all provider interface
func AllProvider(c *cli.Context) []PriceProvider {
	return []PriceProvider{
		coinbase.NewCoinBaseFromContext(c),
		coingecko.NewCoinGeckoFromContext(c),
	}
}
