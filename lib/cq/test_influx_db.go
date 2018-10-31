package cq

import (
	"fmt"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	testInfluxAddress = "http://127.0.0.1:8086"
	testDBName        = "test_cq"
	testMSName        = "test_measurement"
	nTestRecord       = 100
	recordInterval    = 1000
	tags              = "abc"
	timePrecision     = "ms"
)

//setupTestInfluxClient return a http influxClient and create test Database
func setupTestInfluxClient() (*client.Client, error) {
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: testInfluxAddress,
	})
	if err != nil {
		return nil, err
	}
	if _, err := queryDB(influxClient, fmt.Sprintf("CREATE DATABASE %s", testDBName)); err != nil {
		return nil, err
	}
	if err := addTestData(influxClient); err != nil {
		return nil, err
	}
	return &influxClient, nil
}

// addTestData will systematically add test Data to fulfill how many records are needed and interval between them
func addTestData(c client.Client) error {
	bp, err := client.NewBatchPoints(
		client.BatchPointsConfig{
			Database:  testDBName,
			Precision: timePrecision,
		},
	)
	if err != nil {
		return err
	}
	ts := time.Now()
	// The amount is inserted so that every interval, it guaranteed to yield a different sum
	for i := 0; i < nTestRecord; i++ {
		tags := map[string]string{
			"nameTag": string(tags[i%len(tags)]),
		}
		fields := map[string]interface{}{
			"amount": (i%7 + 1),
		}
		pt, err := client.NewPoint(testMSName, tags, fields, ts)
		if err != nil {
			return err
		}
		ts = ts.Add(time.Second)
		bp.AddPoint(pt)
	}
	return c.Write(bp)
}

// queryDB convenience function to query the database
func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: testDBName,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
