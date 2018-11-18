package cq

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	libcq "github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	firstTradedSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/first_traded"
	tradeSumSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/trade_summary"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
)

const (
	//TradeSummaryMeasurement is the measurement to store trade summary
	TradeSummaryMeasurement = "trade_summary"
	//BurnFeeSummaryMeasurement is the measurement to storee burnfee summary
	BurnFeeSummaryMeasurement = "burn_fee_summary"
)

// CreateSummaryCqs return a set of cqs required for trade Summary aggregation
func CreateSummaryCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var result []*libcq.ContinuousQuery

	uniqueAddrCqs, err := libcq.NewContinuousQuery(
		"unique_addr",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf(
			"SELECT COUNT(record) AS %[1]s INTO %[2]s FROM (SELECT SUM(%[3]s) AS record FROM %[4]s GROUP BY %[5]s)",
			tradeSumSchema.UniqueAddresses.String(),
			TradeSummaryMeasurement,
			logSchema.EthAmount.String(),
			common.TradeLogMeasurementName,
			logSchema.UserAddr.String(),
		),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, uniqueAddrCqs)

	volCqs, err := libcq.NewContinuousQuery(
		"summary_volume",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf(
			"SELECT SUM(%[1]s) AS %[2]s, SUM(%[3]s) AS %[4]s, COUNT(%[1]s) AS %[5]s, MEAN(%[3]s) AS %[6]s, MEAN(%[1]s) AS %[7]s INTO %[8]s FROM (SELECT %[1]s, %[1]s*%[9]s AS %[3]s FROM %[10]s WHERE (%[11]s!='%[12]s' AND %[13]s!='%[14]s') OR (%[11]s!='%[14]s' AND %[13]s!='%[12]s'))",
			logSchema.EthAmount.String(),
			tradeSumSchema.TotalETHVolume.String(),
			logSchema.FiatAmount.String(),
			tradeSumSchema.TotalUSDAmount.String(),
			tradeSumSchema.TotalTrade.String(),
			tradeSumSchema.USDPerTrade.String(),
			tradeSumSchema.ETHPerTrade.String(),
			TradeSummaryMeasurement,
			logSchema.EthUSDRate.String(),
			common.TradeLogMeasurementName,
			logSchema.SrcAddr.String(),
			core.ETHToken.Address,
			logSchema.DstAddr.String(),
			core.WETHToken.Address,
		),
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

	newUnqAddressCq, err := libcq.NewContinuousQuery(
		"new_unique_addr",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf(
			"SELECT COUNT(%[1]s) as %[2]s INTO %[3]s FROM %[4]s",
			firstTradedSchema.Traded.String(),
			tradeSumSchema.NewUniqueAddresses.String(),
			TradeSummaryMeasurement,
			common.FirstTradedMeasurementName,
		),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, newUnqAddressCq)

	kyced, err := libcq.NewContinuousQuery(
		"kyced",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT COUNT(kyced) as kyced INTO trade_summary FROM (SELECT DISTINCT(kyced) AS kyced FROM kyced GROUP BY user_addr)",
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, kyced)

	return result, nil
}
