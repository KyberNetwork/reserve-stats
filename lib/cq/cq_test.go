package cq

import (
	"fmt"
	"testing"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
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
			queries: []string{
				`CREATE CONTINUOUS QUERY "test_cq_minus30m" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h,30m) END`,
				`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h) END`,
			},
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
				`CREATE CONTINUOUS QUERY "test_cq_minus10m" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h,10m) END`,
				`CREATE CONTINUOUS QUERY "test_cq_minus20m" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h,20m) END`,
				`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h) END`,
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
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	return c, sugar, nil
}

func TestContinuousQuery_Deploy(t *testing.T) {
	c, sugar, err := setupTest()
	require.NoError(t, err)
	//tear down
	defer func() {
		if _, err := influxdb.QueryDB(c, fmt.Sprintf("DROP DATABASE %s", testDBName), testDBName); err != nil {
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
			"test_cq_minus10m": "CREATE CONTINUOUS QUERY test_cq_minus10m ON test_db RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * INTO test_db.autogen.test_queries_minus10m FROM test_db.autogen.super_database GROUP BY email, time(1h, 10m) END",
			"test_cq_minus20m": "CREATE CONTINUOUS QUERY test_cq_minus20m ON test_db RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * INTO test_db.autogen.test_queries_minus20m FROM test_db.autogen.super_database GROUP BY email, time(1h, 20m) END",
		}
	)
	for cqName, cq := range expectedCqs {
		resCq, ok := cqs[cqName]
		if !ok {
			t.Errorf("Result doesn't contain cq %s", cqName)
		}
		assert.Equal(t, cq, resCq)
	}

	//test drop CQ
	assert.NoError(t, cq.Drop(c, sugar))
	cqs, err = cq.GetCurrentCQs(c, sugar)
	require.NoError(t, err)
	for cqName := range expectedCqs {
		if _, ok := cqs[cqName]; ok {
			t.Errorf("expect cq %s to be dropped, yet it is still there", cqName)
		}
	}

	_, err = NewContinuousQuery(
		"test_cq",
		"test_db_2",
		"1h",
		"2h",
		`SELECT * FROM super_database GROUP BY "email"`,
		"1h",
		[]string{"10m", "20m"},
	)
	require.NoError(t, err)

	_, err = NewContinuousQuery(
		"test_cq",
		"test_db",
		"3h",
		"2h",
		`SELECT * FROM super_database GROUP BY "email"`,
		"1h",
		[]string{"10m", "20m"},
	)
	require.NoError(t, err)

	_, err = NewContinuousQuery(
		"test_cq",
		"test_db",
		"1h",
		"4h",
		`SELECT * FROM super_database GROUP BY "email"`,
		"1h",
		[]string{"10m", "20m"},
	)
	require.NoError(t, err)

	_, err = NewContinuousQuery(
		"test_cq",
		"test_db",
		"1h",
		"2h",
		`SELECT * FROM super_database GROUP BY "username"`,
		"1h",
		[]string{"10m", "20m"},
	)
	require.NoError(t, err)

	_, err = NewContinuousQuery(
		"test_cq",
		"test_db",
		"1h",
		"2h",
		`SELECT * FROM super_database GROUP BY "email"`,
		"2h",
		[]string{"10m", "20m"},
	)
	require.NoError(t, err)

	_, err = NewContinuousQuery(
		"test_cq",
		"test_db",
		"1h",
		"2h",
		`SELECT * FROM super_database GROUP BY "email"`,
		"1h",
		[]string{"15m", "25m"},
	)
	require.NoError(t, err)
}

func TestContinuousQuery_Execute(t *testing.T) {
	c, sugar, err := setupTest()
	require.NoError(t, err)
	//tear down
	defer func() {
		if _, err := influxdb.QueryDB(c, fmt.Sprintf("DROP DATABASE %s", testDBName), testDBName); err != nil {
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
	require.NoError(t, err)

	err = cq.Execute(c, sugar)
	require.NoError(t, err)

	resp, err := influxdb.QueryDB(c, "SHOW measurements", cq.Database)
	require.NoError(t, err)
	if len(resp[0].Series) == 0 {
		t.Error("expect valid result, got empty result")
	}
	var (
		expectedMsms = map[string]bool{
			"test_aggregate_minus10s": false,
			"test_aggregate_minus15s": false,
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

func TestHasGroupBy(t *testing.T) {
	var tests = []struct {
		query      string
		hasGroupBy bool
	}{
		{
			query:      "SELECT * FROM super_database",
			hasGroupBy: false,
		},
		{
			query:      "SELECT * FROM super_database GROUP BY time(1h)",
			hasGroupBy: true,
		},
	}

	for _, tc := range tests {
		assert.True(t, hasGroupBy(tc.query) == tc.hasGroupBy)
	}
}
