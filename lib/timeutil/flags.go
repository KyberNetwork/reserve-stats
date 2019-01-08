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

const timeLayout = "2006-Jan-02"

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

//GetFromTimeFromContext return from time from context and error if it's not provide
func GetFromTimeFromContext(c *cli.Context) (time.Time, error) {
	timeString := c.String(fromTimeFlag)
	if timeString == "" {
		return time.Time{}, errors.New("from time flag is not provide")
	}
	return time.Parse(timeLayout, timeString)
}

//GetToTimeFromContextWithDaemon return to time from context. Return err=ErrEmptyFlag to indicate daemon Mode if it's not provide
func GetToTimeFromContextWithDaemon(c *cli.Context) (time.Time, error) {
	timeString := c.String(toTimeFlag)
	if timeString == "" {
		return time.Time{}, ErrEmptyFlag
	}
	return time.Parse(timeLayout, timeString)
}
