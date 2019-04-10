package cq

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	libcq "github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
	volSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/volume"
)

const (
	// the trades from WETH-ETH doesn't count. Hence the select clause skips every trade ETH-WETH or WETH-ETH
	// These trades are excluded by its src_addr and dst_addr, which is 0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE for ETH
	// and 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2 for WETH
	rsvVolHourMsmName = `rsv_volume_hour`
	rsvVolDayMsmName  = `rsv_volume_day`

	rsvVolTemplate = `SELECT SUM({{.AmountType}}) AS {{.TokenVolume}}, SUM({{.ETHAmount}}) AS {{.ETHVolume}}, SUM(usd_amount) AS {{.USDVolume}} ` +
		`INTO {{.MeasurementName}} FROM ` +
		`(SELECT {{.AmountType}}, {{.ETHAmount}}, {{.ETHAmount}}*{{.ETHUSDRate}} AS usd_amount FROM {{.TradeLogMeasurementName}} ` +
		` WHERE(` + ethWETHExcludingTemp + `) AND {{.RsvAddressType}}!='') GROUP BY {{.AddressType}},{{.RsvAddressType}}`
)

func supportedTimeZone() []string {
	timezone := []string{}
	for i := -11; i < 15; i++ {
		if i != 0 {
			timezone = append(timezone, fmt.Sprintf("%dh", i))
		}
	}
	return timezone
}

func executeAssetVolumeTemplate(stringTemplate, amountType, cqMeasurementName, amountTypeAddr string) (string, error) {
	tmpl, err := template.New("queryAssetVolume").Parse(stringTemplate)
	if err != nil {
		return "", err
	}
	var queryBuf bytes.Buffer
	if err = tmpl.Execute(&queryBuf, struct {
		AmountType              string
		TokenVolume             string
		ETHAmount               string
		ETHVolume               string
		USDVolume               string
		VolumeCqMeasurement     string
		ETHUSDRate              string
		TradeLogMeasurementName string
		SrcAddr                 string
		DstAddr                 string
		ETHTokenAddr            string
		WETHTokenAddr           string
		AmountTypeAddr          string
	}{
		AmountType:              amountType,
		TokenVolume:             volSchema.TokenVolume.String(),
		ETHAmount:               logSchema.EthAmount.String(),
		ETHVolume:               volSchema.ETHVolume.String(),
		USDVolume:               volSchema.USDVolume.String(),
		VolumeCqMeasurement:     cqMeasurementName,
		ETHUSDRate:              logSchema.EthUSDRate.String(),
		TradeLogMeasurementName: common.TradeLogMeasurementName,
		SrcAddr:                 logSchema.SrcAddr.String(),
		DstAddr:                 logSchema.DstAddr.String(),
		ETHTokenAddr:            blockchain.ETHAddr.Hex(),
		WETHTokenAddr:           blockchain.WETHAddr.Hex(),
		AmountTypeAddr:          amountTypeAddr,
	}); err != nil {
		return "", err
	}
	return queryBuf.String(), nil
}

// CreateAssetVolumeCqs return a set of cqs required for asset volume aggregation
func CreateAssetVolumeCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var (
		result []*libcq.ContinuousQuery
	)

	assetVolCqsTemplate := `SELECT SUM({{.AmountType}}) AS {{.TokenVolume}}, SUM({{.ETHAmount}}) AS {{.ETHVolume}}, ` +
		`SUM(usd_amount) AS {{.USDVolume}} INTO {{.VolumeCqMeasurement}}  ` +
		`FROM (SELECT {{.AmountType}}, {{.ETHAmount}}, {{.ETHAmount}}*{{.ETHUSDRate}} AS usd_amount FROM {{.TradeLogMeasurementName}} ` +
		`WHERE (` + ethWETHExcludingTemp + `)) GROUP BY {{.AmountTypeAddr}}`

	queryString, err := executeAssetVolumeTemplate(assetVolCqsTemplate, logSchema.DstAmount.String(),
		common.VolumeHourMeasurementName, logSchema.DstAddr.String())
	if err != nil {
		return nil, err
	}

	assetVolDstHourCqs, err := libcq.NewContinuousQuery(
		"asset_volume_dst_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		queryString,
		"1h",
		nil,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolDstHourCqs)

	queryString, err = executeAssetVolumeTemplate(assetVolCqsTemplate, logSchema.SrcAmount.String(),
		common.VolumeHourMeasurementName, logSchema.SrcAddr.String())
	if err != nil {
		return nil, err
	}

	assetVolSrcHourCqs, err := libcq.NewContinuousQuery(
		"asset_volume_src_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		queryString,
		"1h",
		nil,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolSrcHourCqs)

	queryString, err = executeAssetVolumeTemplate(assetVolCqsTemplate, logSchema.DstAmount.String(),
		common.VolumeDayMeasurementName, logSchema.DstAddr.String())
	if err != nil {
		return nil, err
	}

	assetVolDstDayCqs, err := libcq.NewContinuousQuery(
		"asset_volume_dst_day",
		dbName,
		"1h",
		"2d",
		queryString,
		"1d",
		nil,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolDstDayCqs)

	queryString, err = executeAssetVolumeTemplate(assetVolCqsTemplate, logSchema.SrcAmount.String(),
		common.VolumeDayMeasurementName, logSchema.SrcAddr.String())
	if err != nil {
		return nil, err
	}

	assetVolSrcDayCqs, err := libcq.NewContinuousQuery(
		"asset_volume_src_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		queryString,
		"1d",
		nil,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolSrcDayCqs)
	return result, nil
}

func executeUserVolumeTemplate(templateString, measurementName string) (string, error) {
	tmpl, err := template.New("queryUserVolume").Parse(templateString)
	if err != nil {
		return "", err
	}
	var queryBuf bytes.Buffer
	if err = tmpl.Execute(&queryBuf, struct {
		ETHAmount                 string
		ETHVolume                 string
		USDVolume                 string
		UserVolumeMeasurementName string
		ETHUSDRate                string
		TradeLogMeasurementName   string
		UserAddr                  string
	}{
		ETHAmount:                 logSchema.EthAmount.String(),
		ETHVolume:                 volSchema.ETHVolume.String(),
		USDVolume:                 volSchema.USDVolume.String(),
		UserVolumeMeasurementName: measurementName,
		ETHUSDRate:                logSchema.EthUSDRate.String(),
		TradeLogMeasurementName:   common.TradeLogMeasurementName,
		UserAddr:                  logSchema.UserAddr.String(),
	}); err != nil {
		return "", err
	}
	return queryBuf.String(), nil
}

//CreateUserVolumeCqs continueous query for aggregate user volume
func CreateUserVolumeCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var (
		result []*libcq.ContinuousQuery
	)
	userVolumeCqsQueryTemplate := `SELECT SUM({{.ETHAmount}}) AS {{.ETHVolume}}, SUM(usd_amount) AS {{.USDVolume}} ` +
		`INTO {{.UserVolumeMeasurementName}} FROM (SELECT {{.ETHAmount}}, {{.ETHAmount}}*{{.ETHUSDRate}} AS usd_amount ` +
		`FROM {{.TradeLogMeasurementName}}) GROUP BY {{.UserAddr}}`

	queryString, err := executeUserVolumeTemplate(userVolumeCqsQueryTemplate, common.UserVolumeDayMeasurementName)
	if err != nil {
		return nil, err
	}

	userVolumeDayCqs, err := libcq.NewContinuousQuery(
		"user_volume_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		queryString,
		"1d",
		nil,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, userVolumeDayCqs)

	queryString, err = executeUserVolumeTemplate(userVolumeCqsQueryTemplate, common.UserVolumeHourMeasurementName)
	if err != nil {
		return nil, err
	}

	userVolumeHourCqs, err := libcq.NewContinuousQuery(
		"user_volume_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		queryString,
		"1h",
		nil,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, userVolumeHourCqs)
	return result, nil
}

// RsvFieldsType declare the set of names requires to completed a reserveVolume Cqs
type RsvFieldsType struct {
	// AmountType: it can be dst_amount or src_amount
	AmountType string
	// RsvAddressType: it can be dst_rsv_amount or src_rsv_amount
	RsvAddressType string
	// Addresstype: it can be dst_addr or src_addr
	AddressType string
}

func renderRsvCqFromTemplate(tmpl *template.Template, mName string, types RsvFieldsType) (string, error) {
	var query bytes.Buffer
	err := tmpl.Execute(&query, struct {
		TokenVolume             string
		ETHAmount               string
		ETHVolume               string
		USDVolume               string
		ETHUSDRate              string
		TradeLogMeasurementName string
		SrcAddr                 string
		ETHTokenAddr            string
		DstAddr                 string
		WETHTokenAddr           string
		RsvFieldsType
		MeasurementName string
	}{
		TokenVolume:             volSchema.TokenVolume.String(),
		ETHAmount:               logSchema.EthAmount.String(),
		ETHVolume:               volSchema.ETHVolume.String(),
		USDVolume:               volSchema.USDVolume.String(),
		ETHUSDRate:              logSchema.EthUSDRate.String(),
		TradeLogMeasurementName: common.TradeLogMeasurementName,
		SrcAddr:                 logSchema.SrcAddr.String(),
		ETHTokenAddr:            blockchain.ETHAddr.Hex(),
		DstAddr:                 logSchema.DstAddr.String(),
		WETHTokenAddr:           blockchain.WETHAddr.Hex(),
		RsvFieldsType:           types,
		MeasurementName:         mName,
	})
	if err != nil {
		return "", err
	}
	return query.String(), nil
}

// CreateReserveVolumeCqs return a set of cqs required for asset volume aggregation
func CreateReserveVolumeCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var (
		result     []*libcq.ContinuousQuery
		cqsGroupBY = map[string]RsvFieldsType{
			"rsv_volume_src_src": {
				AmountType:     logSchema.SrcAmount.String(),
				RsvAddressType: logSchema.SrcReserveAddr.String(),
				AddressType:    logSchema.SrcAddr.String()},
			"rsv_volume_src_dst": {
				AmountType:     logSchema.SrcAmount.String(),
				RsvAddressType: logSchema.DstReserveAddr.String(),
				AddressType:    logSchema.SrcAddr.String()},
			"rsv_volume_dst_src": {
				AmountType:     logSchema.DstAmount.String(),
				RsvAddressType: logSchema.SrcReserveAddr.String(),
				AddressType:    logSchema.DstAddr.String()},
			"rsv_volume_dst_dst": {
				AmountType:     logSchema.DstAmount.String(),
				RsvAddressType: logSchema.DstReserveAddr.String(),
				AddressType:    logSchema.DstAddr.String()},
		}
	)

	tpml, err := template.New("cq.CreateReserveVolumeCqs").Parse(rsvVolTemplate)
	if err != nil {
		return nil, err
	}

	for name, types := range cqsGroupBY {
		query, err := renderRsvCqFromTemplate(tpml, rsvVolHourMsmName, types)
		if err != nil {
			return nil, err
		}
		hourCQ, err := libcq.NewContinuousQuery(
			name+"_hour",
			dbName,
			hourResampleInterval,
			hourResampleFor,
			query,
			"1h",
			nil,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, hourCQ)

		query, err = renderRsvCqFromTemplate(tpml, rsvVolDayMsmName, types)
		if err != nil {
			return nil, err
		}
		dayCQ, err := libcq.NewContinuousQuery(
			name+"_day",
			dbName,
			dayResampleInterval,
			dayResampleFor,
			query,
			"1d",
			nil,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, dayCQ)
	}

	return result, nil
}
