package common

const (
	// DatabaseName is the InfluxDB database name to store trade events.
	DatabaseName = "trade_logs"
	// UnknownCountry is the special code for unknown country queries
	UnknownCountry = "UNKNOWN"
	//TradeLogMeasurementName is the measurement for trade log
	TradeLogMeasurementName = "trades"
	//BurnFeeMeasurementName is the measurement for burn fee
	BurnFeeMeasurementName = "burn_fees"
	//WalletMeasurementName is the measurement for wallet fee
	WalletMeasurementName = "wallet_fees"
	//FirstTradedMeasurementName is the measurement name for first traded records
	FirstTradedMeasurementName = "first_trades"
	//CountryStatsMeasurementName is the measurement name for country stat records
	CountryStatsMeasurementName = "country_stats"
	//KYCedMeasurementName is the measurement for kyc status
	KYCedMeasurementName = "kyced"
	// BurnFeeVolumeDayMeasurement is the measure to store aggregatedBurnFee in Day Frequency
	BurnFeeVolumeDayMeasurement = "burn_fee_day"
	// BurnFeeVolumeHourMeasurement is the measure to store aggregatedBurnFee in Hour Frequency
	BurnFeeVolumeHourMeasurement = "burn_fee_hour"
	//IntegrationVolumeMeasurement is the name for Integration Volume Measurement
	IntegrationVolumeMeasurement = "integration_volume"
	//TradeSummaryMeasurement is the measurement to store trade summary
	TradeSummaryMeasurement = "trade_summary"
	//BurnFeeSummaryMeasurement is the measurement to store burnfee summary
	BurnFeeSummaryMeasurement = "burn_fee_summary"
	//HeatMapMeasurement is the measurement to store country heatmap stats
	HeatMapMeasurement = "volume_country_stats"
	//WalletStatsMeasurement is the measurement name to which wallet stats is stored
	WalletStatsMeasurement = "wallet_stats"
	//WalletFeeVolumeMeasurementDay is the measurement to which wallet fee volume daily is stored
	WalletFeeVolumeMeasurementDay = "wallet_fee_day"
	//WalletFeeVolumeMeasurementHour is the measurement to which wallet fee volume hourly is stored
	WalletFeeVolumeMeasurementHour = "wallet_fee_hour"
	//VolumeHourMeasurementName is the measurement to which volume hourly is stored
	VolumeHourMeasurementName = "volume_hour"
	//VolumeDayMeasurementName is the measurement to which volume daily is stored
	VolumeDayMeasurementName = "volume_day"
	//UserVolumeDayMeasurementName is the measurement to which uservolume daily is stored
	UserVolumeDayMeasurementName = "user_volume_day"
	//UserVolumeHourMeasurementName is the measurement to which uservolume hourly is stored
	UserVolumeHourMeasurementName = "user_volume_hour"
)
