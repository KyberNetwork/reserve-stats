package binance

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/urfave/cli"
)

const (
	binanceRequestPerSecond       = "binance-requests-per-second"
	binanceClientValidationFlag   = "binance-client-validation"
	binanceAccountsConfigFileFlag = "binance-account-config-file"
)

//NewCliFlags return cli flags to configure cex-trade client
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.Float64Flag{
			Name:   binanceRequestPerSecond,
			Usage:  "binance request limit per second, default to 20 which etherscan's normal rate limit",
			EnvVar: "BINANCE_REQUESTS_PER_SECOND",
			Value:  10,
		},
		cli.BoolTFlag{
			Name:   binanceClientValidationFlag,
			Usage:  "if set to true, the client is validate by calling GetAccounts with its API key",
			EnvVar: "BINANCE_CLIENT_VALIDATION",
		},
		cli.StringFlag{
			Name:   binanceAccountsConfigFileFlag,
			Usage:  "config file",
			EnvVar: "BINANCE_ACCOUNT_CONFIG_FILE",
		},
	}
}

// AccountsFromContext get accounts from file config
func AccountsFromContext(c *cli.Context) ([]common.Account, error) {
	var accounts []common.Account
	configFile := c.String(binanceAccountsConfigFileFlag)
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return accounts, err
	}
	err = json.Unmarshal(data, &accounts)
	return accounts, err
}

// ClientOptionFromContext return options for binance client
func ClientOptionFromContext(c *cli.Context) ([]Option, error) {
	var (
		options []Option
	)
	rps := c.Float64(binanceRequestPerSecond)
	if rps <= 0 {
		return nil, errors.New("rate limit must be greater than 0")
	}

	options = append(options, WithRateLimiter(NewRateLimiter(rps)))
	if validateRequire := c.BoolT(binanceClientValidationFlag); validateRequire {
		options = append(options, WithValidation())
	}
	return options, nil
}
