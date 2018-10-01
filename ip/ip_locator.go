package ip

import (
	"fmt"
	"net"
	"path"

	"github.com/KyberNetwork/reserve-stats/util"
	geoip2 "github.com/oschwald/geoip2-golang"
	"go.uber.org/zap"
)

const geoDBFile = "GeoLite2-Country.mmdb"

// Locator is a resolver that query data of IP from MaxMind's GeoLite2 database.
type Locator struct {
	r     *geoip2.Reader
	sugar *zap.SugaredLogger
}

// NewLocator returns an instance of ipLocator.
func NewLocator(sugar *zap.SugaredLogger) (*Locator, error) {
	dbPath := path.Join(util.CurrentDir(), geoDBFile)
	r, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &Locator{
		sugar: sugar,
		r:     r,
	}, nil
}

// IPToCountry returns the country of given IP address.
func (il *Locator) IPToCountry(ip string) (string, error) {
	IPParsed := net.ParseIP(ip)
	if IPParsed == nil {
		return "", fmt.Errorf("invalid ip %s", ip)
	}
	record, err := il.r.Country(IPParsed)
	if err != nil {
		il.sugar.Infow("failed to query data from geo-database!")
		return "", err
	}

	country := record.Country.IsoCode //iso code of country
	if country == "" {
		return "", fmt.Errorf("Can't find country of the given ip: %s", ip)
	}
	return country, nil
}
