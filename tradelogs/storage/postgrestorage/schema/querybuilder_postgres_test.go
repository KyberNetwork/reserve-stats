package schema

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestGetTimeCondition(t *testing.T) {
	const (
		fromTime  = 1539000000000
		toTime    = 1539250666000
		frequency = "h"
	)
	var (
		from = timeutil.TimestampMsToTime(uint64(fromTime))
		to   = timeutil.TimestampMsToTime(uint64(toTime))
	)
	condition, err := BuildTimeCondition(from, to, frequency)
	require.NoError(t, err)
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	sugar.Infow("condition generated successful", "condition", condition)
}

func TestBuildDateTruncField(t *testing.T) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	sugar.Infow("build date_trunc", "output", BuildDateTruncField("day", 7))
}

func TestRoundHourTime(t *testing.T) {
	const fromTime = 1539250666000
	from := timeutil.TimestampMsToTime(uint64(fromTime))
	from = RoundTime(from, "date", 7)
	require.Equal(t, "2018-10-10 17:00:00", from.UTC().Format(DefaultDateFormat))
}
