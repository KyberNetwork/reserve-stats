package main

import (
	"log"
	"math/big"
	"os"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	blockFlag           = "block"
	reserveAddresesFlag = "reserve-address"
)

func main() {
	app := libapp.NewApp()
	app.Name = "accounting-listed-token-fetcher"
	app.Usage = "get listed token for provided block and reserve address"
	app.Action = run
	app.Flags = append(app.Flags,
		cli.StringSliceFlag{
			Name:   blockFlag,
			EnvVar: "BLOCK",
			Usage:  "block to get listed token",
		},
		cli.StringSliceFlag{
			Name:   reserveAddresesFlag,
			EnvVar: "RESERVE_ADDRESS",
			Usage:  "reserve address to get listed token",
		},
	)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		block       *big.Int
		reserveAddr common.Address
	)
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	if c.String(blockFlag) == "" {
		sugar.Info("no block number provided, get listed token from latest block")
	} else {
		block, err = libapp.ParseBigIntFlag(c, blockFlag)
		if err != nil {
			return err
		}
	}
	ethClient, err := blockchain.NewEthereumClientFromFlag(c)
	if err != nil {
		return err
	}
	return getListedToken(ethClient, block, reserveAddr, sugar)
}

func getListedToken(ethClient *ethclient.Client, block *big.Int, reserveAddr common.Address, sugar *zap.SugaredLogger) error {
	// step 1: get conversionRatesContract address

	// step 2: get listedTokens from conversionRatesContract

	// step 3: use api etherscan to get first transaction timstamp

	return nil
}
