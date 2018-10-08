package ipinfo

import (
	"net"
	"testing"

	"go.uber.org/zap"
)

func newTestLocator() (*Locator, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	return NewLocator(sugar, "testdata")
}

//TestValidIP test when input is a valid IP of US
func TestValidIP(t *testing.T) {
	l, err := newTestLocator()
	if err != nil {
		t.Error("Could not create Locator", "error", err.Error())
	}
	ip := net.ParseIP("81.2.69.142")
	country, err := l.IPToCountry(ip)
	if err != nil {
		t.Error("Get unexpected error when call IPTOCountry", "error", err.Error())
	}
	if country != "GB" {
		t.Error("Get location of ip was incorrect", "ip", ip, "result", country, "expected", "GB")
	}
}

//TestNotFoundIP test when input ip could not be locate
func TestNotFoundIP(t *testing.T) {
	l, err := newTestLocator()
	if err != nil {
		t.Error("Could not create Locator", "error", err.Error())
	}
	ip := net.ParseIP("192.168.0.1")
	country, err := l.IPToCountry(ip)
	if err != nil {
		t.Error("Get unexpected error when call IPTOCountry", "error", err.Error())
	}
	if country != "" {
		t.Error("Get location of ip was incorrect", "ip", ip, "result", country, "expected", "\"\"")
	}
}
