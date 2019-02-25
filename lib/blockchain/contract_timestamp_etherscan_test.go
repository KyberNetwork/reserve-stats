package blockchain

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nanmu42/etherscan-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestEtherscanContractTimestampResolver(t *testing.T) {
	testutil.SkipExternal(t)

	sugar := testutil.MustNewDevelopmentSugaredLogger()
	client := etherscan.New(etherscan.Mainnet, "")
	resolv := NewEtherscanContractTimestampResolver(sugar, client)

	ts, err := resolv.Resolve(common.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"))
	require.NoError(t, err)
	assert.Equal(t, ts.Unix(), int64(1518038157))

	// not a contract
	_, err = resolv.Resolve(common.HexToAddress("0x8007ce15acda724689760B4Ba493d4766F973649"))
	assert.Equal(t, ErrNotAvailble, err)

	// non existing address
	_, err = resolv.Resolve(common.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D80"))
	assert.Equal(t, ErrNotAvailble, err)

}
