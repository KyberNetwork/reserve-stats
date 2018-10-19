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
)
