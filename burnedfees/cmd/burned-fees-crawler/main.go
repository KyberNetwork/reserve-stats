package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/burnedfees/crawler"
	"github.com/KyberNetwork/reserve-stats/burnedfees/storage"
	influxdbstorage "github.com/KyberNetwork/reserve-stats/burnedfees/storage/influxdb"
	"github.com/KyberNetwork/reserve-stats/burnedfees/storage/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
)

const (
	fromBlockFlag    = "from-block"
	defaultFromBlock = 5069586

	toBlockFlag = "to-block"

	maxBlocksFlag    = "max-blocks"
	defaultMaxBlocks = 100000

	dbEngineFlag = "db-engine"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Burned Fees Crawler"
	app.Action = run
	app.Version = "0.0.1"
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   fromBlockFlag,
			Usage:  "Fetch burned fees from block",
			EnvVar: "FROM_BLOCK",
		},
		cli.StringFlag{
			Name:   toBlockFlag,
			Usage:  "Fetch burned fees to block",
			EnvVar: "TO_BLOCK",
		},
		cli.Uint64Flag{
			Name:   maxBlocksFlag,
			Usage:  "The maximum number of block on each query",
			EnvVar: "MAX_BLOCKS",
			Value:  defaultMaxBlocks,
		},
		cli.StringFlag{
			Name:   dbEngineFlag,
			Usage:  "database engine to store burned fee",
			EnvVar: "DB_ENGINE",
			Value:  "postgres",
		},
	)

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(storage.PostgresDefaultDb)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		fromBlock *big.Int
		toBlock   *big.Int
		daemon    bool
	)

	if err := libapp.Validate(c); err != nil {
		return err
	}
	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	ethClient, err := blockchain.NewEthereumClientFromFlag(c)
	if err != nil {
		return err
	}

	burners := append(
		contracts.OldBurnerContractAddress().MustGetFromContext(c),
		contracts.BurnerContractAddress().MustGetOneFromContext(c))

	influxClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	blkTimeRsv, err := blockchain.NewBlockTimeResolver(sugar, ethClient)
	if err != nil {
		return err
	}

	amountFmt, err := blockchain.NewTokenAmountFormatter(ethClient)
	if err != nil {
		return err
	}

	dbEngine := c.String(dbEngineFlag)
	var st storage.Interface
	if dbEngine == storage.PostgresDBEngine {
		db, err := libapp.NewDBFromContext(c)
		if err != nil {
			return err
		}
		st, err = postgres.NewPostgresStorage(db, sugar, blkTimeRsv, amountFmt)
		if err != nil {
			return err
		}
	} else {
		st, err = influxdbstorage.NewBurnedFeesStorage(sugar, influxClient, blkTimeRsv, amountFmt)
		if err != nil {
			return err
		}
	}

	cr := crawler.NewBurnedFeesCrawler(sugar, ethClient, st, burners)

	if c.String(fromBlockFlag) != "" {
		fromBlock, err = libapp.ParseBigIntFlag(c, fromBlockFlag)
		if err != nil {
			return err
		}
	}

	if c.String(toBlockFlag) == "" {
		daemon = true
	} else {
		toBlock, err = libapp.ParseBigIntFlag(c, toBlockFlag)
		if err != nil {
			return err
		}
	}

	for {

		if fromBlock == nil {
			lastBlock, fErr := st.LastBlock()
			if fErr != nil {
				return fErr
			}
			if lastBlock != 0 {
				fromBlock = big.NewInt(lastBlock)
			} else {
				fromBlock = big.NewInt(defaultFromBlock)
			}
		}

		if toBlock == nil {
			currentHeader, err := ethClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			toBlock = big.NewInt(0).Add(currentHeader.Number, big.NewInt(1))
		}

		if fErr := cr.Crawl(fromBlock.Uint64(), toBlock.Uint64(), c.Uint64(maxBlocksFlag)); fErr != nil {
			return fErr
		}

		if daemon {
			delayTime := time.Minute
			sugar.Infow("waiting before fetching new trade logs",
				"sleep", delayTime.String())
			fromBlock, toBlock = nil, nil
			time.Sleep(delayTime)
		} else {
			break
		}
	}

	return nil
}
