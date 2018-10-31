package cq

import (
	"fmt"
	"testing"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewContinuousQuery(t *testing.T) {
	var tests = []struct {
		testName              string
		name                  string
		database              string
		resampleEveryInterval string
		resampleForInterval   string
		query                 string
		timeInterval          string
		offsetIntervals       []string
		queries               []string
	}{
		{
			testName:              "simple continuous query",
			name:                  "test_cq",
			database:              "test_db",
			resampleEveryInterval: "",
			resampleForInterval:   "",
			query:                 `SELECT * FROM super_database`,
			timeInterval:          "1h",
			offsetIntervals:       nil,
			queries:               []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" BEGIN SELECT * FROM super_database GROUP BY time(1h) END`},
		},
		{
			testName:              "continuous query with resample every interval",
			name:                  "test_cq",
			database:              "test_db",
			resampleEveryInterval: "2h",
			query:                 `SELECT * FROM super_database`,
			timeInterval:          "1h",
			queries:               []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 2h BEGIN SELECT * FROM super_database GROUP BY time(1h) END`},
		},
		{
			testName:            "continuous query with resample for interval",
			name:                "test_cq",
			database:            "test_db",
			resampleForInterval: "2h",
			query:               `SELECT * FROM super_database`,
			timeInterval:        "1h",
			queries:             []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE FOR 2h BEGIN SELECT * FROM super_database GROUP BY time(1h) END`},
		},
		{
			testName:              "continuous query with both resample every and for intervals",
			name:                  "test_cq",
			database:              "test_db",
			resampleEveryInterval: "1h",
			resampleForInterval:   "2h",
			query:                 `SELECT * FROM super_database`,
			timeInterval:          "1h",
			queries:               []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY time(1h) END`},
		},
		{
			testName:              "continuous query with group by in query clause",
			name:                  "test_cq",
			database:              "test_db",
			resampleEveryInterval: "1h",
			resampleForInterval:   "2h",
			query:                 `SELECT * FROM super_database GROUP BY "email"`,
			timeInterval:          "1h",
			queries:               []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h) END`},
		},
		{
			testName:              "continuous query with one offset interval",
			name:                  "test_cq",
			database:              "test_db",
			resampleEveryInterval: "1h",
			resampleForInterval:   "2h",
			query:                 `SELECT * FROM super_database GROUP BY "email"`,
			timeInterval:          "1h",
			offsetIntervals:       []string{"30m"},
			queries:               []string{`CREATE CONTINUOUS QUERY "test_cq_30m" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h,30m) END`},
		},
		{
			testName:              "continuous query with multiple offset intervals",
			name:                  "test_cq",
			database:              "test_db",
			resampleEveryInterval: "1h",
			resampleForInterval:   "2h",
			query:                 `SELECT * FROM super_database GROUP BY "email"`,
			timeInterval:          "1h",
			offsetIntervals:       []string{"10m", "20m"},
			queries: []string{
				`CREATE CONTINUOUS QUERY "test_cq_10m" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h,10m) END`,
				`CREATE CONTINUOUS QUERY "test_cq_20m" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h,20m) END`,
			},
		},
	}

	for _, tc := range tests {
		cq, err := NewContinuousQuery(
			tc.name,
			tc.database,
			tc.resampleEveryInterval,
			tc.resampleForInterval,
			tc.query,
			tc.timeInterval,
			tc.offsetIntervals,
		)
		require.NoError(t, err, tc.testName)
		assert.Equal(t, cq.queries, tc.queries, tc.testName)
	}
}

func setupTest() (client.Client, *zap.SugaredLogger, error) {
	c, err := setupTestInfluxClient()
	if err != nil {
		return nil, nil, err
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, nil, err
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	return c, sugar, nil
}

func TestContinuousQuery_Deploy(t *testing.T) {
	// TODO: init test eth client, create database, insert sample data if needed
	c, sugar, err := setupTest()
	//tear down
	defer func() {
		if _, err := queryDB(c, fmt.Sprintf("DROP DATABASE %s", testDBName), testDBName); err != nil {
			t.Error(err)
		}
	}()

	cq, err := NewContinuousQuery(
		"test_cq",
		testDBName,
		"1h",
		"2h",
		`SELECT * INTO test_queries FROM super_database GROUP BY "email"`,
		"1h",
		[]string{"10m", "20m"},
	)
	require.NoError(t, err)

	//deploy and test cqs
	assert.NoError(t, cq.Deploy(c, sugar))
	cqs, err := cq.GetCurrentCQs(c, sugar)
	require.NoError(t, err)
	var (
		expectedCqs = map[string]string{
			"test_cq_10m": "CREATE CONTINUOUS QUERY test_cq_10m ON test_db RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * INTO test_db.autogen.test_queries_10m FROM test_db.autogen.super_database GROUP BY email, time(1h, 10m) END",
			"test_cq_20m": "CREATE CONTINUOUS QUERY test_cq_20m ON test_db RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * INTO test_db.autogen.test_queries_20m FROM test_db.autogen.super_database GROUP BY email, time(1h, 20m) END",
		}
	)
	for cqName, cq := range expectedCqs {
		resCq, ok := cqs[cqName]
		if !ok {
			t.Errorf("Result doesn't contain cq %s", cqName)
		}
		assert.Equal(t, cq, resCq)
	}

	// TODO: make sure that deploy can be successfully called second time

	cq, err = NewContinuousQuery(
		"test_cq",
		"test_db_2",
		"1h",
		"2h",
		`SELECT * FROM super_database GROUP BY "email"`,
		"1h",
		[]string{"10m", "20m"},
	)
	require.NoError(t, err)
	// TODO: makes sure that cqs database changed form test_db --> test_db_2

	cq, err = NewContinuousQuery(
		"test_cq",
		"test_db",
		"3h",
		"2h",
		`SELECT * FROM super_database GROUP BY "email"`,
		"1h",
		[]string{"10m", "20m"},
	)
	require.NoError(t, err)
	// TODO: makes sure that cqs resample every interval changed from 1h --> 3h

	cq, err = NewContinuousQuery(
		"test_cq",
		"test_db",
		"1h",
		"4h",
		`SELECT * FROM super_database GROUP BY "email"`,
		"1h",
		[]string{"10m", "20m"},
	)
	require.NoError(t, err)
	// TODO: makes sure that cqs resample for interval changed from 2h --> 4h

	cq, err = NewContinuousQuery(
		"test_cq",
		"test_db",
		"1h",
		"2h",
		`SELECT * FROM super_database GROUP BY "username"`,
		"1h",
		[]string{"10m", "20m"},
	)
	require.NoError(t, err)
	// TODO: makes sure that cqs query updated

	cq, err = NewContinuousQuery(
		"test_cq",
		"test_db",
		"1h",
		"2h",
		`SELECT * FROM super_database GROUP BY "email"`,
		"2h",
		[]string{"10m", "20m"},
	)
	require.NoError(t, err)
	// TODO: makes sure that cqs time interval changed from 1h --> 2h

	cq, err = NewContinuousQuery(
		"test_cq",
		"test_db",
		"1h",
		"2h",
		`SELECT * FROM super_database GROUP BY "email"`,
		"1h",
		[]string{"15m", "25m"},
	)
	require.NoError(t, err)
	// TODO: makes sure that cqs offset interval changed from 10, 20 --> 15, 25

	// TODO: refactors above tests to table test format
}

func TestContinuousQuery_Execute(t *testing.T) {
	c, sugar, err := setupTest()
	//tear down
	defer func() {
		if _, err := queryDB(c, fmt.Sprintf("DROP DATABASE %s", testDBName), testDBName); err != nil {
			t.Error(err)
		}
	}()
	cq, err := NewContinuousQuery(
		"first_test",
		testDBName,
		"1h",
		"2h",
		`SELECT SUM(amount) AS volume INTO test_aggregate FROM test_measurement GROUP BY "nameTag"`,
		"1m",
		[]string{"10s", "15s"},
	)

	err = cq.Execute(c, sugar)
	require.NoError(t, err)

	resp, err := cq.queryDB(c, "SHOW measurements")
	require.NoError(t, err)
	if len(resp[0].Series) == 0 {
		t.Error("expect valid result, got empty result")
	}

	var (
		expectedMsms = map[string]bool{
			"test_aggregate_10s": false,
			"test_aggregate_15s": false,
		}
	)
	for _, v := range resp[0].Series[0].Values {
		x, ok := v[0].(string)
		if !ok {
			t.Errorf("invalid value from result %v", v[0])
		}
		expectedMsms[x] = true
	}
	for k, v := range expectedMsms {
		if !v {
			t.Errorf("result doesn't contain measurement %s", k)
		}
	}
}
