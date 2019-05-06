package common

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestCapFunctions(t *testing.T) {
	logger := testutil.MustNewDevelopmentSugaredLogger()
	logger.Info("test caps configuration functions")

	var (
		nonKYCUserDailyCap = 10000.0
		nonKYCUserTxCap    = 1000.0
		kycUserDailyCap    = 100000.0
		kycUserTxCap       = 10000.0
	)

	userCap := NewUserCapConfiguration(nonKYCUserDailyCap, nonKYCUserTxCap, kycUserDailyCap, kycUserTxCap)

	nonKYCUserCap := userCap.UserCap(false)
	assert.Equal(t, nonKYCUserDailyCap, nonKYCUserCap.DailyLimit)
	assert.Equal(t, nonKYCUserTxCap, nonKYCUserCap.TxLimit)

	assert.True(t, userCap.IsRich(false, 10000.0+1))
	assert.False(t, userCap.IsRich(false, 10000.0-1))
}
