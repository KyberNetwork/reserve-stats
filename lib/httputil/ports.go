package httputil

// HTTPPort define custom type for port
type HTTPPort int

const (
	// IPLocatorPort is port number of IpLocator service
	IPLocatorPort HTTPPort = 8001

	// UsersPort is the port number of Users service.
	UsersPort HTTPPort = 8002

	// ReserveRatesPort is the port number of Reserve Rates API.
	ReserveRatesPort HTTPPort = 8003

	// TradeLogsPort is the port number of TradeLogs service.
	TradeLogsPort HTTPPort = 8004

	// GatewayPort is the port number of API gateway service.
	GatewayPort HTTPPort = 8005

	// PriceAnalytic is the port number of Price Analytics API
	PriceAnalytic HTTPPort = 8006

	// AppNames is the port number for Integration Address to AppNames API
	AppNames HTTPPort = 8007

	// UsersPublicPort is the port number for user stats public service
	UsersPublicPort HTTPPort = 8008

	// AccountingReserveAddressPort is the port number of accounting-reserve-addresses-api service.
	AccountingReserveAddressPort = 8009

	// AccountingCEXTradesPort is the port number of accounting-cex-trade-trades-api service.
	AccountingCEXTradesPort = 8010

	// AccountingTransactionsPort is the port number of accounting-reserve-transactions-api service.
	AccountingTransactionsPort = 8011

	// AccountingWalletErc20Port is the port number of accounting-wallet-erc20-api service
	AccountingWalletErc20Port = 8012

	//AccountingReserveTokensPort is the port number of accounting-reserve-tokens service
	AccountingReserveTokensPort = 8013

	//AccountingCEXWithdrawalsPort is the port number of accounting-listed-token service
	AccountingCEXWithdrawalsPort = 8014

	// AccountingReserveRatesPort is the port number of account-reserve-rates-api service
	AccountingReserveRatesPort = 8015

	// AccountingCEXDepositPort is the port number of accounting-listed-token service
	AccountingCEXDepositPort = 8016
	// Accounting0xTradesPort ...
	Accounting0xTradesPort = 8017
)
