package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	libredis "github.com/KyberNetwork/reserve-stats/lib/redis"
	"github.com/KyberNetwork/reserve-stats/users/common"
	server "github.com/KyberNetwork/reserve-stats/users/public-server"
	"github.com/KyberNetwork/tokenrate/coingecko"
)

const (
	nonKYCDailyLimitFlag         = "non-kyc-daily-limit"
	nonKYCTxLimitFlag            = "non-kyc-tx-limit"
	kycedDailyLimitFlag          = "kyced-daily-limit"
	kycedTxLimitFlag             = "kyced-tx-limit"
	defaultNonKYCDailyLimitValue = 15000
	defaultNonKYCTxLimitValue    = 15000
	defaultKYCEDDailyLimitValue  = 1000000
	defaultKYCEDTxLimitValue     = 200000
)

func main() {
	app := libapp.NewApp()
	app.Name = "User stat public service"
	app.Usage = "Return user stat information from cache"
	app.Action = run
	app.Version = "0.1"

	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.UsersPublicPort)...)
	app.Flags = append(app.Flags, libredis.NewCliFlags()...)
	app.Flags = append(app.Flags,
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
	if err := app.Validate(c); err != nil {
		return err
	}

	sugar, flusher, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flusher()
	sugar.Info("Run user stats public service")

	redisClient, err := libredis.NewClientFromContext(c)
	if err != nil {
		return err
	}

	sugar.Debugw("initiate redis client", "client", redisClient)

	nonKYCDailyLimit := c.Float64(nonKYCDailyLimitFlag)
	nonKYCTxLimit := c.Float64(nonKYCTxLimitFlag)
	kycedDailyLimit := c.Float64(kycedDailyLimitFlag)
	kycedTxLimit := c.Float64(kycedTxLimitFlag)

	nonKYCCap := common.NewUserCap(true, common.WithLimit(nonKYCDailyLimit, nonKYCTxLimit))
	kycedCap := common.NewUserCap(true, common.WithLimit(kycedDailyLimit, kycedTxLimit))

	publicServer := server.NewServer(sugar, httputil.NewHTTPAddressFromContext(c), coingecko.New(), redisClient, kycedCap, nonKYCCap)

	return publicServer.Run()
}
