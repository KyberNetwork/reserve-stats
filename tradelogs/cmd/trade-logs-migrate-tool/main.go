package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

const (
	duneDataPathFlag    = "dune-data-path"
	defaultDuneDataPath = "data.json"

	maxRecordSavePerTimeFlag  = "max-record-save-per-time"
	defaultMaxElemSavePerTime = 50
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs migrate tool"
	app.Usage = "Migrate Trade Log data"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   duneDataPathFlag,
			Usage:  "path to dune data",
			Value:  defaultDuneDataPath,
			EnvVar: "DUNE_DATA_PATH",
		},
		cli.UintFlag{
			Name:   maxRecordSavePerTimeFlag,
			Usage:  "max record can save per time",
			Value:  defaultMaxElemSavePerTime,
			EnvVar: "MAX_RECORD_SAVE_PER_TIME",
		},
	)

	app.Flags = append(app.Flags, storage.NewCliFlags()...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(storage.PostgresDefaultDb)...)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		err                  error
		storageInterface     storage.Interface
		maxRecordSavePerTime = c.Uint(maxRecordSavePerTimeFlag)
	)

	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	tokenAmountFormatter, err := blockchain.NewToKenAmountFormatterFromContext(c)
	if err != nil {
		return err
	}

	storageInterface, err = storage.NewStorageInterfaceFromContext(sugar, c, tokenAmountFormatter)
	if err != nil {
		return err
	}

	duneData, err := duneDataFromFile(c.String(duneDataPathFlag))
	if err != nil {
		sugar.Errorw("cannot parse dune data from file", "err", err)
		return err
	}

	var updatedTradeLogs []common.TradeLog
	for _, r := range duneData.QueryResult.Data.Rows {
		walletAddr := ethereum.HexToAddress(addHexPrefix(r.WalletID))
		if !r.CallSuccess || blockchain.IsZeroAddress(walletAddr) {
			continue
		}
		txHash := ethereum.HexToHash(addHexPrefix(r.CallTxHash))
		tradeLogs, err := storageInterface.LoadTradeLogsByTxHash(txHash)
		if err != nil {
			sugar.Errorw("cannot get trade logs by hash", "tx", txHash.Hex())
			return err
		}
		if len(tradeLogs) == 0 {
			sugar.Debugw("no trade log for given tx hash", "tx", txHash.Hex())
			continue
		}
		for _, tl := range tradeLogs {
			if !blockchain.IsZeroAddress(tl.WalletAddress) {
				continue
			}
			var (
				isSameSrc      = tl.SrcAddress == ethereum.HexToAddress(r.Src)
				isSameDst      = tl.DestAddress == ethereum.HexToAddress(r.Dest)
				isSameReceiver = tl.ReceiverAddress == ethereum.HexToAddress(r.DestAddress)
			)
			if isSameSrc && isSameDst && isSameReceiver {
				tl.WalletAddress = walletAddr
				tl.WalletName = tradelogs.WalletAddrToName(walletAddr)
				updatedTradeLogs = append(updatedTradeLogs, tl)
				if len(updatedTradeLogs) == int(maxRecordSavePerTime) {
					if errS := storageInterface.SaveTradeLogs(updatedTradeLogs); errS != nil {
						return errS
					}
					updatedTradeLogs = nil
				}
			}
		}
	}
	if len(updatedTradeLogs) > 0 {
		return storageInterface.SaveTradeLogs(updatedTradeLogs)
	}
	return nil
}

func duneDataFromFile(path string) (duneQueryData, error) {
	var dqd duneQueryData
	file, err := os.Open(path)
	if err != nil {
		return dqd, err
	}
	bytesValue, err := ioutil.ReadAll(file)
	if err != nil {
		return dqd, err
	}
	if err = json.Unmarshal(bytesValue, &dqd); err != nil {
		return dqd, err
	}
	return dqd, nil
}

func addHexPrefix(s string) string {
	return fmt.Sprintf("0x%s", s)
}
