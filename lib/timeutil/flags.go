package timeutil

import (
	"errors"
	"time"

	"github.com/urfave/cli"
)

const (
	toTimeFlag         = "to"
	fromTimeFlag       = "from"
	fromMillisTimeFlag = "from-millis"
	toMillisTimeFlag   = "to-millis"
)

//ErrEmptyFlag is the error returned when empty flag
var ErrEmptyFlag = errors.New("empty flag")

// NewTimeRangeCliFlags returns cli flags to configure a fromTime-toTime.
func NewTimeRangeCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   fromTimeFlag,
			Usage:  "from time in format YYYY-MM-DD",
			EnvVar: "FROM_TIME",
		},
		cli.StringFlag{
			Name:   toTimeFlag,
			Usage:  "to time in format YYYY-MM-DD. Default is time.Now()",
			EnvVar: "TO_TIME",
			Value:  "",
		},
	}
}

func timeFlagFromContext(c *cli.Context, flag string) (time.Time, error) {
	const shortForm = "2006-01-02"
	timeString := c.String(flag)
	if timeString == "" {
		return time.Time{}, ErrEmptyFlag
	}
	return time.Parse(shortForm, timeString)
}

//FromTimeFromContext return from time from context and error if it's not provide
func FromTimeFromContext(c *cli.Context) (time.Time, error) {
	return timeFlagFromContext(c, fromTimeFlag)
}

//ToTimeFromContext return to time from context. Return err=ErrEmptyFlag to indicate daemon Mode if it's not provide
func ToTimeFromContext(c *cli.Context) (time.Time, error) {
	return timeFlagFromContext(c, toTimeFlag)
}

func millisTimeFlagFromContext(c *cli.Context, flag string) (time.Time, error) {
	timeUint := c.Uint64(flag)
	if timeUint == 0 {
		return time.Time{}, nil
	}
	return TimestampMsToTime(timeUint), nil
}

//FromTimeMillisFromContext return from time from context. Return err=ErrEmptyFlag if no flag is provided
func FromTimeMillisFromContext(c *cli.Context) (time.Time, error) {
	return millisTimeFlagFromContext(c, fromMillisTimeFlag)
}

//ToTimeMillisFromContext return from time from context. Return err=ErrEmptyFlag if no flag is provided
func ToTimeMillisFromContext(c *cli.Context) (time.Time, error) {
	return millisTimeFlagFromContext(c, toMillisTimeFlag)
}

func millisTimestampFlagFromContext(c *cli.Context, flag string) (uint64, error) {
	timeUint := c.Uint64(flag)
	if timeUint == 0 {
		return 0, ErrEmptyFlag
	}
	return timeUint, nil
}

//FromTimestampMillisFromContext return from time (in timestamp) from context. Return err=ErrEmptyFlag if no flag is provided
func FromTimestampMillisFromContext(c *cli.Context) (uint64, error) {
	return millisTimestampFlagFromContext(c, fromMillisTimeFlag)
}

//ToTimestampMillisFromContext return to time (in timestamp) from context. Return err=ErrEmptyFlag if no flag is provided
func ToTimestampMillisFromContext(c *cli.Context) (uint64, error) {
	return millisTimestampFlagFromContext(c, toMillisTimeFlag)
}

//NewMilliTimeRangeCliFlags return clit flags to input from/to time in millisecond from terminal
func NewMilliTimeRangeCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.Uint64Flag{
			Name:   fromMillisTimeFlag,
			Usage:  "From timestamp(millisecond) to query from",
			EnvVar: "FROM_MILLIS",
		},
		cli.Uint64Flag{
			Name:   toMillisTimeFlag,
			Usage:  "To timestamp(millisecond) to query to",
			EnvVar: "TO_MILLIS",
		},
	}
}
