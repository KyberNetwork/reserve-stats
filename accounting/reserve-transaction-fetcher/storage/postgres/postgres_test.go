package postgres

import (
	"math/big"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestNormalTx(t *testing.T) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	db, teardown := testutil.MustNewDevelopmentDB()
	s, err := NewStorage(sugar, db)
	require.NoError(t, err)

	defer func(t *testing.T) {
		require.NoError(t, teardown())
	}(t)

	txTimestamp := timeutil.TimestampMsToTime(1439048640 * 1000).UTC()
	txVal, ok := big.NewInt(0).SetString("11901464239480000000000000", 10)
	require.True(t, ok)
	txGasPrice, ok := big.NewInt(0).SetString("10000000000000", 10)
	require.True(t, ok)

	err = s.StoreReserve(ethereum.HexToAddress("0x5abfec25f74cd88437631a7731906932776356f9"), common.Reserve.String())
	require.NoError(t, err)

	testTxs := []common.NormalTx{
		{
			BlockNumber: 54092,
			Timestamp:   txTimestamp,
			Hash:        "0x9c81f44c29ff0226f835cd0a8a2f2a7eca6db52a711f8211b566fd15d3e0e8d4",
			BlockHash:   "0xd3cabad6adab0b52eb632c386ea194036805713682c62cb589b5abcd76de2159",
			From:        "0x5abfec25f74cd88437631a7731906932776356f9",
			To:          "",
			Value:       txVal,
			Gas:         2000000,
			GasUsed:     1436963,
			GasPrice:    txGasPrice,
			IsError:     0,
		},
		{
			BlockNumber: 54092,
			Timestamp:   txTimestamp.Add(time.Second),
			Hash:        "0x98beb27135aa0a25650557005ad962919d6a278c4b3dde7f4f6a3a1e65aa746c",
			BlockHash:   "0x373d339e45a701447367d7b9c7cef84aab79c2b2714271b908cda0ab3ad0849b",
			From:        "0x3fb1cd2cd96c6d5c0b5eb3322d807b34482481d4",
			To:          "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
			Value:       txVal,
			Gas:         122261,
			GasUsed:     122207,
			GasPrice:    txGasPrice,
			IsError:     0,
		},
	}
	err = s.StoreNormalTx(testTxs, ethereum.HexToAddress("0x5abfec25f74cd88437631a7731906932776356f9"))
	require.NoError(t, err)
	txs, err := s.GetNormalTx(txTimestamp.Add(-time.Second), txTimestamp.Add(time.Second*10))
	require.NoError(t, err)
	assert.Equal(t, testTxs, txs)

	// make sure we can safely insert duplicated transaction with value changed
	testTxs[0].Gas++
	testTxs[1].Gas++
	err = s.StoreNormalTx(testTxs, ethereum.HexToAddress("0x5abfec25f74cd88437631a7731906932776356f9"))
	require.NoError(t, err)
	txs, err = s.GetNormalTx(txTimestamp.Add(-time.Second), txTimestamp.Add(time.Second*10))
	require.NoError(t, err)
	assert.Equal(t, testTxs, txs)

	txs, err = s.GetNormalTx(txTimestamp.Add(time.Second*2), txTimestamp.Add(time.Second*3))
	require.NoError(t, err)
	assert.Len(t, txs, 0)
}

func TestInternalTx(t *testing.T) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	db, teardown := testutil.MustNewDevelopmentDB()
	s, err := NewStorage(sugar, db)
	require.NoError(t, err)

	defer func(t *testing.T) {
		require.NoError(t, teardown())
	}(t)

	txTimestamp := timeutil.TimestampMsToTime(1477837690 * 1000).UTC()

	err = s.StoreReserve(ethereum.HexToAddress("0x5abfec25f74cd88437631a7731906932776356f9"), common.Reserve.String())
	require.NoError(t, err)

	testTxs := []common.InternalTx{
		{
			BlockNumber: 2535368,
			Timestamp:   txTimestamp,
			Hash:        "0x8a1a9989bda84f80143181a68bc137ecefa64d0d4ebde45dd94fc0cf49e70cb6",
			From:        "0x20d42f2e99a421147acf198d775395cac2e8b03d",
			To:          "",
			Value:       big.NewInt(0),
			Gas:         254791,
			GasUsed:     46750,
			IsError:     0,
		},
	}

	err = s.StoreInternalTx(testTxs, ethereum.HexToAddress("0x5abfec25f74cd88437631a7731906932776356f9"))
	require.NoError(t, err)
	txs, err := s.GetInternalTx(txTimestamp.Add(-time.Second), txTimestamp.Add(time.Second*10))
	require.NoError(t, err)
	assert.Equal(t, testTxs, txs)

	err = s.StoreInternalTx(testTxs, ethereum.HexToAddress("0x5abfec25f74cd88437631a7731906932776356f9"))
	require.NoError(t, err)
	txs, err = s.GetInternalTx(txTimestamp.Add(-time.Second), txTimestamp.Add(time.Second*10))
	require.NoError(t, err)
	assert.Equal(t, testTxs, txs)

	txs, err = s.GetInternalTx(txTimestamp.Add(time.Second*2), txTimestamp.Add(time.Second*3))
	require.NoError(t, err)
	assert.Len(t, txs, 0)
}

func TestERC20Transfer(t *testing.T) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	db, teardown := testutil.MustNewDevelopmentDB()
	s, err := NewStorage(sugar, db)
	require.NoError(t, err)

	defer func(t *testing.T) {
		require.NoError(t, teardown())
	}(t)

	txTimestamp := timeutil.TimestampMsToTime(1473433992 * 1000).UTC()
	txVal, ok := big.NewInt(0).SetString("101000000000000000000", 10)
	require.True(t, ok)

	err = s.StoreReserve(ethereum.HexToAddress("0x4e83362442b8d1bec281594cea3050c8eb01311c"), common.Reserve.String())
	require.NoError(t, err)

	testTxs := []common.ERC20Transfer{
		{
			BlockNumber:     2228258,
			Timestamp:       txTimestamp,
			Hash:            ethereum.HexToHash("0x5f2cd76fd3656686e356bc02cc91d8d0726a16936fd08e67ed30467053225a86"),
			From:            ethereum.HexToAddress("0x4e83362442b8d1bec281594cea3050c8eb01311c"),
			ContractAddress: ethereum.HexToAddress("0xecf8f87f810ecf450940c9f60066b4a7a501d6a7"),
			To:              ethereum.HexToAddress("0xac75b73394c329376c214663d92156afa864a77f"),
			Value:           txVal,
			Gas:             1000000,
			GasUsed:         93657,
			GasPrice:        big.NewInt(20000000000),
		},
		{
			BlockNumber:     2228258,
			Timestamp:       txTimestamp,
			Hash:            ethereum.HexToHash("0xf0ddb076798afba4fc507e96333aaa9a610b76b6774fdd94c10b78b06e18f9e6"),
			From:            ethereum.HexToAddress("0x4e83362442b8d1bec281594cea3050c8eb01311c"),
			ContractAddress: ethereum.HexToAddress("0xecf8f87f810ecf450940c9f60066b4a7a501d6a7"),
			To:              ethereum.HexToAddress("0xac75b73394c329376c214663d92156afa864a77f"),
			Value:           txVal,
			Gas:             1000000,
			GasUsed:         93657,
			GasPrice:        big.NewInt(20000000000),
		},
	}

	err = s.StoreERC20Transfer(testTxs, ethereum.HexToAddress("0x4e83362442b8d1bec281594cea3050c8eb01311c"))
	require.NoError(t, err)
	txs, err := s.GetERC20Transfer(txTimestamp.Add(-time.Second), txTimestamp.Add(time.Second*10))
	require.NoError(t, err)
	assert.Equal(t, testTxs, txs)

	err = s.StoreERC20Transfer(testTxs, ethereum.HexToAddress("0x4e83362442b8d1bec281594cea3050c8eb01311c"))
	require.NoError(t, err)
	txs, err = s.GetERC20Transfer(txTimestamp.Add(-time.Second), txTimestamp.Add(time.Second*10))
	require.NoError(t, err)
	assert.Equal(t, testTxs, txs)

	txs, err = s.GetERC20Transfer(txTimestamp.Add(time.Second*2), txTimestamp.Add(time.Second*3))
	require.NoError(t, err)
	assert.Len(t, txs, 0)
}

func TestLastInserted(t *testing.T) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	db, teardown := testutil.MustNewDevelopmentDB()
	s, err := NewStorage(sugar, db)
	require.NoError(t, err)

	defer func(t *testing.T) {
		require.NoError(t, teardown())
	}(t)
	err = s.StoreReserve(ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"), common.Reserve.String())
	require.NoError(t, err)

	var (
		testAddr         = ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F")
		testLastInserted = big.NewInt(7461105)
	)

	lastInserted, err := s.GetLastInserted(testAddr)
	require.NoError(t, err)
	assert.Nil(t, lastInserted)

	err = s.StoreLastInserted(testAddr, testLastInserted)
	require.NoError(t, err)

	lastInserted, err = s.GetLastInserted(testAddr)
	require.NoError(t, err)
	assert.Equal(t, 0, testLastInserted.Cmp(lastInserted))

	testLastInserted.Add(testLastInserted, big.NewInt(100))
	err = s.StoreLastInserted(testAddr, testLastInserted)
	require.NoError(t, err)

	lastInserted, err = s.GetLastInserted(testAddr)
	require.NoError(t, err)
	assert.Equal(t, 0, testLastInserted.Cmp(lastInserted))
}
