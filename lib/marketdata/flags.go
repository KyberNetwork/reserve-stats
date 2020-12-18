package marketdata

import "github.com/urfave/cli"

const (
	marketDataBaseURLFlag    = "market-data-base-url"
	defaultMarketDataBaseURL = "https://staging-market-data.ksntats.com"
)

// NewMarketDataFlags ...
func NewMarketDataFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   marketDataBaseURLFlag,
			Usage:  "base url for market data client",
			EnvVar: "MARKET_DATA_BASE_URL",
			Value:  defaultMarketDataBaseURL,
		},
	}
}

// GetMarketDataBaseURLFromContext ...
func GetMarketDataBaseURLFromContext(c *cli.Context) string {
	return c.String(marketDataBaseURLFlag)
}
