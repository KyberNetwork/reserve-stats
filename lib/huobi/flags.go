package huobi

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/urfave/cli"
	"golang.org/x/time/rate"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
)

const (
	huobiRequestPerSecond      = "huobi-requests-per-second"
	huobiClientValidation      = "huobi-client-validation"
	huobiAccountConfigFileFlag = "huobi-account-config-file"
)

//NewCliFlags return cli flags to configure cex-trade client
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.Float64Flag{
			Name:   huobiRequestPerSecond,
			Usage:  "huobi request limit per second, default to 8 which huobi's tested rate limit",
			EnvVar: "HUOBI_REQUESTS_PER_SECOND",
			Value:  8,
		},
		cli.BoolTFlag{
			Name:   huobiClientValidation,
			Usage:  "if set to true, the client is validate by calling GetAccounts with its API key",
			EnvVar: "HUOBI_CLIENT_VALIDATION",
		},
		cli.StringFlag{
			Name:   huobiAccountConfigFileFlag,
			Usage:  "huobi account config file",
			EnvVar: "HUOBI_ACCOUNT_CONFIG_FILE",
		},
	}
}

// AccountsFromContext get accounts from file config
func AccountsFromContext(c *cli.Context) ([]common.Account, error) {
	var accounts []common.Account
	configFile := c.String(huobiAccountConfigFileFlag)
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return accounts, err
	}
	err = json.Unmarshal(data, &accounts)
	return accounts, err
}

// ClientOptionFromContext ...
func ClientOptionFromContext(c *cli.Context) ([]Option, error) {
	var (
		options []Option
	)

	rps := c.Float64(huobiRequestPerSecond)
	if rps <= 0 {
		return nil, errors.New("request per second must be greater than 0")
	}

	options = append(options, WithRateLimiter(rate.NewLimiter(rate.Limit(rps), 1)))
	if validateRequire := c.BoolT(huobiClientValidation); validateRequire {
		options = append(options, WithValidation())
	}
	return options, nil
}
