package tokenrate

import (
	"testing"

	gocache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestCachedRateProvider_USDRate(t *testing.T) {
	cache := gocache.New(defaultExpire, defaultExpire)
	cache.Set("test", 1, defaultExpire)
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	provider := NewCachedRateProvider(sugar, NewMock(), WithExpires(defaultExpire))
	rate, err := provider.USDRate(timeutil.TimestampMsToTime(1574681881000))
	require.NoError(t, err)
	assert.Equal(t, mockRate, rate)
	_, err = provider.USDRate(timeutil.TimestampMsToTime(1574681881000))
	require.NoError(t, err)

	_, err = provider.USDRate(timeutil.TimestampMsToTime(1574681882000))
	require.NoError(t, err)
	assert.Equal(t, mockRate, rate)
}
