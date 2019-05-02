package common

import "github.com/urfave/cli"

const (
	nonKYCDailyLimitFlag         = "non-kyc-daily-limit"
	defaultNonKYCDailyLimitValue = 15000

	nonKYCTxLimitFlag         = "non-kyc-tx-limit"
	defaultNonKYCTxLimitValue = 15000
)

//UserCap is users transaction cap.
type UserCap struct {
	// DailyLimit is the USD amount if the user is considered rich
	// and will receive different rates.
	DailyLimit float64 `json:"daily_limit"`
	// TxLimit is the maximum value in USD of a transaction an user
	// is allowed to make.
	TxLimit float64 `json:"tx_limit"`
}

// UserCapConfiguration is the cap configuration for KYC and non-KYC users.
type UserCapConfiguration struct {
	NonKYC UserCap
}

// NewUserCapConfiguration creates new instance of UserCapConfiguration from given parameters.
func NewUserCapConfiguration(nonKYCDailyLimit, nonKYCTxLimit float64) *UserCapConfiguration {
	return &UserCapConfiguration{
		NonKYC: UserCap{
			DailyLimit: nonKYCDailyLimit,
			TxLimit:    nonKYCTxLimit,
		},
	}
}

// UserCap returns UserCap of user for either kyced or non kyced.
func (c *UserCapConfiguration) UserCap() UserCap {
	return c.NonKYC
}

// IsRich returns true if user volume is greater or equal to daily limit.
func (c *UserCapConfiguration) IsRich(volume float64) bool {
	return volume >= c.NonKYC.DailyLimit
}

// NewUserCapCliFlags creates new cli configuration flags for user cap.
func NewUserCapCliFlags() []cli.Flag {
	return []cli.Flag{
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
	}
}

// NewUserCapConfigurationFromContext creates new UserCapConfiguration from given cli context.
func NewUserCapConfigurationFromContext(c *cli.Context) *UserCapConfiguration {
	return NewUserCapConfiguration(
		c.Float64(nonKYCDailyLimitFlag),
		c.Float64(nonKYCTxLimitFlag),
	)
}
