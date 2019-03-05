package influxdb

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	influxdbEndpointFlag = "influxdb-endpoint"
	influxdbUsernameFlag = "influxdb-username"
	influxdbPasswordFlag = "influxdb-password"

	influxdbEndpointDefaultValue = "http://127.0.0.1:8086"
)

// NewCliFlags returns cli flags to configure a core client.
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   influxdbEndpointFlag,
			Usage:  "Endpoint to influxdb",
			EnvVar: "INFLUXDB_ENDPOINT",
			Value:  influxdbEndpointDefaultValue,
		},
		cli.StringFlag{
			Name:   influxdbUsernameFlag,
			Usage:  "Influxdb user",
			EnvVar: "INFLUXDB_USERNAME",
		},
		cli.StringFlag{
			Name:   influxdbPasswordFlag,
			Usage:  "Influxdb password",
			EnvVar: "INFLUXDB_PASSWORD",
		},
	}
}

// NewClientFromContext returns new core client from cli flags.
func NewClientFromContext(c *cli.Context) (client.Client, error) {
	var (
		err          error
		influxClient client.Client
	)

	endpoint := c.String(influxdbEndpointFlag)
	err = validation.Validate(endpoint,
		validation.Required,
		is.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid influxdb endpoint: %s", err.Error())
	}

	username := c.String(influxdbUsernameFlag)
	password := c.String(influxdbPasswordFlag)

	influxClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     endpoint,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	// ping to make sure connection success
	if _, _, err = influxClient.Ping(5 * time.Second); err != nil {
		return nil, err
	}

	return influxClient, nil
}
