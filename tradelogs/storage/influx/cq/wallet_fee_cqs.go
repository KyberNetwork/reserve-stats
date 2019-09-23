package cq

import (
	"bytes"
	"text/template"

	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/tradelog"
)

func executeWalletFeeTemplate(templateString, walletFeeAmount, measurementName, reserveAddr string) (string, error) {
	tmpl, err := template.New("burnFee").Parse(templateString)
	if err != nil {
		return "", err
	}
	var queryBuf bytes.Buffer
	if err = tmpl.Execute(&queryBuf, struct {
		WalletFeeAmount      string
		MeasurementName      string
		TradeMeasurementName string
		ReserveAddress       string
		WalletAddress        string
	}{
		WalletFeeAmount:      walletFeeAmount,
		MeasurementName:      measurementName,
		TradeMeasurementName: common.TradeLogMeasurementName,
		ReserveAddress:       reserveAddr,
		WalletAddress:        logSchema.WalletAddress.String(),
	}); err != nil {
		return "", err
	}
	return queryBuf.String(), nil
}

// CreateWalletFeeCqs return a set of cqs required for burnfee aggregation
func CreateWalletFeeCqs(dbName string) ([]*cq.ContinuousQuery, error) {
	var (
		result []*cq.ContinuousQuery
	)
	queryTmpl := `SELECT SUM({{.WalletFeeAmount}}) as sum_amount INTO {{.MeasurementName}} 
	FROM {{.TradeMeasurementName}} WHERE {{.ReserveAddress}} != '' GROUP BY {{.ReserveAddress}}, {{.WalletAddress}}`

	queryString, err := executeWalletFeeTemplate(queryTmpl, logSchema.SourceWalletFeeAmount.String(),
		common.WalletFeeVolumeMeasurementHour, logSchema.SrcReserveAddr.String())

	if err != nil {
		return nil, err
	}

	srcWalletFeeHourCqs, err := cq.NewContinuousQuery(
		"src_wallet_fee_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		queryString,
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, srcWalletFeeHourCqs)

	queryString, err = executeWalletFeeTemplate(queryTmpl, logSchema.DestWalletFeeAmount.String(),
		common.WalletFeeVolumeMeasurementHour, logSchema.DstReserveAddr.String())

	if err != nil {
		return nil, err
	}

	dstWalletFeeHourCqs, err := cq.NewContinuousQuery(
		"dst_wallet_fee_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		queryString,
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, dstWalletFeeHourCqs)

	queryString, err = executeWalletFeeTemplate(queryTmpl, logSchema.SourceWalletFeeAmount.String(),
		common.WalletFeeVolumeMeasurementDay, logSchema.SrcReserveAddr.String())

	if err != nil {
		return nil, err
	}

	srcWalletFeeDayCqs, err := cq.NewContinuousQuery(
		"src_wallet_fee_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		queryString,
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, srcWalletFeeDayCqs)

	queryString, err = executeWalletFeeTemplate(queryTmpl, logSchema.DestWalletFeeAmount.String(),
		common.WalletFeeVolumeMeasurementDay, logSchema.DstReserveAddr.String())

	if err != nil {
		return nil, err
	}

	dstWalletFeeDayCqs, err := cq.NewContinuousQuery(
		"dst_wallet_fee_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		queryString,
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, dstWalletFeeDayCqs)

	return result, nil
}
