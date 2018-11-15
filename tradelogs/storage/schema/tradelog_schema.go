package schema

// TradeLogSchemaFieldName define a list of field names for a TradeLog record
//go:generate stringer -type=TradeLogSchemaFieldName -linecomment
type TradeLogSchemaFieldName int

const (
	//BlockNumber is enumerated field name for TradeLog.BlockNumbers
	BlockNumber TradeLogSchemaFieldName = iota //block_number
	//TxHash is enumerated field name for tradeLog.TxHash
	TxHash //tx_hash
	//EthReceivalSender is enumerated field name for TradeLog.EtherReceivalSender
	EthReceivalSender //eth_receival_sender
	//UserAddr is enumerated field name for TradeLog.UserAddr
	UserAddr //user_addr
	//SrcAddr is enumerated field name for TradeLog.SrcAddr
	SrcAddr //src_addr
	//DstAddr is enumerated field name for TradeLog.DstAddr
	DstAddr //dst_addr
	//Country is enumerated field name for TradeLog.CountryName
	Country //country
	//IP is enumerated field name for TradeLog.IP
	IP //ip
	//EthUSDProvider is the enumerated field name for TradeLog.ETHUSDProvider
	EthUSDProvider //eth_rate_provider
	//DstReserveAddr is enumerated fieldname for destination reserve Address
	DstReserveAddr // dst_rsv_addrs
	//SrcReserveAddr is enumerated field for source reserve Address
	SrcReserveAddr // src_rsv_addrs
	//EthReceivalAmount is the enumerated field for ETHReceivalAmount
	EthReceivalAmount // eth_receival_amount
	//SrcAmount is the enumerated field for source amount
	SrcAmount //src_amount
	//DstAmount is the enumerated field for source amount
	DstAmount //dst_amount
	//EthUSDRate is the enumerated field for ETH-USD rate
	EthUSDRate // eth_usd_rate
	//EthAmount is the enumerated field for ETH Amount
	EthAmount // eth_amount
	//FiatAmount is the enumerated field for fiat amount
	FiatAmount // fiat_amount
	//IntegrationApp is the name of apps integrated kyberswap
	IntegrationApp //integration_app
	//WalletAddress is the address of wallet associated with trade log
	WalletAddress //wallet_addr
	//LogIndex is the index of the log in that block
	LogIndex //log_index
)

//tradeLogSchemaFields translates the stringer of reserveRate fields into its enumerated form
var tradeLogSchemaFields = map[string]TradeLogSchemaFieldName{
	"block_number":        BlockNumber,
	"tx_hash":             TxHash,
	"eth_receival_sender": EthReceivalSender,
	"user_addr":           UserAddr,
	"src_addr":            SrcAddr,
	"dst_addr":            DstAddr,
	"country":             Country,
	"ip":                  IP,
	"eth_rate_provider":   EthUSDProvider,
	"dst_rsv_addr":        DstReserveAddr,
	"src_rsv_addr":        SrcReserveAddr,
	"eth_receival_amount": EthReceivalAmount,
	"src_amount":          SrcAmount,
	"dst_amount":          DstAmount,
	"eth_usd_rate":        EthUSDRate,
	"eth_amount":          EthAmount,
	"fiat_amount":         FiatAmount,
}
