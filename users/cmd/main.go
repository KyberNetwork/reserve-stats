package main

import (
	"fmt"
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/users/cmc"
	"github.com/KyberNetwork/reserve-stats/users/http"
	"github.com/KyberNetwork/reserve-stats/users/stats"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"github.com/urfave/cli"
)

const (
	hostFlag     = "postgres_host"
	userFlag     = "postgres_user"
	passwordFlag = "postgres_password"
	databaseFlag = "postgres_database"
	bindFlag     = "bind"
)

func main() {
	app := app.NewApp()
	app.Name = "User stat module"
	app.Usage = "Store and return user stat information"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   hostFlag,
			Usage:  "Postgresql host to connect",
			EnvVar: "USER_POSTGRES_HOST",
			Value:  "127.0.0.1:5432",
		},
		cli.StringFlag{
			Name:   userFlag,
			Usage:  "Postgresql user to connect",
			EnvVar: "USER_POSTGRES_USER",
			Value:  "",
		},
		cli.StringFlag{
			Name:   passwordFlag,
			Usage:  "Postgresql password to connect",
			EnvVar: "USER_POSTGRES_PASSWORD",
			Value:  "",
		},
		cli.StringFlag{
			Name:   databaseFlag,
			Usage:  "Postgres database to connect",
			EnvVar: "USER_POSTGRES_DATABASE",
			Value:  "",
		},
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func run(c *cli.Context) error {
	logger, err := app.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Info("Run user module")
	// init storage
	userDB := storage.NewDB(
		sugar,
		c.String(hostFlag),
		c.String(userFlag),
		c.String(passwordFlag),
		c.String(databaseFlag),
	)

	// init stats
	cmc := cmc.NewCMCEthUSDRate(sugar)
	userStats := stats.NewUserStats(cmc, userDB)

	// run http server
	servePort := c.Int(bindFlag)
	host := fmt.Sprintf(":%d", servePort)
	server := http.NewServer(sugar, userStats, host)
	server.Run()
	return nil
}
