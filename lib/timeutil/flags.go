package timeutil

import (
	"errors"
	"time"

	"github.com/urfave/cli"
)

const (
	toTimeFlag   = "to"
	fromTimeFlag = "from"
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
