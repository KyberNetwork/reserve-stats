package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	libredis "github.com/KyberNetwork/reserve-stats/lib/redis"
	"github.com/KyberNetwork/reserve-stats/users/cacher"
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/storage"
)

const (
	expireTimeFlag               = "expire-time"
	nonKYCDailyLimitFlag         = "non-kyc-daily-limit"
	nonKYCTxLimitFlag            = "non-kyc-tx-limit"
	kycedDailyLimitFlag          = "kyced-daily-limit"
	kycedTxLimitFlag             = "kyced-tx-limit"
	defaultExpireTime            = 3600 // 1 hour
	defaultNonKYCDailyLimitValue = 15000
	defaultNonKYCTxLimitValue    = 15000
	defaultKYCEDDailyLimitValue  = 1000000
	defaultKYCEDTxLimitValue     = 200000
)

func main() {
	app := libapp.NewApp()
	app.Name = "app for caching user information"
	app.Usage = "cache user info hourly"
	app.Action = run
	app.Version = "0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultDB)...)
	app.Flags = append(app.Flags, libredis.NewCliFlags()...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags,
		cli.IntFlag{
			Name:   expireTimeFlag,
			Usage:  "Time to expire redis cache, count by second",
			EnvVar: "EXPIRE_TIME",
			Value:  defaultExpireTime,
		},
		cli.Float64Flag{
			Name:   nonKYCDailyLimitFlag,
			Usage:  "Daily limit for non kyc user",
			EnvVar: "NON_KYC_DAILY_LIMIT",
			Value:  defaultNonKYCDailyLimitValue,
		},
		cli.Float64Flag{
			Name:   nonKYCTxLimitFlag,
			Usage:  "Tx limit for non kyc user",
			EnvVar: "NON_KYC_TX_LIMIT",
			Value:  defaultNonKYCTxLimitValue,
		},
		cli.Float64Flag{
			Name:   kycedDailyLimitFlag,
			Usage:  "Daily limit for kyced user",
			EnvVar: "KYCED_DAILY_LIMIT",
			Value:  defaultKYCEDDailyLimitValue,
		},
		cli.Float64Flag{
			Name:   kycedTxLimitFlag,
			Usage:  "Tx limit for kyced user",
			EnvVar: "KYCED_TX_LIMIT",
			Value:  defaultKYCEDTxLimitValue,
		},
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
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

	sugar.Info("Run user public cacher")

	postgresDB, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	userDB, err := storage.NewDB(
		sugar,
		postgresDB,
	)
	if err != nil {
		return err
	}
	sugar.Debugw("Initiated postgres client", "client", userDB)

	influxDBClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	redisCacheClient, err := libredis.NewClientFromContext(c)
	if err != nil {
		return err
	}

	nonKYCDailyLimit := c.Float64(nonKYCDailyLimitFlag)
	nonKYCTxLimit := c.Float64(nonKYCTxLimitFlag)
	kycedDailyLimit := c.Float64(kycedDailyLimitFlag)
	kycedTxLimit := c.Float64(kycedTxLimitFlag)

	nonKYCCap := common.NewUserCap(true, common.WithLimit(nonKYCDailyLimit, nonKYCTxLimit))
	kycedCap := common.NewUserCap(true, common.WithLimit(kycedDailyLimit, kycedTxLimit))

	expireTimeSecond := c.Int64(expireTimeFlag)
	expireTime := time.Duration(expireTimeSecond) * time.Second

	sugar.Debugw("Initiated redis cached", "cache", redisCacheClient)

	redisCacher := cacher.NewRedisCacher(sugar, userDB, influxDBClient, redisCacheClient, expireTime, kycedCap, nonKYCCap)

	return redisCacher.CacheUserInfo()
}
