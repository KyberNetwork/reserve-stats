package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"os"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/fetcher"
	"github.com/KyberNetwork/reserve-stats/accounting/wallet-erc20/storage/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/etherscan"
)

const (
	walletAddressesFlag = "wallet-addresses"
	fromBlockFlag       = "from-block"
	toBlockFlag         = "to-block"
	defaultPostGresDB   = common.DefaultDB
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
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultPostGresDB)...)
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
		sugar.Warnf("error in parsing fromblock: %v, fetch from last stored block instead", err)
		fromBlock = nil
	}

	toBlock, err := libapp.ParseBigIntFlag(c, toBlockFlag)
	if err != nil {
		sugar.Warnf("error in parsing toblock: %v, fetch to latest block Instead", err)
		toBlock = nil
	}

	etherscanClient, err := etherscan.NewEtherscanClientFromContext(c)
	if err != nil {
		return err
	}

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	wdb, err := postgres.NewDB(sugar, db)
	if err != nil {
		return err
	}

	f := fetcher.NewEtherscanTransactionFetcher(sugar, etherscanClient)

	//f := fetcher.NewWalletFetcher(sugar, etherscanClient)
	for _, walletAddr := range walletAddrs {
		if fromBlock == nil {
			lastStoredBlock, err := wdb.GetLastStoredBlock(ethereum.HexToAddress(walletAddr))
			switch err {
			case sql.ErrNoRows:
				sugar.Infof("no record for wallet %s yet. fetch from beginning", walletAddr)
			case nil:
				sugar.Infof("found record for wallet %s. fetch from block %d", lastStoredBlock+1)
				fromBlock = big.NewInt(int64(lastStoredBlock + 1))
			default:
				return err
			}
		}
		transfers, err := f.ERC20Transfer(ethereum.HexToAddress(walletAddr), fromBlock, toBlock)
		if err != nil {
			return err
		}
		if err := wdb.UpdateERC20Transfers(transfers); err != nil {
			return err
		}
	}
	return nil
}
