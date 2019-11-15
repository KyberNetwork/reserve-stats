package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tokenprice/common"
	"github.com/KyberNetwork/reserve-stats/tokenprice/provider"
	"github.com/KyberNetwork/reserve-stats/tokenprice/storage"
)

const (
	fromTimeFlag = "from-time"
	toTimeFlag   = "to-time"
)

func main() {
	app := libapp.NewAppWithMode()
	app.Name = "Token price crawler"
	app.Usage = "Crawl token price from other exchanges"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   fromTimeFlag,
			Usage:  "provide from time to crawl token price with format YYYY-MM-DD, e.g: 2019-10-11",
			EnvVar: "FROM_TIME",
		},
		cli.StringFlag{
			Name:   toTimeFlag,
			Usage:  "provide to time to crawl token price wiht format YYYY-MM-DD, e.g: 2019-10-12",
			EnvVar: "TO_TIME",
		},
	)

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(storage.DefaultDB)...)
	app.Flags = append(app.Flags, libapp.NewSentryFlags()...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func validateTime(fromTimeS, toTimeS string) (time.Time, time.Time, error) {
	var (
		fromTime, toTime time.Time
		err              error
	)
	if len(fromTimeS) != 0 {
		fromTime, err = common.DateStringToTime(fromTimeS)
		if err != nil {
			return fromTime, toTime, err
		}
	}
	if len(toTimeS) != 0 {
		if len(fromTimeS) == 0 {
			return fromTime, toTime, errors.New("from-time cannot be blank")
		}
		toTime, err = common.DateStringToTime(toTimeS)
		if err != nil {
			return fromTime, toTime, err
		}
		if toTime.Sub(fromTime) < 0 {
			return fromTime, toTime, errors.New("from-time must be smaller than to-time")
		}
	}
	return fromTime, toTime, nil
}

func run(c *cli.Context) error {
	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	p, err := provider.NewPriceProvider(provider.Coinbase)
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
	)

	logger := sugar.With("token", common.ETHID, "currency", common.USDID)

	fromTime, toTime, err := validateTime(fromTimeS, toTimeS)
	if err != nil {
		return err
	}

	if len(toTimeS) != 0 {
		return crawlTokenPriceWithTimeRange(sugar, fromTime, toTime, p, s)
	}
	if len(fromTimeS) != 0 {
		logger.Info("to-time is blank, get price from from-time and run get price daily...")
		toTime = time.Now().UTC()
		if err := crawlTokenPriceWithTimeRange(sugar, fromTime, toTime, p, s); err != nil {
			logger.Errorw("failed to get rate with time range", "from-time", fromTime, "to-time", toTime)
			return err
		}
	}

	if err := crawlTokenPriceDaily(sugar, common.ETHID, common.USDID, p, s); err != nil {
		logger.Errorw("failed to get rate daily", "error", err)
	}
	cs := make(chan os.Signal, 1)
	signal.Notify(cs, os.Interrupt)
	<-cs
	logger.Info("got interrupt signal, program exited")
	return nil
}

func crawlTokenPriceWithTimeRange(sugar *zap.SugaredLogger, fromTime, toTime time.Time, p provider.PriceProvider, s storage.Storage) error {
	var (
		logger = sugar.With("func", caller.GetCurrentFunctionName())

		timeW8PerRequest = 1 * time.Second
	)
	logger = logger.With("from time", fromTime, "to time", toTime)

	for t := fromTime; t.Sub(toTime) <= 0; t = t.Add(24 * time.Hour) {
		price, err := p.Price(common.ETHID, common.USDID, t)
		if err != nil {
			logger.Errorw("failed to get token price", "error", err)
			return err
		}
		logger.Infow("get token price successfully", "time", t, "price", price)

		if err := s.SaveTokenPrice(common.ETHID, common.USDID, p.Name(), t, price); err != nil {
			logger.Errorw("failed to save data to database", "error", err)
			return err
		}
		logger.Info("save token price successfully")
		// sleep for a second to avoid rate limit
		time.Sleep(timeW8PerRequest)
	}
	return nil
}

func crawlTokenPriceDaily(sugar *zap.SugaredLogger, token, currency string, p provider.PriceProvider, s storage.Storage) error {
	var (
		logger = sugar.With("func", caller.GetCurrentFunctionName(),
			"token", token, "currency", currency, "provider", p.Name())
	)
	job := func() {
		logger.Info("Running job")
		var now = time.Now().UTC()
		price, err := p.Price(token, currency, now)
		if err != nil {
			logger.Errorw("failed to get token price", "error", err)
			return
		}
		logger.Infow("get token price successfully", "time", now, "price", price)
		if err := s.SaveTokenPrice(token, currency, p.Name(), now, price); err != nil {
			logger.Errorw("failed to save data to database", "error", err)
		}
		logger.Info("save token price successfully")
	}
	// get price today
	job()

	// run job get price daily
	l, _ := time.LoadLocation("Local")
	firstMomentOfDay := time.Now().Truncate(24 * time.Hour).In(l).Format("15:04:05")
	if _, err := scheduler.Every().Day().At(firstMomentOfDay).Run(job); err != nil {
		return err
	}
	return nil
}
