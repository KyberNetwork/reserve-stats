package cq

import (
	"bytes"
	"text/template"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	firstTradedSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influxstorage/schema/first_traded"
	kycedschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influxstorage/schema/kyced"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influxstorage/schema/tradelog"
	walletStatSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influxstorage/schema/walletstats"
)

func executeWalletStatsTemplate(templateString string) (string, error) {
	tmpl, err := template.New("walletStats").Parse(templateString)
	if err != nil {
		return "", err
	}
	var queryBuf bytes.Buffer
	if err = tmpl.Execute(&queryBuf, struct {
		UniqueAddresses            string
		WalletStatsMeasurementName string
		ETHAmount                  string
		TradeLogMeasurementName    string
		UserAddr                   string
		WalletAddr                 string
	}{
		UniqueAddresses:            walletStatSchema.UniqueAddresses.String(),
		WalletStatsMeasurementName: common.WalletStatsMeasurement,
		ETHAmount:                  logSchema.EthAmount.String(),
		TradeLogMeasurementName:    common.TradeLogMeasurementName,
		UserAddr:                   logSchema.UserAddr.String(),
		WalletAddr:                 logSchema.WalletAddress.String(),
	}); err != nil {
		return "", err
	}
	return queryBuf.String(), nil
}

//CreateWalletStatsCqs return a new set of cqs required for wallet stats aggregation
func CreateWalletStatsCqs(dbName string) ([]*cq.ContinuousQuery, error) {
	var (
		result []*cq.ContinuousQuery
	)
	walletStatsUniqueAddresses := `SELECT COUNT(record) AS {{.UniqueAddresses}} INTO {{.WalletStatsMeasurementName}} FROM ` +
		`(SELECT SUM({{.ETHAmount}}) AS record FROM {{.TradeLogMeasurementName}} GROUP BY {{.UserAddr}}, {{.WalletAddr}}) GROUP BY {{.WalletAddr}}`

	queryString, err := executeWalletStatsTemplate(walletStatsUniqueAddresses)
	if err != nil {
		return nil, err
	}

	uniqueAddrCqs, err := cq.NewContinuousQuery(
		"wallet_unique_addr",
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
	result = append(result, uniqueAddrCqs)

	walletStatsVolCqs := `SELECT SUM({{.ETHAmount}}) AS {{.ETHVolume}}, SUM(usd_amount) AS {{.USDVolume}}, COUNT({{.ETHAmount}}) AS {{.TotalTrade}}, ` +
		`MEAN(usd_amount) AS {{.USDPerTrade}}, MEAN({{.ETHAmount}}) AS {{.ETHPerTrade}} INTO {{.WalletStatsMeasurementName}} FROM ` +
		`(SELECT {{.ETHAmount}}, {{.ETHAmount}}*{{.ETHUSDRate}} AS usd_amount FROM {{.TradeLogMeasurementName}} ` +
		`WHERE ({{.SrcAddr}}!='{{.ETHTokenAddr}}' AND {{.DstAddr}}!='{{.WETHTokenAddr}}') ` +
		`OR ({{.SrcAddr}}!='{{.WETHTokenAddr}}' AND {{.DstAddr}}!='{{.ETHTokenAddr}}') GROUP BY {{.WalletAddr}}) GROUP BY {{.WalletAddr}}`

	tmpl, err := template.New("walletStatsQuery").Parse(walletStatsVolCqs)
	if err != nil {
		return nil, err
	}

	var queryBuf bytes.Buffer
	if err = tmpl.Execute(&queryBuf, struct {
		ETHAmount                  string
		ETHVolume                  string
		USDVolume                  string
		TotalTrade                 string
		USDPerTrade                string
		ETHPerTrade                string
		WalletStatsMeasurementName string
		ETHUSDRate                 string
		TradeLogMeasurementName    string
		SrcAddr                    string
		ETHTokenAddr               string
		DstAddr                    string
		WETHTokenAddr              string
		WalletAddr                 string
	}{
		ETHAmount:                  logSchema.EthAmount.String(),
		ETHVolume:                  walletStatSchema.ETHVolume.String(),
		USDVolume:                  walletStatSchema.USDVolume.String(),
		TotalTrade:                 walletStatSchema.TotalTrade.String(),
		USDPerTrade:                walletStatSchema.USDPerTrade.String(),
		ETHPerTrade:                walletStatSchema.ETHPerTrade.String(),
		WalletStatsMeasurementName: common.WalletStatsMeasurement,
		ETHUSDRate:                 logSchema.EthUSDRate.String(),
		TradeLogMeasurementName:    common.TradeLogMeasurementName,
		SrcAddr:                    logSchema.SrcAddr.String(),
		ETHTokenAddr:               core.ETHToken.Address,
		DstAddr:                    logSchema.DstAddr.String(),
		WETHTokenAddr:              core.WETHToken.Address,
		WalletAddr:                 logSchema.WalletAddress.String(),
	}); err != nil {
		return nil, err
	}

	volCqs, err := cq.NewContinuousQuery(
		"wallet_summary_volume",
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
	result = append(result, volCqs)

	kycedQueryTemplate := `SELECT COUNT(kyced) AS {{.KYCedAddresses}} INTO {{.WalletStatsMeasurementName}} ` +
		`FROM (SELECT DISTINCT({{.KYCed}}) AS kyced FROM {{.KYCedMeasurementName}} GROUP BY {{.UserAddr}}, ` +
		`{{.WalletAddr}}) GROUP BY {{.WalletAddr}}`

	tmpl, err = template.New("walletStatsQuery").Parse(kycedQueryTemplate)
	if err != nil {
		return nil, err
	}
	var kycedQueryBuf bytes.Buffer
	if err = tmpl.Execute(&kycedQueryBuf, struct {
		KYCedAddresses             string
		WalletStatsMeasurementName string
		KYCed                      string
		KYCedMeasurementName       string
		UserAddr                   string
		WalletAddr                 string
	}{
		KYCedAddresses:             walletStatSchema.KYCedAddresses.String(),
		WalletStatsMeasurementName: common.WalletStatsMeasurement,
		KYCed:                      kycedschema.KYCed.String(),
		KYCedMeasurementName:       common.KYCedMeasurementName,
		UserAddr:                   kycedschema.UserAddress.String(),
		WalletAddr:                 kycedschema.WalletAddress.String(),
	}); err != nil {
		return nil, err
	}

	kyced, err := cq.NewContinuousQuery(
		"wallet_kyced",
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

	newUniqueAddressCqTemplate := `SELECT COUNT({{.Traded}}) AS {{.NewUniqueAddresses}} INTO ` +
		`{{.WalletStatsMeasurementName}} FROM {{.FirstTradeMeasurementName}} GROUP BY {{.WalletAddr}}`
	tmpl, err = template.New("walletStatsQuery").Parse(newUniqueAddressCqTemplate)
	if err != nil {
		return nil, err
	}
	var newUniqueAddressQueryBuf bytes.Buffer
	if err = tmpl.Execute(&newUniqueAddressQueryBuf, struct {
		Traded                     string
		NewUniqueAddresses         string
		WalletStatsMeasurementName string
		FirstTradeMeasurementName  string
		WalletAddr                 string
	}{
		Traded:                     firstTradedSchema.Traded.String(),
		NewUniqueAddresses:         walletStatSchema.NewUniqueAddresses.String(),
		WalletStatsMeasurementName: common.WalletStatsMeasurement,
		FirstTradeMeasurementName:  common.FirstTradedMeasurementName,
		WalletAddr:                 firstTradedSchema.WalletAddress.String(),
	}); err != nil {
		return nil, err
	}

	newUnqAddressCq, err := cq.NewContinuousQuery(
		"wallet_new_unique_addr",
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

	totalBurnFeeCqTemplate := `SELECT SUM({{.SrcBurnAmount}})+SUM({{.DstBurnAmount}}) AS {{.TotalBurnFee}} INTO 
	{{.WalletStatsMeasurementName}} FROM {{.TradeMeasurementName}} GROUP BY {{.WalletAddr}}`
	tmpl, err = template.New("walletStatsQuery").Parse(totalBurnFeeCqTemplate)
	if err != nil {
		return nil, err
	}
	var totalBurnFeeQueryBuf bytes.Buffer
	if err = tmpl.Execute(&totalBurnFeeQueryBuf, struct {
		SrcBurnAmount              string
		DstBurnAmount              string
		TotalBurnFee               string
		WalletStatsMeasurementName string
		TradeMeasurementName       string
		WalletAddr                 string
	}{
		SrcBurnAmount:              logSchema.SourceBurnAmount.String(),
		DstBurnAmount:              logSchema.DestBurnAmount.String(),
		TotalBurnFee:               walletStatSchema.TotalBurnFee.String(),
		WalletStatsMeasurementName: common.WalletStatsMeasurement,
		TradeMeasurementName:       common.TradeLogMeasurementName,
		WalletAddr:                 logSchema.WalletAddress.String(),
	}); err != nil {
		return nil, err
	}

	totalBurnFeeCqs, err := cq.NewContinuousQuery(
		"wallet_total_burn_fee",
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
