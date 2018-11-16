package main

import (
	"log"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/priceanalytics/migration"
	"github.com/urfave/cli"
)

const (
	boltDatabSourceFlag = "legacy-database"
	defaultDB           = "reserve_stats"
)

func main() {
	app := libapp.NewApp()
	app.Name = "price analytic migration"
	app.Usage = "migrate db from bolt to postgres"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDB)...)
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:  boltDatabSourceFlag,
			Usage: "bolt db file",
			Value: "priceanalytic.db",
		})
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}

	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	var legacyDBPath = c.String(boltDatabSourceFlag)

	f, err := os.Open(legacyDBPath)
	if err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}
	dbMigration, err := migration.NewDBMigration(sugar, legacyDBPath, db)
	if err != nil {
		return err
	}
	return dbMigration.Migrate()
}
