package main

import (
	"log"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/users/migration"
	"github.com/urfave/cli"
)

const (
	boltDataSourceFlag = "legacy-database"
	defaultDB          = "users"
)

func main() {
	app := libapp.NewApp()
	app.Name = "db migration"
	app.Usage = "migrate db from bolt to postgresql"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDB)...)
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:  boltDataSourceFlag,
			Usage: "bolt db file",
			Value: "users.db",
		},
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Migration crashed into err: %s", err)
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

	var legacyDBPath = c.String(boltDataSourceFlag)

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
	return dbMigration.Migrate(migration.DefaultUsersTableName, migration.DefaultAddressesTableName)
}
