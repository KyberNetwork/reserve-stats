package cq

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	libcq "github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	countryStatSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/country_stats"
	firstTradedSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/first_traded"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
)

// CreateCountryCqs return a set of cqs required for country trade aggregation
func CreateCountryCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var (
		result []*libcq.ContinuousQuery
	)
	uniqueAddrCqs, err := libcq.NewContinuousQuery(
		"country_unique_addr",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf(`SELECT COUNT(record) AS %[1]s INTO %[2]s FROM `+
			`(SELECT COUNT(%[3]s) AS record FROM %[4]s GROUP BY %[5]s) GROUP BY %[6]s`,
			countryStatSchema.UniqueAddresses.String(),
			common.CountryStatsMeasurementName,
			logSchema.EthAmount.String(),
			common.TradeLogMeasurementName,
			logSchema.UserAddr.String(),
			logSchema.Country.String(),
		),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, uniqueAddrCqs)
	volCqs, err := libcq.NewContinuousQuery(
		"summary_country_volume",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf(`SELECT SUM(%[1]s) AS %[2]s, SUM(usd_amount) AS %[3]s, COUNT(%[1]s) AS %[4]s, `+
			`MEAN(usd_amount) AS %[5]s, MEAN(%[1]s) AS %[6]s INTO %[7]s FROM `+
			`(SELECT %[8]s, %[1]s, %[1]s*%[9]s AS usd_amount FROM %[10]s `+
			`WHERE (%[11]s!='%[12]s' AND %[13]s!='%[14]s') `+
			`OR (%[11]s!='%[14]s' AND %[13]s!='%[12]s')) GROUP BY %[13]s`,
			logSchema.EthAmount.String(),
			countryStatSchema.TotalETHVolume.String(),
			countryStatSchema.TotalUSDAmount.String(),
			countryStatSchema.TotalTrade.String(),
			countryStatSchema.USDPerTrade.String(),
			countryStatSchema.ETHPerTrade.String(),
			common.CountryStatsMeasurementName,
			logSchema.DstAmount.String(),
			logSchema.EthUSDRate.String(),
			common.TradeLogMeasurementName,
			logSchema.SrcAddr.String(),
			core.ETHToken.Address,
			logSchema.DstAddr.String(),
			core.WETHToken.Address,
			logSchema.Country.String(),
		),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, volCqs)
	newUnqAddressCq, err := libcq.NewContinuousQuery(
		"new_country_unique_addr",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf("SELECT COUNT(%[1]s) as %[2]s INTO %[3]s FROM %[4]s GROUP BY %[5]s",
			firstTradedSchema.Traded.String(),
			countryStatSchema.NewUniqueAddresses.String(),
			common.CountryStatsMeasurementName,
			common.FirstTradedMeasurementName,
			logSchema.Country.String(),
		),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, newUnqAddressCq)

	assetVolDstDayCqs, err := libcq.NewContinuousQuery(
		"asset_country_volume_dst_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_country_stats FROM "+
			"(SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE "+
			"((src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR "+
			"(src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE'))) GROUP BY dst_addr, country",
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolDstDayCqs)

	assetVolSrcDayCqs, err := libcq.NewContinuousQuery(
		"asset_country_volume_src_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_country_stats FROM "+
			"(SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE "+
			"((src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR "+
			"(src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE'))) GROUP BY src_addr, country",
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolSrcDayCqs)

	kyced, err := libcq.NewContinuousQuery(
		"country_kyced",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf("SELECT COUNT(kyced) AS %[1]s INTO %[2]s FROM (SELECT DISTINCT(kyced) AS kyced FROM kyced GROUP BY user_addr, country) GROUP BY country",
			countryStatSchema.KYCedAddresses.String(),
			common.CountryStatsMeasurementName,
		),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, kyced)

	totalBurnFeeCqs, err := libcq.NewContinuousQuery(
		"country_total_burn_fee",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(src_burn_amount)+SUM(dst_burn_amount) AS total_burn_fee INTO country_stats FROM trades GROUP BY country",
		"1d",
		supportedTimeZone(),
	)

	if err != nil {
		return nil, err
	}
	result = append(result, totalBurnFeeCqs)

	return result, nil
}
