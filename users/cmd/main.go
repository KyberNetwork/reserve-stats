package main

import (
	"fmt"
	"log"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/ethrate"
	"github.com/KyberNetwork/reserve-stats/users/http"
	"github.com/KyberNetwork/reserve-stats/users/stats"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"github.com/urfave/cli"
)

const (
	defaultDB = "users"
	bindFlag  = "bind"
)

func main() {
	app := libapp.NewApp()
	app.Name = "User stat module"
	app.Usage = "Store and return user stat information"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDB)...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Info("Run user module")
	// init storage
	userDB := storage.NewDB(
		sugar,
		libapp.NewDBFromContext(c),
	)

	// init stats
	cmc := ethrate.NewCMCRate(sugar)
	userStats := stats.NewUserStats(cmc, userDB)

	// run http server
	servePort := c.Int(bindFlag)
	host := fmt.Sprintf(":%d", servePort)
	server := http.NewServer(sugar, userStats, host)
	server.Run()
	return nil
}
