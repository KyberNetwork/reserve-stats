package cq

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	firstTradedSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/first_traded"
	kycedschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/kyced"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
	walletStatSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/walletstats"
)

//CreateWalletStatsCqs return a new set of cqs required for wallet stats aggregation
func CreateWalletStatsCqs(dbName string) ([]*cq.ContinuousQuery, error) {
	var (
		result []*cq.ContinuousQuery
	)
	uniqueAddrCqs, err := cq.NewContinuousQuery(
		"wallet_unique_addr",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf(`SELECT COUNT(record) AS %[1]s INTO %[2]s FROM `+
			`(SELECT SUM(%[3]s) AS record FROM %[4]s GROUP BY %[5]s, %[6]s) GROUP BY %[7]s`,
			walletStatSchema.UniqueAddresses.String(),
			common.WalletStatsMeasurement,
			logSchema.EthAmount.String(),
			common.TradeLogMeasurementName,
			logSchema.UserAddr.String(),
			logSchema.WalletAddress.String(),
			walletStatSchema.WalletAddress.String(),
		),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, uniqueAddrCqs)
	volCqs, err := cq.NewContinuousQuery(
		"wallet_summary_volume",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf(`SELECT SUM(%[1]s) AS %[2]s, SUM(usd_amount) AS %[3]s, COUNT(%[1]s) AS %[4]s, `+
			`MEAN(usd_amount) AS %[5]s, MEAN(%[1]s) AS %[6]s INTO %[7]s FROM `+
			`(SELECT %[8]s, %[1]s, %[1]s*%[9]s AS usd_amount FROM %[10]s `+
			`WHERE (%[11]s!='%[12]s' AND %[13]s!='%[14]s') `+
			`OR (%[11]s!='%[14]s' AND %[13]s!='%[12]s') GROUP BY %[15]s) GROUP BY %[15]s`,
			logSchema.EthAmount.String(),
			walletStatSchema.ETHVolume.String(),
			walletStatSchema.USDVolume.String(),
			walletStatSchema.TotalTrade.String(),
			walletStatSchema.USDPerTrade.String(),
			walletStatSchema.ETHPerTrade.String(),
			common.WalletStatsMeasurement,
			logSchema.DstAmount.String(),
			logSchema.EthUSDRate.String(),
			common.TradeLogMeasurementName,
			logSchema.SrcAddr.String(),
			core.ETHToken.Address,
			logSchema.DstAddr.String(),
			core.WETHToken.Address,
			logSchema.WalletAddress.String(),
		),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, volCqs)

	kyced, err := cq.NewContinuousQuery(
		"wallet_kyced",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf(`SELECT COUNT(kyced) AS %[1]s INTO %[2]s `+
			`FROM (SELECT DISTINCT(%[3]s) AS kyced FROM %[4]s GROUP BY %[5]s, %[6]s) GROUP BY %[7]s`,
			walletStatSchema.KYCedAddresses.String(),
			common.WalletStatsMeasurement,
			kycedschema.KYCed.String(),
			common.KYCedMeasurementName,
			kycedschema.UserAddress.String(),
			kycedschema.WalletAddress.String(),
			walletStatSchema.WalletAddress.String(),
		),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, kyced)

	newUnqAddressCq, err := cq.NewContinuousQuery(
		"wallet_new_unique_addr",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf(`SELECT COUNT(%[1]s) AS %[2]s INTO %[3]s FROM %[4]s GROUP BY %[5]s`,
			firstTradedSchema.Traded.String(),
			walletStatSchema.NewUniqueAddresses.String(),
			common.WalletStatsMeasurement,
			common.FirstTradedMeasurementName,
			firstTradedSchema.WalletAddress.String(),
		),
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, newUnqAddressCq)

	totalBurnFeeCqs, err := cq.NewContinuousQuery(
		"wallet_total_burn_fee",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(src_burn_amount) + SUM(dst_burn_amount) AS total_burn_fee INTO wallet_stats FROM trades GROUP BY wallet_addr",
		"1d",
		supportedTimeZone(),
	)
	if err != nil {
		return nil, err
	}
	result = append(result, totalBurnFeeCqs)

	return result, nil
}
