package main

import (
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/crawler"
	"github.com/KyberNetwork/reserve-stats/setting"
	cli "gopkg.in/urfave/cli.v1"
)

func startService(c *cli.Context) error {
	return nil
}

func newReserveCrawlerCli() *cli.App {
	app := cli.NewApp()
	app.Name = "reserve-rates-crawler"
	app.Usage = "get the rates of all configured reserves at a certain block"
	var block uint64
	var endpoint string
	var coreURL string
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "addresses",
			EnvVar: "RESERVE_ADDRESSES",
			Usage:  "--addresses=\"0x1111,0x222\" or set env RESERVE_ADDRESSES=\"0x1111,0x222\"",
		},
		cli.Uint64Flag{
			Name:        "block",
			Value:       0,
			Destination: &block,
		},
		cli.StringFlag{
			Name:        "endpoint",
			EnvVar:      "ENDPOINT",
			Usage:       "--endpoint=\"infura.io\" or set env ENDPOINT=\"infura.io\"",
			Destination: &endpoint,
		},
		cli.StringFlag{
			Name:        "coreURL",
			Destination: &coreURL,
			EnvVar:      "CORE_URL",
		},
	}
	app.Action = func(c *cli.Context) error {
		addrs := c.StringSlice("addresses")
		sett, err := setting.NewSettingClient(coreURL)
		if err != nil {
			panic(err)
		}
		reserveRateCrawler, err := crawler.NewReserveRatesCrawler(addrs, endpoint, sett)
		if err != nil {
			panic(err)
		}
		result := reserveRateCrawler.GetReserveRates(block)
		log.Printf("result is %v", result)
		return nil
	}
	return app
}

//reserve-rates-crawler --addresses=0xABCDEF,0xDEFGHI --block 100
func main() {
	app := newReserveCrawlerCli()
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
