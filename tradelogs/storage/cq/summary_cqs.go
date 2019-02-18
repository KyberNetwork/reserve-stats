package cq

import (
	"bytes"
	"text/template"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	libcq "github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	firstTradedSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/first_traded"
	kycedschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/kyced"
	tradeSumSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/trade_summary"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
)

// CreateSummaryCqs return a set of cqs required for trade Summary aggregation
func CreateSummaryCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var result []*libcq.ContinuousQuery

	uniqueAddrCqsTemplate := "SELECT COUNT(record) AS {{.UniqueAddresses}} INTO {{.TradeSummaryMeasurementName}} FROM " +
		"(SELECT SUM({{.ETHAmount}}) AS record FROM {{.TradeLogMeasurementName}} GROUP BY {{.UserAddr}})"

	tmpl, err := template.New("uniqueAddr").Parse(uniqueAddrCqsTemplate)
	if err != nil {
		return nil, err
	}

	var queryBuf bytes.Buffer
	if err = tmpl.Execute(&queryBuf, struct {
		UniqueAddresses             string
		TradeSummaryMeasurementName string
		ETHAmount                   string
		TradeLogMeasurementName     string
		UserAddr                    string
	}{
		UniqueAddresses:             tradeSumSchema.UniqueAddresses.String(),
		TradeSummaryMeasurementName: common.TradeSummaryMeasurement,
		ETHAmount:                   logSchema.EthAmount.String(),
		TradeLogMeasurementName:     common.TradeLogMeasurementName,
		UserAddr:                    logSchema.UserAddr.String(),
	}); err != nil {
		return nil, err
	}

	uniqueAddrCqs, err := libcq.NewContinuousQuery(
		"unique_addr",
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

	volCqsTemplate := "SELECT SUM({{.ETHAmount}}) AS {{.TotalETHVolume}}, SUM({{.FiatAmount}}) AS {{.TotalUSDAmount}}, COUNT({{.ETHAmount}}) AS {{.TotalTrade}}, " +
		"MEAN({{.FiatAmount}}) AS {{.USDPerTrade}}, MEAN({{.ETHAmount}}) AS {{.ETHPerTrade}} INTO {{.TradeSummaryMeasurementName}} " +
		"FROM (SELECT {{.ETHAmount}}, {{.ETHAmount}}*{{.ETHUSDRate}} AS {{.FiatAmount}} FROM {{.TradeLogMeasurementName}} WHERE " +
		"({{.SrcAddr}}!='{{.ETHTokenAddr}}' AND {{.DstAddr}}!='{{.WETHTokenAddr}}') OR ({{.SrcAddr}}!='{{.WETHTokenAddr}}' AND {{.DstAddr}}!='{{.ETHTokenAddr}}'))"
	tmpl, err = template.New("uniqueAddr").Parse(volCqsTemplate)
	if err != nil {
		return nil, err
	}
	var volQueryBuf bytes.Buffer
	if err = tmpl.Execute(&volQueryBuf, struct {
		ETHAmount                   string
		TotalETHVolume              string
		FiatAmount                  string
		TotalUSDAmount              string
		TotalTrade                  string
		USDPerTrade                 string
		ETHPerTrade                 string
		TradeSummaryMeasurementName string
		ETHUSDRate                  string
		TradeLogMeasurementName     string
		SrcAddr                     string
		DstAddr                     string
		ETHTokenAddr                string
		WETHTokenAddr               string
	}{
		ETHAmount:                   logSchema.EthAmount.String(),
		TotalETHVolume:              tradeSumSchema.TotalETHVolume.String(),
		FiatAmount:                  logSchema.FiatAmount.String(),
		TotalUSDAmount:              tradeSumSchema.TotalUSDAmount.String(),
		TotalTrade:                  tradeSumSchema.TotalTrade.String(),
		USDPerTrade:                 tradeSumSchema.USDPerTrade.String(),
		ETHPerTrade:                 tradeSumSchema.ETHPerTrade.String(),
		TradeSummaryMeasurementName: common.TradeSummaryMeasurement,
		ETHUSDRate:                  logSchema.EthUSDRate.String(),
		TradeLogMeasurementName:     common.TradeLogMeasurementName,
		SrcAddr:                     logSchema.SrcAddr.String(),
		DstAddr:                     logSchema.DstAddr.String(),
		ETHTokenAddr:                core.ETHToken.Address,
		WETHTokenAddr:               core.WETHToken.Address,
	}); err != nil {
		return nil, err
	}
	volCqs, err := libcq.NewContinuousQuery(
		"summary_volume",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		volQueryBuf.String(),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, volCqs)

	totalBurnFeeCqs, err := libcq.NewContinuousQuery(
		"summary_total_burn_fee",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(src_burn_amount)+SUM(dst_burn_amount) AS total_burn_fee INTO burn_fee_summary FROM trades",
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, totalBurnFeeCqs)

	newUniqueAddressCqTemplate := "SELECT COUNT({{.Traded}}) as {{.NewUniqueAddresses}} INTO {{.TradeSummaryMeasurementName}} FROM {{.FirstTradeMeasurementName}}"
	tmpl, err = template.New("uniqueAddr").Parse(newUniqueAddressCqTemplate)
	if err != nil {
		return nil, err
	}
	var newUniqueAddressQueryBuf bytes.Buffer
	if err = tmpl.Execute(&newUniqueAddressQueryBuf, struct {
		Traded                      string
		NewUniqueAddresses          string
		TradeSummaryMeasurementName string
		FirstTradeMeasurementName   string
	}{
		Traded:                      firstTradedSchema.Traded.String(),
		NewUniqueAddresses:          tradeSumSchema.NewUniqueAddresses.String(),
		TradeSummaryMeasurementName: common.TradeSummaryMeasurement,
		FirstTradeMeasurementName:   common.FirstTradedMeasurementName,
	}); err != nil {
		return nil, err
	}
	newUnqAddressCq, err := libcq.NewContinuousQuery(
		"new_unique_addr",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		newUniqueAddressQueryBuf.String(),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, newUnqAddressCq)

	kycedTemplate := "SELECT COUNT(kyced) as {{.KYCedAddresses}} INTO {{.TradeSummaryMeasurementName}} FROM " +
		"(SELECT DISTINCT({{.KYCed}}) AS kyced FROM {{.KYCedMeasurementName}} GROUP BY {{.UserAddr}})"
	tmpl, err = template.New("uniqueAddr").Parse(kycedTemplate)
	if err != nil {
		return nil, err
	}
	var kycedQueryBuf bytes.Buffer
	if err = tmpl.Execute(&kycedQueryBuf, struct {
		KYCedAddresses              string
		TradeSummaryMeasurementName string
		KYCed                       string
		KYCedMeasurementName        string
		UserAddr                    string
	}{
		KYCedAddresses:              tradeSumSchema.KYCedAddresses.String(),
		TradeSummaryMeasurementName: common.TradeSummaryMeasurement,
		KYCed:                       kycedschema.KYCed.String(),
		KYCedMeasurementName:        common.KYCedMeasurementName,
		UserAddr:                    kycedschema.UserAddress.String(),
	}); err != nil {
		return nil, err
	}
	kyced, err := libcq.NewContinuousQuery(
		"kyced",
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

	return result, nil
}
