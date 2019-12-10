package influx

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/postprocessor"
	"github.com/stretchr/testify/assert"
)

func aggregationTestData(is *Storage) error {

	cqs, err := tradelogcq.CreateAssetVolumeCqs(is.dbName)
	if err != nil {
		return err
	}
	for _, cq := range cqs {
		err = cq.Execute(is.influxClient, is.sugar)
		if err != nil {
			return err
		}
	}
	return nil
}

func aggregationVolumeTestData(is *Storage) error {
	cqs, err := tradelogcq.CreateReserveVolumeCqs(is.dbName)
	if err != nil {
		return err
	}
	for _, cq := range cqs {
		err = cq.Execute(is.influxClient, is.sugar)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestGetAssetVolume(t *testing.T) {
	const (
		dbName = "test_volume"
		// These params are expected to be change when export.dat changes.
		fromTime    = 1539248043000
		toTime      = 1539248666000
		ethAmount   = 238.33849929550047
		totalVolume = 1.056174642648189277
		freq        = "h"
		timeStamp   = "2018-10-11T09:00:00Z"
		ethAddress  = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	)

	is, err := newTestInfluxStorage(dbName)
	assert.NoError(t, err)

	from := timeutil.TimestampMsToTime(fromTime)
	to := timeutil.TimestampMsToTime(toTime)

	defer func() {
		assert.NoError(t, is.tearDown())
	}()
	assert.NoError(t, loadTestData(dbName))
	assert.NoError(t, aggregationTestData(is))
	volume, err := is.GetAssetVolume(ethereum.HexToAddress(ethAddress), from, to, freq)
	assert.NoError(t, err)

	t.Logf("Volume result %v", volume)

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := volume[timeUint]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}

	require.Equal(t, ethAmount, result.USDAmount)
	require.Equal(t, totalVolume, result.Volume)
}

func TestGetReserveVolume(t *testing.T) {
	const (
		dbName = "test_rsv_volume"

		// These params are expected to be change when export.dat changes.
		fromTime    = 1539248043000
		toTime      = 1539248666000
		ethAmount   = 227.05539848662738
		totalVolume = 1.006174642648189232
		freq        = "h"
		timeStamp   = "2018-10-11T09:00:00Z"
		rsvAddrStr  = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
		ethAddress  = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	)

	is, err := newTestInfluxStorage(dbName)
	defer func() {
		if err := is.tearDown(); err != nil {
			t.Fatal(err)
		}
	}()
	if err := loadTestData(dbName); err != nil {
		t.Fatal(err)
	}
	if err := aggregationVolumeTestData(is); err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	from := timeutil.TimestampMsToTime(fromTime)
	to := timeutil.TimestampMsToTime(toTime)

	volume, err := is.GetReserveVolume(ethereum.HexToAddress(rsvAddrStr), ethereum.HexToAddress(ethAddress), from, to, freq)
	if err != nil {
		t.Fatal(err)
	}
	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	if err != nil {
		t.Fatal(err)
	}
	result, ok := volume[timeutil.TimeToTimestampMs(timeUnix)]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}

	require.Equal(t, ethAmount, result.USDAmount)
	require.Equal(t, totalVolume, result.Volume)

}

func TestGetMonthlyVolume(t *testing.T) {
	var (
		dbName = "test_monthly_volume"
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

	storage, err := NewInfluxStorage(sugar, dbName, influxClient, blockchain.NewMockTokenAmountFormatter(), blockchain.KNCAddr)
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
			ETHUSDRate:        100,
		},
	}

	actualVolume := &common.VolumeStats{
		ETHAmount: 1.9999902753074876,
		USDAmount: 199.99902753074878,
		Volume:    1.9999902753074876,
	}

	require.NoError(t, err)
	require.NoError(t, storage.SaveTradeLogs(tradeLogs))
	p := postprocessor.New(influxClient, sugar, dbName)
	require.NoError(t, p.Run(start, end))
	monthlyVolume, err := storage.GetMonthlyVolume(ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"), start, end)
	require.NoError(t, err)
	require.Equal(t, actualVolume, monthlyVolume[1554076800000])
}
