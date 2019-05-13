package main

import (
	"context"
	"log"
	"math/big"
	"os"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/client"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/fetcher"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/storage"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/storage/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/etherscan"
)

const (
	addressesFlag = "addresses"

	fromBlockFlag              = "from-block"
	toBlockFlag                = "to-block"
	normalOffsetFlag           = "normal-offset"
	internalOffsetFlag         = "internal-offset"
	transferOffsetFlag         = "transfer-offset"
	attemptFlag                = "attempt"
	defaultNormalOffsetValue   = 500
	defaultInternalOffsetValue = 500
	defaultTransferOffsetValue = 200
	defaultAttemptValue        = 4
)

func fetchTx(
	sugar *zap.SugaredLogger,
	f fetcher.TransactionFetcher,
	s storage.ReserveTransactionStorage,
	addr common.ReserveAddress,
	fromBlock, toBlock *big.Int,
	normalOffset, internalOffset, transferOffset int) error {
	var logger = sugar.With(
		"func", "fetchTx",
		"addr", addr.Address.String(),
		"from", fromBlock.String(),
		"to", toBlock.String(),
	)

	if addr.Type == common.Reserve ||
		addr.Type == common.IntermediateOperator ||
		addr.Type == common.PricingOperator ||
		addr.Type == common.SanityOperator ||
		addr.Type == common.DepositOperator {
		logger.Infow("fetching normal transactions")
		normalTxs, err := f.NormalTx(addr.Address, fromBlock, toBlock, normalOffset)
		if err != nil {
			return err
		}

		if len(normalTxs) > 0 {
			logger.Infow("storing normal transactions to database", "transactions", len(normalTxs))
			if err = s.StoreNormalTx(normalTxs, addr.Address); err != nil {
				return err
			}
		}

		logger.Infow("fetching internal transactions")
		internalTxs, err := f.InternalTx(addr.Address, fromBlock, toBlock, internalOffset)
		if err != nil {
			return err
		}

		if len(internalTxs) > 0 {
			logger.Infow("storing internal transactions to database", "transactions", len(internalTxs))
			if err = s.StoreInternalTx(internalTxs, addr.Address); err != nil {
				return err
			}
		}
	}

	// for reserve, intermediate address type, need to fetch: normal, internal, ERC20 transactions.
	// for pricing operator, sanity operator, we need to fetch: normal, internal transactions as they are not supposed
	// to hold any ERC20 tokens.
	if addr.Type == common.Reserve || addr.Type == common.IntermediateOperator || addr.Type == common.CompanyWallet {
		logger.Infow("fetching ERC20 transactions")
		transfers, err := f.ERC20Transfer(addr.Address, fromBlock, toBlock, transferOffset)
		if err != nil {
			return err
		}
		logger.Infow("storing ERC20 transfers to database", "transfers", len(transfers))
		if len(transfers) > 0 {
			if err = s.StoreERC20Transfer(transfers, addr.Address); err != nil {
				return err
			}
		}
	}

	logger.Infow("storing last inserted block to database")
	return s.StoreLastInserted(addr.Address, toBlock)
}

func main() {
	app := libapp.NewApp()
	app.Name = "Accounting Reserve Transaction Fetcher"
	app.Usage = "Accounting Reserve Transaction Fetcher"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags,
		cli.StringSliceFlag{
			Name:   addressesFlag,
			EnvVar: "ADDRESSES",
			Usage:  "list of addresses to fetch transactions, only use in development",
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
		cli.IntFlag{
			Name:   normalOffsetFlag,
			Usage:  "Offset to get normal transactions",
			EnvVar: "NORMAL_OFFSET",
			Value:  defaultNormalOffsetValue,
		},
		cli.IntFlag{
			Name:   internalOffsetFlag,
			Usage:  "Offset to get internal transactions",
			EnvVar: "INTERNAL_OFFSET",
			Value:  defaultInternalOffsetValue,
		},
		cli.IntFlag{
			Name:   transferOffsetFlag,
			Usage:  "Offset to get erc20 transfer transactions",
			EnvVar: "TRANSFER_OFFSET",
			Value:  defaultTransferOffsetValue,
		},
		cli.IntFlag{
			Name:   attemptFlag,
			Usage:  "number of attempt to retry",
			EnvVar: "ATTEMPT",
			Value:  defaultAttemptValue,
		},
	)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultTransactionsDB)...)
	app.Flags = append(app.Flags, etherscan.NewCliFlags()...)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())
	app.Flags = append(app.Flags, client.NewClientFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		fromBlock     *big.Int
		toBlock       *big.Int
		addressClient client.Interface
		addrs         []common.ReserveAddress
	)

	if err := libapp.Validate(c); err != nil {
		return err
	}

	sugar, flusher, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}

	defer flusher()

	devAddrs := c.StringSlice(addressesFlag)
	if len(devAddrs) != 0 {
		sugar.Info("using provided addresses instead of querying from accounting-reserve-addresses service")
		etherscanClient, err := etherscan.NewEtherscanClientFromContext(c)
		if err != nil {
			return err
		}
		resolver := blockchain.NewEtherscanContractTimestampResolver(sugar, etherscanClient)
		addressClient, err = client.NewFixedAddresses(devAddrs, resolver)
		if err != nil {
			return err
		}
	} else {
		sugar.Info("no address provided, look up in address client instead")
		addressClient, err = client.NewClientFromContext(c, sugar)
		if err != nil {
			return err
		}
	}

	addrs, err = addressClient.ReserveAddresses()
	if err != nil {
		return err
	}

	if len(c.String(fromBlockFlag)) != 0 {
		fromBlock, err = libapp.ParseBigIntFlag(c, fromBlockFlag)
		if err != nil {
			return err
		}
	}
	ethClient, err := blockchain.NewEthereumClientFromFlag(c)
	if err != nil {
		return err
	}

	ethClient, err := blockchain.NewEthereumClientFromFlag(c)
	if err != nil {
		return err
	}

	if len(c.String(toBlockFlag)) == 0 {
		header, err := ethClient.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return err
		}
		toBlock = header.Number
	} else {
		toBlock, err = libapp.ParseBigIntFlag(c, toBlockFlag)
		if err != nil {
			return err
		}
	}

	normalOffset := c.Int(normalOffsetFlag)
	internalOffset := c.Int(internalOffsetFlag)
	transferOffset := c.Int(transferOffsetFlag)
	attempt := c.Int(attemptFlag)

	etherscanClient, err := etherscan.NewEtherscanClientFromContext(c)
	if err != nil {
		return err
	}

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	s, err := postgres.NewStorage(sugar, db)
	if err != nil {
		return err
	}

	f := fetcher.NewEtherscanTransactionFetcher(sugar, etherscanClient, ethClient, attempt)
	for _, addr := range addrs {
		fromBlock, toBlock, addr := fromBlock, toBlock, addr
		if err := s.StoreReserve(addr.Address, addr.Type.String()); err != nil {
			return err
		}

		lastInserted, err := s.GetLastInserted(addr.Address)
		if err != nil {
			return err
		}

		if lastInserted != nil {
			fromBlock = lastInserted
			sugar.Infow("starting from last inserted block",
				"address", addr.Address.String(),
				"last_inserted", toBlock.String(),
				"to_block", toBlock,
			)
		}

		if err = fetchTx(sugar, f, s, addr, fromBlock, toBlock, normalOffset, internalOffset, transferOffset); err != nil {
			return err
		}
	}
	return nil
}
