package timeutil

import (
	"errors"
	"time"

	"github.com/urfave/cli"
)

const (
	toTimeFlag        = "to"
	fromTimeFlag      = "from"
	timePrecisionFlag = "time-precision"
)

// NewTimeRangeCliFlags returns cli flags to configure a fromTime-toTime.
func NewTimeRangeCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.Uint64Flag{
			Name:   fromTimeFlag,
			Usage:  "from time. Default is  00:00:00 UTC on 1 January 1970 in millisecond",
			EnvVar: "FROM_TIME",
			Value:  0,
		},
		cli.Uint64Flag{
			Name:   toTimeFlag,
			Usage:  "to time. Default is time.Now() in millisecond",
			EnvVar: "TO_TIME",
			Value:  0,
		},
	}
}

//MustGetFromTimeFromContext return from time from context and error if it's not provide
func MustGetFromTimeFromContext(c *cli.Context) (time.Time, error) {
	fromTime := TimestampMsToTime(c.Uint64(fromTimeFlag))
	if c.Uint64(fromTimeFlag) == 0 {
		return fromTime, errors.New("From time flag is not provide")
	}

	return fromTime, nil
}

//GetToTimeFromContext return totime from context. Return time.Now if it's not provide
func GetToTimeFromContext(c *cli.Context) time.Time {
	if c.Uint64(toTimeFlag) == 0 {
		return time.Now()
	}
	toTime := TimestampMsToTime(c.Uint64(toTimeFlag))

	return toTime
}
