package ipinfo

import (
	"fmt"
	"os"
	"testing"

	"go.uber.org/zap"
)

func createIPLocator() (*Locator, error) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()
	return NewLocator(sugar, ".")
}

//TestValidIP test when input is a valid IP of US
func TestValidIP(t *testing.T) {
	os.Remove("GeoLite2-Country.mmdb")
	l, err := createIPLocator()
	if err != nil {
		t.Errorf(err.Error())
	}
	ip := "8.8.8.8"
	country, err := l.IPToCountry(ip)
	if err != nil {
		t.Errorf(err.Error())
	}
	if country != "US" {
		t.Errorf("Get location of ip %s was incorrect, got %s, want %s", ip, country, "US")
	}
}

//TestValidIP test when input is an invalid IP
func TestInvalidIP(t *testing.T) {
	l, err := createIPLocator()
	if err != nil {
		t.Errorf(err.Error())
	}
	ip := "22"
	_, err = l.IPToCountry(ip)
	if err.Error() != fmt.Sprintf("%s is invalid ip", ip) {
		t.Errorf(err.Error())
	}
}

//TestNotFoundIP test when input ip could not be locate
func TestNotFoundIP(t *testing.T) {
	l, err := createIPLocator()
	if err != nil {
		t.Errorf(err.Error())
	}
	ip := "192.168.0.1"
	country, err := l.IPToCountry(ip)
	if err != nil {
		t.Errorf(err.Error())
	}
	if country != "" {
		t.Errorf("Get location of ip %s was incorrect, got %s, want %s", ip, country, "US")
	}
}
