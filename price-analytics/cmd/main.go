package main

import (
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/urfave/cli"
)

const (
	httpHost             = ":9001"
	postgresHostFlag     = "postgres_host"
	postgresUserFlag     = "postgres_user"
	postgresPasswordFlag = "postgres_password"
	postgresDatabaseFlag = "postgres_database"
)

func main() {
	app := app.NewApp()
	app.Name = "Price analytics data"
	app.Usage = "store price analytic data"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  postgresHostFlag,
			Usage: "host to connect to postgres database",
			Value: "127.0.0.1:5432",
		},
		cli.StringFlag{
			Name:  postgresUserFlag,
			Usage: "user to connect to postgres database",
			Value: "",
		},
		cli.StringFlag{
			Name:  postgresPasswordFlag,
			Usage: "password to connect to postgres database",
			Value: "",
		},
		cli.StringFlag{
			Name:  postgresDatabaseFlag,
			Usage: "database to connect to postgres database",
			Value: "",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	return nil
}
