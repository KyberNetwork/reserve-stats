package main

import (
	"fmt"
	"os"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/users/cmc"
	"github.com/KyberNetwork/reserve-stats/users/http"
	"github.com/KyberNetwork/reserve-stats/users/stats"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	hostFlag     = "postgres_host"
	userFlag     = "postgres_user"
	passwordFlag = "postgres_password"
	databaseFlag = "postgres_database"
	bindFlag     = "bind"
)

func configLog(stdoutLog bool) {
	logConfig := zap.NewDevelopmentConfig()
	if stdoutLog {
		logConfig.OutputPaths = []string{"stdout"}
	}
	logger, err := logConfig.Build()
	if err != nil {
		zap.S().Error("Cannot init zap logger")
	}
	zap.ReplaceGlobals(logger)
}

func main() {
	configLog(true)
	app := app.NewApp()
	app.Name = "User stat module"
	app.Usage = "Store and return user stat information"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   hostFlag,
			Usage:  "Postgresql host to connect",
			EnvVar: "POSTGRES_HOST",
			Value:  "127.0.0.1:5432",
		},
		cli.StringFlag{
			Name:   userFlag,
			Usage:  "Postgresql user to connect",
			EnvVar: "POSTGRES_USER",
			Value:  "",
		},
		cli.StringFlag{
			Name:   passwordFlag,
			Usage:  "Postgresql password to connect",
			EnvVar: "POSTGRES_PASSWORD",
			Value:  "",
		},
		cli.StringFlag{
			Name:   databaseFlag,
			Usage:  "Postgres database to connect",
			EnvVar: "POSTGRES_DATABASE",
			Value:  "",
		},
	)

	if err := app.Run(os.Args); err != nil {
		zap.S().Fatal(err)
	}

}

func run(c *cli.Context) error {
	zap.S().Info("Run user module")
	// init storage
	zap.S().Infof("Postgresql address: %s", c.String(hostFlag))
	userDB := storage.NewDB(
		c.String(hostFlag),
		c.String(userFlag),
		c.String(passwordFlag),
		c.String(databaseFlag),
	)

	// init stats
	cmc := cmc.NewCMCEthUSDRate()
	userStats := stats.NewUserStats(cmc, userDB)

	// run http server
	servePort := c.Int(bindFlag)
	host := fmt.Sprintf(":%d", servePort)
	server := http.NewServer(userStats, host)
	server.Run()
	return nil
}
