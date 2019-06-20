package schema

import (
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/stretchr/testify/require"
	"testing"
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
