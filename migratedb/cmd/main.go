package main

import (
	"log"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/migratedb/migratedb"
	"github.com/urfave/cli"
)

const (
	boltDataSourceFlag = "database"
	defaultDB          = "users"
)

func main() {
	app := libapp.NewApp()
	app.Name = "db migration"
	app.Usage = "migrate db from bolt to postgresql"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDB)...)
	app.Flags = append(app.Flags, []cli.Flag{
		cli.StringFlag{
			Name:  boltDataSourceFlag,
			Usage: "bolt db file",
			Value: "users.db",
		},
	}...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func run(c *cli.Context) error {
	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}
	dbMigrate, err := migratedb.NewMigrateStorage(c.String(boltDataSourceFlag), db)
	if err != nil {
		return err
	}
	return dbMigrate.MigrateDB()
}
