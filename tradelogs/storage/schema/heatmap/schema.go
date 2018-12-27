package heatmap

// FieldName define a list of field names for a token heatmap record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for time
	Time FieldName = iota //time
	//ETHVolume is the field name for eth volume
	ETHVolume //eth_volume
	//TokenVolume is the field name for token volume
	TokenVolume //token_volume
	//USDVolume is the field name for USD volume
	USDVolume //usd_volume
	//Country is the field name for country name
	Country //country
	//DstAddress is the field name for Destination Address
	DstAddress //dst_addr
	//SrcAddrss is the field name for srouce address
	SrcAddress //src_addr
)

//heatMapFields translates the stringer of token heatmap fields into its enumerated form
var heatMapFields = map[string]FieldName{
	"time":         Time,
	"eth_volume":   ETHVolume,
	"token_volume": TokenVolume,
	"usd_volume":   USDVolume,
	"country":      Country,
	"dst_addr":     DstAddress,
	"src_addr":     SrcAddress,
}
