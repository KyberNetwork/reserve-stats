package cq

import (
	"bytes"
	"text/template"

	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
)

const (
	// DayMeasurement is the measure to store aggregatedBurnFee in Day Frequency
	DayMeasurement = "burn_fee_day"
	// HourMeasurement is the measure to store aggregatedBurnFee in Hour Frequency
	HourMeasurement = "burn_fee_hour"
)

func executeBurnFeeTemplate(templateString, measurementName, burnAmount, address string) (string, error) {
	tmpl, err := template.New("burnFee").Parse(templateString)
	if err != nil {
		return "", err
	}
	var queryBuf bytes.Buffer
	if err = tmpl.Execute(&queryBuf, struct {
		BurnAmount           string
		MeasurementName      string
		TradeMeasurementName string
		Address              string
	}{
		BurnAmount:           burnAmount,
		MeasurementName:      measurementName,
		TradeMeasurementName: common.TradeLogMeasurementName,
		Address:              address,
	}); err != nil {
		return "", err
	}
	return queryBuf.String(), nil
}

// CreateBurnFeeCqs return a set of cqs required for burnfee aggregation
func CreateBurnFeeCqs(dbName string) ([]*cq.ContinuousQuery, error) {
	var (
		result []*cq.ContinuousQuery
	)

	queryTmpl := `SELECT SUM({{.BurnAmount}}) as sum_amount INTO {{.MeasurementName}} 
	FROM {{.TradeMeasurementName}} GROUP BY {{.Address}}`

	queryString, err := executeBurnFeeTemplate(queryTmpl, logSchema.SourceBurnAmount.String(),
		HourMeasurement, logSchema.SrcReserveAddr.String())

	if err != nil {
		return nil, err
	}

	srcBurnfeeHourCqs, err := cq.NewContinuousQuery(
		"src_burn_amount_hour",
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
	result = append(result, srcBurnfeeHourCqs)

	queryString, err = executeBurnFeeTemplate(queryTmpl, logSchema.DestBurnAmount.String(),
		HourMeasurement, logSchema.DstReserveAddr.String())

	if err != nil {
		return nil, err
	}

	dstBurnfeedstHourCqs, err := cq.NewContinuousQuery(
		"dst_burn_amount_hour",
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
	result = append(result, dstBurnfeedstHourCqs)

	queryString, err = executeBurnFeeTemplate(queryTmpl, logSchema.SourceBurnAmount.String(),
		DayMeasurement, logSchema.SrcReserveAddr.String())

	if err != nil {
		return nil, err
	}
	srcBurnfeeDayCqs, err := cq.NewContinuousQuery(
		"src_burn_amount_day",
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
	result = append(result, srcBurnfeeDayCqs)

	queryString, err = executeBurnFeeTemplate(queryTmpl, logSchema.DestBurnAmount.String(),
		DayMeasurement, logSchema.DstReserveAddr.String())

	if err != nil {
		return nil, err
	}
	dstBurnfeedstDayCqs, err := cq.NewContinuousQuery(
		"dst_burn_amount_day",
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
	result = append(result, dstBurnfeedstDayCqs)
	return result, nil
}
