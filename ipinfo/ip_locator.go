package ipinfo

import (
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/KyberNetwork/reserve-stats/util"
	geoip2 "github.com/oschwald/geoip2-golang"
	"go.uber.org/zap"
)

const (
	geoDBFile = "GeoLite2-Country.mmdb"
	url       = "https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.mmdb.gz"
)

func getGeoDBFile(sugar *zap.SugaredLogger, dbPath string) error {
	f, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	sugar.Debug("Begin download dbfile from url", url)
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	r, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, r)
	sugar.Debug("Finish download dbfile.")
	return err
}

// Locator is a resolver that query data of IP from MaxMind's GeoLite2 database.
type Locator struct {
	r     *geoip2.Reader
	sugar *zap.SugaredLogger
}

// NewLocator returns an instance of ipLocator.
func NewLocator(sugar *zap.SugaredLogger) (*Locator, error) {
	dbPath := path.Join(util.CurrentDir(), geoDBFile)
	err := getGeoDBFile(sugar, dbPath)
	if err != nil {
		return nil, err
	}
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
		return "", fmt.Errorf("%s is invalid ip", ip)
	}
	record, err := il.r.Country(IPParsed)
	if err != nil {
		il.sugar.Infow("failed to query data from geo-database!", "error", err)
		return "", err
	}

	country := record.Country.IsoCode //iso code of country
	if country == "" {
		il.sugar.Debug("Could not find country code of ", ip)
	}
	return country, nil
}
