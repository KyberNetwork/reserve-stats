package main

import (
	"fmt"
	"log"
	"os"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"

	fetcher "github.com/KyberNetwork/reserve-stats/accounting/wallet-erc20/fetcher"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/etherscan"
)

const (
	walletAddressesFlag = "wallet-addresses"
	fromBlockFlag       = "from-block"
	toBlockFlag         = "to-block"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Accounting Wallet ERC20 Transaction Fetcher"
	app.Usage = "Accounting Wallet ERC20 Transaction Fetcher"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags,
		cli.StringSliceFlag{
			Name:   walletAddressesFlag,
			EnvVar: "WALLET_ADDRESSES",
			Usage:  "list of wallet addresses to fetch transactions",
		},
		cli.StringFlag{
			Name:   fromBlockFlag,
			Usage:  "Fetch transactions from block",
			EnvVar: "FROM_BLOCK",
		},
		cli.StringFlag{
			Name:   toBlockFlag,
			Usage:  "Fetch transactions to block",
			EnvVar: "TO_BLOCK",
		},
	)
	app.Flags = append(app.Flags, etherscan.NewCliFlags()...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func validateAddrsList(addrs []string) error {
	if len(addrs) == 0 {
		return fmt.Errorf("no address provided")
	}

	for _, addr := range addrs {
		if !ethereum.IsHexAddress(addr) {
			return fmt.Errorf("invalid address provided: address=%s", addr)
		}
	}
	return nil
}

func run(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}

	sugar, flusher, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}

	defer flusher()

	walletAddrs := c.StringSlice(walletAddressesFlag)
	if err := validateAddrsList(walletAddrs); err != nil {
		sugar.Errorf("error in wallet addresses input %v", err)
		return err
	}

	fromBlock, err := libapp.ParseBigIntFlag(c, fromBlockFlag)
	if err != nil {
		return err
	}

	toBlock, err := libapp.ParseBigIntFlag(c, toBlockFlag)
	if err != nil {
		return err
	}

	etherscanClient, err := etherscan.NewEtherscanClientFromContext(c)
	if err != nil {
		return err
	}

	f := fetcher.NewWalletFetcher(sugar, etherscanClient)
	for _, walletAddr := range walletAddrs {
		transfers, err := f.Fetch(ethereum.HexToAddress(walletAddr), fromBlock, toBlock)
		if err != nil {
			return err
		}
		sugar.Infow("fetched ERC20 transactions",
			"wallet addr", walletAddr,
			"txs", transfers,
		)
	}
	return nil
}
