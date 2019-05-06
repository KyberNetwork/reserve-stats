package fetcher

import (
	"math/big"
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/nanmu42/etherscan-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestFetcher(t *testing.T) {
	testutil.SkipExternal(t)

	var (
		sugar     = testutil.MustNewDevelopmentSugaredLogger()
		testAddr1 = ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F")
		//testAddr2 = ethereum.HexToAddress("0xdd974d5c2e2928dea5f71b9825b8b646686bd200") // KNC ERC20 contract
		client = etherscan.New(etherscan.Mainnet, "")
		offset = 500
	)

	f := NewEtherscanTransactionFetcher(sugar, client, 2)

	// the result should not included the to block: 7358394
	expected := []int{7356442, 7356961, 7357002, 7357451, 7357872, 7358016, 7358169, 7358208}
	normalTxs, err := f.NormalTx(testAddr1, big.NewInt(7356442), big.NewInt(7358394), offset)
	require.NoError(t, err)
	require.Len(t, normalTxs, len(expected))
	for i := range normalTxs {
		assert.Equal(t, expected[i], normalTxs[i].BlockNumber)
	}

	// this test is to proved that Etherscan API can returns > 100000 normal transactions.
	// It is disabled as it will take a very long time to finish.
	//normalTxs, err = f.Fetch(testAddr2, nil, nil)
	//require.NoError(t, err)
	//assert.True(t, len(normalTxs) > 400000)

	internalTxs, err := f.InternalTx(testAddr1, big.NewInt(7356442), big.NewInt(7356500), offset)
	require.NoError(t, err)
	expectedHashes := []string{
		"0x111580b2c03d6bae12e9113ff4fc46da3c38daf98f2b282c0eb1ebc0be57c870",
		"0x0c64431234e94af1073dd812bd38c518acd333fcb2599c56284c8367ae4b11d4",
		"0x485210e2464218ade460fb2e6a392b9445b3bfb5d3cd65f33f8e8edd876d810c",
		"0xc8100083fb30da902290310b53c6abd4d2d0f5762ccfb0a7830a66d138c63393",
		"0x9781a916fd91fe7fc9505b0d8eac48cfed1aea660ae87b3f24a75f8918f85413",
		"0x6ef3406b11f61040e156e1fc4bdcee23865ce63431da9906438110200862c05d",
		"0xd25a43b41787c61056aca478b771ecf19147a2577403dad0c2e9276c62490ca9",
		"0xe1dbffc4838aa2e434ffe5fdf8a766568d92158358cadc6e650da84933801b86",
		"0x35fe5bfff1127f9694785a000515c1291fe26965997a357759d508c43ada0488",
	}
	for i := range internalTxs {
		assert.Equal(t, expectedHashes[i], internalTxs[i].Hash)
	}

	transfers, err := f.ERC20Transfer(testAddr1, big.NewInt(7356442), big.NewInt(7356500), offset)
	require.NoError(t, err)
	expectedHashes = []string{
		"0x3dbb05df251ee6c5fe4a4334baa05dfcf9ef85295487b1fc2f9fe7b72a8b7b5f",
		"0x111580b2c03d6bae12e9113ff4fc46da3c38daf98f2b282c0eb1ebc0be57c870",
		"0x0c64431234e94af1073dd812bd38c518acd333fcb2599c56284c8367ae4b11d4",
		"0x485210e2464218ade460fb2e6a392b9445b3bfb5d3cd65f33f8e8edd876d810c",
		"0xc8100083fb30da902290310b53c6abd4d2d0f5762ccfb0a7830a66d138c63393",
		"0x9781a916fd91fe7fc9505b0d8eac48cfed1aea660ae87b3f24a75f8918f85413",
		"0x6ef3406b11f61040e156e1fc4bdcee23865ce63431da9906438110200862c05d",
		"0xd25a43b41787c61056aca478b771ecf19147a2577403dad0c2e9276c62490ca9",
		"0xe1dbffc4838aa2e434ffe5fdf8a766568d92158358cadc6e650da84933801b86",
		"0x35fe5bfff1127f9694785a000515c1291fe26965997a357759d508c43ada0488",
	}
	for i := range transfers {
		assert.Equal(t, expectedHashes[i], transfers[i].Hash.Hex())
	}

	// this test is to proved that Etherscan API can returns > 100000 ERC20 transfers.
	// It is disabled as it will take a very long time to finish.
	//transfers, err = f.ERC20Transfer(testAddr1, nil, nil)
	//require.NoError(t, err)
	//assert.True(t, len(normalTxs) > 100000)
}
