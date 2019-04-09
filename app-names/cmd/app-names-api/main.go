package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/app-names/http"
	"github.com/KyberNetwork/reserve-stats/app-names/storage"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

const (
	defaultDB = "app-names"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Integration Application Name manager"
	app.Action = run
	app.Version = "0.0.1"
	app.Flags = append(app.Flags)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.AppName)...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDB)...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}
	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	appNameDB, err := storage.NewAppNameDB(sugar, db)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := db.Close(); cErr != nil {
			sugar.Errorw("failed to close database", err, cErr)
		}
	}()

	server, err := http.NewServer(httputil.NewHTTPAddressFromContext(c), appNameDB, sugar)
	if err != nil {
		return err
	}

	sugar.Info("Run Addr to Appname module")
	if err = server.Run(); err != nil {
		return err
	}
	return db.Close()
}
