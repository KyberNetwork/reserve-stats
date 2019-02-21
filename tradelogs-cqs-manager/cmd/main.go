package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
)

const (
	listAllCqsFlag = "list-all-cqs"

	dropCqFlag = "drop-cq"

	executeCqFlag = "execute-cq"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs cqs manager"
	app.Usage = "Manage trade logs cqs"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.BoolFlag{
			Name:   listAllCqsFlag,
			Usage:  "List all tradelogs cqs",
			EnvVar: "LIST_ALL_CQS",
		},
		cli.StringFlag{
			Name:   dropCqFlag,
			Usage:  "Drop a cq",
			EnvVar: "DROP_CQ",
		},
		cli.StringFlag{
			Name:   executeCqFlag,
			Usage:  "Execute a cq",
			EnvVar: "EXECUTE_CQ",
		},
	)
	app.Flags = append(app.Flags, cq.NewCQFlags()...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func extractActualQuery(query string) (string, error) {
	begin := strings.Index(query, "BEGIN")
	end := strings.Index(query, "END")
	actualQuery := query[begin+len("BEGIN") : end]
	return actualQuery, nil
}

func getallCqs(influxClient client.Client, sugar *zap.SugaredLogger) (map[string]cq.ContinuousQuery, error) {
	var (
		logger = sugar.With(
			"func", "tradelogs-cq-manager/listAllCqs",
		)
		cqs = make(map[string]cq.ContinuousQuery)
	)
	q := fmt.Sprintf("SHOW CONTINUOUS QUERIES")
	res, err := influxdb.QueryDB(influxClient, q, "")
	if err != nil {
		logger.Debugw("influx query run error", "error", err)
		return cqs, err
	}

	// cqs map cq name to cq query
	for _, serie := range res[0].Series {
		for _, value := range serie.Values {
			cqName := value[0].(string)
			cqQuery := value[1].(string)
			actualQuery, err := extractActualQuery(cqQuery)
			if err != nil {
				return cqs, err
			}
			cqs[cqName] = cq.ContinuousQuery{
				Name:     cqName,
				Database: serie.Name,
				Query:    actualQuery,
			}
		}
	}
	return cqs, nil
}

func listAllCqs(influxClient client.Client, sugar *zap.SugaredLogger) error {
	var (
		logger = sugar.With(
			"func", "tradelogs-cq-manager/listaAllCqs",
		)
	)
	cqs, err := getallCqs(influxClient, sugar)
	if err != nil {
		return err
	}
	for cqName := range cqs {
		fmt.Println(cqName)
	}
	logger.Debug("list all cq complete")
	return nil
}

func dropCq(c client.Client, sugar *zap.SugaredLogger, cqName string) error {
	var (
		logger = sugar.With(
			"func", "tradelogs-cq-manager/cq/DropACQ",
		)
	)
	cqs, err := getallCqs(c, sugar)
	if err != nil {
		return err
	}
	cq, exist := cqs[cqName]
	if !exist {
		logger.Debugw("cq name does not exist", "name", cqName)
		return nil
	}
	return cq.Drop(c, sugar)
}

func executeCq(c client.Client, sugar *zap.SugaredLogger, cqName string) error {
	var (
		logger = sugar.With(
			"func", "tradelogs-cq-manager/cq/DropACQ",
		)
	)
	cqs, err := getallCqs(c, sugar)
	if err != nil {
		return err
	}
	cq, exist := cqs[cqName]
	if !exist {
		logger.Debugw("cq name does not exist", "name", cqName)
		return nil
	}

	//execute cq
	logger.Debugw("execute a cq", "cq", cq.Name)
	_, err = influxdb.QueryDB(c, cq.Query, cq.Database)
	logger.Debugw("cq query", "query", cq.Query)
	return err
}

func run(c *cli.Context) error {

	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	influxClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	sugar.Info("initialized influxClient successfully: ", influxClient)

	if c.Bool(listAllCqsFlag) {
		if err := listAllCqs(influxClient, sugar); err != nil {
			return err
		}
	}

	if c.String(dropCqFlag) != "" {
		cqToDrop := c.String(dropCqFlag)
		if err := dropCq(influxClient, sugar, cqToDrop); err != nil {
			return err
		}
	}

	if c.String(executeCqFlag) != "" {
		cqToExecute := c.String(executeCqFlag)
		if err := executeCq(influxClient, sugar, cqToExecute); err != nil {
			return err
		}
	}

	return nil
}
