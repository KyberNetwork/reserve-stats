package main

import (
	"fmt"
	"os"

	"github.com/KyberNetwork/reserve-stats/users/http"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	servePort = 9000
	hostFlag  = "host"
)

func configLog(stdoutLog bool) {
	logConfig := zap.NewDevelopmentConfig()
	//TODO: if stdout true write log to stdout and file else
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
	app := cli.NewApp()
	app.Name = "User stat module"
	app.Usage = "Store and return user stat information"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   hostFlag,
			Usage:  "",
			EnvVar: "POSTGRE_HOST",
			Value:  "",
		},
	}

	if err := app.Run(os.Args); err != nil {
		zap.S().Fatal(err)
	}

}

func run(c *cli.Context) error {
	zap.S().Info("Run user module")
	// init storage
	userDB := storage.NewDB(
		"127.0.0.1:5432",
		"hahoang",
		"",
		"hahoang",
	)

	// run http server
	host := fmt.Sprintf(":%d", servePort)
	server := http.NewServer(userDB, host)
	server.Run()
	return nil
}
