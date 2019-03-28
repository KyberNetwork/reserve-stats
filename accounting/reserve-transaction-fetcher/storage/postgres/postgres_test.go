package postgres

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/common"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestNormalTx(t *testing.T) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	_, db := testutil.MustNewDevelopmentDB()
	s, err := NewStorage(sugar, db, WithTableName(&tableNames{
		Normal:   "normal_test_normal_tx",
		Internal: "internal_test_normal_tx",
		ERC20:    "erc20_test_normal_tx",
	}))
	require.NoError(t, err)

	defer func(t *testing.T) {
		require.NoError(t, s.TearDown())
	}(t)

	txTimestamp := timeutil.TimestampMsToTime(1439048640 * 1000).UTC()
	txVal, ok := big.NewInt(0).SetString("11901464239480000000000000", 10)
	require.True(t, ok)
	txGasPrice, ok := big.NewInt(0).SetString("10000000000000", 10)
	require.True(t, ok)

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
	err = s.StoreNormalTx(testTxs)
	require.NoError(t, err)
	txs, err := s.GetNormalTx(txTimestamp.Add(-time.Second), txTimestamp.Add(time.Second*10))
	require.NoError(t, err)
	assert.Equal(t, testTxs, txs)

	// make sure we can safely insert duplicated transaction with value changed
	testTxs[0].Gas++
	testTxs[1].Gas++
	err = s.StoreNormalTx(testTxs)
	require.NoError(t, err)
	txs, err = s.GetNormalTx(txTimestamp.Add(-time.Second), txTimestamp.Add(time.Second*10))
	require.NoError(t, err)
	assert.Equal(t, testTxs, txs)

	txs, err = s.GetNormalTx(txTimestamp.Add(time.Second*2), txTimestamp.Add(time.Second*3))
	require.NoError(t, err)
	assert.Len(t, txs, 0)
}
