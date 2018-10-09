package main

import (
	"log"
	"os"

	kyberApp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/crawler"
	"github.com/KyberNetwork/reserve-stats/setting"
	cli "github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	addressesFlag = "addresses"
	blockFlag     = "block"
	coreFlag      = "coreURL"
)

func newReserveCrawlerCli() *cli.App {
	app := kyberApp.NewApp()
	app.Name = "reserve-rates-crawler"
	app.Usage = "get the rates of all configured reserves at a certain block"
	var block uint64
	var coreURL string
	app.Flags = append(app.Flags,
		cli.StringSliceFlag{
			Name:   addressesFlag,
			EnvVar: "RESERVE_ADDRESSES",
			Usage:  "list of reserve contract addresses. Example: --addresses={\"0x1111\",\"0x222\"}",
		},
		cli.Uint64Flag{
			Name:        blockFlag,
			Value:       0,
			Usage:       "block from which rate is queried. Default value is 0, in which case the latest rate is returned",
			Destination: &block,
		},
		cli.StringFlag{
			Name:        coreFlag,
			Destination: &coreURL,
			EnvVar:      "CORE_URL",
		},
		kyberApp.NewEthereumNodeFlags(""),
	)
	app.Action = func(c *cli.Context) error {
		addrs := c.StringSlice(addressesFlag)
		sett, err := setting.NewSettingClient(coreURL)
		if err != nil {
			return err
		}
		client, err := kyberApp.NewEthereumClientFromFlag(c)
		if err != nil {
			return err
		}
		logger, err := kyberApp.NewLogger(c)
		if err != nil {
			return err
		}
		reserveRateCrawler, err := crawler.NewReserveRatesCrawler(addrs, client, sett, logger.Sugar())
		if err != nil {
			return err
		}
		result, err := reserveRateCrawler.GetReserveRates(block)
		if err != nil {
			return err
		}
		logger.Info("rate result is", zap.Reflect("rates", result))
		return nil
	}
	return app
}

//reserve-rates-crawler --addresses=0xABCDEF,0xDEFGHI --block 100
func main() {
	app := newReserveCrawlerCli()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
