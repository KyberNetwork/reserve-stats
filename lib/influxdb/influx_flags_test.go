package influxdb

import (
	"os"
	"testing"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
)

const (
	//influxEndpoint correct is port 8086
	influxEndpointNotExist = "http://127.0.0.1:8087"
	influxdbUsername       = "INFLUXDB_USERNAME"
	influxdbPassword       = "INFLUXDB_PASSWORD"
)

func TestNewInfluxDB(t *testing.T) {
	username := os.Getenv(influxdbUsername)
	password := os.Getenv(influxdbPassword)

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxEndpointNotExist,
		Username: username,
		Password: password,
	})
	assert.Nil(t, err)
	_, _, err = influxClient.Ping(5 * time.Second)
	assert.NotNil(t, err)
}
