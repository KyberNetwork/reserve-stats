package influx

import (
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/stretchr/testify/require"

	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/cq"
)

func getUserListTestData(is *Storage) error {
	cqs, err := tradelogcq.CreateUserVolumeCqs(is.dbName)
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

func TestInfluxStorage_GetUserList(t *testing.T) {
	const (
		dbName   = "test_user_list"
		fromTime = 1539248043000
		toTime   = 1539248666000
		timezone = 1
	)

	is, err := newTestInfluxStorage(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, is.tearDown())
	}()
	if err := loadTestData(dbName); err != nil {
		t.Fatal(err)
	}

	require.NoError(t, loadTestData(dbName))
	require.NoError(t, getUserListTestData(is))

	from := timeutil.TimestampMsToTime(fromTime)
	to := timeutil.TimestampMsToTime(toTime)

	users, err := is.GetUserList(from, to, 0)
	require.NoError(t, err)
	require.Contains(t, users, common.UserInfo{
		Addr:      "0x8fA07F46353A2B17E92645592a94a0Fc1CEb783F",
		ETHVolume: 0.003353633937958248, // change from 0.0022361552369382478 because eth_amount is doubled for burnable token trade
		USDVolume: 0.7567877959608185,   // change from 0.5046152992532744 because eth_amount id double for burnable token trade
	})

	users, err = is.GetUserList(from, to, timezone)
	require.NoError(t, err)
	for _, user := range users {
		t.Log(user)
	}
	//require.Contains(t, users, common.UserInfo{
	//	Addr: "0x8fA07F46353A2B17E92645592a94a0Fc1CEb783F",
	//	ETHVolume: 0.0022361552369382478,
	//	USDVolume: 0.5046152992532744,
	//})
}
