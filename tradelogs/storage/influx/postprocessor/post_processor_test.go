package postprocessor

import (
	"fmt"
	"math/big"
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx"
)

func TestPostProcessor_Run(t *testing.T) {
	var (
		dbName = "test_post_processor"
		start  = timeutil.TimestampMsToTime(1554076800000) // 2019-04-01 00:00:00 GMT
		end    = timeutil.TimestampMsToTime(1556668800000) // 2019-05-01 00:00:00 GMT
	)
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://127.0.0.1:8086",
	})
	require.NoError(t, err)
	defer func() {
		_, err := influxdb.QueryDB(influxClient, fmt.Sprintf("DROP DATABASE %s", dbName), dbName)
		assert.NoError(t, err)
	}()

	storage, err := influx.NewInfluxStorage(sugar, dbName, influxClient, blockchain.NewMockTokenAmountFormatter())
	require.NoError(t, err)

	tradeLogs := []common.TradeLog{
		{
			Timestamp:       timeutil.TimestampMsToTime(1554353231000),
			BlockNumber:     uint64(6100010),
			TransactionHash: ethereum.HexToHash("0x33dcdbed63556a1d90b7e0f626bfaf20f6f532d2ae8bf24c22abb15c4e1fff01"),
			TxSender:        ethereum.HexToAddress("0x63825c174ab367968ec60f061753d3bbd36a0d8f"),
			UserAddress:     ethereum.HexToAddress("0x85c5c26dc2af5546341fc1988b9d178148b4838b"),
			SrcAddress:      ethereum.HexToAddress("0xd0a4b8946cb52f0661273bfbc6fd0e0c75fc6433"),
			DestAddress:     ethereum.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"),
			SrcAmount:       big.NewInt(421371814779117936),
			DestAmount:      big.NewInt(999995137653743773),
			FiatAmount:      0,
			BurnFees: []common.BurnFee{
				{
					ReserveAddress: ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"),
					Amount:         big.NewInt(1427493059000719235),
				},
			},
			WalletFees:        nil,
			EthAmount:         big.NewInt(999995137653743773),
			OriginalEthAmount: big.NewInt(999995137653743773),
			SrcReserveAddress: ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"),
			DstReserveAddress: ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"),
		},
	}

	require.NoError(t, err)
	require.NoError(t, storage.SaveTradeLogs(tradeLogs))

	p := New(influxClient, sugar, dbName)

	volume, err := p.getVolumeData(start, end)
	require.NoError(t, err)
	assert.Contains(t, volume, "0x63825c174ab367968EC60f061753D3bbD36A0D8F")
	assert.NoError(t, p.writeReserveVolumeMonthly(start, volume))

	fee, err := p.getFeeData(start, end)
	require.NoError(t, err)
	assert.Contains(t, fee, "0x63825c174ab367968EC60f061753D3bbD36A0D8F")
	assert.NoError(t, p.writeReserveFeeMonthly(start, fee))

}
