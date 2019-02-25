package integrationvolume

// FieldName define a list of field names for a token heatmap record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for time
	Time FieldName = iota //time
	//KyberSwapVolume is the field name for eth volume from KyberSwap
	KyberSwapVolume //kyber_swap_volume
	//NonKyberSwapVolume is the field name for eth volume from other app than KyberSwap
	NonKyberSwapVolume //non_kyber_swap_volume
)

//integrationVolumeFields translates the stringer of token integration volume fields into its enumerated form
var integrationVolumeFields = map[string]FieldName{
	"time":                  Time,
	"kyber_swap_volume":     KyberSwapVolume,
	"non_kyber_swap_volume": NonKyberSwapVolume,
}
