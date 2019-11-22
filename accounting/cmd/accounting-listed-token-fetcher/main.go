package main

import (
	"context"
	"log"
	"math/big"
	"os"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/listed-tokens/fetcher"
	"github.com/KyberNetwork/reserve-stats/accounting/listed-tokens/storage"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/client"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/etherscan"
)

const (
	blockFlag          = "block"
	reserveAddressFlag = "reserve-address"
)

func main() {
	app := libapp.NewApp()
	app.Name = "accounting-listed-token-fetcher"
	app.Usage = "get listed token for provided block and reserve address"
	app.Action = run
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   blockFlag,
			EnvVar: "BLOCK",
			Usage:  "block to get listed token",
		},
		cli.StringSliceFlag{
			Name:   reserveAddressFlag,
			EnvVar: "RESERVE_ADDRESS",
			Usage:  "reserve address to get listed token",
		},
	)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())
	app.Flags = append(app.Flags, etherscan.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultListedTokenDB)...)
	app.Flags = append(app.Flags, client.NewClientFlags()...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		block         *big.Int
		addressClient client.Interface
	)
	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	ethClient, err := blockchain.NewEthereumClientFromFlag(c)
	if err != nil {
		return err
	}

	if c.String(blockFlag) == "" {
		sugar.Info("no block number provided, get listed token from latest block")
		header, err := ethClient.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
		block = header.Number
	} else {
		block, err = libapp.ParseBigIntFlag(c, blockFlag)
		if err != nil {
			return err
		}
	}

	addrs := c.StringSlice(reserveAddressFlag)
	if len(addrs) != 0 {
		sugar.Infow("using provided addresses instead of querying from accounting-reserve-addresses service")
		etherscanClient, err := etherscan.NewEtherscanClientFromContext(c)
		if err != nil {
			return err
		}
		resolver := blockchain.NewEtherscanContractTimestampResolver(sugar, etherscanClient)
		addressClient, err = client.NewFixedAddresses(addrs, resolver)
		if err != nil {
			return err
		}
	} else {
		addressClient, err = client.NewClientFromContext(c, sugar)
		if err != nil {
			return err
		}
	}

	tokenSymbol, err := blockchain.NewTokenInfoGetterFromContext(c, nil)
	if err != nil {
		return err
	}

	etherscanClient, err := etherscan.NewEtherscanClientFromContext(c)
	if err != nil {
		return err
	}

	resolv := blockchain.NewEtherscanContractTimestampResolver(sugar, etherscanClient)

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}
	listedTokenStorage, err := storage.NewDB(sugar, db)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := listedTokenStorage.Close(); cErr != nil {
			sugar.Errorw("Close database error", "error", cErr)
		}
	}()

	reserveAddrs, err := addressClient.ReserveAddresses(common.Reserve)
	if err != nil {
		return err
	}

	f := fetcher.NewListedTokenFetcher(ethClient, resolv, sugar)
	for _, addr := range reserveAddrs {
		listedTokens, err := f.GetListedToken(block, addr.Address, tokenSymbol)
		if err != nil {
			return err
		}

		if err = listedTokenStorage.CreateOrUpdate(listedTokens, block, addr.Address); err != nil {
			return err
		}
	}

	return listedTokenStorage.Close()
}
