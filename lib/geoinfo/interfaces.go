package geoinfo

// Interface represents a client o interact with Geoinfo APIs.
type Interface interface {
	GetTxInfo(tx string) (string, string, error)
}
