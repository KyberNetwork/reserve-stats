package schema

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

// DateFunctionParams params for frequency
var DateFunctionParams = map[string]string{
	"h": "hour",
	"d": "day",
}

// BuildDateTruncField for aggregated query
func BuildDateTruncField(dateTruncParam string, timeZone int8) string {
	if timeZone != 0 && dateTruncParam == "day" {
		var intervalParse = fmt.Sprintf("interval '%d hour'", timeZone)
		return "date_trunc('" + dateTruncParam + "', timestamp + " + intervalParse + ") - " + intervalParse
	}
	return `date_trunc('` + dateTruncParam + `', timestamp)`
}

// RoundTime returns time is rounded by day or hour
// if time is rounded by day, it also use time zone.
func RoundTime(t time.Time, freq string, timeZone int8) time.Time {
	if freq == "hour" {
		return t.Truncate(time.Hour)
	}
	return timeutil.Midnight(t.In(time.FixedZone("", int(timeZone)*60*60)))
}
