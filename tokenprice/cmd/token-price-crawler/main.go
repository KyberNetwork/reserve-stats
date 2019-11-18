package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tokenprice/common"
	"github.com/KyberNetwork/reserve-stats/tokenprice/provider"
	"github.com/KyberNetwork/reserve-stats/tokenprice/storage"
)

const (
	fromTimeFlag       = "from-time"
	toTimeFlag         = "to-time"
	jobRunningTimeFlag = "job-running-time"
	sourceFlag         = "source"

	defaultJobRunningTime = "07:00:00"
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
		cli.StringFlag{
			Name:   jobRunningTimeFlag,
			Usage:  "crawler will fetch the price daily at this time, e.g: 07:00:00",
			EnvVar: "JOB_RUNNING_TIME",
			Value:  defaultJobRunningTime,
		},
		cli.StringFlag{
			Name:   sourceFlag,
			Usage:  "provide source to get price",
			EnvVar: "SOURCE",
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
		currentTime      = time.Now().UTC()
	)
	if len(fromTimeS) != 0 {
		fromTime, err = common.YYYYMMDDToTime(fromTimeS)
		if err != nil {
			return fromTime, toTime, err
		}
	} else {
		fromTime = currentTime
	}
	if len(toTimeS) != 0 {
		toTime, err = common.YYYYMMDDToTime(toTimeS)
		if err != nil {
			return fromTime, toTime, err
		}
		if toTime.Sub(currentTime) > 0 {
			toTime = currentTime
		}
		if toTime.Sub(fromTime) < 0 {
			return fromTime, toTime, errors.New("from-time must be smaller than to-time")
		}
	} else {
		toTime = currentTime
	}
	return fromTime, toTime, nil
}

func run(c *cli.Context) error {
	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	var (
		source = c.String(sourceFlag)
		ps     []provider.PriceProvider
	)

	if len(source) != 0 {
		p, err := provider.NewPriceProvider(c, source)
		if err != nil {
			sugar.Errorw("failed to init provider", "error", err)
			return err
		}
		ps = append(ps, p)
	} else {
		ps = provider.AllProvider(c)
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
		return errors.Wrap(err, "invalid time")
	}

	if len(toTimeS) != 0 {
		return crawlTokenPriceWithTimeRange(sugar, fromTime, toTime, ps, s)
	}
	logger.Info("to-time is blank, get history price from from-time and run get price daily...")
	if err := crawlTokenPriceWithTimeRange(sugar, fromTime, toTime, ps, s); err != nil {
		logger.Errorw("failed to get rate with time range", "from-time", fromTime, "to-time", toTime)
		return err
	}
	if err := crawlTokenPriceDaily(sugar, ps, s, c.String(jobRunningTimeFlag)); err != nil {
		logger.Errorw("failed to get rate daily", "error", err)
	}
	cs := make(chan os.Signal, 1)
	signal.Notify(cs, os.Interrupt)
	<-cs
	logger.Info("got interrupt signal, program exited")
	return nil
}

func crawlTokenPriceWithTimeRange(
	sugar *zap.SugaredLogger,
	fromTime, toTime time.Time,
	ps []provider.PriceProvider,
	s storage.Storage) error {
	var (
		logger = sugar.With("func", caller.GetCurrentFunctionName(),
			"from time", fromTime, "to time", toTime)
	)
	eg, _ := errgroup.WithContext(context.Background())
	for _, p := range ps {
		var (
			p       = p
			pLogger = logger.With("provider", p.Name())
		)
		eg.Go(func() error {
			for t := fromTime; t.Sub(toTime) <= 0; t = t.Add(24 * time.Hour) {
				pLogger.Infof("provider = %s", p.Name())
				price, err := p.ETHPrice(t)
				if err != nil {
					pLogger.Errorw("failed to get token price", "error", err)
					return err
				}
				pLogger.Infow("get token price successfully", "time", t, "price", price)

				if err := s.SaveTokenPrice(common.ETHID, common.USDID, p.Name(), t, price); err != nil {
					pLogger.Errorw("failed to save data to database", "error", err)
					return err
				}
				pLogger.Info("save token price successfully")
				// avoid rate limit
				p.Wait()
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "failed to get token price")
	}
	return nil
}

func crawlTokenPriceDaily(sugar *zap.SugaredLogger, ps []provider.PriceProvider, s storage.Storage, jobRunningTime string) error {
	var (
		logger = sugar.With("func", caller.GetCurrentFunctionName(),
			"token", common.ETHID,
			"currency", common.USDID,
			"job running time", jobRunningTime)
	)
	if _, err := time.Parse("15:04:05", jobRunningTime); err != nil {
		return err
	}
	job := func() {
		logger.Info("Running job")
		var now = time.Now().UTC()
		for _, p := range ps {
			logger.Infof("provider = %s", p.Name())
			price, err := p.ETHPrice(now)
			if err != nil {
				logger.Errorw("failed to get token price", "error", err)
				return
			}
			logger.Infow("get token price successfully", "time", now, "price", price)
			if err := s.SaveTokenPrice(common.ETHID, common.USDID, p.Name(), now, price); err != nil {
				logger.Errorw("failed to save data to database", "error", err)
			}
			logger.Info("save token price successfully")
		}
	}
	// run job get price daily
	if _, err := scheduler.Every().Day().At(jobRunningTime).Run(job); err != nil {
		return errors.Wrap(err, "failed to run daily job")
	}
	return nil
}
