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

//MustGetFromTimeFromContext return from time from context and error if it's not provide
func MustGetFromTimeFromContext(c *cli.Context) (time.Time, error) {
	timeString := c.String(fromTimeFlag)
	if timeString == "" {
		return time.Time{}, errors.New("From time flag is not provide")
	}
	const shortForm = "2006-Jan-02"
	return time.Parse(shortForm, timeString)
}

//GetToTimeFromContext return totime from context. Return time.Now if it's not provide
func GetToTimeFromContext(c *cli.Context) (time.Time, error) {
	timeString := c.String(toTimeFlag)
	if timeString == "" {
		return time.Now(), nil
	}
	const shortForm = "2006-Jan-02"
	return time.Parse(shortForm, timeString)
}
