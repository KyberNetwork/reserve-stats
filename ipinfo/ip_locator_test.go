package ipinfo

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
)

func createIPLocator() (*Locator, error) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()
	return NewLocator(sugar, "test/")
}

//TestValidIP test when input is a valid IP of US
func TestValidIP(t *testing.T) {
	l, err := createIPLocator()
	if err != nil {
		t.Error("Could not create Locator", "error", err.Error())
	}
	ip := "81.2.69.142"
	country, err := l.IPToCountry(ip)
	if err != nil {
		t.Error("Get unexpected error when call IPTOCountry", "error", err.Error())
	}
	if country != "GB" {
		t.Error("Get location of ip was incorrect", "ip", ip, "result", country, "expected", "GB")
	}
}

//TestValidIP test when input is an invalid IP
func TestInvalidIP(t *testing.T) {
	l, err := createIPLocator()
	if err != nil {
		t.Error("Could not create Locator", "error", err.Error())
	}
	ip := "22"
	_, err = l.IPToCountry(ip)
	if err.Error() != fmt.Sprintf("%s is invalid ip", ip) {
		t.Error("Get unexpected error when call IPTOCountry", "error", err.Error())
	}
}

//TestNotFoundIP test when input ip could not be locate
func TestNotFoundIP(t *testing.T) {
	l, err := createIPLocator()
	if err != nil {
		t.Error("Could not create Locator", "error", err.Error())
	}
	ip := "192.168.0.1"
	country, err := l.IPToCountry(ip)
	if err != nil {
		t.Error("Get unexpected error when call IPTOCountry", "error", err.Error())
	}
	if country != "" {
		t.Error("Get location of ip was incorrect", "ip", ip, "result", country, "expected", "\"\"")
	}
}
