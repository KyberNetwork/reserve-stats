package influxdb

import (
	"github.com/influxdata/influxdb/client/v2"
)

//QueryDB is function to run a command to influx db
func QueryDB(clnt client.Client, cmd, database string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: database,
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
