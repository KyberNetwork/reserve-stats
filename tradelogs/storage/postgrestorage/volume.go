package postgrestorage

import (
	"bytes"
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"
	"html/template"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
)

// Get aggregated Burn fee by hour or day
func (tldb *TradeLogDB) GetAssetVolume(token ethereum.Address, fromTime, toTime time.Time,
	frequency string) (map[uint64]*common.VolumeStats, error) {
	var (
		resultBuffer              bytes.Buffer
		ok                        bool
		dateParam                 string
		timeCondition             string
		ethWETHExcludingCondition string
		queryTmpl                 = `SELECT time, SUM(token_volume) token_volume, SUM(eth_amount) eth_volume,
		SUM(eth_amount * eth_usd_rate) usd_volume
		FROM (
		SELECT date_trunc('{{.DateParam}}',"timestamp") AS time, src_amount token_volume, eth_amount, eth_usd_rate
		FROM "` + schema.TradeLogsTableName + `" 
		WHERE EXISTS (SELECT NULL FROM ` + schema.TokenTableName + ` WHERE address = $1 AND id=src_address_id)
		AND {{.TimeCondition}} AND {{.EthWETHExcludingCondition}}
		UNION ALL
		SELECT date_trunc('{{.DateParam}}',"timestamp") AS time, dst_amount token_volume, eth_amount, eth_usd_rate
		FROM "` + schema.TradeLogsTableName + `"
		WHERE EXISTS (SELECT NULL FROM ` + schema.TokenTableName + ` WHERE address = $1 AND id=dst_address_id)
		AND {{.TimeCondition}} AND {{.EthWETHExcludingCondition}}
		) a GROUP BY time`
	)
	logger := tldb.sugar.With("func", "tradelogs/storage/postgrestorage/TradeLogDB.GetAssetVolume",
		"from", fromTime, "to", toTime, "frequency", frequency)
	tmpl, err := template.New("asset volume template").Parse(queryTmpl)
	if err != nil {
		return nil, err
	}
	if timeCondition, err = schema.BuildTimeCondition(fromTime, toTime, frequency); err != nil {
		return nil, err
	}
	if ethWETHExcludingCondition, err = schema.BuildEthWethExcludingCondition(); err != nil {
		return nil, err
	}
	if dateParam, ok = schema.DateFunctionParams[frequency]; !ok {
		return nil, fmt.Errorf("invalid frequency %s", frequency)
	}

	err = tmpl.Execute(&resultBuffer, struct {
		DateParam                 string
		TimeCondition             template.HTML
		EthWETHExcludingCondition template.HTML
	}{
		DateParam:                 dateParam,
		TimeCondition:             template.HTML(timeCondition),
		EthWETHExcludingCondition: template.HTML(ethWETHExcludingCondition),
	})
	if err != nil {
		return nil, err
	}

	var datas []struct {
		TokenVolume float64   `db:"token_volume"`
		EthVolume   float64   `db:"eth_volume"`
		USDVolume   float64   `db:"usd_volume"`
		Time        time.Time `db:"time"`
	}
	logger.Debugw("execute template successful", "prepare statement", resultBuffer.String())

	//fmt.Println(len(datas), "\t", timeCondition, resultBuffer.String())

	err = tldb.db.Select(&datas, resultBuffer.String(), token.Hex())
	if err != nil {
		return nil, err
	}

	if len(datas) == 0 {
		logger.Debugw("return empty result", "prepare statement", resultBuffer.String())
		return nil, nil
	}
	result := make(map[uint64]*common.VolumeStats)
	for _, data := range datas {
		fmt.Println(data.TokenVolume, " ", data.USDVolume)
		result[timeutil.TimeToTimestampMs(data.Time)] = &common.VolumeStats{
			Volume:    data.TokenVolume,
			ETHAmount: data.EthVolume,
			USDAmount: data.USDVolume,
		}
	}
	return result, nil
}

//TODO implement this
func (tldb *TradeLogDB) GetReserveVolume(rsvAddr ethereum.Address, token ethereum.Address,
	fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error) {
	var (
		resultBuffer              bytes.Buffer
		ok                        bool
		dateParam                 string
		timeCondition             string
		ethWETHExcludingCondition string
	)

	logger := tldb.sugar.With("func", "tradelogs/storage/postgrestorage/TradeLogDB.GetReserveVolume",
		"from", fromTime, "to", toTime, "frequency", frequency)

	tmpl, err := template.New("asset volume template").Parse(reserveQueryTmpl)
	if err != nil {
		return nil, err
	}
	if timeCondition, err = schema.BuildTimeCondition(fromTime, toTime, frequency); err != nil {
		return nil, err
	}
	if ethWETHExcludingCondition, err = schema.BuildEthWethExcludingCondition(); err != nil {
		return nil, err
	}
	if dateParam, ok = schema.DateFunctionParams[frequency]; !ok {
		return nil, fmt.Errorf("invalid frequency %s", frequency)
	}

	err = tmpl.Execute(&resultBuffer, struct {
		DateParam                 string
		TimeCondition             template.HTML
		EthWETHExcludingCondition template.HTML
	}{
		DateParam:                 dateParam,
		TimeCondition:             template.HTML(timeCondition),
		EthWETHExcludingCondition: template.HTML(ethWETHExcludingCondition),
	})
	if err != nil {
		return nil, err
	}

	var datas []struct {
		TokenVolume float64   `db:"token_volume"`
		EthVolume   float64   `db:"eth_volume"`
		USDVolume   float64   `db:"usd_volume"`
		Time        time.Time `db:"time"`
	}
	logger.Debugw("execute template successful", "prepare statement", resultBuffer.String())

	err = tldb.db.Select(&datas, resultBuffer.String(), token.Hex(), rsvAddr.Hex())
	if err != nil {
		return nil, err
	}

	if len(datas) == 0 {
		logger.Debugw("return empty result", "prepare statement", resultBuffer.String())
		return nil, nil
	}
	result := make(map[uint64]*common.VolumeStats)
	for _, data := range datas {
		fmt.Println(data.TokenVolume, " ", data.USDVolume)
		result[timeutil.TimeToTimestampMs(data.Time)] = &common.VolumeStats{
			Volume:    data.TokenVolume,
			ETHAmount: data.EthVolume,
			USDAmount: data.USDVolume,
		}
	}
	return result, nil
}

const reserveQueryTmpl = `SELECT time, SUM(token_volume) token_volume, SUM(eth_amount) eth_volume,
		SUM(eth_amount * eth_usd_rate) usd_volume
		FROM (
		SELECT date_trunc('{{.DateParam}}',"timestamp") AS time, src_amount token_volume, eth_amount, eth_usd_rate
		FROM "` + schema.TradeLogsTableName + `" 
		WHERE EXISTS (SELECT NULL FROM ` + schema.TokenTableName + ` WHERE address = $1 AND id=src_address_id)
		AND EXISTS (SELECT NULL FROM ` + schema.ReserveTableName + ` WHERE address = $2 AND (id= src_reserve_address_id OR id = dst_reserve_address_id))
		AND {{.TimeCondition}} AND {{.EthWETHExcludingCondition}}
		UNION ALL
		SELECT date_trunc('{{.DateParam}}',"timestamp") AS time, dst_amount token_volume, eth_amount, eth_usd_rate
		FROM "` + schema.TradeLogsTableName + `"
		WHERE EXISTS (SELECT NULL FROM ` + schema.TokenTableName + ` WHERE address = $1 AND id=dst_address_id)
		AND EXISTS (SELECT NULL FROM ` + schema.ReserveTableName + ` WHERE address = $2 AND (id= src_reserve_address_id OR id = dst_reserve_address_id))
		AND {{.TimeCondition}} AND {{.EthWETHExcludingCondition}}
		) a GROUP BY time`
