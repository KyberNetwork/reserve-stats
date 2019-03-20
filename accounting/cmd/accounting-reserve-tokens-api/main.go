package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

func main() {
	app := libapp.NewApp()
	app.Name = "User stat module"
	app.Usage = "Store and return user stat information"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultDB)...)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.UsersPort)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	sugar, flusher, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flusher()

	sugar.Info("Run reserve token api")

	// db, err := libapp.NewDBFromContext(c)
	// if err != nil {
	// 	return err
	// }

	// listedTokenStorage, err := storage.NewDB(
	// 	sugar,
	// 	db,
	// )
	// if err != nil {
	// 	return err
	// }

	// defer listedTokenStorage.Close()

	return nil
}
