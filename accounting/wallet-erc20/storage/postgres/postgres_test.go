package postgres

import (
	"math/big"
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestSaveAndGetAccountingRates(t *testing.T) {
	var (
		testWalletAddr   = ethereum.HexToAddress("0xbc33a1f908612640f2849b56b67a4de4d179c151")
		testContractAddr = ethereum.HexToAddress("0xdd974d5c2e2928dea5f71b9825b8b646686bd200")
		testToAddr       = ethereum.HexToAddress("0xc065900403f9ff07dfaecc0d08978c2e0dee1578")
		testHash         = ethereum.HexToHash("0x0ad1ba536ae3140769fb03715eb093daba2776e55e80b8f2783b0d16ae07f61d")
		testData         = []common.ERC20Transfer{
			{
				BlockNumber:     5155487,
				Timestamp:       timeutil.TimestampMsToTime(1553757140000),
				Hash:            testHash,
				From:            testWalletAddr,
				ContractAddress: testContractAddr,
				To:              testToAddr,
				Value:           big.NewInt(1111),
				Gas:             2222,
				GasUsed:         3333,
				GasPrice:        big.NewInt(4444),
			},
		}
	)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	sugar := logger.Sugar()

	_, db := testutil.MustNewDevelopmentDB()
	//Test if the constraint unique works...
	for i := 0; i < 10; i++ {
		td := testData[0]

		testData = append(testData, td)
	}

	wdb, err := NewDB(sugar, db, WithTableName("test_wallet_erc20"))
	require.NoError(t, err)

	defer func() {
		err := wdb.TearDown()
		require.NoError(t, err)
		err = wdb.Close()
		require.NoError(t, err)
	}()

	//TODO: update here

	data, err := wdb.GetERC20Transfers(testWalletAddr, testContractAddr, timeutil.TimestampMsToTime(1553757130000), timeutil.TimestampMsToTime(1553757150000))
	require.NoError(t, err)
	assert.Equal(t, 1, len(data))
	assert.Equal(t, testData[0].Gas, data[0].Gas)
}
