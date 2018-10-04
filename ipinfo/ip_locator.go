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

	"github.com/oschwald/geoip2-golang"
	"go.uber.org/zap"
)

const (
	geoDBFile = "GeoLite2-Country.mmdb"
	url       = "https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.mmdb.gz"
)

func getGeoDBFile(sugar *zap.SugaredLogger, dbPath string) error {
	const timeout = time.Minute * 5

	f, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if os.IsExist(err) {
		sugar.Debugw("db file already exists", "db_path", dbPath)
		return nil
	} else if err != nil {
		return err
	}
	defer f.Close()

	sugar.Debugw("begin downloading db file from url", "url", url)
	client := &http.Client{Timeout: timeout}
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
	sugar.Debugw("finished downloading db file to", "db_path", dbPath)
	return err
}

// Locator is a resolver that query data of IP from MaxMind's GeoLite2 database.
type Locator struct {
	r     *geoip2.Reader
	sugar *zap.SugaredLogger
}

// NewLocator returns an instance of ipLocator.
func NewLocator(sugar *zap.SugaredLogger, dataDir string) (*Locator, error) {
	dbPath := path.Join(dataDir, geoDBFile)
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
	ipParsed := net.ParseIP(ip)
	if ipParsed == nil {
		return "", fmt.Errorf("%s is invalid ip", ip)
	}
	record, err := il.r.Country(ipParsed)
	if err != nil {
		il.sugar.Infow("failed to query data from geo-database!", "error", err)
		return "", err
	}

	country := record.Country.IsoCode //iso code of country
	if country == "" {
		il.sugar.Debugw("could not find country code of given IP", "ip", ip)
	}
	return country, nil
}
