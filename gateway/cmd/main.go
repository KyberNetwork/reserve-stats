package main

import (
	"fmt"
	"github.com/KyberNetwork/reserve-stats/gateway/http"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"log"
	"os"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
)

const (
	tradeLogsAPIURLFlag = "trade-logs-url"
)

var (
	defaultTradeLogsAPIURLValue = fmt.Sprintf("http://127.0.0.1:%d", httputil.TradeLogsPort)
)

func main() {
	app := libapp.NewApp()
	app.Name = "Gateway"
	app.Usage = "Reserve Stats API Gateway"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   tradeLogsAPIURLFlag,
			Usage:  "Trade Logs API URL",
			Value:  defaultTradeLogsAPIURLValue,
			EnvVar: "TRADE_LOGS_API_URL",
		},
		// TODO: add flag for reserve-rates-api
		// TODO: add flag for users
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

	svr, err := http.NewServer(httputil.NewHTTPAddressFromContext(c), c.String(tradeLogsAPIURLFlag))
	if err != nil {
		return err
	}
	return svr.Start()
}
