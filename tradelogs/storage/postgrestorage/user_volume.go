package postgrestorage

import (
	"bytes"
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"
	"html/template"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
)

func (tldb *TradeLogDB) GetUserVolume(userAddress ethereum.Address, from, to time.Time, freq string) (map[uint64]common.UserVolume, error) {
	var (
		dateFunctionParam string
		ok                bool
		queryBuffer       bytes.Buffer
		timeCondition     string
	)

	logger := tldb.sugar.With("from", from, "to", to, "userAddress", userAddress, "freq", freq)
	if dateFunctionParam, ok = schema.DateFunctionParams[strings.ToLower(freq)]; !ok {
		return nil, fmt.Errorf("invalid burn fee frequency %s", freq)
	}

	tpl, err := template.New("user volume template").Parse(userVolumeQueryTemplate)
	if err != nil {
		return nil, err
	}
	if timeCondition, err = schema.BuildTimeCondition(from, to, freq); err != nil {
		return nil, err
	}

	err = tpl.Execute(&queryBuffer, struct {
		DateParam     string
		TimeCondition template.HTML
	}{
		DateParam:     dateFunctionParam,
		TimeCondition: template.HTML(timeCondition),
	})
	if err != nil {
		return nil, err
	}
	logger.Debugw("prepare statement", "stmt", queryBuffer.String())

	var datas []struct {
		Time      time.Time `db:"time"`
		EthAmount float64   `db:"eth_volume"`
		UsdAmount float64   `db:"usd_volume"`
	}
	if err = tldb.db.Select(&datas, queryBuffer.String(), userAddress.Hex()); err != nil {
		return nil, err
	}

	result := make(map[uint64]common.UserVolume)
	for _, data := range datas {
		key := timeutil.TimeToTimestampMs(data.Time)
		result[key] = common.UserVolume{
			ETHAmount: data.EthAmount,
			USDAmount: data.UsdAmount,
		}
	}
	return result, nil
}

const userVolumeQueryTemplate = `
	SELECT date_trunc('{{.DateParam}}', "timestamp")as time, SUM(eth_amount) eth_volume,
		SUM(eth_amount * eth_usd_rate) usd_volume
	FROM "` + schema.TradeLogsTableName + `" a
	WHERE timestamp >= $1 AND timestamp <= $2 
	AND EXISTS (SELECT NULL FROM "` + schema.UserTableName + `" WHERE user_address_id=id and address = $1)
	GROUP BY time
`
