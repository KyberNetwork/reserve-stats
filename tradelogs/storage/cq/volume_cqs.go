package cq

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/KyberNetwork/reserve-stats/lib/core"
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
)

var rsvVolTemplate = fmt.Sprintf(`SELECT SUM({{.AmountType}}) AS %[1]s, SUM(%[2]s) AS %[3]s, SUM(usd_amount) AS %[4]s `+
	`INTO {{.MeasurementName}} FROM `+
	`(SELECT {{.AmountType}}, %[2]s, %[2]s*%[5]s AS usd_amount FROM %[6]s WHERE`+
	`((%[7]s!='%[8]s' AND %[9]s!='%[10]s') OR `+
	`(%[7]s!='%[10]s' AND %[9]s!='%[8]s')) `+
	`AND {{.RsvAddressType}}!='') GROUP BY {{.AddressType}},{{.RsvAddressType}}`,
	volSchema.TokenVolume.String(),
	logSchema.EthAmount.String(),
	volSchema.ETHVolume.String(),
	volSchema.USDVolume.String(),
	logSchema.EthUSDRate.String(),
	common.TradeLogMeasurementName,
	logSchema.SrcAddr.String(),
	core.ETHToken.Address,
	logSchema.DstAddr.String(),
	core.WETHToken.Address,
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

// CreateAssetVolumeCqs return a set of cqs required for asset volume aggregation
func CreateAssetVolumeCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var (
		result []*libcq.ContinuousQuery
	)
	assetVolDstHourCqs, err := libcq.NewContinuousQuery(
		"asset_volume_dst_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		fmt.Sprintf(`SELECT SUM(%[1]s) AS %[2]s, SUM(%[3]s) AS %[4]s, SUM(usd_amount) AS %[5]s INTO %[6]s `+
			`FROM (SELECT %[1]s, %[3]s, %[3]s*%[7]s AS usd_amount FROM %[8]s WHERE `+
			`((%[9]s!='%[10]s' AND %[11]s!='%[12]s') OR `+
			`(%[9]s!='%[12]s' AND %[11]s!='%[10]s'))) GROUP BY %[11]s`,
			logSchema.DstAmount.String(),
			volSchema.TokenVolume.String(),
			logSchema.EthAmount.String(),
			volSchema.ETHVolume.String(),
			volSchema.USDVolume.String(),
			common.VolumeHourMeasurementName,
			logSchema.EthUSDRate.String(),
			common.TradeLogMeasurementName,
			logSchema.SrcAddr.String(),
			core.ETHToken.Address,
			logSchema.DstAddr.String(),
			core.WETHToken.Address,
		),
		"1h",
		nil,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolDstHourCqs)
	assetVolSrcHourCqs, err := libcq.NewContinuousQuery(
		"asset_volume_src_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		fmt.Sprintf(`SELECT SUM(%[1]s) AS %[2]s, SUM(%[3]s) AS %[4]s, SUM(usd_amount) AS %[5]s INTO %[6]s `+
			`FROM (SELECT %[1]s, %[3]s, %[3]s*%[7]s AS usd_amount FROM %[8]s WHERE `+
			`((%[9]s!='%[10]s' AND %[11]s!='%[12]s') OR `+
			`(%[9]s!='%[12]s' AND %[11]s!='%[10]s'))) GROUP BY %[9]s`,
			logSchema.SrcAmount.String(),
			volSchema.TokenVolume.String(),
			logSchema.EthAmount.String(),
			volSchema.ETHVolume.String(),
			volSchema.USDVolume.String(),
			common.VolumeHourMeasurementName,
			logSchema.EthUSDRate.String(),
			common.TradeLogMeasurementName,
			logSchema.SrcAddr.String(),
			core.ETHToken.Address,
			logSchema.DstAddr.String(),
			core.WETHToken.Address,
		),

		"1h",
		nil,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolSrcHourCqs)
	assetVolDstDayCqs, err := libcq.NewContinuousQuery(
		"asset_volume_dst_day",
		dbName,
		"1h",
		"2d",
		fmt.Sprintf(`SELECT SUM(%[1]s) AS %[2]s, SUM(%[3]s) AS %[4]s, SUM(usd_amount) AS %[5]s INTO %[6]s `+
			`FROM (SELECT %[1]s, %[3]s, %[3]s*%[7]s AS usd_amount FROM %[8]s WHERE `+
			`((%[9]s!='%[10]s' AND %[11]s!='%[12]s') OR `+
			`(%[9]s!='%[12]s' AND %[11]s!='%[10]s'))) GROUP BY %[11]s`,
			logSchema.DstAmount.String(),
			volSchema.TokenVolume.String(),
			logSchema.EthAmount.String(),
			volSchema.ETHVolume.String(),
			volSchema.USDVolume.String(),
			common.VolumeDayMeasurementName,
			logSchema.EthUSDRate.String(),
			common.TradeLogMeasurementName,
			logSchema.SrcAddr.String(),
			core.ETHToken.Address,
			logSchema.DstAddr.String(),
			core.WETHToken.Address,
		),
		"1d",
		nil,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolDstDayCqs)

	assetVolSrcDayCqs, err := libcq.NewContinuousQuery(
		"asset_volume_src_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf(`SELECT SUM(%[1]s) AS %[2]s, SUM(%[3]s) AS %[4]s, SUM(usd_amount) AS %[5]s INTO %[6]s `+
			`FROM (SELECT %[1]s, %[3]s, %[3]s*%[7]s AS usd_amount FROM %[8]s WHERE `+
			`((%[9]s!='%[10]s' AND %[11]s!='%[12]s') OR `+
			`(%[9]s!='%[12]s' AND %[11]s!='%[10]s'))) GROUP BY %[9]s`,
			logSchema.SrcAmount.String(),
			volSchema.TokenVolume.String(),
			logSchema.EthAmount.String(),
			volSchema.ETHVolume.String(),
			volSchema.USDVolume.String(),
			common.VolumeDayMeasurementName,
			logSchema.EthUSDRate.String(),
			common.TradeLogMeasurementName,
			logSchema.SrcAddr.String(),
			core.ETHToken.Address,
			logSchema.DstAddr.String(),
			core.WETHToken.Address,
		),
		"1d",
		nil,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolSrcDayCqs)
	return result, nil
}

//CreateUserVolumeCqs continueous query for aggregate user volume
func CreateUserVolumeCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var (
		result []*libcq.ContinuousQuery
	)
	userVolumeDayCqs, err := libcq.NewContinuousQuery(
		"user_volume_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume "+
			"INTO user_volume_day FROM (SELECT eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades) GROUP BY user_addr",
		"1d",
		nil,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, userVolumeDayCqs)
	userVolumeHourCqs, err := libcq.NewContinuousQuery(
		"user_volume_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(eth_amount) as eth_volume, SUM(usd_amount) as usd_volume "+
			"INTO user_volume_hour FROM (SELECT eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades) GROUP BY user_addr",
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
		RsvFieldsType
		MeasurementName string
	}{
		RsvFieldsType:   types,
		MeasurementName: mName,
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
				AmountType:     "src_amount",
				RsvAddressType: "src_rsv_addr",
				AddressType:    "src_addr"},
			"rsv_volume_src_dst": {
				AmountType:     "src_amount",
				RsvAddressType: "dst_rsv_addr",
				AddressType:    "src_addr"},
			"rsv_volume_dst_src": {
				AmountType:     "dst_amount",
				RsvAddressType: "src_rsv_addr",
				AddressType:    "dst_addr"},
			"rsv_volume_dst_dst": {
				AmountType:     "dst_amount",
				RsvAddressType: "dst_rsv_addr",
				AddressType:    "dst_addr"},
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
