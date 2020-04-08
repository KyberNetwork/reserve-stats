package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/appnames"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/userprofile"
	"github.com/KyberNetwork/reserve-stats/tradelogs/http"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs HTTP Api"
	app.Usage = "Serve trade logs data"
	app.Version = "0.0.1"
	app.Action = func(c *cli.Context) error {
		if err := libapp.Validate(c); err != nil {
			return err
		}

		sugar, flush, err := libapp.NewSugaredLogger(c)
		if err != nil {
			return err
		}
		defer flush()

		tokenAmountFormatter, err := blockchain.NewToKenAmountFormatterFromContext(c)
		if err != nil {
			return err
		}

		storageInterface, err := storage.NewStorageInterfaceFromContext(sugar, c, tokenAmountFormatter)
		if err != nil {
			return err
		}

		var options []http.ServerOption
		addrToAppName, err := appnames.NewClientFromContext(sugar, c)
		if err != nil {
			return err
		}
		if addrToAppName != nil {
			options = append(options, http.WithApplicationNames(addrToAppName))
		}

		userClient, err := userprofile.NewClientFromContext(sugar, c)
		if err != nil {
			return err
		}

		cachedUserClient := userprofile.NewCachedClientFromContext(userClient, c)
		if cachedUserClient != nil {
			options = append(options, http.WithUserProfile(cachedUserClient))
		}

		symbolResolver, err := blockchain.NewTokenInfoGetterFromContext(c, storageInterface)
		if err != nil {
			return err
		}

		api := http.NewServer(storageInterface, httputil.NewHTTPAddressFromContext(c),
			sugar, symbolResolver, options...)
		err = api.Start()
		if err != nil {
			return err
		}

		if err = api.Start(); err != nil {
			return err
		}

		return nil
	}

	app.Flags = append(app.Flags, storage.NewCliFlags()...)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.TradeLogsPort)...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(storage.PostgresDefaultDB)...)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())
	app.Flags = append(app.Flags, appnames.NewCliFlags()...)
	app.Flags = append(app.Flags, userprofile.NewCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
