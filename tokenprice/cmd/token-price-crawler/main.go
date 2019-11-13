package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tokenprice/common"
	"github.com/KyberNetwork/reserve-stats/tokenprice/provider"
	"github.com/KyberNetwork/reserve-stats/tokenprice/storage"
	"github.com/KyberNetwork/reserve-stats/tokenprice/storage/postgres"
)

const (
	fromTimeFlag = "from-time"
	toTimeFlag   = "to-time"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Token rate crawler"
	app.Usage = "Crawl token rate from some EX"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   fromTimeFlag,
			Usage:  "provide from time to crawl token rate",
			EnvVar: "FROM_TIME",
		},
		cli.StringFlag{
			Name:   toTimeFlag,
			Usage:  "provide to time to crawl token rate",
			EnvVar: "TO_TIME",
		},
	)

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(storage.DefaultDB)...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func validateTime(fromTimeS, toTimeS string) (time.Time, time.Time, error) {
	var (
		fromTime, toTime time.Time
		err              error
	)
	if len(fromTimeS) == 0 {
		return fromTime, toTime, errors.New("from-time cannot be blank")
	}
	if len(toTimeS) == 0 {
		toTimeS = common.TimeToDateString(time.Now().UTC())
	}

	fromTime, err = common.DateStringToTime(fromTimeS)
	if err != nil {
		return fromTime, toTime, err
	}
	toTime, err = common.DateStringToTime(toTimeS)
	if err != nil {
		return fromTime, toTime, err
	}
	if toTime.Sub(fromTime) < 0 {
		return fromTime, toTime, errors.New("from-time must be smaller than to-time")
	}
	return fromTime, toTime, nil
}

func run(c *cli.Context) error {
	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	p, err := provider.NewProvider(provider.Coinbase)
	if err != nil {
		sugar.Errorw("failed to init provider", "error", err)
		return err
	}
	s, err := storage.NewStorageFromContext(sugar, c)
	if err != nil {
		sugar.Errorw("failed to init storage", "error", err)
		return err
	}

	var (
		fromTimeS = c.String(fromTimeFlag)
		toTimeS   = c.String(toTimeFlag)

		ethID = "ETH"
		usdID = "USD"

		timeW8PerRequest = 1 * time.Second
	)

	logger := sugar.With("token", ethID, "currency", usdID)

	if len(fromTimeS) == 0 && len(toTimeS) == 0 {
		logger.Info("from-time and to-time are blank, run get token rate daily...")
		if err := crawlTokenRateDaily(sugar, ethID, usdID, p, s); err != nil {
			logger.Errorw("failed to crawl token rate daily", "error", err)
			return err
		}
	}

	logger = logger.With("from time", fromTimeS, "to time", toTimeS)

	fromeTime, toTime, err := validateTime(fromTimeS, toTimeS)
	if err != nil {
		logger.Errorw("times are invalids", "error", err)
		return err
	}

	for t := fromeTime; t.Sub(toTime) <= 0; t = t.Add(24 * time.Hour) {
		rate, err := p.Rate(ethID, usdID, t)
		if err != nil {
			logger.Errorw("failed to get token rate", "error", err)
			return err
		}
		logger.Infow("get token rate successfully", "time", t, "rate", rate)
		if err := s.SaveTokenRate(ethID, usdID, p.Name(), t, rate); err != nil && err != postgres.ErrExists {
			logger.Errorw("failed to save data to database", "error", err)
			return err
		}
		logger.Info("save token rate successfully")
		// sleep for a second to avoid rate limit
		time.Sleep(timeW8PerRequest)
	}
	return nil
}

func crawlTokenRateDaily(sugar *zap.SugaredLogger, token, currency string, p provider.Provider, s storage.Storage) error {
	var (
		errCh  = make(chan error, 2)
		logger = sugar.With("func", caller.GetCurrentFunctionName(),
			"token", token, "currency", currency, "provider", p.Name())
	)
	job := func() {
		logger.Info("Running job")
		var now = time.Now().UTC()
		rate, err := p.Rate(token, currency, now)
		if err != nil {
			logger.Errorw("failed to get token rate", "error", err)
			errCh <- err
			return
		}
		logger.Infow("get token rate successfully", "time", now, "rate", rate)
		if err := s.SaveTokenRate(token, currency, p.Name(), now, rate); err != nil && err != postgres.ErrExists {
			logger.Errorw("failed to save data to database", "error", err)
			errCh <- err
		}
		logger.Info("save token rate successfully")
	}
	// get rate today
	job()

	if _, err := scheduler.Every().Day().At("00:00:01").Run(job); err != nil {
		return err
	}
	if err := <-errCh; err != nil {
		return err
	}
	return nil
}
