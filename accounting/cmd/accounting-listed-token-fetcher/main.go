package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"reflect"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/listed-tokens/fetcher"
	"github.com/KyberNetwork/reserve-stats/accounting/listed-tokens/storage"
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
		cli.StringFlag{
			Name:   reserveAddressFlag,
			EnvVar: "RESERVE_ADDRESS",
			Usage:  "reserve address to get listed token",
		},
	)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())
	app.Flags = append(app.Flags, etherscan.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultDB)...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		block       *big.Int
		reserveAddr ethereum.Address
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
	if c.String(reserveAddressFlag) == "" {
		return fmt.Errorf("reserve address is required")
	}
	reserveAddrStr := c.String(reserveAddressFlag)
	sugar.Debug(reserveAddrStr)
	reserveAddr = ethereum.HexToAddress(reserveAddrStr)

	tokenSymbol, err := blockchain.NewTokenInfoGetterFromContext(c)
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
	listedTokenStorage, err := storage.NewDB(sugar, db, common.ListedTokenTable)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := listedTokenStorage.Close(); cErr != nil {
			sugar.Errorw("Close database error", "error", cErr)
		}
	}()

	f := fetcher.NewListedTokenFetcher(ethClient, resolv, sugar)
	listedTokens, err := f.GetListedToken(block, reserveAddr, tokenSymbol)
	if err != nil {
		return err
	}

	storedListedToken, err := listedTokenStorage.GetTokens()
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(storedListedToken, listedTokens) {
		if err = listedTokenStorage.CreateOrUpdate(listedTokens, block); err != nil {
			return err
		}
	}

	return listedTokenStorage.Close()
}
