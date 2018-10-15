package main

import (
	"context"
	"log"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/crawler"
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/storage/influx"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	addressesFlag = "addresses"
	blockFlag     = "block"
	dbURLFlag     = "dbURL"
	dbUNameFlag   = "dbUname"
	dbPwdFlag     = "dbPwd"
	defaultDBURL  = "http://localhost:8086/"
)

func newRateStorage(c *cli.Context) (*influx.RateStorage, error) {
	url := c.GlobalString(dbURLFlag)
	uname := c.GlobalString(dbUNameFlag)
	pwd := c.GlobalString(dbPwdFlag)
	return influx.NewRateInfluxDBStorage(url, uname, pwd)
}

func newReserveCrawlerCli() *cli.App {
	app := libapp.NewApp()
	app.Name = "reserve-rates-crawler"
	app.Usage = "get the rates of all configured reserves at a certain block"
	var block uint64
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
		libapp.NewEthereumNodeFlags(),
		cli.StringFlag{
			Name:  dbURLFlag,
			Value: defaultDBURL,
			Usage: "url to InfluxDB server",
		},
		cli.StringFlag{
			Name:   dbUNameFlag,
			Usage:  "userName for InfluxDB server",
			EnvVar: "INFLUX_UNAME",
			Value:  "",
		}, cli.StringFlag{
			Name:   dbPwdFlag,
			Usage:  "url to InfluxDB server",
			EnvVar: "INFLUX_PWD",
			Value:  "",
		},
	)
	app.Flags = append(app.Flags, core.NewCliFlags()...)
	app.Action = func(c *cli.Context) error {
		addrs := c.StringSlice(addressesFlag)
		client, err := libapp.NewEthereumClientFromFlag(c)
		if err != nil {
			return err
		}
		logger, err := libapp.NewLogger(c)
		if err != nil {
			return err
		}
		blockTimeResolver, err := blockchain.NewBlockTimeResolver(logger.Sugar(), client)
		if err != nil {
			return err
		}
		coreClient, err := core.NewClientFromContext(logger.Sugar(), c)
		if err != nil {
			return err
		}
		rateStorage, err := newRateStorage(c)
		if err != nil {
			return err
		}
		reserveRateCrawler, err := crawler.NewReserveRatesCrawler(addrs, client, coreClient, logger.Sugar(), blockTimeResolver, rateStorage)
		if err != nil {
			return err
		}

		if block == 0 {
			currentBlock, err := client.BlockByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			block = currentBlock.Number().Uint64()
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
