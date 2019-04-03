package main

import (
	"fmt"
	"log"
	"os"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/gateway/http"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

const (
	writeAccessKeyFlag         = "write-access-key"
	writeSecretKeyFlag         = "write-secret-key"
	readAccessKeyFlag          = "read-access-key"
	readSecretKeyFlag          = "read-secret-key"
	cexTradeAPIURLFlag         = "cex-trade-url"
	reserveAddressesAPIURLFlag = "reserve-addresses-url"
	cexWithdrawalURLFlag       = "cex-withdrawal-url"
	reserveTokenURLFlag        = "reserve-token-url"
	reserveTransactionURLFlag  = "reserve-transaction-url"
	erc20APIURLFlag            = "erc20-api-url"
)

var (
	defaultCexTradeAPIValue           = fmt.Sprintf("http://127.0.0.1:%d", httputil.AccountingCEXTradesPort)
	defaultReserveAddressAPIValue     = fmt.Sprintf("http://127.0.0.1:%d", httputil.AccountingReserveAddressPort)
	defaultCexWithdrawalAPIValue      = fmt.Sprintf("http://127.0.0.1:%d", httputil.AccountingCexWithdrawalPort)
	defaultReserveTokenAPIValue       = fmt.Sprintf("http://127.0.0.1:%d", httputil.AccountingListedTokenPort)
	defaultReserveTransactionAPIValue = fmt.Sprintf("http://127.0.0.1:%d", httputil.AccountingTransactionsPort)
	defaultERC20APIValue              = fmt.Sprintf("http://127.0.0.1:%d", httputil.AccountingWalletErc20Port)
)

func main() {
	app := libapp.NewApp()
	app.Name = "gateway"
	app.Usage = "Accounting API Gateway"
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
			Name:   cexTradeAPIURLFlag,
			Usage:  "cex trade api url",
			Value:  defaultCexTradeAPIValue,
			EnvVar: "CEX_TRADE_URL",
		},
		cli.StringFlag{
			Name:   reserveAddressesAPIURLFlag,
			Usage:  "reserve addresses api url",
			Value:  defaultReserveAddressAPIValue,
			EnvVar: "RESERVE_ADDRESS_URL",
		},
		cli.StringFlag{
			Name:   cexWithdrawalURLFlag,
			Usage:  "cex withdrawal api url",
			Value:  defaultReserveAddressAPIValue,
			EnvVar: "CEX_WITHDRAWAL_URL",
		},
		cli.StringFlag{
			Name:   reserveTokenURLFlag,
			Usage:  "reserve token api url",
			Value:  defaultReserveAddressAPIValue,
			EnvVar: "RESERVE_TOKEN_URL",
		},
		cli.StringFlag{
			Name:   reserveTransactionURLFlag,
			Usage:  "reserve transaction api url",
			Value:  defaultReserveAddressAPIValue,
			EnvVar: "RESERVE_TRANSACTION_URL",
		},
		cli.StringFlag{
			Name:   erc20APIURLFlag,
			Usage:  "erc20 api url",
			Value:  defaultERC20APIValue,
			EnvVar: "ERC20_API_URL",
		},
	)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.GatewayPort)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer libapp.NewFlusher(logger)()

	err = validation.Validate(c.String(cexTradeAPIURLFlag),
		validation.Required,
		is.URL)
	if err != nil {
		return fmt.Errorf("invalid cex trade API URL: %s", c.String(cexTradeAPIURLFlag))
	}

	err = validation.Validate(c.String(reserveAddressesAPIURLFlag),
		validation.Required,
		is.URL)
	if err != nil {
		return fmt.Errorf("invalid reserve address API URL: %s", c.String(reserveAddressesAPIURLFlag))
	}

	err = validation.Validate(c.String(cexWithdrawalURLFlag),
		validation.Required,
		is.URL)
	if err != nil {
		return fmt.Errorf("invalid cex withdrawal API URL: %s", c.String(cexWithdrawalURLFlag))
	}

	err = validation.Validate(c.String(reserveTokenURLFlag),
		validation.Required,
		is.URL)
	if err != nil {
		return fmt.Errorf("invalid reserve token API URL: %s", c.String(reserveTokenURLFlag))
	}

	err = validation.Validate(c.String(reserveTransactionURLFlag),
		validation.Required,
		is.URL)
	if err != nil {
		return fmt.Errorf("invalid reserve transaction API URL: %s", c.String(reserveTransactionURLFlag))
	}

	err = validation.Validate(c.String(erc20APIURLFlag),
		validation.Required,
		is.URL)
	if err != nil {
		return fmt.Errorf("invalid erc20 transaction API URL: %s", c.String(erc20APIURLFlag))
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
		auth,
		perm,
		logger,
		http.WithCexTradesURL(c.String(cexTradeAPIURLFlag)),
		http.WithResreveAddressesURL(c.String(reserveAddressesAPIURLFlag)),
		http.WithCexWithdrawalURL(c.String(cexWithdrawalURLFlag)),
		http.WithReserveTokenURL(c.String(reserveTokenURLFlag)),
		http.WithReserveTransactionURL(c.String(reserveTransactionURLFlag)),
		http.WithERC20APIURL(c.String(erc20APIURLFlag)),
	)
	if err != nil {
		return err
	}
	return svr.Start()
}
