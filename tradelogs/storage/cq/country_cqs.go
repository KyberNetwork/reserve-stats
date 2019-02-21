package cq

import (
	"bytes"
	"text/template"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	libcq "github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	countryStatSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/country_stats"
	firstTradedSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/first_traded"
	heatMapSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/heatmap"
	kycedschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/kyced"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
)

func executeCountryVolumeTemplate(templateString, amountType, amountTypeAddr string) (string, error) {
	var (
		queryBuf bytes.Buffer
	)
	tmpl, err := template.New("countryVolume").Parse(templateString)
	if err != nil {
		return "", err
	}
	if err = tmpl.Execute(&queryBuf, struct {
		AmountType                    string
		TokenVolume                   string
		ETHAmount                     string
		ETHVolume                     string
		USDVolume                     string
		VolumeCountryStatsMeasurement string
		ETHUSDRate                    string
		TradeLogMeasurementName       string
		SrcAddr                       string
		ETHTokenAddr                  string
		DstAddr                       string
		WETHTokenAddr                 string
		AmountTypeAddr                string
		Country                       string
	}{
		AmountType:                    amountType,
		TokenVolume:                   heatMapSchema.TokenVolume.String(),
		ETHAmount:                     logSchema.EthAmount.String(),
		ETHVolume:                     heatMapSchema.ETHVolume.String(),
		USDVolume:                     heatMapSchema.USDVolume.String(),
		VolumeCountryStatsMeasurement: common.VolumeCountryStatsMeasurement,
		ETHUSDRate:                    logSchema.EthUSDRate.String(),
		TradeLogMeasurementName:       common.TradeLogMeasurementName,
		SrcAddr:                       logSchema.SrcAddr.String(),
		ETHTokenAddr:                  core.ETHToken.Address,
		DstAddr:                       logSchema.DstAddr.String(),
		WETHTokenAddr:                 core.WETHToken.Address,
		AmountTypeAddr:                amountTypeAddr,
		Country:                       logSchema.Country.String(),
	}); err != nil {
		return "", err
	}
	return queryBuf.String(), nil
}

// CreateCountryCqs return a set of cqs required for country trade aggregation
func CreateCountryCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var (
		result []*libcq.ContinuousQuery
	)
	uniqueAddrCqsTemplate := `SELECT COUNT(record) AS {{.UniqueAddresses}} INTO {{.CountryStatsMeasurementName}} FROM ` +
		`(SELECT COUNT({{.ETHAmount}}) AS record FROM {{.TradeLogMeasurementName}} GROUP BY {{.UserAddr}}) ` +
		`GROUP BY {{.Country}}`

	tmpl, err := template.New("uniqueAddr").Parse(uniqueAddrCqsTemplate)
	if err != nil {
		return nil, err
	}
	var queryBuf bytes.Buffer
	if err = tmpl.Execute(&queryBuf, struct {
		UniqueAddresses             string
		CountryStatsMeasurementName string
		ETHAmount                   string
		TradeLogMeasurementName     string
		UserAddr                    string
		Country                     string
	}{
		UniqueAddresses:             countryStatSchema.UniqueAddresses.String(),
		CountryStatsMeasurementName: common.CountryStatsMeasurementName,
		ETHAmount:                   logSchema.EthAmount.String(),
		TradeLogMeasurementName:     common.TradeLogMeasurementName,
		UserAddr:                    logSchema.UserAddr.String(),
		Country:                     logSchema.Country.String(),
	}); err != nil {
		return nil, err
	}

	uniqueAddrCqs, err := libcq.NewContinuousQuery(
		"country_unique_addr",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		queryBuf.String(),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, uniqueAddrCqs)

	volCqsTemplate := `SELECT SUM({{.ETHAmount}}) AS {{.TotalETHVolume}}, SUM(usd_amount) AS {{.TotalUSDAmount}}, ` +
		`COUNT({{.ETHAmount}}) AS {{.TotalTrade}}, MEAN(usd_amount) AS {{.USDPerTrade}}, ` +
		`MEAN({{.ETHAmount}}) AS {{.ETHPerTrade}} INTO {{.CountryStatsMeasurementName}} FROM ` +
		`(SELECT {{.ETHAmount}}, {{.ETHAmount}}*{{.ETHUSDRate}} AS usd_amount FROM {{.TradeLogMeasurementName}} ` +
		`WHERE (` + ethwethExcludingTemp + `)) GROUP BY {{.Country}}`

	tmpl, err = template.New("volCqs").Parse(volCqsTemplate)
	if err != nil {
		return nil, err
	}
	var volCqsQueryBuf bytes.Buffer
	if err = tmpl.Execute(&volCqsQueryBuf, struct {
		ETHAmount                   string
		TotalETHVolume              string
		TotalUSDAmount              string
		TotalTrade                  string
		USDPerTrade                 string
		ETHPerTrade                 string
		CountryStatsMeasurementName string
		ETHUSDRate                  string
		TradeLogMeasurementName     string
		SrcAddr                     string
		DstAddr                     string
		ETHTokenAddr                string
		WETHTokenAddr               string
		Country                     string
	}{
		ETHAmount:                   logSchema.EthAmount.String(),
		TotalETHVolume:              countryStatSchema.TotalETHVolume.String(),
		TotalUSDAmount:              countryStatSchema.TotalUSDAmount.String(),
		TotalTrade:                  countryStatSchema.TotalTrade.String(),
		USDPerTrade:                 countryStatSchema.USDPerTrade.String(),
		ETHPerTrade:                 countryStatSchema.ETHPerTrade.String(),
		CountryStatsMeasurementName: common.CountryStatsMeasurementName,
		ETHUSDRate:                  logSchema.EthUSDRate.String(),
		TradeLogMeasurementName:     common.TradeLogMeasurementName,
		SrcAddr:                     logSchema.SrcAddr.String(),
		DstAddr:                     logSchema.DstAddr.String(),
		ETHTokenAddr:                core.ETHToken.Address,
		WETHTokenAddr:               core.WETHToken.Address,
		Country:                     logSchema.Country.String(),
	}); err != nil {
		return nil, err
	}
	volCqs, err := libcq.NewContinuousQuery(
		"summary_country_volume",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		volCqsQueryBuf.String(),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, volCqs)

	newUniqueAddressCqTemplate := `SELECT COUNT({{.Traded}}) as {{.NewUniqueAddresses}} INTO ` +
		`{{.CountryStatsMeasurementName}} FROM {{.FirstTradeMeasurementName}} GROUP BY {{.Country}}`

	tmpl, err = template.New("newUniqueAddr").Parse(newUniqueAddressCqTemplate)
	if err != nil {
		return nil, err
	}
	var newUniqueAddressCqQueryBuf bytes.Buffer
	if err = tmpl.Execute(&newUniqueAddressCqQueryBuf, struct {
		Traded                      string
		NewUniqueAddresses          string
		CountryStatsMeasurementName string
		FirstTradeMeasurementName   string
		Country                     string
	}{
		Traded:                      firstTradedSchema.Traded.String(),
		NewUniqueAddresses:          countryStatSchema.NewUniqueAddresses.String(),
		CountryStatsMeasurementName: common.CountryStatsMeasurementName,
		FirstTradeMeasurementName:   common.FirstTradedMeasurementName,
		Country:                     firstTradedSchema.Country.String(),
	}); err != nil {
		return nil, err
	}
	newUnqAddressCq, err := libcq.NewContinuousQuery(
		"new_country_unique_addr",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		newUniqueAddressCqQueryBuf.String(),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, newUnqAddressCq)

	assetVolCqsTemplate := `SELECT SUM({{.AmountType}}) AS {{.TokenVolume}}, SUM({{.ETHAmount}}) AS {{.ETHVolume}}, ` +
		`SUM(usd_amount) AS {{.USDVolume}} INTO {{.VolumeCountryStatsMeasurement}} FROM ` +
		`(SELECT {{.AmountType}}, {{.ETHAmount}}, {{.ETHAmount}}*{{.ETHUSDRate}} AS usd_amount FROM ` +
		`{{.TradeLogMeasurementName}} WHERE (` + ethwethExcludingTemp + `)) GROUP BY {{.AmountTypeAddr}}, {{.Country}}`

	queryString, err := executeCountryVolumeTemplate(assetVolCqsTemplate, logSchema.DstAmount.String(),
		logSchema.DstAddr.String())
	if err != nil {
		return nil, err
	}

	assetVolDstDayCqs, err := libcq.NewContinuousQuery(
		"asset_country_volume_dst_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		queryString,
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolDstDayCqs)

	queryString, err = executeCountryVolumeTemplate(assetVolCqsTemplate, logSchema.SrcAmount.String(),
		logSchema.SrcAddr.String())
	if err != nil {
		return nil, err
	}

	assetVolSrcDayCqs, err := libcq.NewContinuousQuery(
		"asset_country_volume_src_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		queryString,
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolSrcDayCqs)

	kycedTemplate := "SELECT COUNT(kyced) AS {{.KYCedAddresses}} INTO {{.CountryStatsMeasurementName}} FROM " +
		"(SELECT DISTINCT({{.KYCed}}) AS kyced FROM {{.KYCedMeasurementName}} GROUP BY {{.UserAddr}}, {{.Country}}) " +
		"GROUP BY {{.Country}}"

	tmpl, err = template.New("kyced").Parse(kycedTemplate)
	if err != nil {
		return nil, err
	}
	var kycedQueryBuf bytes.Buffer
	if err = tmpl.Execute(&kycedQueryBuf, struct {
		KYCedAddresses              string
		CountryStatsMeasurementName string
		KYCed                       string
		KYCedMeasurementName        string
		UserAddr                    string
		Country                     string
	}{
		KYCedAddresses:              countryStatSchema.KYCedAddresses.String(),
		CountryStatsMeasurementName: common.CountryStatsMeasurementName,
		KYCed:                       kycedschema.KYCed.String(),
		KYCedMeasurementName:        common.KYCedMeasurementName,
		UserAddr:                    kycedschema.UserAddress.String(),
		Country:                     kycedschema.Country.String(),
	}); err != nil {
		return nil, err
	}

	kyced, err := libcq.NewContinuousQuery(
		"country_kyced",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		kycedQueryBuf.String(),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, kyced)

	totalBurnFeeTemplate := `SELECT SUM({{.SrcBurnAmount}})+SUM({{.DstBurnAmount}}) AS {{.TotalBurnFee}} INTO 
	{{.CountryStatMeasurement}} FROM {{.TradeMeasurementName}} GROUP BY {{.Country}}`

	tmpl, err = template.New("totalBurnFee").Parse(totalBurnFeeTemplate)
	if err != nil {
		return nil, err
	}

	var totalBurnFeeQueryBuf bytes.Buffer
	if err = tmpl.Execute(&totalBurnFeeQueryBuf, struct {
		SrcBurnAmount          string
		DstBurnAmount          string
		TotalBurnFee           string
		CountryStatMeasurement string
		TradeMeasurementName   string
		Country                string
	}{
		SrcBurnAmount:          logSchema.SourceBurnAmount.String(),
		DstBurnAmount:          logSchema.DestBurnAmount.String(),
		TotalBurnFee:           countryStatSchema.TotalBurnFee.String(),
		CountryStatMeasurement: common.CountryStatsMeasurementName,
		TradeMeasurementName:   common.TradeLogMeasurementName,
		Country:                logSchema.Country.String(),
	}); err != nil {
		return nil, err
	}

	totalBurnFeeCqs, err := libcq.NewContinuousQuery(
		"country_total_burn_fee",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		totalBurnFeeQueryBuf.String(),
		"1d",
		supportedTimeZone(),
	)

	if err != nil {
		return nil, err
	}
	result = append(result, totalBurnFeeCqs)

	return result, nil
}
