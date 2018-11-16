package common

const (
	// DatabaseName is the InfluxDB database name to store trade events.
	DatabaseName = "trade_logs"
	//TradeLogMeasurementName is the measurement for trade log
	TradeLogMeasurementName = "trades"
	//BurnfeeMeasurementName is the measurement for burn fee
	BurnfeeMeasurementName = "burn_fees"
	//WalletMeasurementName is the measurement for wallet fee
	WalletMeasurementName = "wallet_fees"
)
