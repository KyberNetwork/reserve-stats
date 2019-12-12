package common

import (
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestBlacklist(t *testing.T) {
	logger := testutil.MustNewDevelopmentSugaredLogger()
	bl, err := newBlacklist(logger, "../tests/blacklist.json")
	require.NoError(t, err)

	require.True(t, bl.IsBanned(ethereum.HexToAddress("0xa09871AEadF4994Ca12f5c0b6056BBd1d343c029")))
	require.False(t, bl.IsBanned(ethereum.HexToAddress("0xa09871AEadF4994Ca12f5c0b6056B")))
}
