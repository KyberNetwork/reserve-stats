package main

import (
	"fmt"
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/gateway/http"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
)

const (
	tradeLogsAPIURLFlag    = "trade-logs-url"
	reserveRatesAPIURLFlag = "reserve-rate-url"
	userAPIURLFlag         = "user-url"
	priceAnalyticURLFlag   = "price-analytic-url"
	writeAccessKeyFlag     = "write-access-key"
	writeSecretKeyFlag     = "write-secret-key"
	readAccessKeyFlag      = "read-access-key"
	readSecretKeyFlag      = "read-secret-key"
)

var (
	defaultTradeLogsAPIURLValue  = fmt.Sprintf("http://127.0.0.1:%d", httputil.TradeLogsPort)
	defaultReserveRatesAPIValue  = fmt.Sprintf("http://127.0.0.1:%d", httputil.ReserveRatesPort)
	defaultUserAPIValue          = fmt.Sprintf("http://127.0.0.1:%d", httputil.UsersPort)
	defaultPriceAnalyticAPIValue = fmt.Sprintf("http://127.0.0.1:%d", httputil.PriceAnalytic)
)

func main() {
	app := libapp.NewApp()
	app.Name = "Gateway"
	app.Usage = "Reserve Stats API Gateway"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   writeAccessKeyFlag,
			Usage:  "key for access POST/GET api",
			EnvVar: "WRITE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   writeSecretKeyFlag,
			Usage:  "seceret key for POST/GET api",
			EnvVar: "WRITE_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   readAccessKeyFlag,
			Usage:  "key for access GET api",
			EnvVar: "READ_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   readSecretKeyFlag,
			Usage:  "seceret key for GET api",
			EnvVar: "READ_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   tradeLogsAPIURLFlag,
			Usage:  "Trade Logs API URL",
			Value:  defaultTradeLogsAPIURLValue,
			EnvVar: "TRADE_LOGS_API_URL",
		},
		cli.StringFlag{
			Name:   reserveRatesAPIURLFlag,
			Usage:  "Reserve Rates API URL",
			Value:  defaultReserveRatesAPIValue,
			EnvVar: "RESERVE_RATES_API_URL",
		},
		cli.StringFlag{
			Name:   userAPIURLFlag,
			Usage:  "User API URL",
			Value:  defaultUserAPIValue,
			EnvVar: "USER_API_URL",
		},
		cli.StringFlag{
			Name:   priceAnalyticURLFlag,
			Usage:  "Price analytic API URL",
			Value:  defaultPriceAnalyticAPIValue,
			EnvVar: "PRICE_ANALYTIC_API_URL",
		},
	)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.GatewayPort)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	err := validation.Validate(c.String(tradeLogsAPIURLFlag),
		validation.Required,
		is.URL)
	if err != nil {
		return fmt.Errorf("invalid trades log API URL: %s", c.String(tradeLogsAPIURLFlag))
	}

	err = validation.Validate(c.String(reserveRatesAPIURLFlag),
		validation.Required,
		is.URL)
	if err != nil {
		return fmt.Errorf("invalid reserve rates API URL: %s", c.String(reserveRatesAPIURLFlag))
	}

	err = validation.Validate(c.String(userAPIURLFlag),
		validation.Required,
		is.URL)
	if err != nil {
		return fmt.Errorf("invalid user API URL: %s", c.String(userAPIURLFlag))
	}

	err = validation.Validate(c.String(priceAnalyticURLFlag),
		validation.Required,
		is.URL)
	if err != nil {
		return fmt.Errorf("invalid price analytic API URL: %s", c.String(priceAnalyticURLFlag))
	}

	if err := validation.Validate(c.String(writeAccessKeyFlag), validation.Required); err != nil {
		return fmt.Errorf("access key error: %s", err.Error())
	}

	if err := validation.Validate(c.String(writeSecretKeyFlag), validation.Required); err != nil {
		return fmt.Errorf("secret key error: %s", err.Error())
	}
	auth, err := http.NewAuthenticator(c.String(readAccessKeyFlag), c.String(readSecretKeyFlag),
		c.String(writeAccessKeyFlag), c.String(writeSecretKeyFlag),
	)
	if err != nil {
		return fmt.Errorf("authentication object creation error: %s", err)
	}
	perm, err := http.NewPermissioner(c.String(readAccessKeyFlag), c.String(writeAccessKeyFlag))
	if err != nil {
		return fmt.Errorf("permission object creation error: %s", err)
	}
	svr, err := http.NewServer(httputil.NewHTTPAddressFromContext(c),
		c.String(tradeLogsAPIURLFlag),
		c.String(reserveRatesAPIURLFlag),
		c.String(userAPIURLFlag),
		c.String(priceAnalyticURLFlag),
		auth,
		perm,
	)
	if err != nil {
		return err
	}
	return svr.Start()
}
